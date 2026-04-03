package search

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/models"
)

// SearchParameter represents a FHIR search parameter
type SearchParameter struct {
	Name       string
	Type       string // string, token, number, date, quantity, reference, composite, uri
	Expression string
	Modifier   string
	Prefix     string
	Value      string
	Resource   string
}

// SearchQuery represents a parsed FHIR search query
type SearchQuery struct {
	ResourceType string
	Parameters   []SearchParameter
	Sort         []SortParameter
	Page         int
	Count        int
	Total        string // none, accurate, estimate
	Include      []IncludeParameter
	RevInclude   []RevIncludeParameter
}

// SortParameter represents a sort specification
type SortParameter struct {
	Name  string
	Order string // asc, desc
}

// IncludeParameter represents an _include specification
type IncludeParameter struct {
	Source string
	Path   string
	Target string
}

// RevIncludeParameter represents a _revinclude specification
type RevIncludeParameter struct {
	Source string
	Path   string
	Target string
}

// SearchEngine handles FHIR search operations
type SearchEngine struct {
	store ResourceStore
}

// ResourceStore defines the storage interface for search operations
type ResourceStore interface {
	Search(ctx context.Context, query SearchQuery) (*models.Bundle, error)
	SearchByParams(ctx context.Context, resourceType string, params map[string][]string) (*models.Bundle, error)
}

// NewSearchEngine creates a new search engine
func NewSearchEngine(store ResourceStore) *SearchEngine {
	return &SearchEngine{
		store: store,
	}
}

// ParseSearchQuery parses a raw URL query into a structured SearchQuery
func (se *SearchEngine) ParseSearchQuery(resourceType string, queryParams url.Values) (*SearchQuery, error) {
	query := &SearchQuery{
		ResourceType: resourceType,
		Parameters:   make([]SearchParameter, 0),
		Sort:         make([]SortParameter, 0),
		Include:      make([]IncludeParameter, 0),
		RevInclude:   make([]RevIncludeParameter, 0),
		Page:         1,
		Count:        20, // Default page size
		Total:        "none",
	}

	for key, values := range queryParams {
		if len(values) == 0 {
			continue
		}

		switch key {
		case "_sort":
			query.Sort = se.parseSort(values[0])
		case "_page":
			if page, err := strconv.Atoi(values[0]); err == nil && page > 0 {
				query.Page = page
			}
		case "_count":
			if count, err := strconv.Atoi(values[0]); err == nil && count > 0 {
				query.Count = count
			}
		case "_total":
			if values[0] == "accurate" || values[0] == "estimate" || values[0] == "none" {
				query.Total = values[0]
			}
		case "_include":
			query.Include = se.parseInclude(values)
		case "_revinclude":
			query.RevInclude = se.parseRevInclude(values)
		default:
			// Regular search parameters
			for _, value := range values {
				param := se.parseSearchParameter(resourceType, key, value)
				query.Parameters = append(query.Parameters, param)
			}
		}
	}

	return query, nil
}

// parseSearchParameter parses a single search parameter
func (se *SearchEngine) parseSearchParameter(resourceType, key, value string) SearchParameter {
	param := SearchParameter{
		Name:     key,
		Value:    value,
		Resource: resourceType,
	}

	// Check for modifier (e.g., name:exact, identifier:text)
	if idx := strings.Index(key, ":"); idx != -1 {
		param.Name = key[:idx]
		param.Modifier = key[idx+1:]
	}

	// Check for prefix on number/date/quantity parameters (e.g., gt, lt, eq, ne, ge, le, sa, eb)
	if isPrefixSearch(param.Name) {
		param.Prefix = se.extractPrefix(value)
		if param.Prefix != "" {
			param.Value = value[len(param.Prefix):]
		}
	}

	// Determine parameter type based on name and resource
	param.Type = se.determineParameterType(resourceType, param.Name)

	return param
}

// parseSort parses sort parameters
func (se *SearchEngine) parseSort(sortValue string) []SortParameter {
	sorts := make([]SortParameter, 0)
	fields := strings.Split(sortValue, ",")

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		sort := SortParameter{
			Name:  field,
			Order: "asc",
		}

		// Check for descending order (prefix with -)
		if strings.HasPrefix(field, "-") {
			sort.Name = field[1:]
			sort.Order = "desc"
		}

		sorts = append(sorts, sort)
	}

	return sorts
}

