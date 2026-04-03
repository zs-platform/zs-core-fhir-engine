package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// PostgreSQLTokenStore implements TokenStore using PostgreSQL
type PostgreSQLTokenStore struct {
	db *sql.DB
}

// NewPostgreSQLTokenStore creates a new PostgreSQL token store
func NewPostgreSQLTokenStore(db *sql.DB) *PostgreSQLTokenStore {
	return &PostgreSQLTokenStore{
		db: db,
	}
}

// StoreToken stores a token in the database
func (p *PostgreSQLTokenStore) StoreToken(ctx context.Context, tokenID string, token *TokenResponse) error {
	query := `
		INSERT INTO auth.tokens (id, access_token, token_type, expires_in, scope, refresh_token, id_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		ON CONFLICT (id) DO UPDATE SET
			access_token = EXCLUDED.access_token,
			token_type = EXCLUDED.token_type,
			expires_in = EXCLUDED.expires_in,
			scope = EXCLUDED.scope,
			refresh_token = EXCLUDED.refresh_token,
			id_token = EXCLUDED.id_token,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := p.db.ExecContext(ctx, query,
		tokenID,
		token.AccessToken,
		token.TokenType,
		token.ExpiresIn,
		token.Scope,
		token.RefreshToken,
		token.IDToken,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	log.Debugf("Stored token %s in database", tokenID)
	return nil
}

// GetToken retrieves a token by ID
func (p *PostgreSQLTokenStore) GetToken(ctx context.Context, tokenID string) (*TokenResponse, error) {
	query := `
		SELECT access_token, token_type, expires_in, scope, refresh_token, id_token, created_at, updated_at
		FROM auth.tokens
		WHERE id = $1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`

	var token TokenResponse
	var createdAt, updatedAt time.Time
	var refreshToken, idToken sql.NullString

	err := p.db.QueryRowContext(ctx, query, tokenID).Scan(
		&token.AccessToken,
		&token.TokenType,
		&token.ExpiresIn,
		&token.Scope,
		&refreshToken,
		&idToken,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token not found or expired")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	if refreshToken.Valid {
		token.RefreshToken = refreshToken.String
	}
	if idToken.Valid {
		token.IDToken = idToken.String
	}

	log.Debugf("Retrieved token %s from database", tokenID)
	return &token, nil
}

// RevokeToken removes a token from the database
func (p *PostgreSQLTokenStore) RevokeToken(ctx context.Context, tokenID string) error {
	query := `
		UPDATE auth.tokens
		SET expires_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := p.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	log.Infof("Revoked token %s", tokenID)
	return nil
}

// CleanupExpiredTokens removes expired tokens
func (p *PostgreSQLTokenStore) CleanupExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM auth.tokens
		WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := p.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Infof("Cleaned up %d expired tokens", rowsAffected)
	return nil
}

// InMemoryTokenStore implements TokenStore using in-memory storage
type InMemoryTokenStore struct {
	tokens map[string]*TokenResponse
}

// NewInMemoryTokenStore creates a new in-memory token store
func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		tokens: make(map[string]*TokenResponse),
	}
}

// StoreToken stores a token in memory
func (i *InMemoryTokenStore) StoreToken(ctx context.Context, tokenID string, token *TokenResponse) error {
	i.tokens[tokenID] = token
	log.Debugf("Stored token %s in memory", tokenID)
	return nil
}

// GetToken retrieves a token by ID from memory
func (i *InMemoryTokenStore) GetToken(ctx context.Context, tokenID string) (*TokenResponse, error) {
	token, exists := i.tokens[tokenID]
	if !exists {
		return nil, fmt.Errorf("token not found")
	}

	// Check expiration (simple check based on expires_in)
	if token.ExpiresIn > 0 {
		// This is a simplified check - in production, store actual expiration time
		log.Debugf("Retrieved token %s from memory", tokenID)
		return token, nil
	}

	return nil, fmt.Errorf("token expired")
}

// RevokeToken removes a token from memory
func (i *InMemoryTokenStore) RevokeToken(ctx context.Context, tokenID string) error {
	delete(i.tokens, tokenID)
	log.Infof("Revoked token %s from memory", tokenID)
	return nil
}

