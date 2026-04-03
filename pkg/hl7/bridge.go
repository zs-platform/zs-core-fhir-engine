package hl7

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Bridge handles conversion between HL7 v2 and FHIR
type Bridge struct {
	config      BridgeConfig
	transformer *Transformer
	store       MessageStore
}

// BridgeConfig contains bridge configuration
type BridgeConfig struct {
	Enabled           bool
	SupportedVersions []string
	DefaultMapping    string
	ValidateMessages  bool
	StoreRawHL7       bool
}

// MessageStore stores HL7 messages
type MessageStore interface {
	StoreMessage(ctx context.Context, msg *HL7Message) error
	GetMessage(ctx context.Context, messageID string) (*HL7Message, error)
	ListMessages(ctx context.Context, options ListOptions) ([]*HL7Message, error)
}

// HL7Message represents an HL7 v2 message
type HL7Message struct {
	ID             string                 `json:"id"`
	Raw            string                 `json:"raw"`
	Type           string                 `json:"type"` // ADT, ORM, ORU, MDM, etc.
	Version        string                 `json:"version"`
	Segments       []Segment              `json:"segments"`
	Status         string                 `json:"status"` // received, processing, transformed, error
	Error          string                 `json:"error,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	ProcessedAt    *time.Time             `json:"processedAt,omitempty"`
	FHIRResourceID string                 `json:"fhirResourceId,omitempty"`
	TenantID       string                 `json:"tenantId"`
	Source         string                 `json:"source"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Segment represents an HL7 segment
type Segment struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
}

// ListOptions contains options for listing messages
type ListOptions struct {
	Type     string
	Status   string
	TenantID string
	From     time.Time
	To       time.Time
	Limit    int
	Offset   int
}

// NewBridge creates a new HL7 bridge
func NewBridge(config BridgeConfig, store MessageStore) *Bridge {
	return &Bridge{
		config:      config,
		transformer: NewTransformer(),
		store:       store,
	}
}

// TransformMessage transforms an HL7 message to FHIR
func (b *Bridge) TransformMessage(ctx context.Context, msg *HL7Message) (*TransformationResult, error) {
	if !b.config.Enabled {
		return nil, fmt.Errorf("HL7 bridge is disabled")
	}

	// Validate message
	if b.config.ValidateMessages {
		if err := b.validateMessage(msg); err != nil {
			msg.Status = "error"
			msg.Error = fmt.Sprintf("validation failed: %v", err)
			return nil, err
		}
	}

	// Parse segments
	segments, err := b.parseSegments(msg.Raw)
	if err != nil {
		msg.Status = "error"
		msg.Error = fmt.Sprintf("parse failed: %v", err)
		return nil, err
	}

	msg.Segments = segments
	msg.Type = b.determineMessageType(segments)

	// Transform to FHIR
	result, err := b.transformer.Transform(ctx, msg)
	if err != nil {
		msg.Status = "error"
		msg.Error = fmt.Sprintf("transformation failed: %v", err)
		return nil, err
	}

	msg.Status = "transformed"
	now := time.Now()
	msg.ProcessedAt = &now
	msg.FHIRResourceID = result.ResourceID

	return result, nil
}

// validateMessage validates an HL7 message
func (b *Bridge) validateMessage(msg *HL7Message) error {
	if msg.Raw == "" {
		return fmt.Errorf("message is empty")
	}

	// Check for valid HL7 structure
	if !strings.HasPrefix(msg.Raw, "MSH|") {
		return fmt.Errorf("message must start with MSH segment")
	}

	return nil
}

