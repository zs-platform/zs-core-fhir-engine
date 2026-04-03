package versioning

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// HTTPHandler handles HTTP requests for resource versioning
type HTTPHandler struct {
	versionManager *VersionManager
}

// NewHTTPHandler creates a new versioning HTTP handler
func NewHTTPHandler(versionManager *VersionManager) *HTTPHandler {
	return &HTTPHandler{
		versionManager: versionManager,
	}
}

// RegisterRoutes registers versioning endpoints with the router
func (h *HTTPHandler) RegisterRoutes(router chi.Router) {
	// Resource history endpoint
	router.Get("/fhir/R5/{resourceType}/{resourceID}/_history", h.handleGetHistory)
	router.Get("/fhir/R5/{resourceType}/{resourceID}/_history/{versionID}", h.handleGetVersion)
	
	// Version restoration endpoint
	router.Put("/fhir/R5/{resourceType}/{resourceID}/_history/{versionID}/_restore", h.handleRestoreVersion)
	
	// Version deletion endpoint (for data retention policies)
	router.Delete("/fhir/R5/{resourceType}/{resourceID}/_history/{versionID}", h.handleDeleteVersion)
	
	// Resource version metadata
	router.Get("/fhir/R5/{resourceType}/{resourceID}/_history/meta", h.handleGetHistoryMetadata)
}

// handleGetHistory retrieves version history for a resource
func (h *HTTPHandler) handleGetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")
	resourceID := chi.URLParam(r, "resourceID")
	
	// Parse query parameters
	options := HistoryOptions{
		Page:  1,
		Count: 20,
	}
	
	if page := r.URL.Query().Get("_page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			options.Page = p
		}
	}
	
	if count := r.URL.Query().Get("_count"); count != "" {
		if c, err := strconv.Atoi(count); err == nil && c > 0 {
			options.Count = c
		}
	}
	
	if since := r.URL.Query().Get("_since"); since != "" {
		if t, err := time.Parse(time.RFC3339, since); err == nil {
			options.Since = &t
		}
	}
	
	if until := r.URL.Query().Get("_until"); until != "" {
		if t, err := time.Parse(time.RFC3339, until); err == nil {
			options.Until = &t
		}
	}
	
	if operation := r.URL.Query().Get("_operation"); operation != "" {
		options.Operation = operation
	}
	
	if userID := r.URL.Query().Get("_user"); userID != "" {
		options.UserID = userID
	}
	
	// Get history
	bundle, err := h.versionManager.GetHistory(ctx, resourceType, resourceID, options)
	if err != nil {
		h.writeError(w, "history_error", err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return bundle
	h.writeJSON(w, bundle, http.StatusOK)
}

// handleGetVersion retrieves a specific version
func (h *HTTPHandler) handleGetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")
	resourceID := chi.URLParam(r, "resourceID")
	versionIDStr := chi.URLParam(r, "versionID")
	
	// Parse version ID
	versionID, err := strconv.Atoi(versionIDStr)
	if err != nil {
		h.writeError(w, "invalid_version", "Invalid version ID", http.StatusBadRequest)
		return
	}
	
	// Get version
	version, err := h.versionManager.GetVersion(ctx, resourceType, resourceID, versionID)
	if err != nil {
		h.writeError(w, "version_not_found", err.Error(), http.StatusNotFound)
		return
	}
	
	// Return version
	h.writeJSON(w, version, http.StatusOK)
}

// handleRestoreVersion restores a resource to a specific version
func (h *HTTPHandler) handleRestoreVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")
	resourceID := chi.URLParam(r, "resourceID")
	versionIDStr := chi.URLParam(r, "versionID")
	
	// Parse version ID
	versionID, err := strconv.Atoi(versionIDStr)
	if err != nil {
		h.writeError(w, "invalid_version", "Invalid version ID", http.StatusBadRequest)
		return
	}
	
	// Restore version
	restoredVersion, err := h.versionManager.RestoreVersion(ctx, resourceType, resourceID, versionID)
	if err != nil {
		h.writeError(w, "restore_error", err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return restored version
	h.writeJSON(w, restoredVersion, http.StatusOK)
}

// handleDeleteVersion deletes a specific version
func (h *HTTPHandler) handleDeleteVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")
	resourceID := chi.URLParam(r, "resourceID")
	versionIDStr := chi.URLParam(r, "versionID")
	
	// Parse version ID
	versionID, err := strconv.Atoi(versionIDStr)
	if err != nil {
		h.writeError(w, "invalid_version", "Invalid version ID", http.StatusBadRequest)
		return
	}
	
	// Delete version
	if err := h.versionManager.DeleteVersion(ctx, resourceType, resourceID, versionID); err != nil {
		h.writeError(w, "delete_error", err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// handleGetHistoryMetadata retrieves metadata about version history
func (h *HTTPHandler) handleGetHistoryMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")
	resourceID := chi.URLParam(r, "resourceID")
	
	// Get all versions
	versions, err := h.versionManager.GetAllVersions(ctx, resourceType, resourceID)
	if err != nil {
		h.writeError(w, "history_error", err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Calculate metadata
	metadata := map[string]interface{}{
		"resourceType": resourceType,
		"resourceId":   resourceID,
		"totalVersions": len(versions),
		"currentVersion": 0,
		"firstVersion":   nil,
		"lastVersion":    nil,
		"operations":     make(map[string]int),
	}
	
	if len(versions) > 0 {
		metadata["currentVersion"] = versions[len(versions)-1].VersionID
		metadata["firstVersion"] = map[string]interface{}{
			"versionId": versions[0].VersionID,
			"timestamp": versions[0].Timestamp,
			"operation": versions[0].Operation,
		}
		metadata["lastVersion"] = map[string]interface{}{
			"versionId": versions[len(versions)-1].VersionID,
			"timestamp": versions[len(versions)-1].Timestamp,
			"operation": versions[len(versions)-1].Operation,
		}
		
		// Count operations
		operations := make(map[string]int)
		for _, v := range versions {
			operations[v.Operation]++
		}
		metadata["operations"] = operations
	}
	
	// Return metadata
	h.writeJSON(w, metadata, http.StatusOK)
}

// writeJSON writes a JSON response
func (h *HTTPHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Fallback error response
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes a FHIR OperationOutcome error
func (h *HTTPHandler) writeError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(status)
	
	outcome := map[string]interface{}{
		"resourceType": "OperationOutcome",
		"issue": []map[string]interface{}{
			{
				"severity":    "error",
				"code":        code,
				"diagnostics": message,
			},
		},
	}
	
	json.NewEncoder(w).Encode(outcome)
}
