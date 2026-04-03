-- ZarishSphere FHIR Engine - PostgreSQL Schema
-- This schema creates the necessary tables for FHIR R5 resource storage
-- Follows ADR-0003: PostgreSQL 18.3 as the only database

-- Enable UUID v7 generation (PostgreSQL 18+)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create FHIR resources table
CREATE TABLE IF NOT EXISTS fhir.resources (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type   TEXT NOT NULL,
    fhir_id         TEXT NOT NULL,
    version_id      INTEGER NOT NULL DEFAULT 1,
    resource        JSONB NOT NULL,
    tenant_id       TEXT NOT NULL DEFAULT 'default',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ NULL,
    
    -- Unique constraint for active resources
    UNIQUE(resource_type, fhir_id, tenant_id),
    
    -- Index for fast JSON searches
    INDEX idx_resource_type (resource_type),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_fhir_id (fhir_id),
    INDEX idx_created_at (created_at),
    
    -- GIN index for JSON content
    INDEX idx_resource_gin (resource)
);

-- Create audit events table for compliance
CREATE TABLE IF NOT EXISTS fhir.audit_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type   TEXT NOT NULL,
    resource_id     TEXT NOT NULL,
    action          TEXT NOT NULL,
    user_id         TEXT,
    user_role       TEXT,
    timestamp       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    details         JSONB,
    tenant_id       TEXT NOT NULL DEFAULT 'default',
    
    -- Indexes for audit queries
    INDEX idx_audit_resource (resource_type, resource_id),
    INDEX idx_audit_timestamp (timestamp),
    INDEX idx_audit_tenant (tenant_id),
    INDEX idx_audit_user (user_id)
);

-- Create terminology cache table
CREATE TABLE IF NOT EXISTS fhir.terminology_cache (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system          TEXT NOT NULL,
    code            TEXT NOT NULL,
    display         TEXT NOT NULL,
    definition      JSONB,
    version         TEXT NOT NULL DEFAULT '1.0',
    expires_at      TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '24 hours'),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Indexes for terminology lookups
    INDEX idx_terminology_system_code (system, code),
    INDEX idx_terminology_expires (expires_at)
);

-- Row Level Security for multi-tenancy
ALTER TABLE fhir.resources ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only access their own tenant's data
CREATE POLICY tenant_isolation ON fhir.resources
    USING (tenant_id = current_setting('app.tenant_id'));

-- Policy: Audit events are read-only for regular users
CREATE POLICY audit_readonly ON fhir.audit_events
    USING (current_setting('app.user_role') = 'admin');

-- Grant necessary permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON fhir.resources TO web_app_role;
GRANT SELECT ON fhir.audit_events TO web_app_role;
GRANT SELECT, INSERT ON fhir.terminology_cache TO web_app_role;
