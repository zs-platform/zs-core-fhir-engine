# API Overview

> **Version:** 2.0.0 | **Updated:** 2026-04-03
> **Standards:** FHIR R5 (ADR-0002), Go 1.26.1 Backend (ADR-0001)

The ZarishSphere FHIR Engine provides a comprehensive, production-ready API for healthcare data management with enterprise-grade features.

---

## 🏥 API Components

### 1. **FHIR R5 REST Server** (ADR-0002)
Full RESTful API for FHIR R5 resources with complete CRUD operations, search, and history.

### 2. **SMART on FHIR 2.1 Authentication**
OAuth2/OpenID Connect authentication with Keycloak integration.

### 3. **FHIR Subscriptions** (ADR-0004)
Real-time event subscriptions backed by NATS 2.12.5 JetStream.

### 4. **Terminology Server** (ADR-0002)
ValueSet expansion and code lookup with ICD-11 support.

### 5. **Analytics Dashboard** (ADR-0001)
Usage metrics, resource statistics, and performance monitoring.

### 6. **HL7 v2 Bridge**
Legacy system integration with message transformation.

### 7. **DICOM Integration**
Medical imaging storage with WADO-RS support.

### 8. **AI/ML Services** (ADR-0001)
Model registry, clinical decision support, and predictive analytics.

### 9. **Offline Sync** (ADR-0012)
PowerSync-based bi-directional mobile synchronization.

### 10. **Population Health**
Care gap analysis, risk stratification, and health reporting.

## 🔐 Authentication & Authorization

### SMART on FHIR 2.1 Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/auth/authorize` | OAuth2 authorization endpoint |
| `POST` | `/auth/token` | Token exchange |
| `POST` | `/auth/revoke` | Token revocation |
| `GET` | `/auth/userinfo` | User information |
| `GET` | `/auth/.well-known/openid_configuration` | OIDC discovery |
| `GET` | `/auth/.well-known/jwks.json` | JWKS endpoint |

### Authentication Flow

```bash
# 1. Authorization request
curl "http://localhost:8080/auth/authorize?response_type=code&client_id=zs-fhir-engine&redirect_uri=http://localhost:3000/callback&scope=openid profile patient/*.read"

# 2. Token exchange
curl -X POST http://localhost:8080/auth/token \
  -d "grant_type=authorization_code" \
  -d "code=AUTH_CODE" \
  -d "client_id=zs-fhir-engine"
```

---

## 🚀 Health & Monitoring Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/healthz` | Health check (no auth) |
| `GET` | `/readyz` | Readiness check |
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/version` | API version info |

### Health Check Response

```json
{
  "status": "ok",
  "timestamp": "2026-04-03T10:00:01Z",
  "uptime_seconds": 3600,
  "version": "2.0.0"
}
```

---

## 🏥 FHIR R5 REST Server

### Resource Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/fhir/R5/metadata` | CapabilityStatement |
| `POST` | `/fhir/R5/{resourceType}` | Create resource |
| `GET` | `/fhir/R5/{resourceType}` | Search resources |
| `GET` | `/fhir/R5/{resourceType}/{id}` | Read resource |
| `PUT` | `/fhir/R5/{resourceType}/{id}` | Update resource |
| `DELETE` | `/fhir/R5/{resourceType}/{id}` | Soft delete |
| `GET` | `/fhir/R5/{resourceType}/{id}/_history` | Resource history |
| `GET` | `/fhir/R5/{resourceType}/{id}/_history/{versionId}` | Specific version |

### Supported Resource Types

- **Patient** - Demographics with Bangladesh identifiers (NID, BRN, UHID)
- **Observation** - Clinical measurements
- **Encounter** - Patient visits
- **Condition** - Diagnoses with ICD-11
- **Medication** - Medication information
- **AllergyIntolerance** - Allergies
- **Procedure** - Medical procedures
- **DiagnosticReport** - Lab and imaging reports
- **ImagingStudy** - DICOM imaging references
- **CarePlan** - Patient care plans
- **CareTeam** - Healthcare teams
- **Subscription** - Real-time subscriptions
- **And more** - Full FHIR R5 resource support

### Example: Create Patient

```bash
curl -X POST http://localhost:8080/fhir/R5/Patient \
  -H "Content-Type: application/fhir+json" \
  -H "X-Tenant-ID: hospital-abc" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "resourceType": "Patient",
    "active": true,
    "name": [{"family": "Chowdhury", "given": ["Rahima"]}],
    "gender": "female",
    "birthDate": "1990-05-15",
    "identifier": [{
      "system": "http://bangladesh.gov/nid",
      "value": "1234567890"
    }]
  }'
```

### Example: Search with Parameters

