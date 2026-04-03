# REST API Design Guide

> **Document:** REST-DESIGN-GUIDE.md | **Version:** 2.0.0
> **Standards:** FHIR R5 (ADR-0002), Go 1.26.1 Backend (ADR-0001)

---

## URL Patterns

### FHIR R5 Resources (ADR-0002)

```
# Standard CRUD operations
GET    /fhir/R5/{ResourceType}                         # Search (FHIR R5)
POST   /fhir/R5/{ResourceType}                         # Create (FHIR R5)
GET    /fhir/R5/{ResourceType}/{id}                    # Read (FHIR R5)
PUT    /fhir/R5/{ResourceType}/{id}                    # Update (FHIR R5)
DELETE /fhir/R5/{ResourceType}/{id}                    # Soft delete (FHIR R5)

# Resource versioning & history (FHIR R5)
GET    /fhir/R5/{ResourceType}/{id}/_history           # Get history
GET    /fhir/R5/{ResourceType}/{id}/_history/{versionId} # Specific version
PUT    /fhir/R5/{ResourceType}/{id}/_history/{versionId}/_restore  # Restore version

# Batch/transaction
POST   /fhir/R5/                                       # Batch/transaction bundle
```

### FHIR R5 Operations (ADR-0002)

```
# Capability & validation
GET    /fhir/R5/metadata                               # CapabilityStatement
POST   /fhir/R5/{ResourceType}/$validate               # Validate resource

# Terminology services (ICD-11 per ADR-0002)
GET    /fhir/R5/ValueSet/$expand?url={url}             # Expand ValueSet
POST   /fhir/R5/ValueSet/$validate-code                # Validate code
GET    /fhir/R5/CodeSystem/$lookup                     # Lookup code

# Bulk data export
POST   /fhir/R5/Patient/$export                        # Bulk export
GET    /fhir/R5/Group/{id}/$export                     # Group export
```

### Authentication (SMART on FHIR 2.1)

```
# OAuth2 endpoints
GET    /auth/authorize                                   # OAuth2 authorization
POST   /auth/token                                       # Token exchange
POST   /auth/revoke                                      # Token revocation
GET    /auth/userinfo                                    # User information
GET    /auth/.well-known/openid_configuration           # OIDC discovery
GET    /auth/.well-known/jwks.json                       # JWKS endpoint
```

### FHIR Subscriptions (NATS JetStream - ADR-0004)

```
GET    /fhir/R5/Subscription                            # List subscriptions
POST   /fhir/R5/Subscription                             # Create subscription
GET    /fhir/R5/Subscription/{id}                       # Read subscription
PUT    /fhir/R5/Subscription/{id}                      # Update subscription
DELETE /fhir/R5/Subscription/{id}                      # Delete subscription
```

### Analytics Dashboard (ADR-0001)

```
GET    /analytics/overview                               # System overview
GET    /analytics/resources                              # Resource statistics
GET    /analytics/trends                                 # Usage trends
GET    /analytics/users                                  # Top users
```

### HL7 v2 Bridge (Legacy Integration)

```
POST   /hl7/v2/message                                   # Receive HL7 message
GET    /hl7/v2/messages                                  # List HL7 messages
GET    /hl7/v2/messages/{id}                            # Get HL7 message
POST   /hl7/v2/transform                                 # Transform to FHIR
GET    /hl7/v2/config                                    # Bridge configuration
```

### DICOM Integration

```
POST   /dicom/studies                                    # Store DICOM study
GET    /dicom/studies/{studyUID}                         # Get study
GET    /dicom/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}  # Get instance
GET    /dicom/search                                     # Search DICOM
GET    /wado/rs/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}  # WADO-RS
GET    /dicom/config                                     # DICOM configuration
```

### AI/ML Services (ADR-0001)

```
GET    /ml/models                                        # List models
POST   /ml/models                                        # Register model
GET    /ml/models/{id}                                   # Get model
DELETE /ml/models/{id}                                   # Delete model
POST   /ml/models/{id}/predict                           # Make prediction
POST   /ml/models/{id}/train                             # Train model
POST   /ml/models/{id}/validate                          # Validate model
POST   /ml/nlp/extract                                   # NLP entity extraction
POST   /ml/cds/analyze                                   # CDS analysis
POST   /ml/cds/diagnose                                  # Diagnosis suggestions
```

### Offline Sync (PowerSync - ADR-0012)

```
POST   /sync                                             # Process sync
GET    /sync/status/{sessionId}                          # Sync status
GET    /sync/pending                                     # Pending changes
POST   /sync/resolve                                     # Resolve conflicts
GET    /sync/config                                      # Sync configuration
```

### Population Health

```
GET    /population/registry                              # Patient registry
GET    /population/conditions                            # Condition report
GET    /population/utilization                           # Utilization report
GET    /population/quality                                # Quality measures
POST   /population/report                                  # Generate report
GET    /population/caregaps                              # Care gaps
GET    /population/risk                                   # Risk scores
```

### System Health & Metrics

