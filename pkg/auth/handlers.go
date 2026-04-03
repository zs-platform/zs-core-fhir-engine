package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

// HTTPHandler handles SMART on FHIR HTTP endpoints
type HTTPHandler struct {
	authService AuthService
	auditLogger AuditLogger
	config      *AuthConfig
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	BaseURL     string `json:"base_url"`
	Realm       string `json:"realm"`
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	JWKSURL     string `json:"jwks_url"`
	TokenTTL    int    `json:"token_ttl"`
	RefreshTTL  int    `json:"refresh_ttl"`
}

// NewHTTPHandler creates a new auth HTTP handler
func NewHTTPHandler(authService AuthService, auditLogger AuditLogger, config *AuthConfig) *HTTPHandler {
	return &HTTPHandler{
		authService: authService,
		auditLogger: auditLogger,
		config:      config,
	}
}

// RegisterRoutes registers authentication and authorization routes
func (h *HTTPHandler) RegisterRoutes(router chi.Router) {
	// OAuth2 endpoints
	router.Get("/auth/authorize", h.handleAuthorize)
	router.Post("/auth/token", h.handleToken)
	router.Post("/auth/revoke", h.handleRevoke)
	router.Get("/auth/userinfo", h.handleUserInfo)

	// OpenID Connect endpoints
	router.Get("/auth/.well-known/openid_configuration", h.handleOpenIDConfiguration)
	router.Get("/auth/.well-known/jwks.json", h.handleJWKS)

	// SMART on FHIR specific endpoints
	router.Get("/auth/launch", h.handleLaunch)
	router.Get("/auth/metadata", h.handleMetadata)
}

// handleAuthorize handles the OAuth2 authorization endpoint
func (h *HTTPHandler) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	query := r.URL.Query()
	clientID := query.Get("client_id")
	redirectURI := query.Get("redirect_uri")
	scope := query.Get("scope")
	state := query.Get("state")
	responseType := query.Get("response_type")
	launch := query.Get("launch")

	// Log authorization request
	auditEvent := AuditLog{
		ID:        generateAuditID(),
		Timestamp: time.Now(),
		Action:    "authorize_request",
		ClientID:  clientID,
		Scope:     scope,
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
		Success:   true,
	}

	if err := h.auditLogger.LogAuthEvent(ctx, auditEvent); err != nil {
		log.Warnf("Failed to log auth event: %v", err)
	}

	// Validate parameters
	if responseType != "code" {
		h.writeError(w, "unsupported_response_type", "response_type must be 'code'", http.StatusBadRequest)
		return
	}

	// Generate authorization request
	authReq := AuthRequest{
		ClientID:      clientID,
		RedirectURI:   redirectURI,
		Scope:         scope,
		State:         state,
		ResponseType:  responseType,
		LaunchContext: make(map[string]string),
	}

	if launch != "" {
		authReq.LaunchContext["launch"] = launch
	}

	// Generate authorization URL
	authURL, err := h.authService.GenerateAuthURL(ctx, authReq)
	if err != nil {
		h.writeError(w, "authorization_failed", err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to Keycloak
	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleToken handles the OAuth2 token endpoint
func (h *HTTPHandler) handleToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.writeError(w, "invalid_request", "Failed to parse form data", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")

	switch grantType {
	case "authorization_code":
		h.handleAuthorizationCodeGrant(ctx, w, r)
	case "refresh_token":
		h.handleRefreshTokenGrant(ctx, w, r)
	case "client_credentials":
		h.writeError(w, "unsupported_grant_type", "client_credentials not supported", http.StatusBadRequest)
	default:
		h.writeError(w, "invalid_grant_type", "Invalid grant type", http.StatusBadRequest)
	}
}

// handleAuthorizationCodeGrant handles authorization code grant
func (h *HTTPHandler) handleAuthorizationCodeGrant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	clientID := r.FormValue("client_id")

	// Exchange code for token
	tokenResp, err := h.authService.ExchangeCodeForToken(ctx, code)
	if err != nil {
		h.writeError(w, "invalid_grant", err.Error(), http.StatusBadRequest)
		return
	}

	// Log successful token exchange
	auditEvent := AuditLog{
		ID:        generateAuditID(),
		Timestamp: time.Now(),
		Action:    "token_exchange",
		ClientID:  clientID,
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
		Success:   true,
	}

	if err := h.auditLogger.LogAuthEvent(ctx, auditEvent); err != nil {
		log.Warnf("Failed to log auth event: %v", err)
	}

	// Return token response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	json.NewEncoder(w).Encode(tokenResp)
}

// handleRefreshTokenGrant handles refresh token grant
func (h *HTTPHandler) handleRefreshTokenGrant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refresh_token")

	if refreshToken == "" {
		h.writeError(w, "invalid_request", "refresh_token is required", http.StatusBadRequest)
		return
	}

	// Refresh token
	tokenResp, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		h.writeError(w, "invalid_grant", err.Error(), http.StatusBadRequest)
		return
	}

	// Log successful token refresh
	auditEvent := AuditLog{
		ID:        generateAuditID(),
		Timestamp: time.Now(),
		Action:    "token_refresh",
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
		Success:   true,
	}

	if err := h.auditLogger.LogAuthEvent(ctx, auditEvent); err != nil {
		log.Warnf("Failed to log auth event: %v", err)
	}

	// Return token response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	json.NewEncoder(w).Encode(tokenResp)
}

