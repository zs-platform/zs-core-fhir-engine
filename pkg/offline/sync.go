package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
)

// SyncManager manages offline synchronization for mobile clients
type SyncManager struct {
	store    SyncStore
	config   SyncConfig
	sessions map[string]*SyncSession
	mu       sync.RWMutex
}

// SyncStore defines the storage interface for sync operations
type SyncStore interface {
	GetChanges(ctx context.Context, tenantID string, since time.Time, options SyncOptions) ([]ChangeRecord, error)
	ApplyChanges(ctx context.Context, tenantID string, changes []ChangeRecord) (*SyncResult, error)
	GetLastSyncTime(ctx context.Context, deviceID string) (time.Time, error)
	UpdateLastSyncTime(ctx context.Context, deviceID string, syncTime time.Time) error
	QueueChange(ctx context.Context, tenantID string, change ChangeRecord) error
	GetPendingChanges(ctx context.Context, deviceID string) ([]ChangeRecord, error)
}

// SyncConfig contains sync configuration
type SyncConfig struct {
	BatchSize         int
	ConflictStrategy  string // server-wins, client-wins, manual
	MaxChangesPerSync int
	SyncTimeout       time.Duration
	EnableDeltaSync   bool
}

// SyncOptions contains options for sync operations
type SyncOptions struct {
	ResourceTypes []string
	Limit         int
	Offset        int
}

// SyncSession represents an active sync session
type SyncSession struct {
	ID        string
	DeviceID  string
	UserID    string
	TenantID  string
	StartedAt time.Time
	Status    string // active, completed, failed
	Changes   []ChangeRecord
	mu        sync.Mutex
}

