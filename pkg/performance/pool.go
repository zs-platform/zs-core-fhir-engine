package performance

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// ConnectionPool manages a pool of reusable connections
type ConnectionPool struct {
	maxSize     int
	idleTimeout time.Duration
	connections chan *PooledConnection
	factory     ConnectionFactory
	mu          sync.RWMutex
	stats       PoolStats
}

// ConnectionFactory creates new connections
type ConnectionFactory interface {
	NewConnection(ctx context.Context) (*PooledConnection, error)
	CloseConnection(conn *PooledConnection) error
}

// PooledConnection wraps a connection with metadata
type PooledConnection struct {
	Conn      interface{}
	CreatedAt time.Time
	LastUsed  time.Time
	InUse     bool
}

// PoolStats contains pool statistics
type PoolStats struct {
	TotalConnections  int64
	ActiveConnections int64
	IdleConnections   int64
	WaitCount         int64
	WaitDuration      time.Duration
	MaxWaitDuration   time.Duration
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(maxSize int, idleTimeout time.Duration, factory ConnectionFactory) *ConnectionPool {
	return &ConnectionPool{
		maxSize:     maxSize,
		idleTimeout: idleTimeout,
		connections: make(chan *PooledConnection, maxSize),
		factory:     factory,
	}
}

// Acquire gets a connection from the pool
func (p *ConnectionPool) Acquire(ctx context.Context) (*PooledConnection, error) {
	start := time.Now()

	select {
	case conn := <-p.connections:
		if time.Since(conn.LastUsed) > p.idleTimeout {
			// Connection expired, create new one
			p.factory.CloseConnection(conn)
			return p.factory.NewConnection(ctx)
		}

		conn.InUse = true
		conn.LastUsed = time.Now()

		p.mu.Lock()
		p.stats.ActiveConnections++
		p.stats.IdleConnections--
		p.mu.Unlock()

		log.Debugf("Acquired connection from pool (idle: %d)", len(p.connections))
		return conn, nil

	case <-ctx.Done():
		p.mu.Lock()
		p.stats.WaitCount++
		waitTime := time.Since(start)
		p.stats.WaitDuration += waitTime
		if waitTime > p.stats.MaxWaitDuration {
			p.stats.MaxWaitDuration = waitTime
		}
		p.mu.Unlock()

		return nil, ctx.Err()

	default:
		// Pool empty, create new connection
		p.mu.Lock()
		if p.stats.TotalConnections >= int64(p.maxSize) {
			p.mu.Unlock()
			// Wait for available connection
			select {
			case conn := <-p.connections:
				conn.InUse = true
				conn.LastUsed = time.Now()

				p.mu.Lock()
				p.stats.ActiveConnections++
				p.stats.IdleConnections--
				p.stats.WaitCount++
				waitTime := time.Since(start)
				p.stats.WaitDuration += waitTime
				if waitTime > p.stats.MaxWaitDuration {
					p.stats.MaxWaitDuration = waitTime
				}
				p.mu.Unlock()

				return conn, nil

			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		p.stats.TotalConnections++
		p.stats.ActiveConnections++
		p.mu.Unlock()

		return p.factory.NewConnection(ctx)
	}
}

// Release returns a connection to the pool
func (p *ConnectionPool) Release(conn *PooledConnection) {
	if conn == nil {
		return
	}

	conn.InUse = false
	conn.LastUsed = time.Now()

	p.mu.Lock()
	p.stats.ActiveConnections--
	p.stats.IdleConnections++
	p.mu.Unlock()

	select {
	case p.connections <- conn:
		log.Debugf("Returned connection to pool (idle: %d)", len(p.connections))
	default:
		// Pool full, close connection
		p.factory.CloseConnection(conn)

		p.mu.Lock()
		p.stats.TotalConnections--
		p.stats.IdleConnections--
		p.mu.Unlock()
	}
}

// Stats returns current pool statistics
func (p *ConnectionPool) Stats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stats
}

// Close closes all connections in the pool
func (p *ConnectionPool) Close() error {
	close(p.connections)

	for conn := range p.connections {
		p.factory.CloseConnection(conn)
	}

	return nil
}

// PerformanceMonitor tracks application performance metrics
type PerformanceMonitor struct {
	mu            sync.RWMutex
	metrics       map[string]*Metric
	slowThreshold time.Duration
}

// Metric represents a performance metric
type Metric struct {
	Name        string
	Count       int64
	TotalTime   time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
	LastTime    time.Duration
	AvgTime     time.Duration
	SlowQueries int64
	Errors      int64
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(slowThreshold time.Duration) *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics:       make(map[string]*Metric),
		slowThreshold: slowThreshold,
	}
}

// Record records a performance measurement
func (pm *PerformanceMonitor) Record(name string, duration time.Duration, err error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	metric, exists := pm.metrics[name]
	if !exists {
		metric = &Metric{
			Name:    name,
			MinTime: duration,
			MaxTime: duration,
		}
		pm.metrics[name] = metric
	}

	metric.Count++
	metric.TotalTime += duration
	metric.LastTime = duration

	if duration < metric.MinTime {
		metric.MinTime = duration
	}
	if duration > metric.MaxTime {
		metric.MaxTime = duration
	}

	metric.AvgTime = metric.TotalTime / time.Duration(metric.Count)

	if duration > pm.slowThreshold {
		metric.SlowQueries++
		log.Warnf("Slow %s: %v (threshold: %v)", name, duration, pm.slowThreshold)
	}

	if err != nil {
		metric.Errors++
	}
}

// GetMetric returns a specific metric
func (pm *PerformanceMonitor) GetMetric(name string) (*Metric, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	metric, exists := pm.metrics[name]
	return metric, exists
}

// GetAllMetrics returns all metrics
func (pm *PerformanceMonitor) GetAllMetrics() map[string]*Metric {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*Metric)
	for k, v := range pm.metrics {
		result[k] = v
	}

	return result
}

