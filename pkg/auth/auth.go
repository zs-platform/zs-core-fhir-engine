package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in a SMART on FHIR JWT token
type JWTClaims struct {
	Issuer       string  `json:"iss"`
	Subject      string  `json:"sub"`
	Audience     string  `json:"aud"`
	Expiration   int64   `json:"exp"`
	IssuedAt     int64   `json:"iat"`
	JWTID        string  `json:"jti"`
	Scope        string  `json:"scope"`
	Patient      string  `json:"patient,omitempty"`
	Encounter    string  `json:"encounter,omitempty"`
	Practitioner string  `json:"practitioner,omitempty"`
	Organization string  `json:"organization,omitempty"`
	Context      Context `json:"context,omitempty"`
	jwt.RegisteredClaims
}

// Context represents the launch context in SMART on FHIR
type Context struct {
	Patient      string `json:"patient,omitempty"`
	Encounter    string `json:"encounter,omitempty"`
	Practitioner string `json:"practitioner,omitempty"`
	Organization string `json:"organization,omitempty"`
}

// AuthService interface for SMART on FHIR authentication
type AuthService interface {
	// GenerateAuthURL generates the authorization URL for SMART app launch
	GenerateAuthURL(ctx context.Context, req AuthRequest) (string, error)

	// ExchangeCodeForToken exchanges authorization code for access token
	ExchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error)

	// ValidateToken validates a JWT access token
	ValidateToken(ctx context.Context, token string) (*JWTClaims, error)

	// RefreshToken generates a new access token using refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// RevokeToken revokes an access token
	RevokeToken(ctx context.Context, token string) error

	// GetUserInfo retrieves user information from token
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
}

// AuthRequest represents an authorization request
type AuthRequest struct {
	ClientID      string            `json:"client_id"`
	RedirectURI   string            `json:"redirect_uri"`
	Scope         string            `json:"scope"`
	State         string            `json:"state"`
	ResponseType  string            `json:"response_type"`
	LaunchContext map[string]string `json:"launch_context,omitempty"`
	Audience      string            `json:"aud,omitempty"`
}

// TokenResponse represents the token response from OAuth2 token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// UserInfo represents user information from OpenID Connect
type UserInfo struct {
	Sub                 string            `json:"sub"`
	Name                string            `json:"name,omitempty"`
	GivenName           string            `json:"given_name,omitempty"`
	FamilyName          string            `json:"family_name,omitempty"`
	MiddleName          string            `json:"middle_name,omitempty"`
	Nickname            string            `json:"nickname,omitempty"`
	PreferredUsername   string            `json:"preferred_username,omitempty"`
	Profile             string            `json:"profile,omitempty"`
	Picture             string            `json:"picture,omitempty"`
	Website             string            `json:"website,omitempty"`
	Email               string            `json:"email,omitempty"`
	EmailVerified       bool              `json:"email_verified,omitempty"`
	Gender              string            `json:"gender,omitempty"`
	Birthdate           string            `json:"birthdate,omitempty"`
	Zoneinfo            string            `json:"zoneinfo,omitempty"`
	Locale              string            `json:"locale,omitempty"`
	PhoneNumber         string            `json:"phone_number,omitempty"`
	PhoneNumberVerified bool              `json:"phone_number_verified,omitempty"`
	Address             map[string]string `json:"address,omitempty"`
	UpdatedAt           time.Time         `json:"updated_at,omitempty"`
}

