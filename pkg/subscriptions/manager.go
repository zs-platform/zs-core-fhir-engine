package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
)

// SubscriptionManager manages FHIR subscriptions
type SubscriptionManager struct {
	subscriptions map[string]*Subscription
	handlers      map[string]EventHandler
	mu            sync.RWMutex
	eventBus      EventBus
}

// EventHandler handles subscription events
type EventHandler func(ctx context.Context, event *ResourceEvent) error

// EventBus defines the event bus interface
type EventBus interface {
	Publish(ctx context.Context, topic string, event *ResourceEvent) error
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}

// Subscription represents a FHIR subscription
type Subscription struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"` // active, error, off
	Reason      string                 `json:"reason"`
	Criteria    string                 `json:"criteria"`
	Channel     SubscriptionChannel    `json:"channel"`
	End         *time.Time             `json:"end,omitempty"`
	ReasonEnded string                 `json:"reasonEnded,omitempty"`
	Error       string                 `json:"error,omitempty"`
	TenantID    string                 `json:"tenantId"`
	UserID      string                 `json:"userId"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	EventCount  int                    `json:"eventCount"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SubscriptionChannel defines how notifications are sent
type SubscriptionChannel struct {
	Type    string            `json:"type"`    // rest-hook, websocket, email, sms, message
	Endpoint string           `json:"endpoint,omitempty"`
	Payload  string           `json:"payload"` // empty, fhir-json
	Header   map[string]string `json:"header,omitempty"`
}