// handleRevoke handles the OAuth2 revoke endpoint
func (h *HTTPHandler) handleRevoke(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.writeError(w, "invalid_request", "Failed to parse form data", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")
	tokenTypeHint := r.FormValue("token_type_hint")

	if token == "" {
		h.writeError(w, "invalid_request", "token is required", http.StatusBadRequest)
		return
	}

	// Revoke token
	err := h.authService.RevokeToken(ctx, token)
	if err != nil {
		h.writeError(w, "invalid_request", err.Error(), http.StatusBadRequest)
		return
	}

	// Log successful token revocation
	auditEvent := AuditLog{
		ID:        generateAuditID(),
		Timestamp: time.Now(),
		Action:    "token_revoke",
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
		Success:   true,
	}

	if tokenTypeHint != "" {
		auditEvent.Error = "token_type_hint: " + tokenTypeHint
	}

	if err := h.auditLogger.LogAuthEvent(ctx, auditEvent); err != nil {
		log.Warnf("Failed to log auth event: %v", err)
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}

// handleUserInfo handles the OpenID Connect userinfo endpoint
func (h *HTTPHandler) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.writeError(w, "invalid_token", "Authorization header is required", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		h.writeError(w, "invalid_token", "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Get user info
	userInfo, err := h.authService.GetUserInfo(ctx, token)
	if err != nil {
		h.writeError(w, "invalid_token", err.Error(), http.StatusUnauthorized)
		return
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// handleOpenIDConfiguration handles OpenID Connect discovery
func (h *HTTPHandler) handleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"issuer":                   fmt.Sprintf("%s/realms/%s", h.config.BaseURL, h.config.Realm),
		"authorization_endpoint":   fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", h.config.BaseURL, h.config.Realm),
		"token_endpoint":           fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", h.config.BaseURL, h.config.Realm),
		"userinfo_endpoint":        fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", h.config.BaseURL, h.config.Realm),
		"jwks_uri":                 fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", h.config.BaseURL, h.config.Realm),
		"end_session_endpoint":     fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", h.config.BaseURL, h.config.Realm),
		"response_types_supported": []string{"code"},
		"grant_types_supported":    []string{"authorization_code", "refresh_token"},
		"scopes_supported": []string{
			"openid", "profile", "email",
			"patient/Patient.read", "patient/Patient.write",
			"patient/Observation.read", "patient/Observation.write",
			"system/Patient.read", "system/Patient.write",
			"system/Observation.read", "system/Observation.write",
		},
		"token_endpoint_auth_methods_supported": []string{"client_secret_post"},
		"code_challenge_methods_supported":      []string{"S256"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// handleJWKS handles JSON Web Key Set endpoint
func (h *HTTPHandler) handleJWKS(w http.ResponseWriter, r *http.Request) {
	// This should return the public keys used to verify JWT tokens
	// For now, return a placeholder
	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"kid": "1",
				"n":   "placeholder_modulus",
				"e":   "AQAB",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

// handleLaunch handles SMART on FHIR app launch
func (h *HTTPHandler) handleLaunch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse launch parameters
	query := r.URL.Query()
	iss := query.Get("iss")
	launch := query.Get("launch")

	if iss == "" {
		h.writeError(w, "invalid_launch", "iss parameter is required", http.StatusBadRequest)
		return
	}

	// Create launch context
	launchContext := &SMARTLaunchContext{
		Patient: launch, // Simplified - in real implementation, this would be parsed from JWT
	}

	// Generate authorization request for launch
	launchManager := NewLaunchManager(h.authService, nil)
	authReq, err := launchManager.HandleLaunch(ctx, launchContext)
	if err != nil {
		h.writeError(w, "launch_failed", err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate authorization URL
	authURL, err := h.authService.GenerateAuthURL(ctx, *authReq)
	if err != nil {
		h.writeError(w, "launch_failed", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return launch response
	response := map[string]interface{}{
		"authorization_url": authURL,
		"launch_context":    launchContext,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMetadata handles SMART on FHIR capability statement
func (h *HTTPHandler) handleMetadata(w http.ResponseWriter, r *http.Request) {
	// Return SMART on FHIR capability statement
	capabilityStatement := map[string]interface{}{
		"resourceType": "CapabilityStatement",
		"status":       "active",
		"date":         time.Now().Format(time.RFC3339),
		"fhirVersion":  "5.0.0",
		"format":       []string{"application/fhir+json", "application/fhir+xml"},
		"smart_on_fhir": map[string]interface{}{
			"capabilities": []string{
				"launch-standalone",
				"launch-ehr",
				"client-public",
				"client-confidential-symmetric",
				"sso-openid-connect",
				"context-ehr-patient",
				"context-ehr-encounter",
				"context-ehr-practitioner",
				"permission-patient",
				"permission-v1",
			},
			"authorization_endpoint": fmt.Sprintf("%s/auth/authorize", h.config.BaseURL),
			"token_endpoint":         fmt.Sprintf("%s/auth/token", h.config.BaseURL),
			"userinfo_endpoint":      fmt.Sprintf("%s/auth/userinfo", h.config.BaseURL),
			"jwks_uri":               fmt.Sprintf("%s/auth/.well-known/jwks.json", h.config.BaseURL),
			"registration_endpoint":  fmt.Sprintf("%s/auth/register", h.config.BaseURL),
		},
		"rest": map[string]interface{}{
			"mode": "server",
			"security": []map[string]interface{}{
				{
					"cors": true,
					"service": []interface{}{
						map[string]interface{}{
							"coding": []interface{}{
								map[string]interface{}{
									"system":  "http://terminology.hl7.org/CodeSystem/v3-ObservationValue",
									"code":    "SECUR",
									"display": "Security",
								},
							},
							"text": "OAuth2 using SMART on FHIR",
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(capabilityStatement)
}

// writeError writes an OAuth2 error response
func (h *HTTPHandler) writeError(w http.ResponseWriter, errorCode, description string, statusCode int) {
	errorResp := map[string]interface{}{
		"error":             errorCode,
		"error_description": description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}

// generateAuditID generates a unique audit ID
func generateAuditID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// Middleware for authentication
func (h *HTTPHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Skip authentication for certain endpoints
		if h.isPublicEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.writeError(w, "missing_token", "Authorization header is required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			h.writeError(w, "invalid_token", "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := h.authService.ValidateToken(ctx, token)
		if err != nil {
			h.writeError(w, "invalid_token", err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx = context.WithValue(ctx, "claims", claims)

		// Continue with next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isPublicEndpoint checks if an endpoint should skip authentication
func (h *HTTPHandler) isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/fhir/R5/metadata",
		"/auth/authorize",
		"/auth/token",
		"/auth/.well-known/openid_configuration",
		"/auth/.well-known/jwks.json",
		"/auth/launch",
		"/auth/metadata",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath {
			return true
		}
	}

	return false
}
