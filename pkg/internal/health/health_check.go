package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HealthChecker handles health check endpoints
type HealthChecker struct {
	checks       map[string]HealthCheck
	mu           sync.RWMutex
	startTime    time.Time
	version      string
	buildInfo    BuildInfo
}

// BuildInfo contains application build information
type BuildInfo struct {
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildTime   string `json:"build_time"`
	GoVersion   string `json:"go_version"`
	Environment string `json:"environment"`
}

// HealthCheck represents a health check function
type HealthCheck interface {
	Name() string
	Check(ctx context.Context) HealthStatus
}

// HealthStatus represents the status of a health check
type HealthStatus struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Duration  time.Duration          `json:"duration"`
}

// DatabaseHealthCheck checks database connectivity
type DatabaseHealthCheck struct {
	db Database
}

// NewDatabaseHealthCheck creates a new database health check
func NewDatabaseHealthCheck(db Database) *DatabaseHealthCheck {
	return &DatabaseHealthCheck{db: db}
}

// Name returns the name of the health check
func (dhc *DatabaseHealthCheck) Name() string {
	return "database"
}

// Check performs the database health check
func (dhc *DatabaseHealthCheck) Check(ctx context.Context) HealthStatus {
	start := time.Now()
	
	// Check database connectivity
	if err := dhc.db.Ping(ctx); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Database connection failed: %v", err),
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		}
	}

	// Check database response time
	pingStart := time.Now()
	if err := dhc.db.Ping(ctx); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Database ping failed: %v", err),
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		}
	}
	pingDuration := time.Since(pingStart)

	status := "healthy"
	if pingDuration > 100*time.Millisecond {
		status = "degraded"
	}

	return HealthStatus{
		Status:    status,
		Message:   "Database is accessible",
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"ping_duration_ms": pingDuration.Milliseconds(),
			"max_connections":  dhc.db.MaxConnections(),
			"active_connections": dhc.db.ActiveConnections(),
		},
		Duration: time.Since(start),
	}
}

// Database interface for health checking
type Database interface {
	Ping(ctx context.Context) error
	MaxConnections() int
	ActiveConnections() int
}

// RedisHealthCheck checks Redis connectivity
type RedisHealthCheck struct {
	redis Redis
}

// NewRedisHealthCheck creates a new Redis health check
func NewRedisHealthCheck(redis Redis) *RedisHealthCheck {
	return &RedisHealthCheck{redis: redis}
}

// Name returns the name of the health check
func (rhc *RedisHealthCheck) Name() string {
	return "redis"
}

// Check performs the Redis health check
func (rhc *RedisHealthCheck) Check(ctx context.Context) HealthStatus {
	start := time.Now()
	
	// Check Redis connectivity
	if err := rhc.redis.Ping(ctx); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Redis connection failed: %v", err),
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		}
	}

	// Check Redis response time
	pingStart := time.Now()
	if err := rhc.redis.Ping(ctx); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Redis ping failed: %v", err),
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		}
	}
	pingDuration := time.Since(pingStart)

	status := "healthy"
	if pingDuration > 50*time.Millisecond {
		status = "degraded"
	}

	return HealthStatus{
		Status:    status,
		Message:   "Redis is accessible",
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"ping_duration_ms": pingDuration.Milliseconds(),
			"memory_usage":     rhc.redis.MemoryUsage(),
			"connected_clients": rhc.redis.ConnectedClients(),
		},
		Duration: time.Since(start),
	}
}

// Redis interface for health checking
type Redis interface {
	Ping(ctx context.Context) error
	MemoryUsage() int64
	ConnectedClients() int
}

// NATSHealthCheck checks NATS connectivity
type NATSHealthCheck struct {
	nats NATS
}

// NewNATSHealthCheck creates a new NATS health check
func NewNATSHealthCheck(nats NATS) *NATSHealthCheck {
	return &NATSHealthCheck{nats: nats}
}

// Name returns the name of the health check
func (nhc *NATSHealthCheck) Name() string {
	return "nats"
}