// KeycloakConfig represents Keycloak configuration
type KeycloakConfig struct {
	BaseURL     string `json:"base_url"`
	Realm       string `json:"realm"`
	ClientID    string `json:"client_id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

// TokenStore interface for storing and managing tokens
type TokenStore interface {
	// StoreToken stores a token
	StoreToken(ctx context.Context, tokenID string, token *TokenResponse) error

	// GetToken retrieves a token by ID
	GetToken(ctx context.Context, tokenID string) (*TokenResponse, error)

	// RevokeToken removes a token
	RevokeToken(ctx context.Context, tokenID string) error

	// CleanupExpiredTokens removes expired tokens
	CleanupExpiredTokens(ctx context.Context) error
}

// JWTManager manages JWT token generation and validation
type JWTManager struct {
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	issuer      string
	keycloakURL string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, issuer string, keycloakURL string) *JWTManager {
	return &JWTManager{
		privateKey:  privateKey,
		publicKey:   publicKey,
		issuer:      issuer,
		keycloakURL: keycloakURL,
	}
}

// GenerateToken generates a new JWT token
func (j *JWTManager) GenerateToken(claims *JWTClaims) (string, error) {
	now := time.Now()

	// Set default values
	if claims.Expiration == 0 {
		claims.Expiration = now.Add(time.Hour).Unix()
	}
	if claims.IssuedAt == 0 {
		claims.IssuedAt = now.Unix()
	}
	if claims.Issuer == "" {
		claims.Issuer = j.issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ScopeValidator validates SMART on FHIR scopes
type ScopeValidator struct {
	allowedScopes map[string]bool
}

// NewScopeValidator creates a new scope validator
func NewScopeValidator() *ScopeValidator {
	return &ScopeValidator{
		allowedScopes: map[string]bool{
			// Patient scopes
			"patient/Patient.read":            true,
			"patient/Patient.write":           true,
			"patient/Observation.read":        true,
			"patient/Observation.write":       true,
			"patient/Condition.read":          true,
			"patient/Condition.write":         true,
			"patient/MedicationRequest.read":  true,
			"patient/MedicationRequest.write": true,

			// System scopes
			"system/Patient.read":            true,
			"system/Patient.write":           true,
			"system/Observation.read":        true,
			"system/Observation.write":       true,
			"system/Condition.read":          true,
			"system/Condition.write":         true,
			"system/MedicationRequest.read":  true,
			"system/MedicationRequest.write": true,

			// Launch scopes
			"launch":              true,
			"launch/patient":      true,
			"launch/encounter":    true,
			"launch/practitioner": true,
			"launch/organization": true,

			// OpenID Connect scopes
			"openid":  true,
			"profile": true,
			"email":   true,
			"phone":   true,
			"address": true,

			// Offline access
			"offline_access": true,
		},
	}
}

// ValidateScope validates if the requested scopes are allowed
func (s *ScopeValidator) ValidateScope(scopes string) error {
	// Parse scopes (space-separated)
	scopeList := parseScopes(scopes)

	for _, scope := range scopeList {
		if !s.allowedScopes[scope] {
			return fmt.Errorf("invalid scope: %s", scope)
		}
	}

	return nil
}

// parseScopes parses space-separated scopes into a slice
func parseScopes(scopes string) []string {
	var result []string
	start := 0

	for i, char := range scopes {
		if char == ' ' {
			if start < i {
				result = append(result, scopes[start:i])
			}
			start = i + 1
		}
	}

	if start < len(scopes) {
		result = append(result, scopes[start:])
	}

	return result
}

// BangladeshSpecificClaims represents Bangladesh-specific claims
type BangladeshSpecificClaims struct {
	NID          string `json:"nid,omitempty"`           // National ID
	BRN          string `json:"brn,omitempty"`           // Birth Registration Number
	UHID         string `json:"uhid,omitempty"`          // Unique Health ID
	FCN          string `json:"fcn,omitempty"`           // Family Counting Number (Rohingya)
	ProgressID   string `json:"progress_id,omitempty"`   // Progress ID (Rohingya)
	CampLocation string `json:"camp_location,omitempty"` // Camp Location (Rohingya)
	Division     string `json:"division,omitempty"`      // Administrative Division
	District     string `json:"district,omitempty"`      // Administrative District
	Upazila      string `json:"upazila,omitempty"`       // Administrative Upazila
	Union        string `json:"union,omitempty"`         // Administrative Union
	Ward         string `json:"ward,omitempty"`          // Administrative Ward
}

// AddBangladeshClaims adds Bangladesh-specific claims to JWT claims
func (claims *JWTClaims) AddBangladeshClaims(bdClaims BangladeshSpecificClaims) {
	// Add Bangladesh-specific claims as extensions
	if claims.Context.Organization == "" {
		claims.Context.Organization = bdClaims.Division
	}
}
