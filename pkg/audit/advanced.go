package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

// AuditEvent represents a comprehensive audit event
type AuditEvent struct {
	ID              string                 `json:"id"`
	Timestamp       time.Time              `json:"timestamp"`
	EventType       string                 `json:"eventType"` // create, read, update, delete, login, logout, export
	ResourceType    string                 `json:"resourceType"`
	ResourceID      string                 `json:"resourceId,omitempty"`
	Action          string                 `json:"action"`
	Outcome         string                 `json:"outcome"` // success, failure, partial
	User            UserInfo               `json:"user"`
	Tenant          TenantInfo             `json:"tenant"`
	Client          ClientInfo             `json:"client"`
	Request         RequestInfo            `json:"request"`
	Response        ResponseInfo           `json:"response"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Changes         *ResourceChanges       `json:"changes,omitempty"`
	Classification  string                 `json:"classification"` // phi, restricted, internal, public
	ComplianceFlags []string               `json:"complianceFlags,omitempty"`
}

// UserInfo contains user information
type UserInfo struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email,omitempty"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions,omitempty"`
}

// TenantInfo contains tenant information
type TenantInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ClientInfo contains client information
type ClientInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"userAgent"`
	DeviceID  string `json:"deviceId,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
}

// RequestInfo contains request information
type RequestInfo struct {
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Path      string            `json:"path"`
	Query     map[string]string `json:"query,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	BodySize  int               `json:"bodySize,omitempty"`
	RequestID string            `json:"requestId"`
}

// ResponseInfo contains response information
type ResponseInfo struct {
	StatusCode int   `json:"statusCode"`
	BodySize   int   `json:"bodySize,omitempty"`
	Duration   int64 `json:"durationMs"`
}

// ResourceChanges tracks resource modifications
type ResourceChanges struct {
	Before map[string]interface{} `json:"before,omitempty"`
	After  map[string]interface{} `json:"after,omitempty"`
	Fields []string               `json:"fields,omitempty"`
}

// AuditLogger defines the audit logging interface
type AuditLogger interface {
	Log(ctx context.Context, event *AuditEvent) error
	Query(ctx context.Context, query AuditQuery) ([]*AuditEvent, error)
	Export(ctx context.Context, options ExportOptions) ([]byte, error)
}

// AuditQuery contains query parameters for audit events
type AuditQuery struct {
	EventTypes    []string
	ResourceTypes []string
	UserID        string
	TenantID      string
	From          time.Time
	To            time.Time
	Outcome       string
	Limit         int
	Offset        int
}

// ExportOptions contains options for exporting audit logs
type ExportOptions struct {
	Format     string // json, csv, fhir-audit
	From       time.Time
	To         time.Time
	TenantID   string
	EventTypes []string
}

// AdvancedAuditLogger provides comprehensive audit logging
type AdvancedAuditLogger struct {
	store  AuditStore
	config AuditConfig
}

// AuditStore defines the storage interface for audit events
type AuditStore interface {
	Store(ctx context.Context, event *AuditEvent) error
	Query(ctx context.Context, query AuditQuery) ([]*AuditEvent, error)
	Export(ctx context.Context, options ExportOptions) ([]byte, error)
}

// AuditConfig contains audit configuration
type AuditConfig struct {
	Enabled            bool
	LogPHI             bool
	AsyncLogging       bool
	BufferSize         int
	RetentionDays      int
	ExportEnabled      bool
	AlertOnAnomalies   bool
	ComplianceStandard string // HIPAA, GDPR, Bangladesh_DGHS
}

// NewAdvancedAuditLogger creates a new advanced audit logger
func NewAdvancedAuditLogger(store AuditStore, config AuditConfig) *AdvancedAuditLogger {
	return &AdvancedAuditLogger{
		store:  store,
		config: config,
	}
}

// Log logs an audit event
func (al *AdvancedAuditLogger) Log(ctx context.Context, event *AuditEvent) error {
	if !al.config.Enabled {
		return nil
	}

	if event.ID == "" {
		event.ID = generateEventID()
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Check PHI logging restrictions
	if event.Classification == "phi" && !al.config.LogPHI {
		// Mask sensitive data
		event.Data = map[string]interface{}{"masked": true}
		event.Changes = nil
	}

	// Add compliance flags
	event.ComplianceFlags = al.getComplianceFlags(event)

	return al.store.Store(ctx, event)
}

// Query queries audit events
func (al *AdvancedAuditLogger) Query(ctx context.Context, query AuditQuery) ([]*AuditEvent, error) {
	return al.store.Query(ctx, query)
}

// Export exports audit events
func (al *AdvancedAuditLogger) Export(ctx context.Context, options ExportOptions) ([]byte, error) {
	if !al.config.ExportEnabled {
		return nil, fmt.Errorf("export is not enabled")
	}

	return al.store.Export(ctx, options)
}

// LogResourceAccess logs a resource access event
func (al *AdvancedAuditLogger) LogResourceAccess(ctx context.Context, user UserInfo, tenant TenantInfo, resourceType, resourceID, action string, success bool) {
	event := &AuditEvent{
		EventType:      "read",
		ResourceType:   resourceType,
		ResourceID:     resourceID,
		Action:         action,
		Outcome:        map[bool]string{true: "success", false: "failure"}[success],
		User:           user,
		Tenant:         tenant,
		Classification: "phi",
	}

	if err := al.Log(ctx, event); err != nil {
		log.Errorf("Failed to log audit event: %v", err)
	}
}

// LogResourceChange logs a resource modification event
func (al *AdvancedAuditLogger) LogResourceChange(ctx context.Context, user UserInfo, tenant TenantInfo, resourceType, resourceID, action string, before, after map[string]interface{}, success bool) {
	changes := &ResourceChanges{
		Before: before,
		After:  after,
	}

	// Identify changed fields
	if before != nil && after != nil {
		for key := range after {
			if before[key] != after[key] {
				changes.Fields = append(changes.Fields, key)
			}
		}
	}

	event := &AuditEvent{
		EventType:      action,
		ResourceType:   resourceType,
		ResourceID:     resourceID,
		Action:         action,
		Outcome:        map[bool]string{true: "success", false: "failure"}[success],
		User:           user,
		Tenant:         tenant,
		Changes:        changes,
		Classification: "phi",
	}

	if err := al.Log(ctx, event); err != nil {
		log.Errorf("Failed to log audit event: %v", err)
	}
}

// LogAuthentication logs authentication events
func (al *AdvancedAuditLogger) LogAuthentication(ctx context.Context, user UserInfo, tenant TenantInfo, client ClientInfo, action string, success bool, details map[string]interface{}) {
	event := &AuditEvent{
		EventType:      action,
		ResourceType:   "User",
		ResourceID:     user.ID,
		Action:         action,
		Outcome:        map[bool]string{true: "success", false: "failure"}[success],
		User:           user,
		Tenant:         tenant,
		Client:         client,
		Data:           details,
		Classification: "restricted",
	}

	if err := al.Log(ctx, event); err != nil {
		log.Errorf("Failed to log audit event: %v", err)
	}
}

// getComplianceFlags returns compliance flags for an event
func (al *AdvancedAuditLogger) getComplianceFlags(event *AuditEvent) []string {
	flags := make([]string, 0)

	switch al.config.ComplianceStandard {
	case "HIPAA":
		if event.Classification == "phi" {
			flags = append(flags, "hipaa-protected")
		}
		if event.ResourceType == "Patient" || event.ResourceType == "Observation" {
			flags = append(flags, "phi-access")
		}

	case "GDPR":
		if event.EventType == "export" {
			flags = append(flags, "data-export")
		}
		if event.EventType == "delete" {
			flags = append(flags, "right-to-erasure")
		}

	case "Bangladesh_DGHS":
		if event.ResourceType == "Patient" {
			flags = append(flags, "dghs-regulated")
		}
	}

	return flags
}

// AuditMiddleware creates HTTP middleware for audit logging
type AuditMiddleware struct {
	logger *AdvancedAuditLogger
}

// NewAuditMiddleware creates new audit middleware
func NewAuditMiddleware(logger *AdvancedAuditLogger) *AuditMiddleware {
	return &AuditMiddleware{logger: logger}
}

// Handler returns the HTTP handler
func (am *AuditMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response wrapper
		rw := &auditResponseWriter{ResponseWriter: w, statusCode: 200}

		// Extract user info from context
		user := extractUserInfo(r)
		tenant := extractTenantInfo(r)
		client := extractClientInfo(r)

		// Process request
		next.ServeHTTP(rw, r)

		// Build audit event
		duration := time.Since(start)

		event := &AuditEvent{
			EventType:    getEventType(r.Method),
			ResourceType: chi.URLParam(r, "resourceType"),
			ResourceID:   chi.URLParam(r, "resourceID"),
			Action:       r.Method,
			Outcome:      getOutcome(rw.statusCode),
			User:         user,
			Tenant:       tenant,
			Client:       client,
			Request: RequestInfo{
				Method:    r.Method,
				URL:       r.URL.String(),
				Path:      r.URL.Path,
				Query:     flattenQuery(r.URL.Query()),
				RequestID: r.Header.Get("X-Request-ID"),
			},
			Response: ResponseInfo{
				StatusCode: rw.statusCode,
				Duration:   duration.Milliseconds(),
			},
			Classification: "phi",
		}

		// Log event
		if err := am.logger.Log(r.Context(), event); err != nil {
			log.Errorf("Failed to log audit event: %v", err)
		}
	})
}

// auditResponseWriter wraps http.ResponseWriter to capture status code
type auditResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *auditResponseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// InMemoryAuditStore implements AuditStore with in-memory storage
type InMemoryAuditStore struct {
	events []*AuditEvent
	mu     sync.RWMutex
}

// NewInMemoryAuditStore creates a new in-memory audit store
func NewInMemoryAuditStore() *InMemoryAuditStore {
	return &InMemoryAuditStore{
		events: make([]*AuditEvent, 0),
	}
}

// Store implements AuditStore
func (s *InMemoryAuditStore) Store(ctx context.Context, event *AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)

	// Keep only last 10000 events in memory
	if len(s.events) > 10000 {
		s.events = s.events[len(s.events)-10000:]
	}

	return nil
}

// Query implements AuditStore
func (s *InMemoryAuditStore) Query(ctx context.Context, query AuditQuery) ([]*AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*AuditEvent

	for _, event := range s.events {
		// Apply filters
		if len(query.EventTypes) > 0 && !contains(query.EventTypes, event.EventType) {
			continue
		}
		if len(query.ResourceTypes) > 0 && !contains(query.ResourceTypes, event.ResourceType) {
			continue
		}
		if query.UserID != "" && event.User.ID != query.UserID {
			continue
		}
		if query.TenantID != "" && event.Tenant.ID != query.TenantID {
			continue
		}
		if !query.From.IsZero() && event.Timestamp.Before(query.From) {
			continue
		}
		if !query.To.IsZero() && event.Timestamp.After(query.To) {
			continue
		}
		if query.Outcome != "" && event.Outcome != query.Outcome {
			continue
		}

		results = append(results, event)
	}

	// Apply pagination
	if query.Offset >= len(results) {
		return []*AuditEvent{}, nil
	}

	end := query.Offset + query.Limit
	if end > len(results) || query.Limit == 0 {
		end = len(results)
	}

	return results[query.Offset:end], nil
}

// Export implements AuditStore
func (s *InMemoryAuditStore) Export(ctx context.Context, options ExportOptions) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []*AuditEvent

	for _, event := range s.events {
		// Apply filters
		if !options.From.IsZero() && event.Timestamp.Before(options.From) {
			continue
		}
		if !options.To.IsZero() && event.Timestamp.After(options.To) {
			continue
		}
		if options.TenantID != "" && event.Tenant.ID != options.TenantID {
			continue
		}
		if len(options.EventTypes) > 0 && !contains(options.EventTypes, event.EventType) {
			continue
		}

		events = append(events, event)
	}

	switch options.Format {
	case "json":
		return json.MarshalIndent(events, "", "  ")
	default:
		return json.MarshalIndent(events, "", "  ")
	}
}

// Helper functions

func generateEventID() string {
	return fmt.Sprintf("audit-%d", time.Now().UnixNano())
}

func getEventType(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

func getOutcome(statusCode int) string {
	if statusCode >= 200 && statusCode < 300 {
		return "success"
	} else if statusCode >= 400 && statusCode < 500 {
		return "failure"
	}
	return "partial"
}

func extractUserInfo(r *http.Request) UserInfo {
	// Extract from context or headers
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	return UserInfo{
		ID:       userID,
		Username: r.Header.Get("X-User-Name"),
		Role:     r.Header.Get("X-User-Role"),
	}
}

func extractTenantInfo(r *http.Request) TenantInfo {
	return TenantInfo{
		ID:   r.Header.Get("X-Tenant-ID"),
		Name: r.Header.Get("X-Tenant-Name"),
	}
}

func extractClientInfo(r *http.Request) ClientInfo {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ClientInfo{
		IP:        ip,
		UserAgent: r.UserAgent(),
		SessionID: r.Header.Get("X-Session-ID"),
	}
}

func flattenQuery(query map[string][]string) map[string]string {
	result := make(map[string]string)
	for k, v := range query {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
