# FHIR AuditEvent Policy

> **Document:** FHIR-AUDIT-POLICY.md | **Version:** 2.0.0
> **Standard:** HIPAA § 164.312(b) | **FHIR Version:** R5

---

## Mandatory AuditEvents

The following operations MUST generate a FHIR AuditEvent:

| Operation | Trigger | HIPAA Rationale |
|-----------|---------|----------------|
| Read PHI resource | GET /fhir/R5/{ResourceType}/{id} | Track who viewed patient data |
| Search PHI | GET /fhir/R5/{ResourceType}?... | Track who searched patient records |
| Create PHI | POST /fhir/R5/{ResourceType} | Track who created patient records |
| Update PHI | PUT /fhir/R5/{ResourceType}/{id} | Track who modified patient records |
| Delete PHI | DELETE /fhir/R5/{ResourceType}/{id} | Track who deleted patient records |
| Export PHI | POST /fhir/R5/$export | Track bulk data exports |
| Failed authentication | POST /auth/token → 401 | Track unauthorized access attempts |

PHI resources: Patient, Encounter, Observation, Condition, MedicationRequest, Procedure, AllergyIntolerance, DiagnosticReport, DocumentReference, CarePlan, Consent, NutritionOrder

Non-PHI resources (Practitioner, Organization, ValueSet, etc.) do not require AuditEvents for reads.

---

## AuditEvent Implementation (Go)

```go
// zs-pkg-go-audit implements HIPAA-compliant audit logging
// Used by every FHIR handler after successful operation

type AuditLogger interface {
    LogRead(ctx context.Context, resourceType, resourceID, userID, tenantID string) error
    LogCreate(ctx context.Context, resourceType, resourceID, userID, tenantID string) error
    LogUpdate(ctx context.Context, resourceType, resourceID, userID, tenantID string) error
    LogDelete(ctx context.Context, resourceType, resourceID, userID, tenantID string) error
    LogSearch(ctx context.Context, resourceType, query, userID, tenantID string) error
    LogExport(ctx context.Context, userID, tenantID string) error
}

// Implementation writes async to:
// 1. audit.events PostgreSQL table (primary record)
// 2. NATS subject zs.audit.events (for streaming consumers)
```

---

## Audit Log Retention

- Minimum retention: **7 years** (HIPAA requirement)
- Storage: PostgreSQL `audit.events` table with partition pruning after 7 years
- Backup: Included in regular PostgreSQL backup to Cloudflare R2
- Access: Read-only for audit reviewers; no delete permission for app users

---

## Audit Log Review Schedule

| Review | Frequency | Who |
|--------|-----------|-----|
| Failed authentication review | Weekly | `@DevOps-Ariful-Islam` |
| Unusual access pattern review | Monthly | `@arwa-zarish` |
| Full audit log review | Quarterly | All owners |
