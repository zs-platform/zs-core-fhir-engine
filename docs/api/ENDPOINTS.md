# REST Endpoints

Complete reference for all FHIR REST API endpoints and platform services.

## Base URLs

### FHIR API
```
http://localhost:8080/fhir/R5
```

### Platform Services
```
http://localhost:8080
```

## Authentication Endpoints (SMART on FHIR 2.1)

As per ADR-0001 (Go Backend) and ADR standards, all authentication follows SMART on FHIR 2.1 specification.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/auth/authorize` | GET | OAuth2 authorization endpoint |
| `/auth/token` | POST | Token exchange endpoint |
| `/auth/revoke` | POST | Token revocation |
| `/auth/userinfo` | GET | User information |
| `/auth/.well-known/openid_configuration` | GET | OIDC discovery document |
| `/auth/.well-known/jwks.json` | GET | JSON Web Key Set |

### Example Authentication Flow

```bash
# 1. Redirect user to authorization endpoint
curl "http://localhost:8080/auth/authorize?response_type=code&client_id=zs-fhir-engine&redirect_uri=http://localhost:3000/callback&scope=openid profile patient/*.read"

# 2. Exchange code for token
curl -X POST http://localhost:8080/auth/token \
  -d "grant_type=authorization_code" \
  -d "code=AUTH_CODE" \
  -d "client_id=zs-fhir-engine" \
  -d "redirect_uri=http://localhost:3000/callback"
```

## FHIR Resource Operations (FHIR R5 - ADR-0002)

As per ADR-0002, all FHIR operations comply with FHIR R5 (5.0.0) specification.

### Create Resource

Create a new FHIR resource.

**Request**

```http
POST /fhir/R5/{resourceType}
Content-Type: application/fhir+json
X-Tenant-ID: {tenant-id}

{
  "resourceType": "Patient",
  "active": true,
  "name": [{
    "family": "Chowdhury",
    "given": ["Rahima"]
  }]
}
```

**Response**

```http
HTTP/1.1 201 Created
Content-Type: application/fhir+json
Location: /fhir/R5/Patient/550e8400-e29b-41d4-a716-446655440000

{
  "resourceType": "Patient",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "meta": {
    "versionId": "1",
    "lastUpdated": "2026-03-24T10:00:00Z"
  },
  "active": true,
  "name": [{
    "family": "Chowdhury",
    "given": ["Rahima"]
  }]
}
```

---

### Read Resource

Read a specific resource by ID.

**Request**

```http
GET /fhir/R5/{resourceType}/{id}
X-Tenant-ID: {tenant-id}
```

**Response**

```http
HTTP/1.1 200 OK
Content-Type: application/fhir+json

{
  "resourceType": "Patient",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "meta": {
    "versionId": "1",
    "lastUpdated": "2026-03-24T10:00:00Z"
  },
  "active": true,
  ...
}
```

If the resource is not found:

```http
HTTP/1.1 404 Not Found
Content-Type: application/fhir+json

{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "not-found",
    "diagnostics": "Resource not found"
  }]
}
```

---

### Update Resource

Update an existing resource.

**Request**

```http
PUT /fhir/R5/{resourceType}/{id}
Content-Type: application/fhir+json
X-Tenant-ID: {tenant-id}

{
  "resourceType": "Patient",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "meta": {
    "versionId": "1"
  },
  "active": true,
  "name": [{
    "family": "Chowdhury",
    "given": ["Rahima", "Begum"]
  }]
}
```

---

### Delete Resource

Delete a resource (soft delete with history preservation).

**Request**

```http
DELETE /fhir/R5/{resourceType}/{id}
X-Tenant-ID: {tenant-id}
```

**Response**

```http
HTTP/1.1 204 No Content
```

---

### Search Resources

Search for resources with full FHIR R5 search parameter support.

**Request**

```http
GET /fhir/R5/{resourceType}?{searchParameters}
X-Tenant-ID: {tenant-id}
```

**Response**

```http
HTTP/1.1 200 OK
Content-Type: application/fhir+json