// Check performs the NATS health check
func (nhc *NATSHealthCheck) Check(ctx context.Context) HealthStatus {
	start := time.Now()
	
	// Check NATS connectivity
	if err := nhc.nats.Ping(ctx); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("NATS connection failed: %v", err),
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		}
	}

	return HealthStatus{
		Status:    "healthy",
		Message:   "NATS is accessible",
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"connected_servers": nhc.nats.ConnectedServers(),
			"jetstream_enabled":  nhc.nats.IsJetStreamEnabled(),
		},
		Duration: time.Since(start),
	}
}

// NATS interface for health checking
type NATS interface {
	Ping(ctx context.Context) error
	ConnectedServers() int
	IsJetStreamEnabled() bool
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(version string, buildInfo BuildInfo) *HealthChecker {
	return &HealthChecker{
		checks:    make(map[string]HealthCheck),
		startTime: time.Now(),
		version:   version,
		buildInfo: buildInfo,
	}
}

// RegisterCheck registers a health check
func (hc *HealthChecker) RegisterCheck(check HealthCheck) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[check.Name()] = check
}

// CheckHealth performs all registered health checks
func (hc *HealthChecker) CheckHealth(ctx context.Context) HealthStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	start := time.Now()
	checks := make(map[string]HealthStatus)
	overallStatus := "healthy"

	for name, check := range hc.checks {
		status := check.Check(ctx)
		checks[name] = status

		if status.Status == "unhealthy" {
			overallStatus = "unhealthy"
		} else if status.Status == "degraded" && overallStatus == "healthy" {
			overallStatus = "degraded"
		}
	}

	return HealthStatus{
		Status:    overallStatus,
		Message:   fmt.Sprintf("System is %s", overallStatus),
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"uptime_seconds": time.Since(hc.startTime).Seconds(),
			"version":        hc.version,
			"checks":         checks,
		},
		Duration: time.Since(start),
	}
}

// CheckHealthSimple performs a simple health check
func (hc *HealthChecker) CheckHealthSimple(ctx context.Context) map[string]interface{} {
	status := hc.CheckHealth(ctx)
	
	return map[string]interface{}{
		"status":    status.Status,
		"timestamp": status.Timestamp,
		"version":   hc.version,
		"uptime":    time.Since(hc.startTime).Seconds(),
	}
}

// GetBuildInfo returns build information
func (hc *HealthChecker) GetBuildInfo() BuildInfo {
	return hc.buildInfo
}

// HealthHandler handles HTTP health check requests
type HealthHandler struct {
	healthChecker *HealthChecker
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthChecker *HealthChecker) *HealthHandler {
	return &HealthHandler{
		healthChecker: healthChecker,
	}
}

// HandleHealth handles the /health endpoint
func (hh *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	healthStatus := hh.healthChecker.CheckHealthSimple(ctx)
	
	w.Header().Set("Content-Type", "application/json")
	if healthStatus["status"] == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else if healthStatus["status"] == "degraded" {
		w.WriteHeader(http.StatusOK) // Still 200 for degraded
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	json.NewEncoder(w).Encode(healthStatus)
}

// HandleHealthDetailed handles the /health/detailed endpoint
func (hh *HealthHandler) HandleHealthDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	healthStatus := hh.healthChecker.CheckHealth(ctx)
	
	w.Header().Set("Content-Type", "application/json")
	if healthStatus.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else if healthStatus.Status == "degraded" {
		w.WriteHeader(http.StatusOK) // Still 200 for degraded
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	json.NewEncoder(w).Encode(healthStatus)
}

// HandleReadiness handles the /ready endpoint
func (hh *HealthHandler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	healthStatus := hh.healthChecker.CheckHealth(ctx)
	
	w.Header().Set("Content-Type", "application/json")
	if healthStatus.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "not_ready",
			"timestamp": time.Now(),
		})
	}
}

// HandleLiveness handles the /live endpoint
func (hh *HealthHandler) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	// Liveness probe - just check if the application is running
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}