// ChangeRecord represents a single change to be synced
type ChangeRecord struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	ResourceType string                 `json:"resourceType"`
	ResourceID   string                 `json:"resourceId"`
	Action       string                 `json:"action"` // create, update, delete
	Resource     *fhir.Resource         `json:"resource,omitempty"`
	Previous     *fhir.Resource         `json:"previous,omitempty"`
	Checksum     string                 `json:"checksum"`
	DeviceID     string                 `json:"deviceId"`
	UserID       string                 `json:"userId"`
	Sequence     int64                  `json:"sequence"`
	Conflicts    []Conflict            `json:"conflicts,omitempty"`
	Resolved     bool                   `json:"resolved"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// Conflict represents a sync conflict
type Conflict struct {
	Field     string      `json:"field"`
	ServerValue interface{} `json:"serverValue"`
	ClientValue interface{} `json:"clientValue"`
	Resolved  bool        `json:"resolved"`
	Resolution string     `json:"resolution,omitempty"` // server, client, merged
}

// SyncResult contains the result of a sync operation
type SyncResult struct {
	Success        bool            `json:"success"`
	Applied        int             `json:"applied"`
	Rejected       int             `json:"rejected"`
	Conflicts      int             `json:"conflicts"`
	ServerChanges  []ChangeRecord  `json:"serverChanges,omitempty"`
	NextSyncToken  string          `json:"nextSyncToken,omitempty"`
	Timestamp      time.Time       `json:"timestamp"`
	Errors         []string        `json:"errors,omitempty"`
}

// SyncRequest represents a client sync request
type SyncRequest struct {
	DeviceID     string         `json:"deviceId"`
	LastSyncTime time.Time      `json:"lastSyncTime"`
	ClientChanges []ChangeRecord `json:"clientChanges,omitempty"`
	ResourceTypes []string       `json:"resourceTypes,omitempty"`
	SyncToken    string         `json:"syncToken,omitempty"`
}

// SyncResponse represents the server response to a sync request
type SyncResponse struct {
	Success       bool           `json:"success"`
	ServerChanges []ChangeRecord `json:"serverChanges,omitempty"`
	Result        *SyncResult    `json:"result,omitempty"`
	SyncToken     string         `json:"syncToken"`
	Timestamp     time.Time      `json:"timestamp"`
	NextSyncTime  time.Time      `json:"nextSyncTime"`
}

// NewSyncManager creates a new sync manager
func NewSyncManager(store SyncStore, config SyncConfig) *SyncManager {
	return &SyncManager{
		store:    store,
		config:   config,
		sessions: make(map[string]*SyncSession),
	}
}

// InitiateSync initiates a sync session
func (sm *SyncManager) InitiateSync(ctx context.Context, req SyncRequest) (*SyncSession, error) {
	// Create new sync session
	session := &SyncSession{
		ID:        generateSessionID(),
		DeviceID:  req.DeviceID,
		StartedAt: time.Now(),
		Status:    "active",
		Changes:   make([]ChangeRecord, 0),
	}
	
	sm.mu.Lock()
	sm.sessions[session.ID] = session
	sm.mu.Unlock()
	
	log.Infof("Started sync session %s for device %s", session.ID, req.DeviceID)
	
	return session, nil
}

// ProcessSync processes a sync request and returns server changes
func (sm *SyncManager) ProcessSync(ctx context.Context, tenantID, userID string, req SyncRequest) (*SyncResponse, error) {
	// Start sync session
	session, err := sm.InitiateSync(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// Process client changes first
	var result *SyncResult
	if len(req.ClientChanges) > 0 {
		result, err = sm.ApplyClientChanges(ctx, tenantID, userID, req.DeviceID, req.ClientChanges)
		if err != nil {
			session.Status = "failed"
			return nil, fmt.Errorf("failed to apply client changes: %w", err)
		}
	} else {
		result = &SyncResult{Success: true}
	}
	
	// Get server changes since client's last sync
	serverChanges, err := sm.GetServerChanges(ctx, tenantID, req.LastSyncTime, SyncOptions{
		ResourceTypes: req.ResourceTypes,
		Limit:         sm.config.MaxChangesPerSync,
	})
	if err != nil {
		session.Status = "failed"
		return nil, fmt.Errorf("failed to get server changes: %w", err)
	}
	
	// Update session
	session.Changes = serverChanges
	session.Status = "completed"
	
	// Update last sync time for device
	if err := sm.store.UpdateLastSyncTime(ctx, req.DeviceID, time.Now()); err != nil {
		log.Warnf("Failed to update last sync time: %v", err)
	}
	
	// Generate sync token for next sync
	syncToken := sm.generateSyncToken(tenantID, req.DeviceID)
	
	return &SyncResponse{
		Success:       true,
		ServerChanges: serverChanges,
		Result:        result,
		SyncToken:     syncToken,
		Timestamp:     time.Now(),
		NextSyncTime:  time.Now().Add(5 * time.Minute),
	}, nil
}

// ApplyClientChanges applies changes from the client
func (sm *SyncManager) ApplyClientChanges(ctx context.Context, tenantID, userID, deviceID string, changes []ChangeRecord) (*SyncResult, error) {
	result := &SyncResult{
		Success:   true,
		Timestamp: time.Now(),
		Errors:    make([]string, 0),
	}
	
	for i := range changes {
		change := &changes[i]
		change.DeviceID = deviceID
		change.UserID = userID
		change.Timestamp = time.Now()
		change.Sequence = int64(i)
		
		// Check for conflicts
		conflicts, err := sm.checkConflicts(ctx, tenantID, change)
		if err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to check conflicts for %s/%s: %v", 
				change.ResourceType, change.ResourceID, err))
			continue
		}
		
		if len(conflicts) > 0 {
			change.Conflicts = conflicts
			result.Conflicts++
			
			// Resolve conflicts based on strategy
			sm.resolveConflicts(change)
		}
		
		// Queue the change for application
		if err := sm.store.QueueChange(ctx, tenantID, *change); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to queue change for %s/%s: %v", 
				change.ResourceType, change.ResourceID, err))
			continue
		}
		
		result.Applied++
	}
	
	// Apply all queued changes
	finalResult, err := sm.store.ApplyChanges(ctx, tenantID, changes)
	if err != nil {
		return nil, fmt.Errorf("failed to apply changes: %w", err)
	}
	
	// Merge results
	result.Applied = finalResult.Applied
	result.Rejected = finalResult.Rejected
	result.Conflicts = finalResult.Conflicts
	
	return result, nil
}

// GetServerChanges retrieves changes from the server for the client
func (sm *SyncManager) GetServerChanges(ctx context.Context, tenantID string, since time.Time, options SyncOptions) ([]ChangeRecord, error) {
	return sm.store.GetChanges(ctx, tenantID, since, options)
}

// checkConflicts checks if a change conflicts with server state
func (sm *SyncManager) checkConflicts(ctx context.Context, tenantID string, change *ChangeRecord) ([]Conflict, error) {
	conflicts := make([]Conflict, 0)
	
	// In a real implementation, this would compare the change with current server state
	// For now, we assume no conflicts for simplicity
	
	return conflicts, nil
}

// resolveConflicts resolves conflicts using the configured strategy
func (sm *SyncManager) resolveConflicts(change *ChangeRecord) {
	switch sm.config.ConflictStrategy {
	case "server-wins":
		for i := range change.Conflicts {
			change.Conflicts[i].Resolved = true
			change.Conflicts[i].Resolution = "server"
		}
		change.Resolved = true
		
	case "client-wins":
		for i := range change.Conflicts {
			change.Conflicts[i].Resolved = true
			change.Conflicts[i].Resolution = "client"
		}
		change.Resolved = true
		
	case "manual":
		// Leave conflicts unresolved for manual resolution
		change.Resolved = false
		
	default:
		// Default to server-wins
		for i := range change.Conflicts {
			change.Conflicts[i].Resolved = true
			change.Conflicts[i].Resolution = "server"
		}
		change.Resolved = true
	}
}

// GetSyncStatus returns the status of a sync session
func (sm *SyncManager) GetSyncStatus(ctx context.Context, sessionID string) (*SyncSession, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sync session not found: %s", sessionID)
	}
	
	return session, nil
}

// ResolveConflict manually resolves a conflict
func (sm *SyncManager) ResolveConflict(ctx context.Context, changeID string, field string, resolution string) error {
	// In a real implementation, this would update the specific conflict
	log.Infof("Resolved conflict %s for field %s with resolution %s", changeID, field, resolution)
	return nil
}

// GetPendingChanges returns pending changes for a device
func (sm *SyncManager) GetPendingChanges(ctx context.Context, deviceID string) ([]ChangeRecord, error) {
	return sm.store.GetPendingChanges(ctx, deviceID)
}

// generateSyncToken generates a token for the next sync
func (sm *SyncManager) generateSyncToken(tenantID, deviceID string) string {
	return fmt.Sprintf("sync:%s:%s:%d", tenantID, deviceID, time.Now().Unix())
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return fmt.Sprintf("sync-session-%d", time.Now().UnixNano())
}

// SyncHandler handles HTTP requests for offline sync
type SyncHandler struct {
	manager *SyncManager
}

// NewSyncHandler creates a new sync HTTP handler
func NewSyncHandler(manager *SyncManager) *SyncHandler {
	return &SyncHandler{
		manager: manager,
	}
}

// RegisterRoutes registers sync endpoints
func (sh *SyncHandler) RegisterRoutes(router chi.Router) {
	router.Post("/sync", sh.handleSync)
	router.Get("/sync/status/{sessionID}", sh.handleGetSyncStatus)
	router.Get("/sync/pending", sh.handleGetPendingChanges)
	router.Post("/sync/resolve", sh.handleResolveConflict)
	router.Get("/sync/config", sh.handleGetSyncConfig)
}

// handleSync handles POST /sync
func (sh *SyncHandler) handleSync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Extract tenant and user from context
	tenantID := r.Header.Get("X-Tenant-ID")
	userID := r.Header.Get("X-User-ID")
	
	if tenantID == "" {
		http.Error(w, "Tenant ID required", http.StatusBadRequest)
		return
	}
	
	response, err := sh.manager.ProcessSync(ctx, tenantID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetSyncStatus handles GET /sync/status/{sessionID}
func (sh *SyncHandler) handleGetSyncStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "sessionID")
	
	session, err := sh.manager.GetSyncStatus(ctx, sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// handleGetPendingChanges handles GET /sync/pending
func (sh *SyncHandler) handleGetPendingChanges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceID := r.Header.Get("X-Device-ID")
	
	if deviceID == "" {
		http.Error(w, "Device ID required", http.StatusBadRequest)
		return
	}
	
	changes, err := sh.manager.GetPendingChanges(ctx, deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(changes)
}

// handleResolveConflict handles POST /sync/resolve
func (sh *SyncHandler) handleResolveConflict(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		ChangeID   string `json:"changeId"`
		Field      string `json:"field"`
		Resolution string `json:"resolution"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := sh.manager.ResolveConflict(ctx, req.ChangeID, req.Field, req.Resolution); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// handleGetSyncConfig handles GET /sync/config