```bash
# Search by name
curl "http://localhost:8080/fhir/R5/Patient?family=Chowdhury&given=Rahima" \
  -H "X-Tenant-ID: hospital-abc" \
  -H "Authorization: Bearer {token}"

# Search with date range
curl "http://localhost:8080/fhir/R5/Observation?date=gt2024-01-01&date=lt2024-12-31" \
  -H "X-Tenant-ID: hospital-abc"

# Search with identifier
curl "http://localhost:8080/fhir/R5/Patient?identifier=NID|1234567890" \
  -H "X-Tenant-ID: hospital-abc"

# Paginated search
curl "http://localhost:8080/fhir/R5/Patient?_count=20&_offset=0&_sort=family" \
  -H "X-Tenant-ID: hospital-abc"
```

---

## 📡 FHIR Subscriptions (NATS JetStream - ADR-0004)

### Subscription Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/fhir/R5/Subscription` | List subscriptions |
| `POST` | `/fhir/R5/Subscription` | Create subscription |
| `GET` | `/fhir/R5/Subscription/{id}` | Read subscription |
| `PUT` | `/fhir/R5/Subscription/{id}` | Update subscription |
| `DELETE` | `/fhir/R5/Subscription/{id}` | Delete subscription |

### Create Subscription

```bash
curl -X POST http://localhost:8080/fhir/R5/Subscription \
  -H "Content-Type: application/fhir+json" \
  -H "X-Tenant-ID: hospital-abc" \
  -H "Authorization: Bearer {token}" \
  -d '{
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
  }'
```

### Supported Channel Types

- `rest-hook` - HTTP POST callback
- `websocket` - WebSocket connection
- `email` - Email notification
- `sms` - SMS notification
- `message` - NATS message queue

---

## 📊 Analytics Dashboard (ADR-0001)

### Analytics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/analytics/overview` | System overview |
| `GET` | `/analytics/resources` | Resource statistics |
| `GET` | `/analytics/trends` | Usage trends |
| `GET` | `/analytics/users` | Top users |

### Example Queries

```bash
# System overview
curl "http://localhost:8080/analytics/overview?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z" \
  -H "X-Tenant-ID: hospital-abc" \
  -H "Authorization: Bearer {token}"

# Usage trends
curl "http://localhost:8080/analytics/trends?interval=1h" \
  -H "X-Tenant-ID: hospital-abc"
```

---

## 🔗 HL7 v2 Bridge

### Bridge Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/hl7/v2/message` | Receive HL7 message |
| `GET` | `/hl7/v2/messages` | List messages |
| `GET` | `/hl7/v2/messages/{id}` | Get message |
| `POST` | `/hl7/v2/transform` | Transform to FHIR |
| `GET` | `/hl7/v2/config` | Bridge configuration |

### Send HL7 Message

```bash
curl -X POST http://localhost:8080/hl7/v2/message \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: hospital-abc" \
  -d '{
    "message": "MSH|^~\\&|ADT1|...",
    "source": "hospital-his"
  }'
```

---

## 🖼️ DICOM Integration

### DICOM Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/dicom/studies` | Store DICOM study |
| `GET` | `/dicom/studies/{studyUID}` | Get study |
| `GET` | `/dicom/search` | Search DICOM |
| `GET` | `/wado/rs/studies/{studyUID}/series/{seriesUID}/instances/{instanceUID}` | WADO-RS |
| `GET` | `/dicom/config` | DICOM configuration |

---

## 🤖 AI/ML Services (ADR-0001)

### ML Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/ml/models` | List models |
| `POST` | `/ml/models` | Register model |
| `GET` | `/ml/models/{id}` | Get model |
| `DELETE` | `/ml/models/{id}` | Delete model |
| `POST` | `/ml/models/{id}/predict` | Make prediction |
| `POST` | `/ml/models/{id}/train` | Train model |
| `POST` | `/ml/models/{id}/validate` | Validate model |
| `POST` | `/ml/nlp/extract` | NLP extraction |
| `POST` | `/ml/cds/analyze` | CDS analysis |
| `POST` | `/ml/cds/diagnose` | Diagnosis suggestions |

### Make Prediction

```bash
curl -X POST http://localhost:8080/ml/models/risk-model/predict \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: hospital-abc" \
  -d '{
    "data": {
      "age": 65,
      "conditions": ["diabetes", "hypertension"],
      "medications": ["metformin"]
    }
  }'
```

---

## 🔄 Offline Sync (PowerSync - ADR-0012)

### Sync Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/sync` | Process sync |
| `GET` | `/sync/status/{sessionId}` | Sync status |
| `GET` | `/sync/pending` | Pending changes |
| `POST` | `/sync/resolve` | Resolve conflicts |
| `GET` | `/sync/config` | Sync configuration |

### Initiate Sync

```bash
curl -X POST http://localhost:8080/sync \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: hospital-abc" \
  -H "X-Device-ID: mobile-123" \
  -d '{
    "lastSyncTime": "2024-01-01T00:00:00Z",
    "clientChanges": [
      {
        "resourceType": "Patient",
        "action": "create",
        "resource": { ... }
      }
    ]
  }'
```

---

## 🏥 Population Health