// QueryOptimizer provides query optimization suggestions
type QueryOptimizer struct {
	monitor *PerformanceMonitor
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer(monitor *PerformanceMonitor) *QueryOptimizer {
	return &QueryOptimizer{
		monitor: monitor,
	}
}

// AnalyzeQuery analyzes a query and provides optimization suggestions
func (qo *QueryOptimizer) AnalyzeQuery(resourceType string, params map[string][]string) []string {
	suggestions := make([]string, 0)

	// Check for common performance issues
	if len(params) == 0 {
		suggestions = append(suggestions, "Add search parameters to limit results")
	}

	// Check for expensive operations
	if _, hasInclude := params["_include"]; hasInclude {
		suggestions = append(suggestions, "Consider limiting _include depth for better performance")
	}

	if _, hasRevInclude := params["_revinclude"]; hasRevInclude {
		suggestions = append(suggestions, "_revinclude can be expensive on large datasets")
	}

	// Check pagination
	if count, hasCount := params["_count"]; hasCount {
		if len(count) > 0 {
			// Large page sizes can impact performance
			suggestions = append(suggestions, "Consider using smaller _count values for better response times")
		}
	}

	return suggestions
}

// IndexAdvisor provides index recommendations
type IndexAdvisor struct {
	queryPatterns map[string]int
	mu            sync.RWMutex
}

// NewIndexAdvisor creates a new index advisor
func NewIndexAdvisor() *IndexAdvisor {
	return &IndexAdvisor{
		queryPatterns: make(map[string]int),
	}
}

// RecordQueryPattern records a query pattern for analysis
func (ia *IndexAdvisor) RecordQueryPattern(resourceType, param string) {
	ia.mu.Lock()
	defer ia.mu.Unlock()

	key := resourceType + ":" + param
	ia.queryPatterns[key]++
}

// GetRecommendations returns index recommendations
func (ia *IndexAdvisor) GetRecommendations() map[string][]string {
	ia.mu.RLock()
	defer ia.mu.RUnlock()

	recommendations := make(map[string][]string)

	// Find frequently used parameters that might benefit from indexing
	for pattern, count := range ia.queryPatterns {
		if count > 100 { // Threshold for recommendation
			parts := strings.Split(pattern, ":")
			if len(parts) == 2 {
				resourceType := parts[0]
				param := parts[1]

				if _, exists := recommendations[resourceType]; !exists {
					recommendations[resourceType] = make([]string, 0)
				}

				recommendations[resourceType] = append(recommendations[resourceType], param)
			}
		}
	}

	return recommendations
}
