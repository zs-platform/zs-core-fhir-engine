package analytics

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Dashboard provides analytics and metrics for the FHIR server
type Dashboard struct {
	metricsStore MetricsStore
	config       DashboardConfig
}

// MetricsStore defines the storage interface for metrics
type MetricsStore interface {
	StoreMetric(ctx context.Context, metric Metric) error
	QueryMetrics(ctx context.Context, query MetricQuery) ([]Metric, error)
	GetAggregates(ctx context.Context, tenantID string, period TimePeriod) (*Aggregates, error)
}

// Metric represents a single metric data point
type Metric struct {
	Timestamp    time.Time              `json:"timestamp"`
	Name         string                 `json:"name"`
	Value        float64                `json:"value"`
	Type         string                 `json:"type"` // counter, gauge, histogram
	TenantID     string                 `json:"tenantId"`
	UserID       string                 `json:"userId,omitempty"`
	ResourceType string                 `json:"resourceType,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// MetricQuery contains query parameters for metrics
type MetricQuery struct {
	Names        []string
	TenantID     string
	ResourceType string
	From         time.Time
	To           time.Time
	Aggregation  string // sum, avg, min, max
	Interval     time.Duration
	Limit        int
}

// TimePeriod represents a time period for aggregation
type TimePeriod struct {
	Start time.Time
	End   time.Time
}

// Aggregates contains aggregated metrics
type Aggregates struct {
	Period            TimePeriod       `json:"period"`
	TenantID          string           `json:"tenantId"`
	TotalRequests     int64            `json:"totalRequests"`
	TotalResources    int64            `json:"totalResources"`
	ResourceBreakdown map[string]int64 `json:"resourceBreakdown"`
	ResponseTimeAvg   float64          `json:"responseTimeAvg"`
	ResponseTimeP95   float64          `json:"responseTimeP95"`
	ResponseTimeP99   float64          `json:"responseTimeP99"`
	ErrorRate         float64          `json:"errorRate"`
	ActiveUsers       int64            `json:"activeUsers"`
	TopEndpoints      []EndpointStat   `json:"topEndpoints"`
	UsageByHour       map[int]int64    `json:"usageByHour"`
}

// EndpointStat represents endpoint usage statistics
type EndpointStat struct {
	Endpoint    string  `json:"endpoint"`
	Method      string  `json:"method"`
	Count       int64   `json:"count"`
	AvgDuration float64 `json:"avgDuration"`
	ErrorRate   float64 `json:"errorRate"`
}

// DashboardConfig contains dashboard configuration
type DashboardConfig struct {
	DefaultPeriod     time.Duration
	RetentionDays     int
	AggregationWindow time.Duration
}

// NewDashboard creates a new analytics dashboard
func NewDashboard(store MetricsStore, config DashboardConfig) *Dashboard {
	return &Dashboard{
		metricsStore: store,
		config:       config,
	}
}

// RecordMetric records a metric data point
func (d *Dashboard) RecordMetric(ctx context.Context, metric Metric) error {
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}

	return d.metricsStore.StoreMetric(ctx, metric)
}

// GetOverview returns a high-level dashboard overview
func (d *Dashboard) GetOverview(ctx context.Context, tenantID string, period TimePeriod) (*DashboardOverview, error) {
	aggregates, err := d.metricsStore.GetAggregates(ctx, tenantID, period)
	if err != nil {
		return nil, err
	}

	overview := &DashboardOverview{
		Period:      period,
		TenantID:    tenantID,
		Summary:     aggregates,
		GeneratedAt: time.Now(),
	}

	return overview, nil
}

// GetResourceStats returns resource-specific statistics
func (d *Dashboard) GetResourceStats(ctx context.Context, tenantID string, period TimePeriod) (*ResourceStats, error) {
	query := MetricQuery{
		TenantID: tenantID,
		From:     period.Start,
		To:       period.End,
	}

	metrics, err := d.metricsStore.QueryMetrics(ctx, query)
	if err != nil {
		return nil, err
	}

	stats := &ResourceStats{
		Period:    period,
		TenantID:  tenantID,
		Resources: make(map[string]ResourceMetrics),
	}

	// Aggregate by resource type
	for _, metric := range metrics {
		if metric.ResourceType == "" {
			continue
		}

		if _, exists := stats.Resources[metric.ResourceType]; !exists {
			stats.Resources[metric.ResourceType] = ResourceMetrics{
				ResourceType: metric.ResourceType,
			}
		}

		resourceMetrics := stats.Resources[metric.ResourceType]

		switch metric.Name {
		case "resource_created":
			resourceMetrics.Created++
		case "resource_read":
			resourceMetrics.Read++
		case "resource_updated":
			resourceMetrics.Updated++
		case "resource_deleted":
			resourceMetrics.Deleted++
		}

		stats.Resources[metric.ResourceType] = resourceMetrics
	}

	return stats, nil
}

// GetUsageTrends returns usage trends over time
func (d *Dashboard) GetUsageTrends(ctx context.Context, tenantID string, period TimePeriod, interval time.Duration) ([]TrendPoint, error) {
	query := MetricQuery{
		Names:    []string{"api_request"},
		TenantID: tenantID,
		From:     period.Start,
		To:       period.End,
		Interval: interval,
	}

	metrics, err := d.metricsStore.QueryMetrics(ctx, query)
	if err != nil {
		return nil, err
	}

	// Group by interval
	buckets := make(map[time.Time][]Metric)

	for _, metric := range metrics {
		bucket := metric.Timestamp.Truncate(interval)
		buckets[bucket] = append(buckets[bucket], metric)
	}

	// Convert to trend points
	var trends []TrendPoint
	for bucket, bucketMetrics := range buckets {
		point := TrendPoint{
			Timestamp: bucket,
			Count:     int64(len(bucketMetrics)),
		}

		// Calculate average response time
		var totalDuration float64
		for _, m := range bucketMetrics {
			if duration, ok := m.Metadata["duration_ms"].(float64); ok {
				totalDuration += duration
			}
		}

		if len(bucketMetrics) > 0 {
			point.AvgDuration = totalDuration / float64(len(bucketMetrics))
		}

		trends = append(trends, point)
	}

	// Sort by timestamp
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Timestamp.Before(trends[j].Timestamp)
	})

	return trends, nil
}

// GetTopUsers returns top users by activity
func (d *Dashboard) GetTopUsers(ctx context.Context, tenantID string, period TimePeriod, limit int) ([]UserStat, error) {
	query := MetricQuery{
		TenantID: tenantID,
		From:     period.Start,
		To:       period.End,
	}

	metrics, err := d.metricsStore.QueryMetrics(ctx, query)
	if err != nil {
		return nil, err
	}

	// Aggregate by user
	userStats := make(map[string]*UserStat)

	for _, metric := range metrics {
		if metric.UserID == "" {
			continue
		}

		if _, exists := userStats[metric.UserID]; !exists {
			userStats[metric.UserID] = &UserStat{
				UserID:            metric.UserID,
				ResourcesAccessed: make(map[string]int64),
			}
		}

		userStats[metric.UserID].RequestCount++

		if resourceType, ok := metric.Labels["resource_type"]; ok {
			userStats[metric.UserID].ResourcesAccessed[resourceType]++
		}
	}

	// Convert to slice and sort
	var results []UserStat
	for _, stat := range userStats {
		results = append(results, *stat)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].RequestCount > results[j].RequestCount
	})

	// Apply limit
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	return results, nil
}

// DashboardOverview represents a dashboard overview
type DashboardOverview struct {
	Period      TimePeriod  `json:"period"`
	TenantID    string      `json:"tenantId"`
	Summary     *Aggregates `json:"summary"`
	GeneratedAt time.Time   `json:"generatedAt"`
}

// ResourceStats represents resource statistics
type ResourceStats struct {
	Period    TimePeriod                 `json:"period"`
	TenantID  string                     `json:"tenantId"`
	Resources map[string]ResourceMetrics `json:"resources"`
}

// ResourceMetrics contains metrics for a specific resource type
type ResourceMetrics struct {
	ResourceType string `json:"resourceType"`
	Created      int64  `json:"created"`
	Read         int64  `json:"read"`
	Updated      int64  `json:"updated"`
	Deleted      int64  `json:"deleted"`
}

// TrendPoint represents a single point in a trend
type TrendPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Count       int64     `json:"count"`
	AvgDuration float64   `json:"avgDuration"`
}

// UserStat represents user statistics
type UserStat struct {
	UserID            string           `json:"userId"`
	RequestCount      int64            `json:"requestCount"`
	ResourcesAccessed map[string]int64 `json:"resourcesAccessed"`
}

// DashboardHandler handles HTTP requests for the dashboard
type DashboardHandler struct {
	dashboard *Dashboard
}

// NewDashboardHandler creates a new dashboard HTTP handler
func NewDashboardHandler(dashboard *Dashboard) *DashboardHandler {
	return &DashboardHandler{
		dashboard: dashboard,
	}
}

// RegisterRoutes registers dashboard endpoints
func (dh *DashboardHandler) RegisterRoutes(router chi.Router) {
	router.Get("/analytics/overview", dh.handleGetOverview)
	router.Get("/analytics/resources", dh.handleGetResourceStats)
	router.Get("/analytics/trends", dh.handleGetTrends)
	router.Get("/analytics/users", dh.handleGetTopUsers)
	router.Get("/analytics/metrics", dh.handleGetMetrics)
	router.Post("/analytics/metrics", dh.handleRecordMetric)
}

// handleGetOverview handles GET /analytics/overview
func (dh *DashboardHandler) handleGetOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")

	period := dh.parsePeriod(r)

	overview, err := dh.dashboard.GetOverview(ctx, tenantID, period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

// handleGetResourceStats handles GET /analytics/resources
func (dh *DashboardHandler) handleGetResourceStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")

	period := dh.parsePeriod(r)

	stats, err := dh.dashboard.GetResourceStats(ctx, tenantID, period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleGetTrends handles GET /analytics/trends
func (dh *DashboardHandler) handleGetTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")

	period := dh.parsePeriod(r)
	interval := dh.parseInterval(r)

	trends, err := dh.dashboard.GetUsageTrends(ctx, tenantID, period, interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trends)
}

// handleGetTopUsers handles GET /analytics/users
func (dh *DashboardHandler) handleGetTopUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")

	period := dh.parsePeriod(r)
	limit := 10

	users, err := dh.dashboard.GetTopUsers(ctx, tenantID, period, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// handleGetMetrics handles GET /analytics/metrics
func (dh *DashboardHandler) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")

	query := MetricQuery{
		TenantID: tenantID,
		From:     dh.parsePeriod(r).Start,
		To:       dh.parsePeriod(r).End,
	}

	metrics, err := dh.dashboard.metricsStore.QueryMetrics(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// handleRecordMetric handles POST /analytics/metrics
func (dh *DashboardHandler) handleRecordMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract tenant from header if not provided
	if metric.TenantID == "" {
		metric.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if err := dh.dashboard.RecordMetric(ctx, metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// parsePeriod parses time period from request
func (dh *DashboardHandler) parsePeriod(r *http.Request) TimePeriod {
	now := time.Now()
	period := TimePeriod{
		Start: now.Add(-24 * time.Hour), // Default: last 24 hours
		End:   now,
	}

	if from := r.URL.Query().Get("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			period.Start = t
		}
	}

	if to := r.URL.Query().Get("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			period.End = t
		}
	}

	return period
}

// parseInterval parses interval from request
func (dh *DashboardHandler) parseInterval(r *http.Request) time.Duration {
	interval := 1 * time.Hour // Default: 1 hour

	if i := r.URL.Query().Get("interval"); i != "" {
		if d, err := time.ParseDuration(i); err == nil {
			interval = d
		}
	}

	return interval
}

// InMemoryMetricsStore implements MetricsStore with in-memory storage
type InMemoryMetricsStore struct {
	metrics    []Metric
	aggregates map[string]map[time.Time]*Aggregates
	mu         sync.RWMutex
}

// NewInMemoryMetricsStore creates a new in-memory metrics store
func NewInMemoryMetricsStore() *InMemoryMetricsStore {
	return &InMemoryMetricsStore{
		metrics:    make([]Metric, 0),
		aggregates: make(map[string]map[time.Time]*Aggregates),
	}
}

// StoreMetric implements MetricsStore
func (s *InMemoryMetricsStore) StoreMetric(ctx context.Context, metric Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics = append(s.metrics, metric)

	// Keep only last 100000 metrics in memory
	if len(s.metrics) > 100000 {
		s.metrics = s.metrics[len(s.metrics)-100000:]
	}

	return nil
}

// QueryMetrics implements MetricsStore
func (s *InMemoryMetricsStore) QueryMetrics(ctx context.Context, query MetricQuery) ([]Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []Metric

	for _, metric := range s.metrics {
		// Apply filters
		if query.TenantID != "" && metric.TenantID != query.TenantID {
			continue
		}
		if query.ResourceType != "" && metric.ResourceType != query.ResourceType {
			continue
		}
		if !query.From.IsZero() && metric.Timestamp.Before(query.From) {
			continue
		}
		if !query.To.IsZero() && metric.Timestamp.After(query.To) {
			continue
		}

		if len(query.Names) > 0 {
			found := false
			for _, name := range query.Names {
				if metric.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		results = append(results, metric)
	}

	return results, nil
}

// GetAggregates implements MetricsStore
func (s *InMemoryMetricsStore) GetAggregates(ctx context.Context, tenantID string, period TimePeriod) (*Aggregates, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	aggregates := &Aggregates{
		Period:            period,
		TenantID:          tenantID,
		ResourceBreakdown: make(map[string]int64),
		UsageByHour:       make(map[int]int64),
	}

	// Calculate aggregates from metrics
	var totalDuration float64
	var errorCount int64

	for _, metric := range s.metrics {
		if metric.TenantID != tenantID {
			continue
		}
		if metric.Timestamp.Before(period.Start) || metric.Timestamp.After(period.End) {
			continue
		}

		switch metric.Name {
		case "api_request":
			aggregates.TotalRequests++

			if duration, ok := metric.Metadata["duration_ms"].(float64); ok {
				totalDuration += duration
			}

			if status, ok := metric.Metadata["status_code"].(float64); ok {
				if status >= 400 {
					errorCount++
				}
			}

		case "resource_created":
			aggregates.TotalResources++
			aggregates.ResourceBreakdown[metric.ResourceType]++
		}

		// Track usage by hour
		hour := metric.Timestamp.Hour()
		aggregates.UsageByHour[hour]++
	}

	// Calculate averages
	if aggregates.TotalRequests > 0 {
		aggregates.ResponseTimeAvg = totalDuration / float64(aggregates.TotalRequests)
		aggregates.ErrorRate = float64(errorCount) / float64(aggregates.TotalRequests)
	}

	return aggregates, nil
}