{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 2,
  "link": [
    {"relation": "self", "url": "..."},
    {"relation": "next", "url": "..."}
  ],
  "entry": [
    {
      "resource": { ... }
    },
    {
      "resource": { ... }
    }
  ]
}
```

#### Search Parameters

- **_id**: Resource ID
- **_lastUpdated**: Last modified date
- **_tag**: Tag filter
- **_profile**: Profile filter
- **Resource-specific**: Varies by resource type

#### Search Modifiers

- `:exact` - Exact match
- `:contains` - Contains substring  
- `:text` - Text search
- `:not` - Negation
- `:in` - In set
- `:not-in` - Not in set

#### Search Prefixes (for dates/numbers)

- `eq` - Equal (default)
- `ne` - Not equal
- `gt` - Greater than
- `lt` - Less than
- `ge` - Greater or equal
- `le` - Less or equal
- `sa` - Starts after
- `eb` - Ends before
- `ap` - Approximately

#### Pagination

- `_count`: Number of results per page (default: 20, max: 200)
- `_offset`: Offset for pagination
- `_page`: Page number

#### Sorting

- `_sort`: Sort criteria (e.g., `_sort=family`, `_sort=-_lastUpdated`)

---

## Resource History & Versioning

As per ADR-0002 and FHIR R5 specification, full resource versioning is supported.

### Get Resource History

```http
GET /fhir/R5/{resourceType}/{id}/_history
X-Tenant-ID: {tenant-id}
```

### Get Specific Version

```http
GET /fhir/R5/{resourceType}/{id}/_history/{versionId}
X-Tenant-ID: {tenant-id}
```

### Restore Version

```http
PUT /fhir/R5/{resourceType}/{id}/_history/{versionId}/_restore
X-Tenant-ID: {tenant-id}
```

---

## FHIR Subscriptions (NATS JetStream - ADR-0004)

As per ADR-0004, FHIR R5 Subscriptions are backed by NATS 2.12.5 JetStream.

### Subscription Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/fhir/R5/Subscription` | GET, POST | List/Create subscriptions |
| `/fhir/R5/Subscription/{id}` | GET, PUT, DELETE | Manage subscription |

### Create Subscription

```http
POST /fhir/R5/Subscription
Content-Type: application/fhir+json
X-Tenant-ID: {tenant-id}

{
  "resourceType": "Subscription",
  "status": "active",
  "reason": "Monitor new patients",
  "criteria": "Patient?",
  "channel": {
    "type": "rest-hook",
    "endpoint": "https://myapp.example.com/webhook",
    "payload": "application/fhir+json",
    "header": ["Authorization: Bearer token"]
  }
}
```

### Supported Channel Types

- `rest-hook` - HTTP POST callback
- `websocket` - WebSocket connection
- `email` - Email notification
- `sms` - SMS notification  
- `message` - NATS message queue

---

## Analytics Dashboard (ADR-0001)

As per ADR-0001 Go backend standards, analytics are built-in.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/analytics/overview` | GET | System overview metrics |
| `/analytics/resources` | GET | Resource statistics |
| `/analytics/trends` | GET | Usage trends |
| `/analytics/users` | GET | Top users |

### Example Analytics Query

```bash
# Get system overview
curl "http://localhost:8080/analytics/overview?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z" \
  -H "X-Tenant-ID: {tenant-id}" \
  -H "Authorization: Bearer {token}"

# Get usage trends with hourly interval
curl "http://localhost:8080/analytics/trends?interval=1h" \
  -H "X-Tenant-ID: {tenant-id}"
```

---

## HL7 v2 Bridge

Legacy system integration supporting ADT, ORU, MDM message types.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hl7/v2/message` | POST | Receive HL7 message |
| `/hl7/v2/messages` | GET | List processed messages |
| `/hl7/v2/messages/{id}` | GET | Get specific message |
| `/hl7/v2/transform` | POST | Transform HL7 to FHIR |
| `/hl7/v2/config` | GET | Bridge configuration |

### Send HL7 Message

```http
POST /hl7/v2/message
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "message": "MSH|^~\\&|ADT1|...",
  "source": "hospital-his"
}
```

---

## DICOM Integration

Medical imaging storage with WADO-RS support.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/dicom/studies` | POST | Store DICOM study |
| `/dicom/studies/{studyUID}` | GET | Get study FHIR resources |
| `/dicom/search` | GET | Search DICOM images |
| `/wado/rs/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}` | GET | WADO-RS endpoint |
| `/dicom/config` | GET | DICOM configuration |

---

## AI/ML Services (ADR-0001)

Machine learning model registry and clinical decision support.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/ml/models` | GET, POST | List/Register models |
| `/ml/models/{id}` | GET, DELETE | Get/Delete model |
| `/ml/models/{id}/predict` | POST | Make prediction |
| `/ml/models/{id}/train` | POST | Train model |
| `/ml/models/{id}/validate` | POST | Validate model |
| `/ml/nlp/extract` | POST | Extract entities from text |
| `/ml/cds/analyze` | POST | Analyze patient for CDS |
| `/ml/cds/diagnose` | POST | Suggest diagnoses |

### Make Prediction

```http
POST /ml/models/risk-model/predict
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "data": {
    "age": 65,
    "conditions": ["diabetes", "hypertension"],
    "medications": ["metformin"]
  }
}
```