// parseInclude parses _include parameters
func (se *SearchEngine) parseInclude(values []string) []IncludeParameter {
	includes := make([]IncludeParameter, 0)

	for _, value := range values {
		parts := strings.Split(value, ":")
		if len(parts) >= 2 {
			include := IncludeParameter{
				Source: parts[0],
				Path:   parts[1],
			}
			if len(parts) >= 3 {
				include.Target = parts[2]
			}
			includes = append(includes, include)
		}
	}

	return includes
}

// parseRevInclude parses _revinclude parameters
func (se *SearchEngine) parseRevInclude(values []string) []RevIncludeParameter {
	revIncludes := make([]RevIncludeParameter, 0)

	for _, value := range values {
		parts := strings.Split(value, ":")
		if len(parts) >= 2 {
			revInclude := RevIncludeParameter{
				Source: parts[0],
				Path:   parts[1],
			}
			if len(parts) >= 3 {
				revInclude.Target = parts[2]
			}
			revIncludes = append(revIncludes, revInclude)
		}
	}

	return revIncludes
}

// extractPrefix extracts the prefix from a search value
func (se *SearchEngine) extractPrefix(value string) string {
	prefixes := []string{"eq", "ne", "gt", "lt", "ge", "le", "sa", "eb", "ap"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return prefix
		}
	}

	return ""
}

// isPrefixSearch determines if a parameter supports prefix matching
func isPrefixSearch(paramName string) bool {
	prefixParams := []string{"birthdate", "death-date", "date", "effective", "issued",
		"value-date", "occurrence", " onset", "abatement", "recorded", "authored", "when", "start", "end"}

	for _, p := range prefixParams {
		if strings.Contains(paramName, p) {
			return true
		}
	}

	return false
}

// determineParameterType determines the FHIR search parameter type
func (se *SearchEngine) determineParameterType(resourceType, paramName string) string {
	// Common parameter types across resources
	commonParams := map[string]string{
		"_id":          "token",
		"_lastUpdated": "date",
		"_tag":         "token",
		"_profile":     "uri",
		"_security":    "token",
		"_source":      "uri",
		"_text":        "string",
		"_content":     "string",
		"_list":        "reference",
		"_has":         "reference",
		"_type":        "token",
	}

	if paramType, ok := commonParams[paramName]; ok {
		return paramType
	}

	// Resource-specific parameter types
	resourceParams := se.getResourceSearchParams(resourceType)
	if paramType, ok := resourceParams[paramName]; ok {
		return paramType
	}

	// Default to string if unknown
	return "string"
}

// getResourceSearchParams returns search parameters for specific resource types
func (se *SearchEngine) getResourceSearchParams(resourceType string) map[string]string {
	params := map[string]map[string]string{
		"Patient": {
			"identifier":           "token",
			"name":                 "string",
			"family":               "string",
			"given":                "string",
			"birthdate":            "date",
			"gender":               "token",
			"address":              "string",
			"address-city":         "string",
			"address-state":        "string",
			"address-postalcode":   "string",
			"address-country":      "string",
			"phone":                "token",
			"email":                "token",
			"active":               "token",
			"deceased":             "token",
			"organization":         "reference",
			"general-practitioner": "reference",
			"link":                 "reference",
		},
		"Observation": {
			"identifier":               "token",
			"status":                   "token",
			"category":                 "token",
			"code":                     "token",
			"date":                     "date",
			"subject":                  "reference",
			"patient":                  "reference",
			"encounter":                "reference",
			"performer":                "reference",
			"value-string":             "string",
			"value-quantity":           "quantity",
			"value-concept":            "token",
			"value-date":               "date",
			"component-code":           "token",
			"component-value-quantity": "quantity",
			"component-value-concept":  "token",
			"has-member":               "reference",
			"derived-from":             "reference",
		},
		"Condition": {
			"identifier":          "token",
			"clinical-status":     "token",
			"verification-status": "token",
			"category":            "token",
			"severity":            "token",
			"code":                "token",
			"body-site":           "token",
			"subject":             "reference",
			"patient":             "reference",
			"encounter":           "reference",
			"asserter":            "reference",
			"onset-date":          "date",
			"abatement-date":      "date",
			"recorded-date":       "date",
			"stage":               "token",
			"stage-type":          "token",
		},
		"MedicationRequest": {
			"identifier":       "token",
			"status":           "token",
			"intent":           "token",
			"category":         "token",
			"priority":         "token",
			"medication":       "reference",
			"medication.code":  "token",
			"subject":          "reference",
			"patient":          "reference",
			"encounter":        "reference",
			"authoredon":       "date",
			"requester":        "reference",
			"performer":        "reference",
			"recorder":         "reference",
			"reason-code":      "token",
			"reason-reference": "reference",
			"based-on":         "reference",
			"group-identifier": "token",
		},
	}

	if resourceParams, ok := params[resourceType]; ok {
		return resourceParams
	}

	return make(map[string]string)
}