```
GET    /healthz                                          # Health check (no auth)
GET    /readyz                                           # Readiness check
GET    /metrics                                          # Prometheus metrics
GET    /version                                          # API version info
```

---

## Request Headers

### Required Headers

```
Authorization: Bearer {jwt}                    # Required on all PHI endpoints (SMART on FHIR 2.1)
Content-Type: application/fhir+json           # Required for FHIR resource bodies
X-Tenant-ID: bgd-cxb-camp-1w                   # Required per ADR-0003 (PostgreSQL RLS)
```

### Optional Headers

```
Accept: application/fhir+json               # Request FHIR JSON response format
X-Device-ID: mobile-123                       # Required for offline sync (ADR-0012)
X-Request-ID: {uuid}                         # Optional - for request tracing
If-Match: W/"{versionId}"                      # For conditional updates (optimistic locking)
If-None-Match: *                              # For conditional creates
```

---

## Response Headers

### Standard FHIR R5 Headers

```
Content-Type: application/fhir+json; charset=utf-8
Location: /fhir/R5/Patient/{id}              # On 201 Created
ETag: W/"version-id"                          # For conditional updates
Last-Modified: Thu, 24 Mar 2026 ...            # RFC 7231 format
X-Request-ID: {uuid}                          # Echo the request ID
X-Tenant-ID: {tenant-id}                      # Echo tenant ID
```

### Rate Limiting Headers (ADR-0001)

```
X-RateLimit-Limit: 100                         # Maximum requests allowed
X-RateLimit-Remaining: 95                      # Remaining requests in window
X-RateLimit-Reset: 1700000000                  # Unix timestamp when limit resets
Retry-After: 60                                # Seconds to wait (on 429)
```

---

## Pagination (FHIR R5)

Use FHIR Bundle links for pagination:

```
GET /fhir/R5/Patient?_count=20&_offset=0
→ Bundle.link[relation=next].url = /fhir/R5/Patient?_count=20&_offset=20
→ Bundle.link[relation=prev].url = /fhir/R5/Patient?_count=20&_offset=0
→ Bundle.total = 47
```

- Maximum `_count`: 200 per ADR-0001 performance targets
- Default `_count`: 20
- Use `_offset` for offset-based pagination
- Alternative: Use `_page` for page-based pagination

---

## Search Parameters (FHIR R5)

### Common Parameters

- `_id` - Resource ID
- `_lastUpdated` - Last modified date
- `_tag` - Tag filter
- `_profile` - Profile filter
- `_security` - Security label

### Search Modifiers (FHIR R5)

- `:exact` - Exact match
- `:contains` - Contains substring
- `:text` - Text search
- `:not` - Negation
- `:in` - In set
- `:not-in` - Not in set
- `:below` - Hierarchy below
- `:above` - Hierarchy above

### Search Prefixes (FHIR R5)

- `eq` - Equal (default)
- `ne` - Not equal
- `gt` - Greater than
- `lt` - Less than
- `ge` - Greater or equal
- `le` - Less or equal
- `sa` - Starts after
- `eb` - Ends before
- `ap` - Approximately

---

## Versioning

API version is embedded in the URL: `/fhir/R5/` (FHIR version, not service version).

Service versions are tracked via SemVer in the Docker image tag and Helm chart version.

Platform endpoints (analytics, sync, etc.) do not include FHIR version in URL.

---

## Multi-Tenancy (ADR-0003)

Per ADR-0003 (PostgreSQL-only with RLS), all requests must include tenant context:

1. **Primary**: `X-Tenant-ID` header
2. **Fallback**: JWT claim `tenant_id`
3. **Default**: `default` tenant (for single-tenant deployments)

All data is isolated per tenant using PostgreSQL Row-Level Security (RLS).

---

## Error Handling

All errors return FHIR `OperationOutcome` resources with proper HTTP status codes:

- `400` - Bad Request (invalid syntax or parameters)
- `401` - Unauthorized (missing/invalid authentication)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource doesn't exist)
- `409` - Conflict (version mismatch or duplicate)
- `422` - Unprocessable Entity (FHIR validation failure)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

---

## AsyncAPI & NATS Events (ADR-0004)

For real-time events, use NATS JetStream per ADR-0004:

```
# Resource events
zs.fhir.{ResourceType}.created
zs.fhir.{ResourceType}.updated
zs.fhir.{ResourceType}.deleted

# Subscription events
zs.fhir.Subscription.triggered

# Audit events
zs.audit.event

# Platform events
zs.platform.service.health
```

See `ASYNCAPI-CONVENTIONS.md` for complete event schema documentation.

---

## Performance Targets (ADR-0001)

Per ADR-0001 Go Backend standards:

- Response time: < 50ms for 95% of requests
- Concurrent users: 1000+
- Database connections: 1000+
- Rate limit: 100 req/min per client

---

## Content Types

- **FHIR Resources**: `application/fhir+json`
- **FHIR Patch**: `application/json-patch+json`
- **Binary**: `application/octet-stream`
- **Form Upload**: `multipart/form-data` (DICOM)