// CleanupExpiredTokens removes expired tokens from memory
func (i *InMemoryTokenStore) CleanupExpiredTokens(ctx context.Context) error {
	// In a real implementation, this would check actual expiration times
	// For now, just log the operation
	log.Info("Cleaning up expired tokens from memory")
	return nil
}

// CreateTokensTable creates the tokens table if it doesn't exist
func CreateTokensTable(db *sql.DB) error {
	query := `
		CREATE SCHEMA IF NOT EXISTS auth;
		
		CREATE TABLE IF NOT EXISTS auth.tokens (
			id VARCHAR(255) PRIMARY KEY,
			access_token TEXT NOT NULL,
			token_type VARCHAR(50) NOT NULL,
			expires_in BIGINT,
			scope TEXT,
			refresh_token TEXT,
			id_token TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP WITH TIME ZONE
		);

		CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON auth.tokens(expires_at);
		CREATE INDEX IF NOT EXISTS idx_tokens_created_at ON auth.tokens(created_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tokens table: %w", err)
	}

	log.Info("Created auth.tokens table")
	return nil
}

// AuditLog represents an authentication audit log entry
type AuditLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"` // login, logout, token_refresh, token_revoke
	UserID    string    `json:"user_id"`
	ClientID  string    `json:"client_id"`
	Scope     string    `json:"scope"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// AuditLogger interface for authentication audit logging
type AuditLogger interface {
	// LogAuthEvent logs an authentication event
	LogAuthEvent(ctx context.Context, event AuditLog) error

	// GetAuthLogs retrieves authentication logs for a user
	GetAuthLogs(ctx context.Context, userID string, limit int) ([]AuditLog, error)
}

// PostgreSQLAuditLogger implements AuditLogger using PostgreSQL
type PostgreSQLAuditLogger struct {
	db *sql.DB
}

// NewPostgreSQLAuditLogger creates a new PostgreSQL audit logger
func NewPostgreSQLAuditLogger(db *sql.DB) *PostgreSQLAuditLogger {
	return &PostgreSQLAuditLogger{
		db: db,
	}
}

// LogAuthEvent logs an authentication event
func (p *PostgreSQLAuditLogger) LogAuthEvent(ctx context.Context, event AuditLog) error {
	query := `
		INSERT INTO auth.audit_logs (id, timestamp, action, user_id, client_id, scope, ip_address, user_agent, success, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := p.db.ExecContext(ctx, query,
		event.ID,
		event.Timestamp,
		event.Action,
		event.UserID,
		event.ClientID,
		event.Scope,
		event.IPAddress,
		event.UserAgent,
		event.Success,
		event.Error,
	)

	if err != nil {
		return fmt.Errorf("failed to log auth event: %w", err)
	}

	log.Infof("Logged auth event: %s for user %s", event.Action, event.UserID)
	return nil
}

// GetAuthLogs retrieves authentication logs for a user
func (p *PostgreSQLAuditLogger) GetAuthLogs(ctx context.Context, userID string, limit int) ([]AuditLog, error) {
	query := `
		SELECT id, timestamp, action, user_id, client_id, scope, ip_address, user_agent, success, error
		FROM auth.audit_logs
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := p.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		var errorMsg sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.Action,
			&log.UserID,
			&log.ClientID,
			&log.Scope,
			&log.IPAddress,
			&log.UserAgent,
			&log.Success,
			&errorMsg,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan auth log: %w", err)
		}

		if errorMsg.Valid {
			log.Error = errorMsg.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// CreateAuditLogsTable creates the audit logs table if it doesn't exist
func CreateAuditLogsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS auth.audit_logs (
			id VARCHAR(255) PRIMARY KEY,
			timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			action VARCHAR(50) NOT NULL,
			user_id VARCHAR(255),
			client_id VARCHAR(255),
			scope TEXT,
			ip_address INET,
			user_agent TEXT,
			success BOOLEAN NOT NULL,
			error TEXT
		);

		CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON auth.audit_logs(user_id);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON auth.audit_logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON auth.audit_logs(action);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create audit logs table: %w", err)
	}

	log.Info("Created auth.audit_logs table")
	return nil
}

// BangladeshProfile represents Bangladesh-specific profile data
type BangladeshProfile struct {
	ID           string `json:"id"`
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
	ProfileType  string `json:"profile_type"`            // patient, practitioner, organization
}

// MockBangladeshProfileService implements BangladeshProfileService for testing
type MockBangladeshProfileService struct {
	profiles map[string]*BangladeshProfile
}

// NewMockBangladeshProfileService creates a new mock profile service
func NewMockBangladeshProfileService() *MockBangladeshProfileService {
	profiles := make(map[string]*BangladeshProfile)

	// Add sample profiles
	profiles["patient-001"] = &BangladeshProfile{
		ID:          "patient-001",
		NID:         "1234567890123",
		BRN:         "BRN20234567890",
		UHID:        "UHID1234567890",
		Division:    "Dhaka",
		District:    "Dhaka",
		Upazila:     "Dhanmondi",
		Union:       "Kashimpur",
		Ward:        "1",
		ProfileType: "patient",
	}

	profiles["practitioner-001"] = &BangladeshProfile{
		ID:          "practitioner-001",
		NID:         "9876543210987",
		Division:    "Dhaka",
		District:    "Dhaka",
		Upazila:     "Dhanmondi",
		ProfileType: "practitioner",
	}

	profiles["organization-001"] = &BangladeshProfile{
		ID:          "organization-001",
		Division:    "Dhaka",
		District:    "Dhaka",
		Upazila:     "Dhanmondi",
		ProfileType: "organization",
	}

	return &MockBangladeshProfileService{
		profiles: profiles,
	}
}

// GetPatientProfile retrieves Bangladesh patient profile
func (m *MockBangladeshProfileService) GetPatientProfile(ctx context.Context, patientID string) (*BangladeshSpecificClaims, error) {
	profile, exists := m.profiles[patientID]
	if !exists || profile.ProfileType != "patient" {
		return nil, fmt.Errorf("patient profile not found")
	}

	return &BangladeshSpecificClaims{
		NID:      profile.NID,
		BRN:      profile.BRN,
		UHID:     profile.UHID,
		Division: profile.Division,
		District: profile.District,
		Upazila:  profile.Upazila,
		Union:    profile.Union,
		Ward:     profile.Ward,
	}, nil
}

// GetPractitionerProfile retrieves Bangladesh practitioner profile
func (m *MockBangladeshProfileService) GetPractitionerProfile(ctx context.Context, practitionerID string) (*BangladeshSpecificClaims, error) {
	profile, exists := m.profiles[practitionerID]
	if !exists || profile.ProfileType != "practitioner" {
		return nil, fmt.Errorf("practitioner profile not found")
	}

	return &BangladeshSpecificClaims{
		NID:      profile.NID,
		Division: profile.Division,
		District: profile.District,
		Upazila:  profile.Upazila,
	}, nil
}

// GetOrganizationProfile retrieves Bangladesh organization profile
func (m *MockBangladeshProfileService) GetOrganizationProfile(ctx context.Context, organizationID string) (*BangladeshSpecificClaims, error) {
	profile, exists := m.profiles[organizationID]
	if !exists || profile.ProfileType != "organization" {
		return nil, fmt.Errorf("organization profile not found")
	}

	return &BangladeshSpecificClaims{
		Division: profile.Division,
		District: profile.District,
		Upazila:  profile.Upazila,
	}, nil
}

// MockAuditLogger implements AuditLogger for testing
type MockAuditLogger struct {
	logs []AuditLog
}

// NewMockAuditLogger creates a new mock audit logger
func NewMockAuditLogger() *MockAuditLogger {
	return &MockAuditLogger{
		logs: make([]AuditLog, 0),
	}
}

// LogAuthEvent logs an authentication event
func (m *MockAuditLogger) LogAuthEvent(ctx context.Context, event AuditLog) error {
	m.logs = append(m.logs, event)
	log.Infof("Mock audit log: %s - %s for client %s", event.Action, event.UserID, event.ClientID)
	return nil
}

// GetAuthLogs retrieves authentication logs for a user
func (m *MockAuditLogger) GetAuthLogs(ctx context.Context, userID string, limit int) ([]AuditLog, error) {
	var userLogs []AuditLog
	for _, log := range m.logs {
		if log.UserID == userID {
			userLogs = append(userLogs, log)
			if len(userLogs) >= limit {
				break
			}
		}
	}
	return userLogs, nil
}