// ExecuteSearch executes a search query and returns results
func (se *SearchEngine) ExecuteSearch(ctx context.Context, query *SearchQuery) (*models.Bundle, error) {
	return se.store.Search(ctx, *query)
}

// ValidateSearchQuery validates a search query for correctness
func (se *SearchEngine) ValidateSearchQuery(query *SearchQuery) error {
	// Validate resource type
	validResources := []string{"Patient", "Observation", "Condition", "MedicationRequest",
		"Encounter", "Procedure", "DiagnosticReport", "Practitioner", "Organization"}

	isValid := false
	for _, resource := range validResources {
		if query.ResourceType == resource {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("unsupported resource type: %s", query.ResourceType)
	}

	// Validate parameters
	for _, param := range query.Parameters {
		if err := se.validateParameter(param); err != nil {
			return err
		}
	}

	return nil
}

// validateParameter validates a single search parameter
func (se *SearchEngine) validateParameter(param SearchParameter) error {
	// Validate based on parameter type
	switch param.Type {
	case "date":
		return se.validateDateParameter(param)
	case "number":
		return se.validateNumberParameter(param)
	case "token":
		return se.validateTokenParameter(param)
	case "quantity":
		return se.validateQuantityParameter(param)
	case "reference":
		return se.validateReferenceParameter(param)
	}

	return nil
}

// validateDateParameter validates date parameters
func (se *SearchEngine) validateDateParameter(param SearchParameter) error {
	validPrefixes := []string{"eq", "ne", "gt", "lt", "ge", "le", "sa", "eb", "ap"}

	if param.Prefix != "" {
		isValid := false
		for _, prefix := range validPrefixes {
			if param.Prefix == prefix {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid date prefix: %s", param.Prefix)
		}
	}

	// Try to parse the date value
	if _, err := time.Parse(time.RFC3339, param.Value); err != nil {
		// Try other date formats
		if _, err := time.Parse("2006-01-02", param.Value); err != nil {
			return fmt.Errorf("invalid date format: %s", param.Value)
		}
	}

	return nil
}

// validateNumberParameter validates number parameters
func (se *SearchEngine) validateNumberParameter(param SearchParameter) error {
	if param.Prefix != "" {
		validPrefixes := []string{"eq", "ne", "gt", "lt", "ge", "le", "ap"}
		isValid := false
		for _, prefix := range validPrefixes {
			if param.Prefix == prefix {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid number prefix: %s", param.Prefix)
		}
	}

	// Try to parse as number
	if _, err := strconv.ParseFloat(param.Value, 64); err != nil {
		return fmt.Errorf("invalid number value: %s", param.Value)
	}

	return nil
}

// validateTokenParameter validates token parameters
func (se *SearchEngine) validateTokenParameter(param SearchParameter) error {
	validModifiers := []string{"", "text", "not", "above", "below", "in", "not-in", "ofType"}

	isValid := false
	for _, mod := range validModifiers {
		if param.Modifier == mod {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid token modifier: %s", param.Modifier)
	}

	return nil
}

// validateQuantityParameter validates quantity parameters
func (se *SearchEngine) validateQuantityParameter(param SearchParameter) error {
	// Format should be: prefixvalue|system|code or prefixvalue|unit
	parts := strings.Split(param.Value, "|")

	if len(parts) < 2 {
		return fmt.Errorf("invalid quantity format: %s", param.Value)
	}

	// Validate the numeric value part
	valuePart := parts[0]
	if param.Prefix != "" {
		valuePart = valuePart[len(param.Prefix):]
	}

	if _, err := strconv.ParseFloat(valuePart, 64); err != nil {
		return fmt.Errorf("invalid quantity value: %s", valuePart)
	}

	return nil
}

// validateReferenceParameter validates reference parameters
func (se *SearchEngine) validateReferenceParameter(param SearchParameter) error {
	validModifiers := []string{"", "identifier", "above", "below"}

	isValid := false
	for _, mod := range validModifiers {
		if param.Modifier == mod {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid reference modifier: %s", param.Modifier)
	}

	return nil
}
