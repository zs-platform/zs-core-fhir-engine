package observability

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// Metrics holds basic application metrics
type Metrics struct {
	requestCount      int64
	errorCount        int64
	activeConnections int64
	lastReset         time.Time
}

// Global metrics instance
var metrics = &Metrics{
	requestCount:      0,
	errorCount:        0,
	activeConnections: 0,
	lastReset:         time.Now(),
}

// IncrementRequests increments the request counter
func IncrementRequests() {
	atomic.AddInt64(&metrics.requestCount, 1)
}

// IncrementErrors increments the error counter
func IncrementErrors() {
	atomic.AddInt64(&metrics.errorCount, 1)
}

// AddConnections increments active connections
func AddConnections() {
	atomic.AddInt64(&metrics.activeConnections, 1)
}

// RemoveConnections decrements active connections
func RemoveConnections() {
	atomic.AddInt64(&metrics.activeConnections, -1)
}

// GetMetrics returns current metrics
func GetMetrics() map[string]any {
	return map[string]any{
		"requests_total":     atomic.LoadInt64(&metrics.requestCount),
		"errors_total":       atomic.LoadInt64(&metrics.errorCount),
		"active_connections": atomic.LoadInt64(&metrics.activeConnections),
		"uptime_seconds":     time.Since(metrics.lastReset).Seconds(),
	}
}

// LogRequest logs a request with basic metrics
func LogRequest(method, resourceType, path, statusCode, duration time.Duration) {
	log.Printf("FHIR %s %s %s -> %d (%v)",
		method, resourceType, path, statusCode, duration.Milliseconds())

	IncrementRequests()
	if statusCode >= 400 {
		IncrementErrors()
	}
}

// LogInfo logs informational message
func LogInfo(message string, args ...any) {
	log.Printf("INFO: %s", fmt.Sprintf(message, args...))
}

// LogError logs error message
func LogError(message string, args ...any) {
	log.Printf("ERROR: %s", fmt.Sprintf(message, args...))
	IncrementErrors()
}

// LogWarning logs warning message
func LogWarning(message string, args ...any) {
	log.Printf("WARNING: %s", fmt.Sprintf(message, args...))
}