// ResourceEvent represents a FHIR resource change event
type ResourceEvent struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	ResourceType string                 `json:"resourceType"`
	ResourceID   string                 `json:"resourceId"`
	Action       string                 `json:"action"` // create, update, delete
	Resource     *fhir.Resource         `json:"resource,omitempty"`
	Previous     *fhir.Resource         `json:"previous,omitempty"`
	Changes      []string               `json:"changes,omitempty"`
	Subscription string                 `json:"subscription,omitempty"`
	TenantID     string                 `json:"tenantId"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(eventBus EventBus) *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]*Subscription),
		handlers:      make(map[string]EventHandler),
		eventBus:      eventBus,
	}
}

// CreateSubscription creates a new subscription
func (sm *SubscriptionManager) CreateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error) {
	// Validate criteria
	if err := sm.validateCriteria(subscription.Criteria); err != nil {
		return nil, fmt.Errorf("invalid criteria: %w", err)
	}
	
	// Validate channel
	if err := sm.validateChannel(subscription.Channel); err != nil {
		return nil, fmt.Errorf("invalid channel: %w", err)
	}
	
	// Generate ID if not provided
	if subscription.ID == "" {
		subscription.ID = generateSubscriptionID()
	}
	
	subscription.Status = "active"
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	subscription.EventCount = 0
	
	// Store subscription
	sm.mu.Lock()
	sm.subscriptions[subscription.ID] = subscription
	sm.mu.Unlock()
	
	// Set up event handler
	topic := sm.criteriaToTopic(subscription.Criteria)
	handler := sm.createEventHandler(subscription)
	
	if err := sm.eventBus.Subscribe(ctx, topic, handler); err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}
	
	sm.mu.Lock()
	sm.handlers[subscription.ID] = handler
	sm.mu.Unlock()
	
	log.Infof("Created subscription: %s for criteria: %s", subscription.ID, subscription.Criteria)
	
	return subscription, nil
}

// GetSubscription retrieves a subscription by ID
func (sm *SubscriptionManager) GetSubscription(ctx context.Context, subscriptionID string) (*Subscription, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	subscription, exists := sm.subscriptions[subscriptionID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}
	
	return subscription, nil
}

// UpdateSubscription updates an existing subscription
func (sm *SubscriptionManager) UpdateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	existing, exists := sm.subscriptions[subscription.ID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscription.ID)
	}
	
	// Update fields
	if subscription.Reason != "" {
		existing.Reason = subscription.Reason
	}
	if subscription.Criteria != "" {
		// Unsubscribe from old topic
		oldTopic := sm.criteriaToTopic(existing.Criteria)
		sm.eventBus.Unsubscribe(ctx, oldTopic)
		
		// Update criteria
		existing.Criteria = subscription.Criteria
		
		// Subscribe to new topic
		newTopic := sm.criteriaToTopic(subscription.Criteria)
		handler := sm.createEventHandler(existing)
		sm.eventBus.Subscribe(ctx, newTopic, handler)
		sm.handlers[existing.ID] = handler
	}
	if subscription.Channel.Type != "" {
		existing.Channel = subscription.Channel
	}
	if subscription.End != nil {
		existing.End = subscription.End
	}
	
	existing.UpdatedAt = time.Now()
	
	return existing, nil
}

// DeleteSubscription deletes a subscription
func (sm *SubscriptionManager) DeleteSubscription(ctx context.Context, subscriptionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	subscription, exists := sm.subscriptions[subscriptionID]
	if !exists {
		return fmt.Errorf("subscription not found: %s", subscriptionID)
	}
	
	// Unsubscribe from event bus
	topic := sm.criteriaToTopic(subscription.Criteria)
	if err := sm.eventBus.Unsubscribe(ctx, topic); err != nil {
		log.Warnf("Failed to unsubscribe: %v", err)
	}
	
	delete(sm.subscriptions, subscriptionID)
	delete(sm.handlers, subscriptionID)
	
	log.Infof("Deleted subscription: %s", subscriptionID)
	
	return nil
}

// ListSubscriptions lists all subscriptions for a tenant
func (sm *SubscriptionManager) ListSubscriptions(ctx context.Context, tenantID string) ([]*Subscription, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	var results []*Subscription
	
	for _, subscription := range sm.subscriptions {
		if subscription.TenantID == tenantID {
			results = append(results, subscription)
		}
	}
	
	return results, nil
}

// PublishEvent publishes a resource event to subscribers
func (sm *SubscriptionManager) PublishEvent(ctx context.Context, event *ResourceEvent) error {
	topic := fmt.Sprintf("fhir.%s.%s", event.ResourceType, event.Action)
	
	return sm.eventBus.Publish(ctx, topic, event)
}

// validateCriteria validates subscription criteria
func (sm *SubscriptionManager) validateCriteria(criteria string) error {
	// Format should be: ResourceType?searchParams
	if criteria == "" {
		return fmt.Errorf("criteria is required")
	}
	
	parts := strings.Split(criteria, "?")
	if len(parts) < 1 {
		return fmt.Errorf("invalid criteria format")
	}
	
	validResourceTypes := []string{"Patient", "Observation", "Condition", "Medication", "Encounter", "Procedure", "DiagnosticReport", "Practitioner", "Organization"}
	
	resourceType := parts[0]
	valid := false
	for _, rt := range validResourceTypes {
		if rt == resourceType {
			valid = true
			break
		}
	}
	
	if !valid {
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	
	return nil
}

// validateChannel validates subscription channel
func (sm *SubscriptionManager) validateChannel(channel SubscriptionChannel) error {
	validTypes := []string{"rest-hook", "websocket", "email", "sms", "message"}
	
	valid := false
	for _, t := range validTypes {
		if t == channel.Type {
			valid = true
			break
		}
	}
	
	if !valid {
		return fmt.Errorf("unsupported channel type: %s", channel.Type)
	}
	
	if channel.Type == "rest-hook" && channel.Endpoint == "" {
		return fmt.Errorf("endpoint is required for rest-hook channel")
	}
	
	return nil
}

// criteriaToTopic converts subscription criteria to event bus topic
func (sm *SubscriptionManager) criteriaToTopic(criteria string) string {
	// Extract resource type from criteria
	parts := strings.Split(criteria, "?")
	resourceType := parts[0]
	
	// Return wildcard topic for all actions on this resource type
	return fmt.Sprintf("fhir.%s.*", resourceType)
}

// createEventHandler creates an event handler for a subscription
func (sm *SubscriptionManager) createEventHandler(subscription *Subscription) EventHandler {
	return func(ctx context.Context, event *ResourceEvent) error {
		// Check if subscription is still active
		if subscription.Status != "active" {
			return nil
		}
		
		// Check if subscription has expired
		if subscription.End != nil && time.Now().After(*subscription.End) {
			subscription.Status = "off"
			subscription.ReasonEnded = "expired"
			return nil
		}
		
		// Match criteria
		if !sm.matchesCriteria(event, subscription.Criteria) {
			return nil
		}
		
		// Add subscription ID to event
		event.Subscription = subscription.ID
		
		// Send notification based on channel type
		switch subscription.Channel.Type {
		case "rest-hook":
			return sm.sendRestHook(ctx, subscription, event)
		case "websocket":
			return sm.sendWebsocket(ctx, subscription, event)
		case "message":
			return sm.sendMessage(ctx, subscription, event)
		default:
			return fmt.Errorf("unsupported channel type: %s", subscription.Channel.Type)
		}
	}
}

// matchesCriteria checks if an event matches subscription criteria
func (sm *SubscriptionManager) matchesCriteria(event *ResourceEvent, criteria string) bool {
	// Parse criteria
	parts := strings.Split(criteria, "?")
	resourceType := parts[0]
	
	// Check resource type match
	if event.ResourceType != resourceType {
		return false
	}
	
	// If no search params, match all
	if len(parts) < 2 {
		return true
	}
	
	// Parse search params
	searchParams := parts[1]
	params := parseSearchParams(searchParams)
	
	// Apply filters (simplified - in production would use FHIR search engine)
	for key, value := range params {
		if !sm.matchesParam(event, key, value) {
			return false
		}
	}
	
	return true
}

// matchesParam checks if event matches a search parameter
func (sm *SubscriptionManager) matchesParam(event *ResourceEvent, key, value string) bool {
	// Simplified parameter matching
	// In production, this would use the FHIR search engine
	
	if key == "_lastUpdated" {
		// Check if event is recent enough
		return true
	}
	
	return true
}

// parseSearchParams parses URL query parameters
func parseSearchParams(params string) map[string]string {
	result := make(map[string]string)
	
	pairs := strings.Split(params, "&")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	
	return result
}

// sendRestHook sends a REST hook notification
func (sm *SubscriptionManager) sendRestHook(ctx context.Context, subscription *Subscription, event *ResourceEvent) error {
	// In production, this would make an HTTP POST to the endpoint
	log.Infof("Sending rest-hook to %s for subscription %s", subscription.Channel.Endpoint, subscription.ID)
	
	// Update event count
	subscription.EventCount++
	subscription.UpdatedAt = time.Now()
	
	return nil
}

// sendWebsocket sends a websocket notification
func (sm *SubscriptionManager) sendWebsocket(ctx context.Context, subscription *Subscription, event *ResourceEvent) error {
	// In production, this would send to websocket connection
	log.Infof("Sending websocket notification for subscription %s", subscription.ID)
	
	subscription.EventCount++
	subscription.UpdatedAt = time.Now()
	
	return nil
}

// sendMessage sends a message notification
func (sm *SubscriptionManager) sendMessage(ctx context.Context, subscription *Subscription, event *ResourceEvent) error {
	// In production, this would send to message queue
	log.Infof("Sending message for subscription %s", subscription.ID)
	
	subscription.EventCount++
	subscription.UpdatedAt = time.Now()
	
	return nil
}

// generateSubscriptionID generates a unique subscription ID
func generateSubscriptionID() string {
	return fmt.Sprintf("sub-%d", time.Now().UnixNano())
}

// SubscriptionHandler handles HTTP requests for subscriptions
type SubscriptionHandler struct {
	manager *SubscriptionManager
}

// NewSubscriptionHandler creates a new subscription HTTP handler
func NewSubscriptionHandler(manager *SubscriptionManager) *SubscriptionHandler {
	return &SubscriptionHandler{
		manager: manager,
	}
}

// RegisterRoutes registers subscription endpoints
func (sh *SubscriptionHandler) RegisterRoutes(router chi.Router) {
	router.Post("/fhir/R5/Subscription", sh.handleCreateSubscription)
	router.Get("/fhir/R5/Subscription/{subscriptionID}", sh.handleGetSubscription)
	router.Put("/fhir/R5/Subscription/{subscriptionID}", sh.handleUpdateSubscription)
	router.Delete("/fhir/R5/Subscription/{subscriptionID}", sh.handleDeleteSubscription)
	router.Get("/fhir/R5/Subscription", sh.handleListSubscriptions)
}

// handleCreateSubscription handles POST /Subscription
func (sh *SubscriptionHandler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var subscription Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Extract tenant and user from context
	subscription.TenantID = r.Header.Get("X-Tenant-ID")
	subscription.UserID = r.Header.Get("X-User-ID")
	
	result, err := sh.manager.CreateSubscription(ctx, &subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// handleGetSubscription handles GET /Subscription/{id}
func (sh *SubscriptionHandler) handleGetSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subscriptionID := chi.URLParam(r, "subscriptionID")
	
	subscription, err := sh.manager.GetSubscription(ctx, subscriptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(subscription)
}

// handleUpdateSubscription handles PUT /Subscription/{id}
func (sh *SubscriptionHandler) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subscriptionID := chi.URLParam(r, "subscriptionID")
	
	var subscription Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	subscription.ID = subscriptionID
	
	result, err := sh.manager.UpdateSubscription(ctx, &subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(result)
}

// handleDeleteSubscription handles DELETE /Subscription/{id}
func (sh *SubscriptionHandler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subscriptionID := chi.URLParam(r, "subscriptionID")
	
	if err := sh.manager.DeleteSubscription(ctx, subscriptionID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// handleListSubscriptions handles GET /Subscription
func (sh *SubscriptionHandler) handleListSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")
	
	subscriptions, err := sh.manager.ListSubscriptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(subscriptions)
}

// InMemoryEventBus implements EventBus with in-memory pub/sub
type InMemoryEventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// NewInMemoryEventBus creates a new in-memory event bus
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Publish implements EventBus
func (b *InMemoryEventBus) Publish(ctx context.Context, topic string, event *ResourceEvent) error {
	b.mu.RLock()
	handlers := b.subscribers[topic]
	b.mu.RUnlock()
	
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				log.Errorf("Event handler error: %v", err)
			}
		}(handler)
	}
	
	return nil
}

// Subscribe implements EventBus
func (b *InMemoryEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.subscribers[topic] = append(b.subscribers[topic], handler)
	return nil
}

// Unsubscribe implements EventBus
func (b *InMemoryEventBus) Unsubscribe(ctx context.Context, topic string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	delete(b.subscribers, topic)
	return nil
}
