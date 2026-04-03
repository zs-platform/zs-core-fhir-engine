package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// KeycloakAuthService implements AuthService using Keycloak
type KeycloakAuthService struct {
	config     *KeycloakConfig
	jwtManager *JWTManager
	tokenStore TokenStore
	httpClient *http.Client
}

// NewKeycloakAuthService creates a new Keycloak auth service
func NewKeycloakAuthService(config *KeycloakConfig, jwtManager *JWTManager, tokenStore TokenStore) *KeycloakAuthService {
	return &KeycloakAuthService{
		config:     config,
		jwtManager: jwtManager,
		tokenStore: tokenStore,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateAuthURL generates the authorization URL for SMART app launch
func (k *KeycloakAuthService) GenerateAuthURL(ctx context.Context, req AuthRequest) (string, error) {
	// Validate required parameters
	if req.ClientID == "" {
		return "", fmt.Errorf("client_id is required")
	}
	if req.RedirectURI == "" {
		return "", fmt.Errorf("redirect_uri is required")
	}
	if req.Scope == "" {
		return "", fmt.Errorf("scope is required")
	}
	if req.ResponseType != "code" {
		return "", fmt.Errorf("response_type must be 'code'")
	}

	// Build authorization URL
	authURL, err := url.Parse(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", k.config.BaseURL, k.config.Realm))
	if err != nil {
		return "", fmt.Errorf("failed to parse auth URL: %w", err)
	}

	// Add query parameters
	params := url.Values{}
	params.Set("client_id", req.ClientID)
	params.Set("redirect_uri", req.RedirectURI)
	params.Set("scope", req.Scope)
	params.Set("state", req.State)
	params.Set("response_type", req.ResponseType)
	params.Set("response_mode", "query")

	if req.Audience != "" {
		params.Set("audience", req.Audience)
	}

	// Add launch context parameters
	for key, value := range req.LaunchContext {
		params.Set("launch_"+key, value)
	}

	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (k *KeycloakAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	// Prepare token request
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", k.config.BaseURL, k.config.Realm)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", k.config.RedirectURI)
	data.Set("client_id", k.config.ClientID)
	data.Set("client_secret", k.config.Secret)

	// Make HTTP request
	resp, err := k.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Store token
	tokenID := generateTokenID()
	if err := k.tokenStore.StoreToken(ctx, tokenID, &tokenResp); err != nil {
		log.Warnf("Failed to store token: %v", err)
	}

	log.Infof("Successfully exchanged authorization code for access token")
	return &tokenResp, nil
}

// ValidateToken validates a JWT access token
func (k *KeycloakAuthService) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate JWT
	claims, err := k.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check expiration
	if time.Now().Unix() > claims.Expiration {
		return nil, fmt.Errorf("token has expired")
	}

	// Check issuer
	if claims.Issuer != fmt.Sprintf("%s/realms/%s", k.config.BaseURL, k.config.Realm) {
		return nil, fmt.Errorf("invalid token issuer")
	}

	return claims, nil
}

// RefreshToken generates a new access token using refresh token
func (k *KeycloakAuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Prepare refresh request
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", k.config.BaseURL, k.config.Realm)

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", k.config.ClientID)
	data.Set("client_secret", k.config.Secret)

	// Make HTTP request
	resp, err := k.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	log.Infof("Successfully refreshed access token")
	return &tokenResp, nil
}

// RevokeToken revokes an access token
func (k *KeycloakAuthService) RevokeToken(ctx context.Context, token string) error {
	// Revoke from Keycloak
	revokeURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/revoke", k.config.BaseURL, k.config.Realm)

	data := url.Values{}
	data.Set("client_id", k.config.ClientID)
	data.Set("client_secret", k.config.Secret)
	data.Set("token", token)
	data.Set("token_type_hint", "access_token")

	resp, err := k.httpClient.PostForm(revokeURL, data)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token revocation failed with status: %d", resp.StatusCode)
	}

	// Remove from local store
	tokenID := generateTokenID()
	if err := k.tokenStore.RevokeToken(ctx, tokenID); err != nil {
		log.Warnf("Failed to revoke token from store: %v", err)
	}

	log.Infof("Successfully revoked access token")
	return nil
}

// GetUserInfo retrieves user information from token
func (k *KeycloakAuthService) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	// Prepare user info request
	userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", k.config.BaseURL, k.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	// Make HTTP request
	resp, err := k.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}

	return &userInfo, nil
}

// generateTokenID generates a unique token ID
func generateTokenID() string {
	return fmt.Sprintf("tok_%d", time.Now().UnixNano())
}

// BangladeshTokenEnhancer enhances tokens with Bangladesh-specific claims
type BangladeshTokenEnhancer struct {
	profileService BangladeshProfileService
}