// MetricsCollector handles Prometheus metrics
type MetricsCollector struct {
	registry *prometheus.Registry
	
	// HTTP metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec
	
	// FHIR metrics
	fhirResourcesTotal   *prometheus.CounterVec
	fhirResourceDuration *prometheus.HistogramVec
	fhirValidationErrors *prometheus.CounterVec
	
	// System metrics
	systemGoroutines prometheus.Gauge
	systemMemory     prometheus.Gauge
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	registry := prometheus.NewRegistry()
	
	mc := &MetricsCollector{
		registry: registry,
		
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		
		httpResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"method", "endpoint"},
		),
		
		fhirResourcesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fhir_resources_total",
				Help: "Total number of FHIR resources processed",
			},
			[]string{"resource_type", "operation"},
		),
		
		fhirResourceDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "fhir_resource_duration_seconds",
				Help:    "FHIR resource processing duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"resource_type", "operation"},
		),
		
		fhirValidationErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fhir_validation_errors_total",
				Help: "Total number of FHIR validation errors",
			},
			[]string{"resource_type", "error_type"},
		),
		
		systemGoroutines: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_goroutines",
				Help: "Number of goroutines",
			},
		),
		
		systemMemory: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_memory_bytes",
				Help: "Memory usage in bytes",
			},
		),
	}
	
	// Register metrics
	registry.MustRegister(mc.httpRequestsTotal)
	registry.MustRegister(mc.httpRequestDuration)
	registry.MustRegister(mc.httpResponseSize)
	registry.MustRegister(mc.fhirResourcesTotal)
	registry.MustRegister(mc.fhirResourceDuration)
	registry.MustRegister(mc.fhirValidationErrors)
	registry.MustRegister(mc.systemGoroutines)
	registry.MustRegister(mc.systemMemory)
	
	return mc
}

// RecordHTTPRequest records an HTTP request
func (mc *MetricsCollector) RecordHTTPRequest(method, endpoint, status string, duration time.Duration, size int) {
	mc.httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	mc.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	mc.httpResponseSize.WithLabelValues(method, endpoint).Observe(float64(size))
}

// RecordFHIRResource records a FHIR resource operation
func (mc *MetricsCollector) RecordFHIRResource(resourceType, operation string, duration time.Duration) {
	mc.fhirResourcesTotal.WithLabelValues(resourceType, operation).Inc()
	mc.fhirResourceDuration.WithLabelValues(resourceType, operation).Observe(duration.Seconds())
}

// RecordFHIRValidationError records a FHIR validation error
func (mc *MetricsCollector) RecordFHIRValidationError(resourceType, errorType string) {
	mc.fhirValidationErrors.WithLabelValues(resourceType, errorType).Inc()
}

// UpdateSystemMetrics updates system metrics
func (mc *MetricsCollector) UpdateSystemMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	mc.systemGoroutines.Set(float64(runtime.NumGoroutine()))
	mc.systemMemory.Set(float64(m.Alloc))
}

// GetRegistry returns the Prometheus registry
func (mc *MetricsCollector) GetRegistry() *prometheus.Registry {
	return mc.registry
}

// MetricsHandler handles Prometheus metrics requests
type MetricsHandler struct {
	collector *MetricsCollector
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(collector *MetricsCollector) *MetricsHandler {
	return &MetricsHandler{
		collector: collector,
	}
}

// HandleMetrics handles the /metrics endpoint
func (mh *MetricsHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	// Update system metrics before serving
	mh.collector.UpdateSystemMetrics()
	
	// Use Prometheus HTTP handler
	promhttp.HandlerFor(mh.collector.GetRegistry(), promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

// Helper function to create health checker with default checks
func CreateHealthChecker(version string, buildInfo BuildInfo, db Database, redis Redis, nats NATS) *HealthChecker {
	hc := NewHealthChecker(version, buildInfo)
	
	// Register default health checks
	if db != nil {
		hc.RegisterCheck(NewDatabaseHealthCheck(db))
	}
	
	if redis != nil {
		hc.RegisterCheck(NewRedisHealthCheck(redis))
	}
	
	if nats != nil {
		hc.RegisterCheck(NewNATSHealthCheck(nats))
	}
	
	return hc
}

// Helper function to create metrics collector
func CreateMetricsCollector() *MetricsCollector {
	return NewMetricsCollector()
}