// parseSegments parses an HL7 message into segments
func (b *Bridge) parseSegments(raw string) ([]Segment, error) {
	segments := make([]Segment, 0)

	// Split by segment delimiter
	segmentTexts := strings.Split(raw, "\r")

	for _, segText := range segmentTexts {
		segText = strings.TrimSpace(segText)
		if segText == "" {
			continue
		}

		// Parse segment
		fields := strings.Split(segText, "|")
		if len(fields) == 0 {
			continue
		}

		segment := Segment{
			Name:   fields[0],
			Fields: fields[1:],
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

// determineMessageType determines the HL7 message type from segments
func (b *Bridge) determineMessageType(segments []Segment) string {
	// Find MSH segment
	for _, seg := range segments {
		if seg.Name == "MSH" && len(seg.Fields) >= 8 {
			// Message type is in field 8 (index 7)
			return seg.Fields[7]
		}
	}

	return "Unknown"
}

// ProcessIncoming processes an incoming HL7 message
func (b *Bridge) ProcessIncoming(ctx context.Context, raw string, tenantID, source string) (*TransformationResult, error) {
	// Create message
	msg := &HL7Message{
		ID:        generateMessageID(),
		Raw:       raw,
		Status:    "received",
		CreatedAt: time.Now(),
		TenantID:  tenantID,
		Source:    source,
	}

	// Store raw message if configured
	if b.config.StoreRawHL7 {
		if err := b.store.StoreMessage(ctx, msg); err != nil {
			return nil, fmt.Errorf("failed to store message: %w", err)
		}
	}

	// Transform message
	result, err := b.TransformMessage(ctx, msg)
	if err != nil {
		// Update status on error
		if b.config.StoreRawHL7 {
			b.store.StoreMessage(ctx, msg)
		}
		return nil, err
	}

	// Update stored message
	if b.config.StoreRawHL7 {
		b.store.StoreMessage(ctx, msg)
	}

	return result, nil
}

// TransformationResult contains the result of a transformation
type TransformationResult struct {
	ResourceType  string                 `json:"resourceType"`
	ResourceID    string                 `json:"resourceId"`
	Resource      map[string]interface{} `json:"resource"`
	MessageType   string                 `json:"messageType"`
	TransformTime int64                  `json:"transformTimeMs"`
	Warnings      []string               `json:"warnings,omitempty"`
}

// Transformer handles HL7 to FHIR transformation
type Transformer struct {
	mappings map[string]MessageMapping
}

// MessageMapping defines how to map an HL7 message to FHIR
type MessageMapping struct {
	MessageType     string
	ResourceType    string
	SegmentMappings map[string]SegmentMapping
}

// SegmentMapping defines how to map an HL7 segment to FHIR
type SegmentMapping struct {
	TargetPath string
	Transform  func(fields []string) (interface{}, error)
}

// NewTransformer creates a new transformer
func NewTransformer() *Transformer {
	t := &Transformer{
		mappings: make(map[string]MessageMapping),
	}

	// Register standard mappings
	t.registerStandardMappings()

	return t
}

// Transform transforms an HL7 message to FHIR
func (t *Transformer) Transform(ctx context.Context, msg *HL7Message) (*TransformationResult, error) {
	start := time.Now()

	mapping, exists := t.mappings[msg.Type]
	if !exists {
		return nil, fmt.Errorf("no mapping found for message type: %s", msg.Type)
	}

	// Create FHIR resource
	resource := make(map[string]interface{})
	resource["resourceType"] = mapping.ResourceType
	resource["id"] = generateResourceID()

	// Apply segment mappings
	warnings := make([]string, 0)

	for _, seg := range msg.Segments {
		if segMapping, exists := mapping.SegmentMappings[seg.Name]; exists {
			value, err := segMapping.Transform(seg.Fields)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("Failed to transform %s: %v", seg.Name, err))
				continue
			}

			// Set value at target path
			t.setAtPath(resource, segMapping.TargetPath, value)
		}
	}

	return &TransformationResult{
		ResourceType:  mapping.ResourceType,
		ResourceID:    resource["id"].(string),
		Resource:      resource,
		MessageType:   msg.Type,
		TransformTime: time.Since(start).Milliseconds(),
		Warnings:      warnings,
	}, nil
}

// setAtPath sets a value at a nested path in a map
func (t *Transformer) setAtPath(resource map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	current := resource

	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
			return
		}

		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return
		}
	}
}

// registerStandardMappings registers standard HL7 to FHIR mappings
func (t *Transformer) registerStandardMappings() {
	// ADT^A01 - Admit/Visit Notification -> Patient + Encounter
	t.mappings["ADT^A01"] = MessageMapping{
		MessageType:  "ADT^A01",
		ResourceType: "Bundle",
		SegmentMappings: map[string]SegmentMapping{
			"PID": {
				TargetPath: "entry.0.resource",
				Transform:  t.transformPID,
			},
			"PV1": {
				TargetPath: "entry.1.resource",
				Transform:  t.transformPV1,
			},
		},
	}

	// ORU^R01 - Observation Result -> Observation
	t.mappings["ORU^R01"] = MessageMapping{
		MessageType:  "ORU^R01",
		ResourceType: "Observation",
		SegmentMappings: map[string]SegmentMapping{
			"OBR": {
				TargetPath: "code",
				Transform:  t.transformOBR,
			},
			"OBX": {
				TargetPath: "valueQuantity",
				Transform:  t.transformOBX,
			},
		},
	}

	// MDM^T02 - Medical Document -> DocumentReference
	t.mappings["MDM^T02"] = MessageMapping{
		MessageType:  "MDM^T02",
		ResourceType: "DocumentReference",
		SegmentMappings: map[string]SegmentMapping{
			"TXA": {
				TargetPath: "type",
				Transform:  t.transformTXA,
			},
		},
	}
}

// Transform functions for different segments

func (t *Transformer) transformPID(fields []string) (interface{}, error) {
	// Simplified PID to Patient transformation
	patient := map[string]interface{}{
		"resourceType": "Patient",
		"identifier": []map[string]interface{}{
			{
				"system": "http://hl7.org/fhir/sid/us-ssn",
				"value":  getField(fields, 2),
			},
		},
		"name": []map[string]interface{}{
			{
				"family": getField(fields, 4),
				"given":  []string{getField(fields, 5)},
			},
		},
		"birthDate": getField(fields, 6),
		"gender":    mapGender(getField(fields, 7)),
	}

	return patient, nil
}

