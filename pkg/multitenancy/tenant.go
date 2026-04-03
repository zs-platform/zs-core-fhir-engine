package multitenancy

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// Tenant represents a healthcare organization/tenant
type Tenant struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"` // active, inactive, suspended
	Plan        string                 `json:"plan"`   // basic, professional, enterprise
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	Settings    TenantSettings         `json:"settings"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TenantSettings contains tenant-specific configuration
type TenantSettings struct {
	MaxUsers          int      `json:"maxUsers"`
	MaxPatients       int      `json:"maxPatients"`
	MaxStorageGB      int      `json:"maxStorageGB"`
	MaxConcurrentReqs int      `json:"maxConcurrentReqs"`
	AllowedRegions    []string `json:"allowedRegions,omitempty"`
	DefaultLanguage   string   `json:"defaultLanguage"`
	Timezone          string   `json:"timezone"`
	DateFormat        string   `json:"dateFormat"`
	EnableAudit       bool     `json:"enableAudit"`
	EnableAnalytics   bool     `json:"enableAnalytics"`
	EnableOffline     bool     `json:"enableOffline"`
	DataRetentionDays int      `json:"dataRetentionDays"`
}

// TenantContext holds tenant information in request context
type TenantContext struct {
	TenantID   string
	TenantName string
	Plan       string
	Settings   TenantSettings
}

// contextKey is a custom type for context keys
type contextKey string

const (
	// TenantContextKey is the key for tenant context in request context
	TenantContextKey contextKey = "tenant_context"
)

// TenantStore defines the storage interface for tenant operations
type TenantStore interface {
	CreateTenant(ctx context.Context, tenant *Tenant) error
	GetTenant(ctx context.Context, tenantID string) (*Tenant, error)
	UpdateTenant(ctx context.Context, tenant *Tenant) error
	DeleteTenant(ctx context.Context, tenantID string) error
	ListTenants(ctx context.Context, options ListOptions) ([]*Tenant, error)
	GetTenantByAPIKey(ctx context.Context, apiKey string) (*Tenant, error)
}

// ListOptions contains options for listing tenants
type ListOptions struct {
	Status string
	Plan   string
	Limit  int
	Offset int
}

// TenantManager manages tenant operations
type TenantManager struct {
	store TenantStore
}

// NewTenantManager creates a new tenant manager
func NewTenantManager(store TenantStore) *TenantManager {
	return &TenantManager{
		store: store,
	}
}

// CreateTenant creates a new tenant
func (tm *TenantManager) CreateTenant(ctx context.Context, tenant *Tenant) error {
	if tenant.ID == "" {
		return fmt.Errorf("tenant ID is required")
	}

	if tenant.Status == "" {
		tenant.Status = "active"
	}

	if tenant.Plan == "" {
		tenant.Plan = "basic"
	}

	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	// Set default settings based on plan
	tm.applyPlanDefaults(tenant)

	log.Infof("Creating tenant: %s (%s)", tenant.Name, tenant.ID)
	return tm.store.CreateTenant(ctx, tenant)
}

// GetTenant retrieves a tenant by ID
func (tm *TenantManager) GetTenant(ctx context.Context, tenantID string) (*Tenant, error) {
	return tm.store.GetTenant(ctx, tenantID)
}

// UpdateTenant updates an existing tenant
func (tm *TenantManager) UpdateTenant(ctx context.Context, tenant *Tenant) error {
	tenant.UpdatedAt = time.Now()
	log.Infof("Updating tenant: %s", tenant.ID)
	return tm.store.UpdateTenant(ctx, tenant)
}

// DeleteTenant deletes a tenant
func (tm *TenantManager) DeleteTenant(ctx context.Context, tenantID string) error {
	log.Infof("Deleting tenant: %s", tenantID)
	return tm.store.DeleteTenant(ctx, tenantID)
}

// ListTenants lists all tenants with optional filtering
func (tm *TenantManager) ListTenants(ctx context.Context, options ListOptions) ([]*Tenant, error) {
	return tm.store.ListTenants(ctx, options)
}

// ValidateTenantAccess checks if a tenant has access to a feature
func (tm *TenantManager) ValidateTenantAccess(ctx context.Context, tenantID string, feature string) error {
	tenant, err := tm.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	if tenant.Status != "active" {
		return fmt.Errorf("tenant is not active: %s", tenant.Status)
	}

	// Check feature access based on plan
	if !tm.isFeatureAllowed(tenant.Plan, feature) {
		return fmt.Errorf("feature %s not allowed for %s plan", feature, tenant.Plan)
	}

	return nil
}

// applyPlanDefaults applies default settings based on tenant plan
func (tm *TenantManager) applyPlanDefaults(tenant *Tenant) {
	switch tenant.Plan {
	case "basic":
		tenant.Settings.MaxUsers = 10
		tenant.Settings.MaxPatients = 1000
		tenant.Settings.MaxStorageGB = 10
		tenant.Settings.MaxConcurrentReqs = 50
		tenant.Settings.EnableAudit = true
		tenant.Settings.EnableAnalytics = false
		tenant.Settings.EnableOffline = false
		tenant.Settings.DataRetentionDays = 90

	case "professional":
		tenant.Settings.MaxUsers = 50
		tenant.Settings.MaxPatients = 10000
		tenant.Settings.MaxStorageGB = 100
		tenant.Settings.MaxConcurrentReqs = 200
		tenant.Settings.EnableAudit = true
		tenant.Settings.EnableAnalytics = true
		tenant.Settings.EnableOffline = true
		tenant.Settings.DataRetentionDays = 365

	case "enterprise":
		tenant.Settings.MaxUsers = 500
		tenant.Settings.MaxPatients = 100000
		tenant.Settings.MaxStorageGB = 1000
		tenant.Settings.MaxConcurrentReqs = 1000
		tenant.Settings.EnableAudit = true
		tenant.Settings.EnableAnalytics = true
		tenant.Settings.EnableOffline = true
		tenant.Settings.DataRetentionDays = 2555 // 7 years
	}

	// Set defaults if not specified
	if tenant.Settings.DefaultLanguage == "" {
		tenant.Settings.DefaultLanguage = "en"
	}
	if tenant.Settings.Timezone == "" {
		tenant.Settings.Timezone = "Asia/Dhaka"
	}
	if tenant.Settings.DateFormat == "" {
		tenant.Settings.DateFormat = "DD/MM/YYYY"
	}
}

// isFeatureAllowed checks if a feature is allowed for a plan
func (tm *TenantManager) isFeatureAllowed(plan, feature string) bool {
	featureMatrix := map[string]map[string]bool{
		"basic": {
			"patient_management": true,
			"observations":       true,
			"medications":        true,
			"appointments":       true,
			"basic_reports":      true,
			"audit_logging":      true,
			"analytics":          false,
			"offline_sync":       false,
			"multi_location":     false,
			"api_access":         false,
		},
		"professional": {
			"patient_management": true,
			"observations":       true,
			"medications":        true,
			"appointments":       true,
			"basic_reports":      true,
			"audit_logging":      true,
			"analytics":          true,
			"offline_sync":       true,
			"multi_location":     true,
			"api_access":         true,
		},
		"enterprise": {
			"patient_management":  true,
			"observations":        true,
			"medications":         true,
			"appointments":        true,
			"basic_reports":       true,
			"audit_logging":       true,
			"analytics":           true,
			"offline_sync":        true,
			"multi_location":      true,
			"api_access":          true,
			"custom_integrations": true,
			"dedicated_support":   true,
			"sla_guarantee":       true,
		},
	}

	if planFeatures, ok := featureMatrix[plan]; ok {
		return planFeatures[feature]
	}

	return false
}

// TenantMiddleware extracts tenant information from requests
func TenantMiddleware(tm *TenantManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract tenant ID from header
			tenantID := r.Header.Get("X-Tenant-ID")
			if tenantID == "" {
				// Try from query parameter
				tenantID = r.URL.Query().Get("tenant")
			}

			if tenantID == "" {
				http.Error(w, "Tenant ID required", http.StatusBadRequest)
				return
			}

			// Get tenant
			tenant, err := tm.GetTenant(r.Context(), tenantID)
			if err != nil {
				http.Error(w, "Invalid tenant", http.StatusUnauthorized)
				return
			}

			if tenant.Status != "active" {
				http.Error(w, "Tenant is not active", http.StatusForbidden)
				return
			}

			// Create tenant context
			tenantCtx := TenantContext{
				TenantID:   tenant.ID,
				TenantName: tenant.Name,
				Plan:       tenant.Plan,
				Settings:   tenant.Settings,
			}

			// Add to request context
			ctx := WithTenantContext(r.Context(), tenantCtx)

			// Continue with tenant context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithTenantContext adds tenant context to a context
func WithTenantContext(ctx context.Context, tenantCtx TenantContext) context.Context {
	return context.WithValue(ctx, TenantContextKey, tenantCtx)
}

// GetTenantContext retrieves tenant context from a context
func GetTenantContext(ctx context.Context) (TenantContext, bool) {
	tenantCtx, ok := ctx.Value(TenantContextKey).(TenantContext)
	return tenantCtx, ok
}

// GetTenantID retrieves tenant ID from a context
func GetTenantID(ctx context.Context) string {
	if tenantCtx, ok := GetTenantContext(ctx); ok {
		return tenantCtx.TenantID
	}
	return ""
}

// InMemoryTenantStore implements TenantStore with in-memory storage
type InMemoryTenantStore struct {
	tenants map[string]*Tenant
	apiKeys map[string]string // apiKey -> tenantID
	mu      sync.RWMutex
}

// NewInMemoryTenantStore creates a new in-memory tenant store
func NewInMemoryTenantStore() *InMemoryTenantStore {
	return &InMemoryTenantStore{
		tenants: make(map[string]*Tenant),
		apiKeys: make(map[string]string),
	}
}

// CreateTenant implements TenantStore
func (s *InMemoryTenantStore) CreateTenant(ctx context.Context, tenant *Tenant) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tenants[tenant.ID]; exists {
		return fmt.Errorf("tenant %s already exists", tenant.ID)
	}

	s.tenants[tenant.ID] = tenant
	return nil
}

// GetTenant implements TenantStore
func (s *InMemoryTenantStore) GetTenant(ctx context.Context, tenantID string) (*Tenant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tenant, exists := s.tenants[tenantID]
	if !exists {
		return nil, fmt.Errorf("tenant %s not found", tenantID)
	}

	return tenant, nil
}

// UpdateTenant implements TenantStore
func (s *InMemoryTenantStore) UpdateTenant(ctx context.Context, tenant *Tenant) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tenants[tenant.ID]; !exists {
		return fmt.Errorf("tenant %s not found", tenant.ID)
	}

	s.tenants[tenant.ID] = tenant
	return nil
}

// DeleteTenant implements TenantStore
func (s *InMemoryTenantStore) DeleteTenant(ctx context.Context, tenantID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tenants, tenantID)

	// Clean up API keys
	for key, tid := range s.apiKeys {
		if tid == tenantID {
			delete(s.apiKeys, key)
		}
	}

	return nil
}

// ListTenants implements TenantStore
func (s *InMemoryTenantStore) ListTenants(ctx context.Context, options ListOptions) ([]*Tenant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Tenant

	for _, tenant := range s.tenants {
		// Apply filters
		if options.Status != "" && tenant.Status != options.Status {
			continue
		}
		if options.Plan != "" && tenant.Plan != options.Plan {
			continue
		}

		result = append(result, tenant)
	}

	// Apply pagination
	if options.Offset >= len(result) {
		return []*Tenant{}, nil
	}

	end := options.Offset + options.Limit
	if end > len(result) || options.Limit == 0 {
		end = len(result)
	}

	return result[options.Offset:end], nil
}

// GetTenantByAPIKey implements TenantStore
func (s *InMemoryTenantStore) GetTenantByAPIKey(ctx context.Context, apiKey string) (*Tenant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tenantID, exists := s.apiKeys[apiKey]
	if !exists {
		return nil, fmt.Errorf("invalid API key")
	}

	tenant, exists := s.tenants[tenantID]
	if !exists {
		return nil, fmt.Errorf("tenant not found for API key")
	}

	return tenant, nil
}
