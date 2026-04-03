package dicom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Store handles DICOM storage and integration with FHIR
type Store struct {
	storage ImageStorage
	indexer Indexer
	config  StoreConfig
}

// ImageStorage defines the storage interface for DICOM images
type ImageStorage interface {
	StoreImage(ctx context.Context, image *DICOMImage) error
	RetrieveImage(ctx context.Context, imageID string) (*DICOMImage, error)
	ListImages(ctx context.Context, options ImageListOptions) ([]*DICOMImage, error)
	DeleteImage(ctx context.Context, imageID string) error
	GenerateWADOURL(imageID string) string
}

// Indexer indexes DICOM metadata for searching
type Indexer interface {
	IndexImage(ctx context.Context, image *DICOMImage) error
	SearchImages(ctx context.Context, query ImageQuery) ([]*DICOMImage, error)
}

// StoreConfig contains DICOM store configuration
type StoreConfig struct {
	Enabled           bool
	StorageBackend    string // s3, filesystem, azure
	WADOEnabled       bool
	WADORoot          string
	MaxFileSize       int64
	AllowedModalities []string
}

// DICOMImage represents a DICOM image with metadata
type DICOMImage struct {
	ID                     string                 `json:"id"`
	StudyInstanceUID       string                 `json:"studyInstanceUid"`
	SeriesInstanceUID      string                 `json:"seriesInstanceUid"`
	SOPInstanceUID         string                 `json:"sopInstanceUid"`
	PatientID              string                 `json:"patientId"`
	PatientName            string                 `json:"patientName"`
	Modality               string                 `json:"modality"`
	StudyDate              string                 `json:"studyDate"`
	StudyDescription       string                 `json:"studyDescription"`
	SeriesDescription      string                 `json:"seriesDescription"`
	BodyPart               string                 `json:"bodyPart"`
	InstanceNumber         int                    `json:"instanceNumber"`
	Rows                   int                    `json:"rows"`
	Columns                int                    `json:"columns"`
	FileSize               int64                  `json:"fileSize"`
	ContentType            string                 `json:"contentType"`
	Status                 string                 `json:"status"` // received, stored, indexed, error
	Error                  string                 `json:"error,omitempty"`
	CreatedAt              time.Time              `json:"createdAt"`
	UpdatedAt              time.Time              `json:"updatedAt"`
	FHIRDiagnosticReportID string                 `json:"fhirDiagnosticReportId,omitempty"`
	FHIRImagingStudyID     string                 `json:"fhirImagingStudyId,omitempty"`
	TenantID               string                 `json:"tenantId"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}

// ImageListOptions contains options for listing images
type ImageListOptions struct {
	PatientID string
	StudyUID  string
	SeriesUID string
	Modality  string
	StudyDate string
	TenantID  string
	Limit     int
	Offset    int
}

// ImageQuery contains search query parameters
type ImageQuery struct {
	PatientID     string
	PatientName   string
	Modality      string
	StudyDateFrom string
	StudyDateTo   string
	BodyPart      string
	TenantID      string
}

// NewStore creates a new DICOM store
func NewStore(storage ImageStorage, indexer Indexer, config StoreConfig) *Store {
	return &Store{
		storage: storage,
		indexer: indexer,
		config:  config,
	}
}

// StoreDICOM stores a DICOM image
func (s *Store) StoreDICOM(ctx context.Context, image *DICOMImage, data []byte) error {
	if !s.config.Enabled {
		return fmt.Errorf("DICOM store is disabled")
	}

	// Validate modality
	if !s.isValidModality(image.Modality) {
		return fmt.Errorf("unsupported modality: %s", image.Modality)
	}

	// Set metadata
	image.ID = generateImageID()
	image.FileSize = int64(len(data))
	image.ContentType = "application/dicom"
	image.Status = "received"
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	// Store image
	if err := s.storage.StoreImage(ctx, image); err != nil {
		image.Status = "error"
		image.Error = err.Error()
		return fmt.Errorf("failed to store image: %w", err)
	}

	image.Status = "stored"

	// Index for searching
	if err := s.indexer.IndexImage(ctx, image); err != nil {
		// Log but don't fail
		image.Error = fmt.Sprintf("indexing failed: %v", err)
	} else {
		image.Status = "indexed"
	}

	image.UpdatedAt = time.Now()
	return nil
}

// GetImage retrieves a DICOM image
func (s *Store) GetImage(ctx context.Context, imageID string) (*DICOMImage, error) {
	return s.storage.RetrieveImage(ctx, imageID)
}

// SearchImages searches for DICOM images
func (s *Store) SearchImages(ctx context.Context, query ImageQuery) ([]*DICOMImage, error) {
	return s.indexer.SearchImages(ctx, query)
}

// GenerateFHIRResources generates FHIR ImagingStudy and DiagnosticReport
func (s *Store) GenerateFHIRResources(ctx context.Context, studyUID string) (*FHIRImagingResources, error) {
	// Search for all images in the study
	images, err := s.storage.ListImages(ctx, ImageListOptions{
		StudyUID: studyUID,
	})
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no images found for study: %s", studyUID)
	}

	// Group by series
	seriesMap := make(map[string][]*DICOMImage)
	for _, img := range images {
		seriesMap[img.SeriesInstanceUID] = append(seriesMap[img.SeriesInstanceUID], img)
	}

	// Create ImagingStudy
	imagingStudy := s.createImagingStudy(images[0], seriesMap)

	// Create DiagnosticReport
	diagnosticReport := s.createDiagnosticReport(images[0], imagingStudy)

	return &FHIRImagingResources{
		ImagingStudy:     imagingStudy,
		DiagnosticReport: diagnosticReport,
	}, nil
}

// createImagingStudy creates a FHIR ImagingStudy resource
func (s *Store) createImagingStudy(firstImage *DICOMImage, seriesMap map[string][]*DICOMImage) map[string]interface{} {
	series := make([]map[string]interface{}, 0)

	for seriesUID, images := range seriesMap {
		instances := make([]map[string]interface{}, 0)

		for _, img := range images {
			instance := map[string]interface{}{
				"uid": img.SOPInstanceUID,
				"sopClass": map[string]interface{}{
					"system": "urn:ietf:rfc:3986",
					"code":   "1.2.840.10008.5.1.4.1.1.2", // CT Image Storage
				},
				"number": img.InstanceNumber,
				"title":  img.SeriesDescription,
			}
			instances = append(instances, instance)
		}

		seriesEntry := map[string]interface{}{
			"uid": seriesUID,
			"modality": map[string]interface{}{
				"system": "http://dicom.nema.org/resources/ontology/DCM",
				"code":   images[0].Modality,
			},
			"description": images[0].SeriesDescription,
			"bodySite": map[string]interface{}{
				"display": images[0].BodyPart,
			},
			"instance": instances,
		}
		series = append(series, seriesEntry)
	}

	return map[string]interface{}{
		"resourceType": "ImagingStudy",
		"id":           generateResourceID(),
		"status":       "available",
		"subject": map[string]interface{}{
			"reference": fmt.Sprintf("Patient/%s", firstImage.PatientID),
		},
		"started":     firstImage.StudyDate,
		"description": firstImage.StudyDescription,
		"series":      series,
	}
}

// createDiagnosticReport creates a FHIR DiagnosticReport resource
func (s *Store) createDiagnosticReport(image *DICOMImage, imagingStudy map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"resourceType": "DiagnosticReport",
		"id":           generateResourceID(),
		"status":       "final",
		"category": []map[string]interface{}{
			{
				"coding": []map[string]interface{}{
					{
						"system": "http://terminology.hl7.org/CodeSystem/v2-0074",
						"code":   image.Modality,
					},
				},
			},
		},
		"code": map[string]interface{}{
			"text": fmt.Sprintf("%s Imaging Report", image.Modality),
		},
		"subject": map[string]interface{}{
			"reference": fmt.Sprintf("Patient/%s", image.PatientID),
		},
		"imagingStudy": []map[string]interface{}{
			{
				"reference": fmt.Sprintf("ImagingStudy/%s", imagingStudy["id"]),
			},
		},
		"effectiveDateTime": image.StudyDate,
	}
}

// isValidModality checks if a modality is supported
func (s *Store) isValidModality(modality string) bool {
	if len(s.config.AllowedModalities) == 0 {
		return true
	}

	for _, m := range s.config.AllowedModalities {
		if m == modality {
			return true
		}
	}

	return false
}

// FHIRImagingResources contains generated FHIR resources
type FHIRImagingResources struct {
	ImagingStudy     map[string]interface{} `json:"imagingStudy"`
	DiagnosticReport map[string]interface{} `json:"diagnosticReport"`
}

// WADOServer provides WADO-RS and WADO-URI services
type WADOServer struct {
	store  ImageStorage
	config WADOConfig
}

// WADOConfig contains WADO configuration
type WADOConfig struct {
	Enabled        bool
	RootURL        string
	SupportRS      bool
	SupportURI     bool
	DefaultQuality int
}

// NewWADOServer creates a new WADO server
func NewWADOServer(store ImageStorage, config WADOConfig) *WADOServer {
	return &WADOServer{
		store:  store,
		config: config,
	}
}

// GetWADOURL generates a WADO URL for an image
func (w *WADOServer) GetWADOURL(imageID string) string {
	return w.store.GenerateWADOURL(imageID)
}

// generateImageID generates a unique image ID
func generateImageID() string {
	return fmt.Sprintf("dcm-%d", time.Now().UnixNano())
}

// generateResourceID generates a FHIR resource ID
func generateResourceID() string {
	return fmt.Sprintf("fhir-%d", time.Now().UnixNano())
}

// DICOMHandler handles HTTP requests for DICOM operations
type DICOMHandler struct {
	store *Store
	wado  *WADOServer
}

// NewDICOMHandler creates a new DICOM HTTP handler
func NewDICOMHandler(store *Store, wado *WADOServer) *DICOMHandler {
	return &DICOMHandler{
		store: store,
		wado:  wado,
	}
}

// RegisterRoutes registers DICOM endpoints
func (dh *DICOMHandler) RegisterRoutes(router chi.Router) {
	router.Post("/dicom/studies", dh.handleStoreDICOM)
	router.Get("/dicom/studies/{studyUID}", dh.handleGetStudy)
	router.Get("/dicom/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}", dh.handleGetInstance)
	router.Get("/dicom/search", dh.handleSearch)
	router.Get("/wado/rs/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}", dh.handleWADORS)
	router.Get("/dicom/config", dh.handleGetConfig)
}

// handleStoreDICOM handles POST /dicom/studies
func (dh *DICOMHandler) handleStoreDICOM(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form for DICOM data
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract metadata from form
	metadata := r.FormValue("metadata")

	var image DICOMImage
	if err := json.Unmarshal([]byte(metadata), &image); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get tenant from header
	image.TenantID = r.Header.Get("X-Tenant-ID")

	// Get file
	file, _, err := r.FormFile("dicom")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	data := make([]byte, 0)
	// In production, would read file content here

	// Store DICOM
	if err := dh.store.StoreDICOM(ctx, &image, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(image)
}

// handleGetStudy handles GET /dicom/studies/{studyUID}
func (dh *DICOMHandler) handleGetStudy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	studyUID := chi.URLParam(r, "studyUID")

	// Generate FHIR resources for the study
	resources, err := dh.store.GenerateFHIRResources(ctx, studyUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

// handleGetInstance handles GET /dicom/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}
func (dh *DICOMHandler) handleGetInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	instanceUID := chi.URLParam(r, "instanceUID")

	image, err := dh.store.GetImage(ctx, instanceUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}

// handleSearch handles GET /dicom/search
func (dh *DICOMHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := ImageQuery{
		PatientID:   r.URL.Query().Get("patientId"),
		PatientName: r.URL.Query().Get("patientName"),
		Modality:    r.URL.Query().Get("modality"),
		BodyPart:    r.URL.Query().Get("bodyPart"),
		TenantID:    r.Header.Get("X-Tenant-ID"),
	}

	images, err := dh.store.SearchImages(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// handleWADORS handles WADO-RS requests
func (dh *DICOMHandler) handleWADORS(w http.ResponseWriter, r *http.Request) {
	instanceUID := chi.URLParam(r, "instanceUID")

	if !dh.wado.config.Enabled {
		http.Error(w, "WADO is not enabled", http.StatusNotImplemented)
		return
	}

	// Return WADO URL
	url := dh.wado.GetWADOURL(instanceUID)

	response := map[string]interface{}{
		"wadoUrl": url,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetConfig handles GET /dicom/config
func (dh *DICOMHandler) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"enabled":           dh.store.config.Enabled,
		"storageBackend":    dh.store.config.StorageBackend,
		"wadoEnabled":       dh.store.config.WADOEnabled,
		"wadoRoot":          dh.store.config.WADORoot,
		"allowedModalities": dh.store.config.AllowedModalities,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// InMemoryImageStorage implements ImageStorage with in-memory storage
type InMemoryImageStorage struct {
	images map[string]*DICOMImage
	data   map[string][]byte
	mu     sync.RWMutex
}

// NewInMemoryImageStorage creates a new in-memory image storage
func NewInMemoryImageStorage() *InMemoryImageStorage {
	return &InMemoryImageStorage{
		images: make(map[string]*DICOMImage),
		data:   make(map[string][]byte),
	}
}

// StoreImage implements ImageStorage
func (s *InMemoryImageStorage) StoreImage(ctx context.Context, image *DICOMImage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.images[image.ID] = image
	return nil
}

// RetrieveImage implements ImageStorage
func (s *InMemoryImageStorage) RetrieveImage(ctx context.Context, imageID string) (*DICOMImage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	image, exists := s.images[imageID]
	if !exists {
		return nil, fmt.Errorf("image not found: %s", imageID)
	}

	return image, nil
}

// ListImages implements ImageStorage
func (s *InMemoryImageStorage) ListImages(ctx context.Context, options ImageListOptions) ([]*DICOMImage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*DICOMImage

	for _, image := range s.images {
		if options.PatientID != "" && image.PatientID != options.PatientID {
			continue
		}
		if options.StudyUID != "" && image.StudyInstanceUID != options.StudyUID {
			continue
		}
		if options.SeriesUID != "" && image.SeriesInstanceUID != options.SeriesUID {
			continue
		}
		if options.Modality != "" && image.Modality != options.Modality {
			continue
		}

		results = append(results, image)
	}

	return results, nil
}

// DeleteImage implements ImageStorage
func (s *InMemoryImageStorage) DeleteImage(ctx context.Context, imageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.images, imageID)
	delete(s.data, imageID)

	return nil
}

// GenerateWADOURL implements ImageStorage
func (s *InMemoryImageStorage) GenerateWADOURL(imageID string) string {
	return fmt.Sprintf("/wado/rs/instances/%s", imageID)
}

// InMemoryIndexer implements Indexer with in-memory indexing
type InMemoryIndexer struct {
	images map[string]*DICOMImage
	mu     sync.RWMutex
}

// NewInMemoryIndexer creates a new in-memory indexer
func NewInMemoryIndexer() *InMemoryIndexer {
	return &InMemoryIndexer{
		images: make(map[string]*DICOMImage),
	}
}

// IndexImage implements Indexer
func (i *InMemoryIndexer) IndexImage(ctx context.Context, image *DICOMImage) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.images[image.ID] = image
	return nil
}

// SearchImages implements Indexer
func (i *InMemoryIndexer) SearchImages(ctx context.Context, query ImageQuery) ([]*DICOMImage, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var results []*DICOMImage

	for _, image := range i.images {
		if query.PatientID != "" && image.PatientID != query.PatientID {
			continue
		}
		if query.PatientName != "" && image.PatientName != query.PatientName {
			continue
		}
		if query.Modality != "" && image.Modality != query.Modality {
			continue
		}
		if query.BodyPart != "" && image.BodyPart != query.BodyPart {
			continue
		}

		results = append(results, image)
	}

	return results, nil
}