func (t *Transformer) transformPV1(fields []string) (interface{}, error) {
	encounter := map[string]interface{}{
		"resourceType": "Encounter",
		"status":       "in-progress",
		"class": map[string]interface{}{
			"system": "http://terminology.hl7.org/CodeSystem/v3-ActCode",
			"code":   getField(fields, 1),
		},
	}

	return encounter, nil
}

func (t *Transformer) transformOBR(fields []string) (interface{}, error) {
	code := map[string]interface{}{
		"coding": []map[string]interface{}{
			{
				"system":  "http://loinc.org",
				"code":    getField(fields, 3),
				"display": getField(fields, 4),
			},
		},
	}

	return code, nil
}

func (t *Transformer) transformOBX(fields []string) (interface{}, error) {
	value := map[string]interface{}{
		"value": getField(fields, 4),
		"unit":  getField(fields, 5),
	}

	return value, nil
}

func (t *Transformer) transformTXA(fields []string) (interface{}, error) {
	docType := map[string]interface{}{
		"coding": []map[string]interface{}{
			{
				"system":  "http://loinc.org",
				"code":    getField(fields, 1),
				"display": getField(fields, 1),
			},
		},
	}

	return docType, nil
}

// Helper functions

func getField(fields []string, index int) string {
	if index < len(fields) {
		return fields[index]
	}
	return ""
}

func mapGender(hl7Gender string) string {
	switch hl7Gender {
	case "M":
		return "male"
	case "F":
		return "female"
	case "U":
		return "unknown"
	default:
		return "unknown"
	}
}

func generateMessageID() string {
	return fmt.Sprintf("hl7-%d", time.Now().UnixNano())
}

func generateResourceID() string {
	return fmt.Sprintf("fhir-%d", time.Now().UnixNano())
}

// BridgeHandler handles HTTP requests for the HL7 bridge
type BridgeHandler struct {
	bridge *Bridge
}

// NewBridgeHandler creates a new bridge HTTP handler
func NewBridgeHandler(bridge *Bridge) *BridgeHandler {
	return &BridgeHandler{
		bridge: bridge,
	}
}

// RegisterRoutes registers HL7 bridge endpoints
func (bh *BridgeHandler) RegisterRoutes(router chi.Router) {
	router.Post("/hl7/v2/message", bh.handleReceiveMessage)
	router.Get("/hl7/v2/messages", bh.handleListMessages)
	router.Get("/hl7/v2/messages/{messageID}", bh.handleGetMessage)
	router.Post("/hl7/v2/transform", bh.handleTransform)
	router.Get("/hl7/v2/config", bh.handleGetConfig)
}

// handleReceiveMessage handles POST /hl7/v2/message
func (bh *BridgeHandler) handleReceiveMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Message string `json:"message"`
		Source  string `json:"source"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")

	result, err := bh.bridge.ProcessIncoming(ctx, req.Message, tenantID, req.Source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// handleListMessages handles GET /hl7/v2/messages
func (bh *BridgeHandler) handleListMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")

	options := ListOptions{
		TenantID: tenantID,
	}

	messages, err := bh.bridge.store.ListMessages(ctx, options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// handleGetMessage handles GET /hl7/v2/messages/{messageID}
func (bh *BridgeHandler) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	messageID := chi.URLParam(r, "messageID")

	message, err := bh.bridge.store.GetMessage(ctx, messageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// handleTransform handles POST /hl7/v2/transform
func (bh *BridgeHandler) handleTransform(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := &HL7Message{
		Raw:       req.Message,
		CreatedAt: time.Now(),
	}

	result, err := bh.bridge.TransformMessage(ctx, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleGetConfig handles GET /hl7/v2/config
func (bh *BridgeHandler) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"enabled":           bh.bridge.config.Enabled,
		"supportedVersions": bh.bridge.config.SupportedVersions,
		"storeRawHL7":       bh.bridge.config.StoreRawHL7,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// InMemoryMessageStore implements MessageStore with in-memory storage
type InMemoryMessageStore struct {
	messages map[string]*HL7Message
	mu       sync.RWMutex
}

// NewInMemoryMessageStore creates a new in-memory message store
func NewInMemoryMessageStore() *InMemoryMessageStore {
	return &InMemoryMessageStore{
		messages: make(map[string]*HL7Message),
	}
}

// StoreMessage implements MessageStore
func (s *InMemoryMessageStore) StoreMessage(ctx context.Context, msg *HL7Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[msg.ID] = msg
	return nil
}

// GetMessage implements MessageStore
func (s *InMemoryMessageStore) GetMessage(ctx context.Context, messageID string) (*HL7Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg, exists := s.messages[messageID]
	if !exists {
		return nil, fmt.Errorf("message not found: %s", messageID)
	}

	return msg, nil
}

// ListMessages implements MessageStore
func (s *InMemoryMessageStore) ListMessages(ctx context.Context, options ListOptions) ([]*HL7Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*HL7Message

	for _, msg := range s.messages {
		if options.TenantID != "" && msg.TenantID != options.TenantID {
			continue
		}
		if options.Type != "" && msg.Type != options.Type {
			continue
		}
		if options.Status != "" && msg.Status != options.Status {
			continue
		}

		results = append(results, msg)
	}

	return results, nil
}