### Population Health Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/population/registry` | Patient registry |
| `GET` | `/population/conditions` | Condition report |
| `GET` | `/population/utilization` | Utilization report |
| `GET` | `/population/quality` | Quality measures |
| `POST` | `/population/report` | Generate report |
| `GET` | `/population/caregaps` | Care gaps |
| `GET` | `/population/risk` | Risk scores |

### Generate Report

```bash
curl -X POST http://localhost:8080/population/report \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: hospital-abc" \
  -d '{
    "from": "2024-01-01T00:00:00Z",
    "to": "2024-12-31T23:59:59Z"
  }'
```

---

## 🏥 Terminology Service (ICD-11 - ADR-0002)

### Terminology Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/fhir/R5/ValueSet/$expand` | Expand ValueSet |
| `POST` | `/fhir/R5/ValueSet/$validate-code` | Validate code |
| `GET` | `/fhir/R5/CodeSystem/$lookup` | Lookup code |

### Example: Expand ICD-11 ValueSet

```bash
curl "http://localhost:8080/fhir/R5/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms" \
  -H "X-Tenant-ID: hospital-abc"

# Filter results
curl "http://localhost:8080/fhir/R5/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms&filter=diabetes" \
  -H "X-Tenant-ID: hospital-abc"
```

---

## 🔧 Configuration & Headers
- `FHIR_DB_MODE` - Storage backend (memory/postgresql)
- `FHIR_LOG_LEVEL` - Logging level (debug/info/warn/error)
- `FHIR_TENANT_ID` - Multi-tenant identifier

### Required Headers

```
Authorization: Bearer {jwt}                    # Required for PHI endpoints
Content-Type: application/fhir+json           # Required for FHIR bodies
X-Tenant-ID: {tenant-id}                      # Required per ADR-0003 (PostgreSQL RLS)
```

### Optional Headers

```
Accept: application/fhir+json               # Request format
X-Device-ID: {device-id}                       # Required for offline sync
X-Request-ID: {uuid}                         # Request tracing
```

### Environment Variables

- `FHIR_SERVER_PORT` - Server port (default: 8080)
- `FHIR_SERVER_HOST` - Bind address (default: 0.0.0.0)
- `FHIR_DB_MODE` - Storage backend (memory/postgresql)
- `FHIR_LOG_LEVEL` - Logging level (debug/info/warn/error)
- `NATS_URL` - NATS server URL (ADR-0004)
- `KEYCLOAK_URL` - Keycloak URL for authentication

### Error Response
```json
{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "invalid",
    "details": {
      "text": "Resource validation failed"
    }
  }]
}
```

### Bundle Response (Search)
```json
{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 1,
  "entry": [{
    "resource": {
      "resourceType": "Patient",
      "id": "patient-123"
    }
  }]
}
```

## 🔍 Search Parameters

Basic search parameters are supported:

| Parameter | Type | Description |
|-----------|------|-------------|
| `_id` | token | Resource ID |
| `_lastUpdated` | date | Last updated date |
| `_tag` | token | Resource tags |
| `_profile` | uri | Profile URL |

### Resource-Specific Parameters
- **Patient**: `family`, `given`, `gender`, `birthDate`
- **Observation**: `code`, `subject`, `date`
- **Condition**: `code`, `subject`, `clinical-status`

## 🚀 Production Features

### Performance
- **Connection Pooling**: PostgreSQL connection management
- **Request Metrics**: Built-in performance monitoring
- **Health Checks**: Automated health monitoring
- **Graceful Shutdown**: Clean server termination

### Security
- **Multi-tenancy**: Tenant isolation
- **Audit Logging**: Request tracking
- **Input Validation**: Resource validation
- **Error Handling**: Secure error responses

### Scalability
- **Database Persistence**: PostgreSQL backend
- **Horizontal Scaling**: Load balancer ready
- **Caching**: Built-in caching layer
- **Monitoring**: Metrics and observability

## 📚 Quick Links

- [Server Configuration](/api/server) - Server setup and options
- [Complete Endpoint Reference](/api/endpoints) - All endpoints documented
- [FHIR Resources](/fhir/overview) - Resource documentation
- [Terminology Service](/terminology/overview) - Terminology details
- [Bangladesh Profiles](/fhir/profiles) - Local healthcare standards

## 🔧 Development Tools

### CLI Commands
```bash
# Start server
./zs-core-fhir-engine serve --port 8080

# Validate resources
./zs-core-fhir-engine validate patient.json

# Start terminology server
./zs-core-fhir-engine terminology --port 8081
```

### Testing
```bash
# Health check
curl http://localhost:8080/healthz

# Server metadata
curl http://localhost:8080/fhir/metadata

# Create test patient
curl -X POST http://localhost:8080/fhir/Patient \
  -H "Content-Type: application/fhir+json" \
  -d @test-patient.json
```

---

**Need help?** Check our [Installation Guide](/guide/installation) or open an issue on GitHub.