// BangladeshProfileService interface for Bangladesh profile data
type BangladeshProfileService interface {
	// GetPatientProfile retrieves Bangladesh patient profile
	GetPatientProfile(ctx context.Context, patientID string) (*BangladeshSpecificClaims, error)

	// GetPractitionerProfile retrieves Bangladesh practitioner profile
	GetPractitionerProfile(ctx context.Context, practitionerID string) (*BangladeshSpecificClaims, error)

	// GetOrganizationProfile retrieves Bangladesh organization profile
	GetOrganizationProfile(ctx context.Context, organizationID string) (*BangladeshSpecificClaims, error)
}

// NewBangladeshTokenEnhancer creates a new token enhancer
func NewBangladeshTokenEnhancer(profileService BangladeshProfileService) *BangladeshTokenEnhancer {
	return &BangladeshTokenEnhancer{
		profileService: profileService,
	}
}

// EnhanceToken adds Bangladesh-specific claims to token
func (b *BangladeshTokenEnhancer) EnhanceToken(ctx context.Context, claims *JWTClaims) error {
	// Add Bangladesh-specific claims based on context
	if claims.Context.Patient != "" {
		if profile, err := b.profileService.GetPatientProfile(ctx, claims.Context.Patient); err == nil {
			claims.AddBangladeshClaims(*profile)
		}
	}

	if claims.Context.Practitioner != "" {
		if profile, err := b.profileService.GetPractitionerProfile(ctx, claims.Context.Practitioner); err == nil {
			claims.AddBangladeshClaims(*profile)
		}
	}

	if claims.Context.Organization != "" {
		if profile, err := b.profileService.GetOrganizationProfile(ctx, claims.Context.Organization); err == nil {
			claims.AddBangladeshClaims(*profile)
		}
	}

	return nil
}

// SMARTLaunchContext represents SMART on FHIR launch context
type SMARTLaunchContext struct {
	Patient           string `json:"patient,omitempty"`
	Encounter         string `json:"encounter,omitempty"`
	Practitioner      string `json:"practitioner,omitempty"`
	Organization      string `json:"organization,omitempty"`
	NeedPatientBanner bool   `json:"need_patient_banner,omitempty"`
	SmartStyleURL     string `json:"smart_style_url,omitempty"`
}

// LaunchManager handles SMART on FHIR app launches
type LaunchManager struct {
	authService AuthService
	enhancer    *BangladeshTokenEnhancer
}

// NewLaunchManager creates a new launch manager
func NewLaunchManager(authService AuthService, enhancer *BangladeshTokenEnhancer) *LaunchManager {
	return &LaunchManager{
		authService: authService,
		enhancer:    enhancer,
	}
}

// HandleLaunch handles SMART app launch
func (l *LaunchManager) HandleLaunch(ctx context.Context, launchContext *SMARTLaunchContext) (*AuthRequest, error) {
	// Generate authorization request with launch context
	authReq := AuthRequest{
		ClientID:      "zs-fhir-engine",
		RedirectURI:   "http://localhost:3000/callback",
		Scope:         "launch/patient patient/*.read openid profile",
		State:         generateState(),
		ResponseType:  "code",
		LaunchContext: map[string]string{},
	}

	// Add launch context parameters
	if launchContext.Patient != "" {
		authReq.LaunchContext["patient"] = launchContext.Patient
		authReq.Scope = "launch/patient patient/*.read openid profile"
	}

	if launchContext.Encounter != "" {
		authReq.LaunchContext["encounter"] = launchContext.Encounter
		authReq.Scope += " encounter/*.read"
	}

	if launchContext.Practitioner != "" {
		authReq.LaunchContext["practitioner"] = launchContext.Practitioner
		authReq.Scope += " practitioner/*.read"
	}

	if launchContext.Organization != "" {
		authReq.LaunchContext["organization"] = launchContext.Organization
		authReq.Scope += " organization/*.read"
	}

	return &authReq, nil
}

// generateState generates a random state parameter
func generateState() string {
	return fmt.Sprintf("state_%d", time.Now().UnixNano())
}

// EHRLaunch represents EHR system launch parameters
type EHRLaunch struct {
	Iss           string `json:"iss"`            // EHR issuer URL
	LaunchContext string `json:"launch_context"` // Launch context JWT
	Audience      string `json:"aud,omitempty"`  // Audience (optional)
}

// StandaloneLaunch represents standalone app launch
type StandaloneLaunch struct {
	ClientID    string `json:"client_id"`     // Client ID
	RedirectURI string `json:"redirect_uri"`  // Redirect URI
	Scope       string `json:"scope"`         // Requested scopes
	State       string `json:"state"`         // State parameter
	Audience    string `json:"aud,omitempty"` // Audience (optional)
}
