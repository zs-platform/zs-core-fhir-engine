package security

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// SecurityConfig contains security configuration options
type SecurityConfig struct {
	TLSEnabled                   bool
	TLSCertFile                  string
	TLSKeyFile                   string
	MinTLSVersion                uint16
	AllowedCipherSuites          []uint16
	HSTSEnabled                  bool
	HSTSMaxAge                   int
	HSTSIncludeSubdomains        bool
	HSTSPreload                  bool
	CSPPolicy                    string
	FrameOptions                 string
	ContentTypeOptions           bool
	XSSProtection                bool
	ReferrerPolicy               string
	PermittedCrossDomainPolicies string
	RateLimitEnabled             bool
	RateLimitRequests            int
	RateLimitWindow              time.Duration
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		TLSEnabled:                   false, // Disabled by default for development
		MinTLSVersion:                tls.VersionTLS12,
		AllowedCipherSuites:          secureCipherSuites(),
		HSTSEnabled:                  true,
		HSTSMaxAge:                   31536000, // 1 year
		HSTSIncludeSubdomains:        true,
		HSTSPreload:                  true,
		CSPPolicy:                    "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self';",
		FrameOptions:                 "DENY",
		ContentTypeOptions:           true,
		XSSProtection:                true,
		ReferrerPolicy:               "strict-origin-when-cross-origin",
		PermittedCrossDomainPolicies: "none",
		RateLimitEnabled:             true,
		RateLimitRequests:            100,
		RateLimitWindow:              time.Minute,
	}
}

// secureCipherSuites returns recommended TLS cipher suites
func secureCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(config *SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Strict-Transport-Security (HSTS)
			if config.HSTSEnabled {
				hstsValue := fmt.Sprintf("max-age=%d", config.HSTSMaxAge)
				if config.HSTSIncludeSubdomains {
					hstsValue += "; includeSubDomains"
				}
				if config.HSTSPreload {
					hstsValue += "; preload"
				}
				w.Header().Set("Strict-Transport-Security", hstsValue)
			}

			// Content-Security-Policy
			if config.CSPPolicy != "" {
				w.Header().Set("Content-Security-Policy", config.CSPPolicy)
			}

			// X-Frame-Options
			if config.FrameOptions != "" {
				w.Header().Set("X-Frame-Options", config.FrameOptions)
			}

			// X-Content-Type-Options
			if config.ContentTypeOptions {
				w.Header().Set("X-Content-Type-Options", "nosniff")
			}

			// X-XSS-Protection
			if config.XSSProtection {
				w.Header().Set("X-XSS-Protection", "1; mode=block")
			}

			// Referrer-Policy
			if config.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", config.ReferrerPolicy)
			}

			// X-Permitted-Cross-Domain-Policies
			if config.PermittedCrossDomainPolicies != "" {
				w.Header().Set("X-Permitted-Cross-Domain-Policies", config.PermittedCrossDomainPolicies)
			}

			// Remove potentially dangerous headers
			w.Header().Del("Server")
			w.Header().Del("X-Powered-By")

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter implements request rate limiting
type RateLimiter struct {
	requests map[string]*clientRequests
	mu       sync.RWMutex
	config   *RateLimitConfig
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
}

// clientRequests tracks requests per client
type clientRequests struct {
	count       int
	windowStart time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*clientRequests),
		config:   config,
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.requests[clientID]

	if !exists || now.Sub(client.windowStart) > rl.config.Window {
		// New window
		rl.requests[clientID] = &clientRequests{
			count:       1,
			windowStart: now,
		}
		return true
	}

	if client.count >= rl.config.MaxRequests {
		log.Warnf("Rate limit exceeded for client: %s", clientID)
		return false
	}

	client.count++
	return true
}

// GetRemaining returns remaining requests for a client
func (rl *RateLimiter) GetRemaining(clientID string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	client, exists := rl.requests[clientID]
	if !exists {
		return rl.config.MaxRequests
	}

	if time.Since(client.windowStart) > rl.config.Window {
		return rl.config.MaxRequests
	}

	remaining := rl.config.MaxRequests - client.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset resets rate limiting for a client
func (rl *RateLimiter) Reset(clientID string) {
	rl.mu.Lock()
	delete(rl.requests, clientID)
	rl.mu.Unlock()
}

// RateLimitMiddleware adds rate limiting to HTTP handler
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := getClientID(r)

			if !limiter.Allow(clientID) {
				w.Header().Set("Content-Type", "application/fhir+json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(limiter.config.Window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)

				outcome := map[string]interface{}{
					"resourceType": "OperationOutcome",
					"issue": []map[string]interface{}{
						{
							"severity":    "error",
							"code":        "throttled",
							"diagnostics": "Rate limit exceeded. Please retry after the specified time.",
						},
					},
				}

				json.NewEncoder(w).Encode(outcome)
				return
			}

			// Add rate limit headers
			remaining := limiter.GetRemaining(clientID)
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

			next.ServeHTTP(w, r)
		})
	}
}

// getClientID extracts a client identifier from the request
func getClientID(r *http.Request) string {
	// Try to get from authenticated user
	if userID := r.Context().Value("user_id"); userID != nil {
		return fmt.Sprintf("%v", userID)
	}

	// Fall back to IP address
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a secure default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"}, // Restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Type", "X-RateLimit-Remaining", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(config *CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if isOriginAllowed(origin, config.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if an origin is in the allowed list
func isOriginAllowed(origin string, allowed []string) bool {
	for _, allowed := range allowed {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

// AuditLogger provides security audit logging
type AuditLogger struct {
	events chan AuditEvent
}

// AuditEvent represents a security audit event
type AuditEvent struct {
	Timestamp time.Time
	EventType string
	UserID    string
	ClientIP  string
	Resource  string
	Action    string
	Outcome   string
	Details   map[string]interface{}
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		events: make(chan AuditEvent, 1000),
	}
}

// Log records an audit event
func (al *AuditLogger) Log(event AuditEvent) {
	select {
	case al.events <- event:
		// Event queued for logging
	default:
		// Channel full, log directly
		logSecurityEvent(event)
	}
}

// logSecurityEvent logs a security event
func logSecurityEvent(event AuditEvent) {
	log.Infof("[SECURITY] %s | User: %s | IP: %s | Resource: %s | Action: %s | Outcome: %s | Details: %v",
		event.Timestamp.Format(time.RFC3339),
		event.UserID,
		event.ClientIP,
		event.Resource,
		event.Action,
		event.Outcome,
		event.Details,
	)
}

// AuditMiddleware adds audit logging to HTTP handlers
func AuditMiddleware(auditLogger *AuditLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create response wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(rw, r)

			// Log the request
			event := AuditEvent{
				Timestamp: start,
				EventType: "http_request",
				UserID:    fmt.Sprintf("%v", r.Context().Value("user_id")),
				ClientIP:  getClientID(r),
				Resource:  r.URL.Path,
				Action:    r.Method,
				Outcome:   fmt.Sprintf("%d", rw.statusCode),
				Details: map[string]interface{}{
					"duration": time.Since(start).Milliseconds(),
					"status":   rw.statusCode,
					"query":    r.URL.RawQuery,
				},
			}

			auditLogger.Log(event)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}