func (sh *SyncHandler) handleGetSyncConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"batchSize":         sh.manager.config.BatchSize,
		"conflictStrategy":  sh.manager.config.ConflictStrategy,
		"maxChangesPerSync": sh.manager.config.MaxChangesPerSync,
		"syncTimeout":       sh.manager.config.SyncTimeout.Seconds(),
		"enableDeltaSync":   sh.manager.config.EnableDeltaSync,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// InMemorySyncStore implements SyncStore with in-memory storage
type InMemorySyncStore struct {
	changes      []ChangeRecord
	lastSyncTimes map[string]time.Time
	pending      map[string][]ChangeRecord
	mu           sync.RWMutex
}

// NewInMemorySyncStore creates a new in-memory sync store
func NewInMemorySyncStore() *InMemorySyncStore {
	return &InMemorySyncStore{
		changes:       make([]ChangeRecord, 0),
		lastSyncTimes: make(map[string]time.Time),
		pending:       make(map[string][]ChangeRecord),
	}
}

// GetChanges implements SyncStore
func (s *InMemorySyncStore) GetChanges(ctx context.Context, tenantID string, since time.Time, options SyncOptions) ([]ChangeRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []ChangeRecord
	
	for _, change := range s.changes {
		if change.Timestamp.After(since) {
			// Filter by resource types if specified
			if len(options.ResourceTypes) > 0 {
				found := false
				for _, rt := range options.ResourceTypes {
					if change.ResourceType == rt {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			
			results = append(results, change)
		}
	}
	
	// Sort by timestamp
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.Before(results[j].Timestamp)
	})
	
	// Apply limit
	if options.Limit > 0 && len(results) > options.Limit {
		results = results[:options.Limit]
	}
	
	return results, nil
}

// ApplyChanges implements SyncStore
func (s *InMemorySyncStore) ApplyChanges(ctx context.Context, tenantID string, changes []ChangeRecord) (*SyncResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	result := &SyncResult{
		Success:   true,
		Timestamp: time.Now(),
	}
	
	for _, change := range changes {
		s.changes = append(s.changes, change)
		
		if len(change.Conflicts) > 0 {
			result.Conflicts++
		} else {
			result.Applied++
		}
	}
	
	return result, nil
}

// GetLastSyncTime implements SyncStore
func (s *InMemorySyncStore) GetLastSyncTime(ctx context.Context, deviceID string) (time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if lastSync, exists := s.lastSyncTimes[deviceID]; exists {
		return lastSync, nil
	}
	
	return time.Time{}, nil
}

// UpdateLastSyncTime implements SyncStore
func (s *InMemorySyncStore) UpdateLastSyncTime(ctx context.Context, deviceID string, syncTime time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.lastSyncTimes[deviceID] = syncTime
	return nil
}

// QueueChange implements SyncStore
func (s *InMemorySyncStore) QueueChange(ctx context.Context, tenantID string, change ChangeRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if change.ID == "" {
		change.ID = fmt.Sprintf("change-%d", time.Now().UnixNano())
	}
	
	return nil
}

// GetPendingChanges implements SyncStore
func (s *InMemorySyncStore) GetPendingChanges(ctx context.Context, deviceID string) ([]ChangeRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.pending[deviceID], nil
}
