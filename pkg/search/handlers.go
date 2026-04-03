package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/models"
)

// SearchHandler handles HTTP requests for FHIR search operations
type SearchHandler struct {
	engine *SearchEngine
}

// NewSearchHandler creates a new search HTTP handler
func NewSearchHandler(engine *SearchEngine) *SearchHandler {
	return &SearchHandler{
		engine: engine,
	}
}

// RegisterRoutes registers search endpoints with the router
func (sh *SearchHandler) RegisterRoutes(router chi.Router) {
	// Resource-specific search endpoints
	router.Get("/fhir/R5/{resourceType}", sh.handleSearch)
	router.Get("/fhir/R5/{resourceType}/_search", sh.handleSearchPOST)
	router.Post("/fhir/R5/{resourceType}/_search", sh.handleSearchPOST)

	// Global search across all resources
	router.Get("/fhir/R5", sh.handleGlobalSearch)
}

// handleSearch handles GET search requests
func (sh *SearchHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")

	// Parse query parameters
	if err := r.ParseForm(); err != nil {
		sh.writeError(w, "invalid_request", "Failed to parse query parameters", http.StatusBadRequest)
		return
	}

	// Build search query
	query, err := sh.engine.ParseSearchQuery(resourceType, r.Form)
	if err != nil {
		sh.writeError(w, "invalid_query", err.Error(), http.StatusBadRequest)
		return
	}

	// Validate query
	if err := sh.engine.ValidateSearchQuery(query); err != nil {
		sh.writeError(w, "validation_error", err.Error(), http.StatusBadRequest)
		return
	}

	// Execute search
	results, err := sh.engine.ExecuteSearch(ctx, query)
	if err != nil {
		sh.writeError(w, "search_error", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	sh.writeBundle(w, results, http.StatusOK)
}

// handleSearchPOST handles POST search requests (for complex queries)
func (sh *SearchHandler) handleSearchPOST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resourceType")

	// Parse form data from POST body
	if err := r.ParseForm(); err != nil {
		sh.writeError(w, "invalid_request", "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Build search query
	query, err := sh.engine.ParseSearchQuery(resourceType, r.Form)
	if err != nil {
		sh.writeError(w, "invalid_query", err.Error(), http.StatusBadRequest)
		return
	}

	// Validate query
	if err := sh.engine.ValidateSearchQuery(query); err != nil {
		sh.writeError(w, "validation_error", err.Error(), http.StatusBadRequest)
		return
	}

	// Execute search
	results, err := sh.engine.ExecuteSearch(ctx, query)
	if err != nil {
		sh.writeError(w, "search_error", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	sh.writeBundle(w, results, http.StatusOK)
}

// handleGlobalSearch handles global search across all resource types
func (sh *SearchHandler) handleGlobalSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	if err := r.ParseForm(); err != nil {
		sh.writeError(w, "invalid_request", "Failed to parse query parameters", http.StatusBadRequest)
		return
	}

	// Check for _type parameter to limit resource types
	resourceTypes := []string{"Patient", "Observation", "Condition", "MedicationRequest"}
	if typeParam := r.FormValue("_type"); typeParam != "" {
		resourceTypes = strings.Split(typeParam, ",")
	}

	// Search across all specified resource types
	allResults := make([]models.Resource, 0)

	for _, resourceType := range resourceTypes {
		query, err := sh.engine.ParseSearchQuery(resourceType, r.Form)
		if err != nil {
			continue // Skip invalid resource types
		}

		results, err := sh.engine.ExecuteSearch(ctx, query)
		if err != nil {
			continue // Skip resources that fail to search
		}

		// Collect entries
		for _, entry := range results.Entry {
			if entry.Resource != nil {
				allResults = append(allResults, *entry.Resource)
			}
		}
	}

	// Create bundle with combined results
	bundle := sh.createBundle(allResults, len(allResults))
	sh.writeBundle(w, bundle, http.StatusOK)
}

// writeBundle writes a FHIR bundle response
func (sh *SearchHandler) writeBundle(w http.ResponseWriter, bundle *models.Bundle, status int) {
	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(status)

	// Set bundle metadata
	bundle.ResourceType = "Bundle"
	bundle.Type = "searchset"
	bundle.Timestamp = time.Now().Format(time.RFC3339)

	// Encode and write
	if err := json.NewEncoder(w).Encode(bundle); err != nil {
		// Fallback error response
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes a FHIR OperationOutcome error
func (sh *SearchHandler) writeError(w http.ResponseWriter, code, message string, status int) {
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

// createBundle creates a FHIR bundle from resources
func (sh *SearchHandler) createBundle(resources []models.Resource, total int) *models.Bundle {
	bundle := &models.Bundle{
		ResourceType: "Bundle",
		Type:         "searchset",
		Total:        &total,
		Entry:        make([]models.BundleEntry, 0, len(resources)),
	}

	for i, resource := range resources {
		entry := models.BundleEntry{
			FullURL: fmt.Sprintf("http://localhost:8080/fhir/R5/%s/%s",
				resource.ResourceType(), resource.ID()),
			Resource: &resource,
			Search: &models.BundleEntrySearch{
				Mode:  "match",
				Score: float64(len(resources)-i) / float64(len(resources)),
			},
		}
		bundle.Entry = append(bundle.Entry, entry)
	}

	return bundle
}

// InMemorySearchStore implements ResourceStore with in-memory storage
type InMemorySearchStore struct {
	resources map[string]map[string]models.Resource // resourceType -> id -> resource
}

// NewInMemorySearchStore creates a new in-memory search store
func NewInMemorySearchStore() *InMemorySearchStore {
	return &InMemorySearchStore{
		resources: make(map[string]map[string]models.Resource),
	}
}

// StoreResource stores a resource for search
func (s *InMemorySearchStore) StoreResource(resource models.Resource) error {
	resourceType := resource.ResourceType()
	id := resource.ID()

	if _, ok := s.resources[resourceType]; !ok {
		s.resources[resourceType] = make(map[string]models.Resource)
	}

	s.resources[resourceType][id] = resource
	return nil
}

// Search implements ResourceStore interface
func (s *InMemorySearchStore) Search(ctx context.Context, query SearchQuery) (*models.Bundle, error) {
	return s.executeSearch(query)
}

// SearchByParams implements ResourceStore interface
func (s *InMemorySearchStore) SearchByParams(ctx context.Context, resourceType string, params map[string][]string) (*models.Bundle, error) {
	query := &SearchQuery{
		ResourceType: resourceType,
		Parameters:   make([]SearchParameter, 0),
		Page:         1,
		Count:        20,
	}

	for key, values := range params {
		for _, value := range values {
			param := SearchParameter{
				Name:     key,
				Value:    value,
				Resource: resourceType,
			}
			query.Parameters = append(query.Parameters, param)
		}
	}

	return s.executeSearch(*query)
}

// executeSearch performs the actual search operation
func (s *InMemorySearchStore) executeSearch(query SearchQuery) (*models.Bundle, error) {
	results := make([]models.Resource, 0)

	// Get resources of the specified type
	resources, ok := s.resources[query.ResourceType]
	if !ok {
		// Return empty bundle
		return s.createBundle(results, 0), nil
	}

	// Filter resources based on search parameters
	for _, resource := range resources {
		if s.matchesSearch(resource, query.Parameters) {
			results = append(results, resource)
		}
	}

	// Apply sorting
	if len(query.Sort) > 0 {
		results = s.sortResults(results, query.Sort)
	}

	// Apply pagination
	total := len(results)
	start := (query.Page - 1) * query.Count
	end := start + query.Count

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedResults := results[start:end]

	return s.createBundle(paginatedResults, total), nil
}

// matchesSearch checks if a resource matches the search parameters
func (s *InMemorySearchStore) matchesSearch(resource models.Resource, params []SearchParameter) bool {
	for _, param := range params {
		if !s.matchesParameter(resource, param) {
			return false
		}
	}
	return true
}

// matchesParameter checks if a resource matches a single search parameter
func (s *InMemorySearchStore) matchesParameter(resource models.Resource, param SearchParameter) bool {
	// Get resource data as map for searching
	resourceMap := resourceToMap(resource)

	switch param.Name {
	case "_id":
		return s.matchToken(resourceMap["id"], param.Value, param.Modifier)
	case "_lastUpdated":
		return s.matchDate(resourceMap["meta"], param.Value, param.Prefix)
	case "identifier":
		return s.matchIdentifier(resourceMap["identifier"], param.Value, param.Modifier)
	case "name":
		return s.matchName(resourceMap["name"], param.Value, param.Modifier)
	case "family":
		return s.matchFamilyName(resourceMap["name"], param.Value, param.Modifier)
	case "given":
		return s.matchGivenName(resourceMap["name"], param.Value, param.Modifier)
	case "birthdate":
		return s.matchBirthDate(resourceMap["birthDate"], param.Value, param.Prefix)
	case "gender":
		return s.matchToken(resourceMap["gender"], param.Value, param.Modifier)
	case "address":
		return s.matchAddress(resourceMap["address"], param.Value, param.Modifier)
	case "address-city":
		return s.matchAddressField(resourceMap["address"], "city", param.Value, param.Modifier)
	case "address-state":
		return s.matchAddressField(resourceMap["address"], "state", param.Value, param.Modifier)
	case "address-postalcode":
		return s.matchAddressField(resourceMap["address"], "postalCode", param.Value, param.Modifier)
	case "address-country":
		return s.matchAddressField(resourceMap["address"], "country", param.Value, param.Modifier)
	case "phone":
		return s.matchContactPoint(resourceMap["telecom"], "phone", param.Value, param.Modifier)
	case "email":
		return s.matchContactPoint(resourceMap["telecom"], "email", param.Value, param.Modifier)
	case "active":
		return s.matchBoolean(resourceMap["active"], param.Value)
	case "code":
		return s.matchCode(resourceMap["code"], param.Value, param.Modifier)
	case "subject", "patient":
		return s.matchReference(resourceMap["subject"], param.Value, param.Modifier) ||
			s.matchReference(resourceMap["patient"], param.Value, param.Modifier)
	case "encounter":
		return s.matchReference(resourceMap["encounter"], param.Value, param.Modifier)
	case "date":
		return s.matchDate(resourceMap["effectiveDateTime"], param.Value, param.Prefix) ||
			s.matchDate(resourceMap["effectivePeriod"], param.Value, param.Prefix)
	case "status":
		return s.matchToken(resourceMap["status"], param.Value, param.Modifier)
	case "category":
		return s.matchCodeableConcept(resourceMap["category"], param.Value, param.Modifier)
	default:
		// For unknown parameters, try to match as string
		return s.matchString(fmt.Sprintf("%v", resourceMap[param.Name]), param.Value, param.Modifier)
	}
}

// Helper functions for matching different parameter types

func (s *InMemorySearchStore) matchToken(value interface{}, searchValue, modifier string) bool {
	if value == nil {
		return false
	}

	strValue := fmt.Sprintf("%v", value)

	switch modifier {
	case "text":
		return strings.Contains(strings.ToLower(strValue), strings.ToLower(searchValue))
	case "not":
		return !strings.EqualFold(strValue, searchValue)
	default:
		return strings.EqualFold(strValue, searchValue)
	}
}

func (s *InMemorySearchStore) matchString(value interface{}, searchValue, modifier string) bool {
	if value == nil {
		return false
	}

	strValue := fmt.Sprintf("%v", value)

	switch modifier {
	case "exact":
		return strValue == searchValue
	case "contains":
		return strings.Contains(strings.ToLower(strValue), strings.ToLower(searchValue))
	default:
		// Default is partial match (starts with)
		return strings.HasPrefix(strings.ToLower(strValue), strings.ToLower(searchValue))
	}
}

func (s *InMemorySearchStore) matchDate(value interface{}, searchValue, prefix string) bool {
	// Simplified date matching - in production would use proper date parsing
	if value == nil {
		return false
	}

	dateStr := fmt.Sprintf("%v", value)
	return strings.HasPrefix(dateStr, searchValue)
}

func (s *InMemorySearchStore) matchIdentifier(identifiers interface{}, searchValue, modifier string) bool {
	if identifiers == nil {
		return false
	}

	// Try to match against identifier array
	idList, ok := identifiers.([]interface{})
	if !ok {
		return false
	}

	for _, id := range idList {
		idMap, ok := id.(map[string]interface{})
		if !ok {
			continue
		}

		value, hasValue := idMap["value"]
		if !hasValue {
			continue
		}

		if s.matchToken(value, searchValue, modifier) {
			return true
		}
	}

	return false
}

func (s *InMemorySearchStore) matchName(names interface{}, searchValue, modifier string) bool {
	if names == nil {
		return false
	}

	nameList, ok := names.([]interface{})
	if !ok {
		return false
	}

	for _, name := range nameList {
		nameMap, ok := name.(map[string]interface{})
		if !ok {
			continue
		}

		// Check text
		if text, hasText := nameMap["text"]; hasText {
			if s.matchString(text, searchValue, modifier) {
				return true
			}
		}

		// Check family name
		if family, hasFamily := nameMap["family"]; hasFamily {
			if s.matchString(family, searchValue, modifier) {
				return true
			}
		}

		// Check given names
		if given, hasGiven := nameMap["given"]; hasGiven {
			if givenList, ok := given.([]interface{}); ok {
				for _, g := range givenList {
					if s.matchString(g, searchValue, modifier) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (s *InMemorySearchStore) matchFamilyName(names interface{}, searchValue, modifier string) bool {
	if names == nil {
		return false
	}

	nameList, ok := names.([]interface{})
	if !ok {
		return false
	}

	for _, name := range nameList {
		nameMap, ok := name.(map[string]interface{})
		if !ok {
			continue
		}

		if family, hasFamily := nameMap["family"]; hasFamily {
			if s.matchString(family, searchValue, modifier) {
				return true
			}
		}
	}

	return false
}

func (s *InMemorySearchStore) matchGivenName(names interface{}, searchValue, modifier string) bool {
	if names == nil {
		return false
	}

	nameList, ok := names.([]interface{})
	if !ok {
		return false
	}

	for _, name := range nameList {
		nameMap, ok := name.(map[string]interface{})
		if !ok {
			continue
		}

		if given, hasGiven := nameMap["given"]; hasGiven {
			if givenList, ok := given.([]interface{}); ok {
				for _, g := range givenList {
					if s.matchString(g, searchValue, modifier) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (s *InMemorySearchStore) matchBirthDate(value interface{}, searchValue, prefix string) bool {
	return s.matchDate(value, searchValue, prefix)
}

func (s *InMemorySearchStore) matchAddress(addresses interface{}, searchValue, modifier string) bool {
	if addresses == nil {
		return false
	}

	addrList, ok := addresses.([]interface{})
	if !ok {
		return false
	}

	for _, addr := range addrList {
		addrMap, ok := addr.(map[string]interface{})
		if !ok {
			continue
		}

		// Search in text, line, city, state, postalCode, country
		fields := []string{"text", "city", "state", "postalCode", "country"}
		for _, field := range fields {
			if value, hasField := addrMap[field]; hasField {
				if s.matchString(value, searchValue, modifier) {
					return true
				}
			}
		}

		// Check address lines
		if lines, hasLines := addrMap["line"]; hasLines {
			if lineList, ok := lines.([]interface{}); ok {
				for _, line := range lineList {
					if s.matchString(line, searchValue, modifier) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (s *InMemorySearchStore) matchAddressField(addresses interface{}, field, searchValue, modifier string) bool {
	if addresses == nil {
		return false
	}

	addrList, ok := addresses.([]interface{})
	if !ok {
		return false
	}

	for _, addr := range addrList {
		addrMap, ok := addr.(map[string]interface{})
		if !ok {
			continue
		}

		if value, hasField := addrMap[field]; hasField {
			if s.matchString(value, searchValue, modifier) {
				return true
			}
		}
	}

	return false
}

func (s *InMemorySearchStore) matchContactPoint(telecoms interface{}, system, searchValue, modifier string) bool {
	if telecoms == nil {
		return false
	}

	telecomList, ok := telecoms.([]interface{})
	if !ok {
		return false
	}

	for _, telecom := range telecomList {
		telecomMap, ok := telecom.(map[string]interface{})
		if !ok {
			continue
		}

		// Check system matches
		if sys, hasSys := telecomMap["system"]; hasSys {
			if fmt.Sprintf("%v", sys) != system {
				continue
			}
		}

		// Check value
		if value, hasValue := telecomMap["value"]; hasValue {
			if s.matchToken(value, searchValue, modifier) {
				return true
			}
		}
	}

	return false
}

func (s *InMemorySearchStore) matchBoolean(value interface{}, searchValue string) bool {
	if value == nil {
		return false
	}

	boolValue := fmt.Sprintf("%v", value)
	return strings.EqualFold(boolValue, searchValue)
}

func (s *InMemorySearchStore) matchCode(code interface{}, searchValue, modifier string) bool {
	if code == nil {
		return false
	}

	// Handle CodeableConcept
	codeMap, ok := code.(map[string]interface{})
	if !ok {
		return s.matchToken(code, searchValue, modifier)
	}

	// Check coding array
	if coding, hasCoding := codeMap["coding"]; hasCoding {
		if codingList, ok := coding.([]interface{}); ok {
			for _, c := range codingList {
				codingMap, ok := c.(map[string]interface{})
				if !ok {
					continue
				}

				// Check code
				if codeValue, hasCode := codingMap["code"]; hasCode {
					if s.matchToken(codeValue, searchValue, modifier) {
						return true
					}
				}

				// Check display (for text modifier)
				if modifier == "text" {
					if display, hasDisplay := codingMap["display"]; hasDisplay {
						if s.matchString(display, searchValue, "") {
							return true
						}
					}
				}
			}
		}
	}

	// Check text
	if modifier == "text" {
		if text, hasText := codeMap["text"]; hasText {
			return s.matchString(text, searchValue, "")
		}
	}

	return false
}

func (s *InMemorySearchStore) matchCodeableConcept(concept interface{}, searchValue, modifier string) bool {
	return s.matchCode(concept, searchValue, modifier)
}

func (s *InMemorySearchStore) matchReference(ref interface{}, searchValue, modifier string) bool {
	if ref == nil {
		return false
	}

	refMap, ok := ref.(map[string]interface{})
	if !ok {
		return false
	}

	// Check reference field
	if reference, hasRef := refMap["reference"]; hasRef {
		refStr := fmt.Sprintf("%v", reference)

		switch modifier {
		case "identifier":
			// Check identifier field instead
			if ident, hasIdent := refMap["identifier"]; hasIdent {
				return s.matchToken(ident, searchValue, "")
			}
			return false
		default:
			// Match reference string (e.g., "Patient/123")
			return strings.EqualFold(refStr, searchValue) ||
				strings.HasSuffix(refStr, "/"+searchValue)
		}
	}

	return false
}

func (s *InMemorySearchStore) sortResults(results []models.Resource, sorts []SortParameter) []models.Resource {
	// Simplified sorting - in production would use more sophisticated sorting
	// For now, just return as-is
	return results
}

func (s *InMemorySearchStore) createBundle(resources []models.Resource, total int) *models.Bundle {
	bundle := &models.Bundle{
		ResourceType: "Bundle",
		Type:         "searchset",
		Total:        &total,
		Entry:        make([]models.BundleEntry, 0, len(resources)),
	}

	for i, resource := range resources {
		entry := models.BundleEntry{
			FullURL: fmt.Sprintf("http://localhost:8080/fhir/R5/%s/%s",
				resource.ResourceType(), resource.ID()),
			Resource: &resource,
			Search: &models.BundleEntrySearch{
				Mode:  "match",
				Score: float64(len(resources)-i) / float64(len(resources)),
			},
		}
		bundle.Entry = append(bundle.Entry, entry)
	}

	return bundle
}

// resourceToMap converts a resource to a map for searching
func resourceToMap(resource models.Resource) map[string]interface{} {
	// This is a simplified version - in production would use proper FHIR model reflection
	result := make(map[string]interface{})

	// Use JSON marshaling as a simple way to convert
	data, _ := json.Marshal(resource)
	json.Unmarshal(data, &result)

	return result
}
