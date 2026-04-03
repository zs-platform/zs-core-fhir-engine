package models

import "encoding/json"

// Resource is the minimal contract used by the search package.
type Resource interface {
	ResourceType() string
	ID() string
}

// Bundle is a lightweight FHIR search bundle.
type Bundle struct {
	ResourceType string        `json:"resourceType"`
	Type         string        `json:"type"`
	Total        *int          `json:"total,omitempty"`
	Timestamp    string        `json:"timestamp,omitempty"`
	Entry        []BundleEntry `json:"entry,omitempty"`
}

// BundleEntry represents a single entry in a search bundle.
type BundleEntry struct {
	FullURL  string             `json:"fullUrl,omitempty"`
	Resource *Resource          `json:"resource,omitempty"`
	Search   *BundleEntrySearch `json:"search,omitempty"`
}

// BundleEntrySearch carries per-entry search metadata.
type BundleEntrySearch struct {
	Mode  string  `json:"mode,omitempty"`
	Score float64 `json:"score,omitempty"`
}

// RawResource is a small JSON-backed Resource implementation for tests and
// in-memory stores.
type RawResource struct {
	ResourceTypeValue string                 `json:"resourceType"`
	IDValue           string                 `json:"id"`
	Fields            map[string]interface{} `json:"-"`
}

func (r RawResource) ResourceType() string { return r.ResourceTypeValue }

func (r RawResource) ID() string { return r.IDValue }

func (r RawResource) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{}, len(r.Fields)+2)
	for k, v := range r.Fields {
		data[k] = v
	}
	data["resourceType"] = r.ResourceTypeValue
	data["id"] = r.IDValue
	return json.Marshal(data)
}
