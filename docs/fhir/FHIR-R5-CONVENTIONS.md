# ZarishSphere FHIR R5 Conventions

> **Document:** FHIR-R5-CONVENTIONS.md | **Version:** 2.0.0
> **FHIR Version:** R5 (5.0.0)

---

## Resource ID Convention

ZarishSphere uses UUID v7 for all FHIR resource IDs:

```go
// UUID v7 is time-ordered, enabling chronological sorting
// PostgreSQL 18.3 provides uuidv7() natively
id := db.UUIDv7()  // Format: 0192fbad-xxxx-7xxx-xxxx-xxxxxxxxxxxx
```

Resource IDs MUST:
- Be UUID v7 format
- Be unique within the tenant
- Never be reused after deletion

---

## Tenant Isolation Pattern

Every FHIR resource stored in ZarishSphere includes a tenant identifier in `meta.tag`:

```json
{
  "resourceType": "Patient",
  "id": "0192fbad-1234-7abc-def0-123456789012",
  "meta": {
    "lastUpdated": "2026-03-24T10:30:00Z",
    "versionId": "1",
    "tag": [
      {
        "system": "https://zarishsphere.com/tags/tenant",
        "code": "bgd-cxb-camp-1w",
        "display": "Cox's Bazar Camp 1W Health Post"
      },
      {
        "system": "https://zarishsphere.com/tags/program",
        "code": "bgd-refugee-response",
        "display": "Bangladesh Refugee Response"
      }
    ]
  }
}
```

---

## Required Metadata on All Resources

```json
{
  "meta": {
    "lastUpdated": "required — set by server on every write",
    "versionId": "required — integer, incremented on update",
    "source": "required — URI identifying the system that created/last updated",
    "tag": "required — tenant and program tags (see above)"
  }
}
```

---

## AuditEvent Requirements

Every read, write, or delete of a FHIR resource containing PHI MUST generate a FHIR AuditEvent with:

```json
{
  "resourceType": "AuditEvent",
  "type": { "system": "http://terminology.hl7.org/CodeSystem/audit-event-type", "code": "rest" },
  "action": "R",  // C=create, R=read, U=update, D=delete
  "recorded": "2026-03-24T10:30:00Z",
  "outcome": "0",  // 0=success, 8=error
  "agent": [{
    "who": { "identifier": { "value": "keycloak-user-uuid" } },
    "requestor": true
  }],
  "entity": [{
    "what": { "reference": "Patient/0192fbad-1234-7abc-def0-123456789012" },
    "type": { "code": "2" }
  }]
}
```

---

## FHIR Search Response Format

All FHIR search responses return a `Bundle` of type `searchset`:

```json
{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 47,
  "link": [
    { "relation": "self", "url": "/fhir/R5/Patient?_count=20&_offset=0" },
    { "relation": "next", "url": "/fhir/R5/Patient?_count=20&_offset=20" }
  ],
  "entry": [...]
}
```

Maximum page size: 200 resources. Default: 20.

---

## FHIR Error Response Format

All errors return a FHIR `OperationOutcome`:

```json
{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "invalid",
    "diagnostics": "Patient.name is required",
    "expression": ["Patient.name"]
  }]
}
```

HTTP status codes:
- `200` — successful read/search
- `201` — successful create (with `Location` header)
- `400` — bad request / validation failure
- `401` — missing or invalid authentication
- `403` — insufficient SMART scope
- `404` — resource not found
- `409` — version conflict (optimistic locking)
- `422` — unprocessable entity (FHIR validation failure)
- `500` — internal server error
