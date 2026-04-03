package versioning

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
)

// VersionManager manages resource versioning and history
type VersionManager struct {
	store VersionStore
}

// VersionStore defines the storage interface for versioning
type VersionStore interface {
	CreateVersion(ctx context.Context, resource fhir.Resource, operation string, userID string) (*ResourceVersion, error)
	GetVersion(ctx context.Context, resourceType, resourceID string, versionID int) (*ResourceVersion, error)
	GetHistory(ctx context.Context, resourceType, resourceID string, options HistoryOptions) (*VersionBundle, error)
	GetAllVersions(ctx context.Context, resourceType, resourceID string) ([]ResourceVersion, error)
	RestoreVersion(ctx context.Context, resourceType, resourceID string, versionID int) (*ResourceVersion, error)
	DeleteVersion(ctx context.Context, resourceType, resourceID string, versionID int) error
}

// ResourceVersion represents a single version of a FHIR resource
type ResourceVersion struct {
	ID           string          `json:"id"`
	ResourceType string          `json:"resourceType"`
	ResourceID   string          `json:"resourceId"`
	VersionID    int             `json:"versionId"`
	Resource     fhir.Resource   `json:"resource"`
	Operation    string          `json:"operation"` // create, update, delete
	UserID       string          `json:"userId"`
	UserName     string          `json:"userName,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
	Provenance   *ProvenanceInfo `json:"provenance,omitempty"`
	Diff         *ResourceDiff   `json:"diff,omitempty"`
}

// ProvenanceInfo contains audit information about the version
type ProvenanceInfo struct {
	Agent     string   `json:"agent"`
	AgentType string   `json:"agentType"`
	Entity    string   `json:"entity"`
	Activity  string   `json:"activity"`
	Location  string   `json:"location,omitempty"`
	Policy    []string `json:"policy,omitempty"`
	Signature string   `json:"signature,omitempty"`
}

// ResourceDiff represents the difference between two versions
type ResourceDiff struct {
	Added    map[string]interface{} `json:"added,omitempty"`
	Removed  map[string]interface{} `json:"removed,omitempty"`
	Modified map[string]interface{} `json:"modified,omitempty"`
}

// HistoryOptions contains options for retrieving history
type HistoryOptions struct {
	Since     *time.Time
	Until     *time.Time
	Count     int
	Page      int
	Operation string // filter by operation type
	UserID    string // filter by user
}

// VersionBundle contains a bundle of resource versions
type VersionBundle struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Total        int            `json:"total"`
	Link         []BundleLink   `json:"link,omitempty"`
	Entry        []VersionEntry `json:"entry"`
}

// BundleLink represents pagination links
type BundleLink struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

// VersionEntry represents a single version entry in a bundle
type VersionEntry struct {
	FullURL  string          `json:"fullUrl"`
	Resource fhir.Resource   `json:"resource,omitempty"`
	Request  *BundleRequest  `json:"request,omitempty"`
	Response *BundleResponse `json:"response,omitempty"`
}

// BundleRequest represents the request that created this version
type BundleRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

// BundleResponse represents the response when the version was created
type BundleResponse struct {
	Status  string    `json:"status"`
	ETag    string    `json:"etag,omitempty"`
	LastMod time.Time `json:"lastModified,omitempty"`
}

// NewVersionManager creates a new version manager
func NewVersionManager(store VersionStore) *VersionManager {
	return &VersionManager{
		store: store,
	}
}

// CreateVersion creates a new version of a resource
func (vm *VersionManager) CreateVersion(ctx context.Context, resource fhir.Resource, operation string, userID string, userName string) (*ResourceVersion, error) {
	return vm.store.CreateVersion(ctx, resource, operation, userID)
}

// GetVersion retrieves a specific version of a resource
func (vm *VersionManager) GetVersion(ctx context.Context, resourceType, resourceID string, versionID int) (*ResourceVersion, error) {
	return vm.store.GetVersion(ctx, resourceType, resourceID, versionID)
}

// GetHistory retrieves the version history for a resource
func (vm *VersionManager) GetHistory(ctx context.Context, resourceType, resourceID string, options HistoryOptions) (*VersionBundle, error) {
	return vm.store.GetHistory(ctx, resourceType, resourceID, options)
}

// GetAllVersions retrieves all versions for a resource.
func (vm *VersionManager) GetAllVersions(ctx context.Context, resourceType, resourceID string) ([]ResourceVersion, error) {
	return vm.store.GetAllVersions(ctx, resourceType, resourceID)
}

// RestoreVersion restores a resource to a specific version
func (vm *VersionManager) RestoreVersion(ctx context.Context, resourceType, resourceID string, versionID int) (*ResourceVersion, error) {
	return vm.store.RestoreVersion(ctx, resourceType, resourceID, versionID)
}

// DeleteVersion deletes a specific version from history (for data retention policies)
func (vm *VersionManager) DeleteVersion(ctx context.Context, resourceType, resourceID string, versionID int) error {
	return vm.store.DeleteVersion(ctx, resourceType, resourceID, versionID)
}

// CalculateDiff calculates the difference between two versions
func (vm *VersionManager) CalculateDiff(version1, version2 *ResourceVersion) (*ResourceDiff, error) {
	return calculateDiff(version1.Resource, version2.Resource)
}

// HistoryHandler handles HTTP requests for resource history
type HistoryHandler struct {
	versionManager *VersionManager
}

// NewHistoryHandler creates a new history HTTP handler
func NewHistoryHandler(versionManager *VersionManager) *HistoryHandler {
	return &HistoryHandler{
		versionManager: versionManager,
	}
}

// InMemoryVersionStore implements VersionStore with in-memory storage
type InMemoryVersionStore struct {
	versions map[string][]ResourceVersion // key: resourceType/resourceID
}

// NewInMemoryVersionStore creates a new in-memory version store
func NewInMemoryVersionStore() *InMemoryVersionStore {
	return &InMemoryVersionStore{
		versions: make(map[string][]ResourceVersion),
	}
}

// CreateVersion implements VersionStore
func (s *InMemoryVersionStore) CreateVersion(ctx context.Context, resource fhir.Resource, operation string, userID string) (*ResourceVersion, error) {
	key := fmt.Sprintf("%s/%s", resource.ResourceType, *resource.ID)

	// Get existing versions
	existing := s.versions[key]
	nextVersionID := 1
	if len(existing) > 0 {
		nextVersionID = existing[len(existing)-1].VersionID + 1
	}

	// Create new version
	version := ResourceVersion{
		ID:           generateVersionID(),
		ResourceType: resource.ResourceType,
		ResourceID:   *resource.ID,
		VersionID:    nextVersionID,
		Resource:     resource,
		Operation:    operation,
		UserID:       userID,
		Timestamp:    time.Now(),
		Provenance: &ProvenanceInfo{
			Agent:    userID,
			Entity:   key,
			Activity: operation,
		},
	}

	// Store version
	s.versions[key] = append(existing, version)

	return &version, nil
}

// GetVersion implements VersionStore
func (s *InMemoryVersionStore) GetVersion(ctx context.Context, resourceType, resourceID string, versionID int) (*ResourceVersion, error) {
	key := fmt.Sprintf("%s/%s", resourceType, resourceID)

	versions, ok := s.versions[key]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", key)
	}

	for _, version := range versions {
		if version.VersionID == versionID {
			return &version, nil
		}
	}

	return nil, fmt.Errorf("version %d not found for %s", versionID, key)
}

// GetHistory implements VersionStore
func (s *InMemoryVersionStore) GetHistory(ctx context.Context, resourceType, resourceID string, options HistoryOptions) (*VersionBundle, error) {
	key := fmt.Sprintf("%s/%s", resourceType, resourceID)

	allVersions, ok := s.versions[key]
	if !ok {
		allVersions = []ResourceVersion{}
	}

	// Filter by options
	filtered := make([]ResourceVersion, 0)
	for _, version := range allVersions {
		// Apply filters
		if options.Since != nil && version.Timestamp.Before(*options.Since) {
			continue
		}
		if options.Until != nil && version.Timestamp.After(*options.Until) {
			continue
		}
		if options.Operation != "" && version.Operation != options.Operation {
			continue
		}
		if options.UserID != "" && version.UserID != options.UserID {
			continue
		}

		filtered = append(filtered, version)
	}

	// Apply pagination
	total := len(filtered)
	page := options.Page
	if page < 1 {
		page = 1
	}
	count := options.Count
	if count < 1 {
		count = 20
	}

	start := (page - 1) * count
	end := start + count
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedVersions := filtered[start:end]

	// Create bundle
	bundle := &VersionBundle{
		ResourceType: "Bundle",
		ID:           generateVersionID(),
		Type:         "history",
		Total:        total,
		Entry:        make([]VersionEntry, 0, len(paginatedVersions)),
	}

	// Add pagination links
	baseURL := fmt.Sprintf("http://localhost:8080/fhir/R5/%s/%s/_history", resourceType, resourceID)

	// Self link
	bundle.Link = append(bundle.Link, BundleLink{
		Relation: "self",
		URL:      fmt.Sprintf("%s?page=%d&_count=%d", baseURL, page, count),
	})

	// First link
	bundle.Link = append(bundle.Link, BundleLink{
		Relation: "first",
		URL:      fmt.Sprintf("%s?page=1&_count=%d", baseURL, count),
	})

	// Previous link
	if page > 1 {
		bundle.Link = append(bundle.Link, BundleLink{
			Relation: "previous",
			URL:      fmt.Sprintf("%s?page=%d&_count=%d", baseURL, page-1, count),
		})
	}

	// Next link
	if end < total {
		bundle.Link = append(bundle.Link, BundleLink{
			Relation: "next",
			URL:      fmt.Sprintf("%s?page=%d&_count=%d", baseURL, page+1, count),
		})
	}

	// Last link
	lastPage := (total + count - 1) / count
	if lastPage < 1 {
		lastPage = 1
	}
	bundle.Link = append(bundle.Link, BundleLink{
		Relation: "last",
		URL:      fmt.Sprintf("%s?page=%d&_count=%d", baseURL, lastPage, count),
	})

	// Add entries
	for _, version := range paginatedVersions {
		entry := VersionEntry{
			FullURL: fmt.Sprintf("http://localhost:8080/fhir/R5/%s/%s/_history/%d",
				version.ResourceType, version.ResourceID, version.VersionID),
			Resource: version.Resource,
			Request: &BundleRequest{
				Method: version.Operation,
				URL:    fmt.Sprintf("%s/%s", version.ResourceType, version.ResourceID),
			},
			Response: &BundleResponse{
				Status:  "200",
				ETag:    fmt.Sprintf("W/\"%d\"", version.VersionID),
				LastMod: version.Timestamp,
			},
		}
		bundle.Entry = append(bundle.Entry, entry)
	}

	return bundle, nil
}

// GetAllVersions implements VersionStore
func (s *InMemoryVersionStore) GetAllVersions(ctx context.Context, resourceType, resourceID string) ([]ResourceVersion, error) {
	key := fmt.Sprintf("%s/%s", resourceType, resourceID)

	versions, ok := s.versions[key]
	if !ok {
		return []ResourceVersion{}, nil
	}

	return versions, nil
}

// RestoreVersion implements VersionStore
func (s *InMemoryVersionStore) RestoreVersion(ctx context.Context, resourceType, resourceID string, versionID int) (*ResourceVersion, error) {
	// Get the version to restore
	version, err := s.GetVersion(ctx, resourceType, resourceID, versionID)
	if err != nil {
		return nil, err
	}

	// Create a new version with the restored content
	restoredVersion, err := s.CreateVersion(ctx, version.Resource, "restore", "system")
	if err != nil {
		return nil, err
	}

	return restoredVersion, nil
}

// DeleteVersion implements VersionStore
func (s *InMemoryVersionStore) DeleteVersion(ctx context.Context, resourceType, resourceID string, versionID int) error {
	key := fmt.Sprintf("%s/%s", resourceType, resourceID)

	versions, ok := s.versions[key]
	if !ok {
		return fmt.Errorf("resource not found: %s", key)
	}

	// Find and remove the version
	newVersions := make([]ResourceVersion, 0, len(versions)-1)
	found := false

	for _, version := range versions {
		if version.VersionID == versionID {
			found = true
			continue
		}
		newVersions = append(newVersions, version)
	}

	if !found {
		return fmt.Errorf("version %d not found for %s", versionID, key)
	}

	s.versions[key] = newVersions
	return nil
}

// Helper functions

func generateVersionID() string {
	return fmt.Sprintf("version-%d", time.Now().UnixNano())
}

func calculateDiff(resource1, resource2 fhir.Resource) (*ResourceDiff, error) {
	diff := &ResourceDiff{
		Added:    make(map[string]interface{}),
		Removed:  make(map[string]interface{}),
		Modified: make(map[string]interface{}),
	}

	// Convert resources to maps for comparison
	map1, err := resourceToMap(resource1)
	if err != nil {
		return nil, err
	}

	map2, err := resourceToMap(resource2)
	if err != nil {
		return nil, err
	}

	// Calculate differences
	calculateMapDiff(map1, map2, "", diff)

	return diff, nil
}

func resourceToMap(resource fhir.Resource) (map[string]interface{}, error) {
	data, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func calculateMapDiff(map1, map2 map[string]interface{}, prefix string, diff *ResourceDiff) {
	// Find added and modified fields
	for key, value2 := range map2 {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		value1, exists := map1[key]
		if !exists {
			// Field was added
			diff.Added[fullKey] = value2
		} else {
			// Field exists in both, check if modified
			switch v1 := value1.(type) {
			case map[string]interface{}:
				if v2, ok := value2.(map[string]interface{}); ok {
					calculateMapDiff(v1, v2, fullKey, diff)
				} else {
					diff.Modified[fullKey] = map[string]interface{}{
						"old": value1,
						"new": value2,
					}
				}
			default:
				if fmt.Sprintf("%v", value1) != fmt.Sprintf("%v", value2) {
					diff.Modified[fullKey] = map[string]interface{}{
						"old": value1,
						"new": value2,
					}
				}
			}
		}
	}

	// Find removed fields
	for key, value1 := range map1 {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if _, exists := map2[key]; !exists {
			diff.Removed[fullKey] = value1
		}
	}
}