---

## Offline Sync (PowerSync - ADR-0012)

As per ADR-0012, mobile offline synchronization with bi-directional sync.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/sync` | POST | Process sync request |
| `/sync/status/{sessionId}` | GET | Get sync status |
| `/sync/pending` | GET | Get pending changes |
| `/sync/resolve` | POST | Resolve conflict |
| `/sync/config` | GET | Sync configuration |

### Initiate Sync

```http
POST /sync
Content-Type: application/json
X-Tenant-ID: {tenant-id}
X-Device-ID: mobile-123

{
  "lastSyncTime": "2024-01-01T00:00:00Z",
  "clientChanges": [
    {
      "resourceType": "Patient",
      "action": "create",
      "resource": { ... }
    }
  ]
}
```

---

## Population Health

Population health management and care gap analysis.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/population/registry` | GET | Get patient registry |
| `/population/conditions` | GET | Condition prevalence report |
| `/population/utilization` | GET | Care utilization report |
| `/population/quality` | GET | Quality measures report |
| `/population/report` | POST | Generate comprehensive report |
| `/population/caregaps` | GET | Identify care gaps |
| `/population/risk` | GET | Calculate risk scores |

### Generate Health Report

```http
POST /population/report
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "from": "2024-01-01T00:00:00Z",
  "to": "2024-12-31T23:59:59Z"
}
```

---

## Terminology Endpoints

### Expand ValueSet

Expand a ValueSet to get all included concepts (ICD-11 support per ADR-0002).

**Request**

```http
GET /fhir/R5/ValueSet/$expand?url={system}
X-Tenant-ID: {tenant-id}
```

**Example**

```bash
# Expand ICD-11 codes
curl "http://localhost:8080/fhir/R5/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms" \
  -H "X-Tenant-ID: {tenant-id}"

# Filter results
curl "http://localhost:8080/fhir/R5/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms&filter=hyper" \
  -H "X-Tenant-ID: {tenant-id}"
```

**Response**

```json
{
  "resourceType": "ValueSet",
  "expansion": {
    "contains": [
      {
        "system": "http://id.who.int/icd/release/11/mms",
        "code": "BA00",
        "display": "Essential hypertension"
      }
    ]
  }
}
```

---

## System & Health Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/healthz` | GET | Health check |
| `/readyz` | GET | Readiness check |
| `/metrics` | GET | Prometheus metrics |
| `/version` | GET | API version info |

---

## Error Handling

All errors return FHIR `OperationOutcome` resources per FHIR R5 specification.

### Bad Request (400)

```http
HTTP/1.1 400 Bad Request
Content-Type: application/fhir+json

{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "invalid",
    "diagnostics": "Invalid request: Missing required field 'resourceType'"
  }]
}
```

### Unauthorized (401)

```http
HTTP/1.1 401 Unauthorized
Content-Type: application/fhir+json
WWW-Authenticate: Bearer

{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "security",
    "diagnostics": "Invalid or missing authentication token"
  }]
}
```

### Forbidden (403)

```http
HTTP/1.1 403 Forbidden
Content-Type: application/fhir+json

{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "security",
    "diagnostics": "Insufficient permissions for this operation"
  }]
}
```

### Not Found (404)

```http
HTTP/1.1 404 Not Found
Content-Type: application/fhir+json

{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "not-found",
    "diagnostics": "Resource not found"
  }]
}
```

### Rate Limited (429)

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/fhir+json
Retry-After: 60

{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "throttled",
    "diagnostics": "Rate limit exceeded. Please retry after 60 seconds."
  }]
}
```

---

## Multi-tenancy Headers

All requests must include tenant identification for proper data isolation (per ADR-0003 PostgreSQL RLS).

| Header | Required | Description |
|--------|----------|-------------|
| `X-Tenant-ID` | Yes | Tenant identifier for data isolation |
| `Authorization` | Yes (for protected endpoints) | Bearer token from SMART on FHIR |
| `X-Device-ID` | No (required for sync) | Device identifier for offline sync |
| `X-Request-ID` | No | Request tracking ID |

---

## Rate Limiting

API rate limits are enforced per client (per ADR-0001 security standards):

- **Default**: 100 requests per minute
- **Burst**: 150 requests
- **Headers**:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests in window
  - `X-RateLimit-Reset`: Time when limit resets
  - `Retry-After`: Seconds to wait when rate limited

---

## Content Types

- **Request/Response**: `application/fhir+json`
- **Patch**: `application/json-patch+json`
- **Binary**: `application/octet-stream`
- **Form Data**: `multipart/form-data` (for DICOM uploads)
