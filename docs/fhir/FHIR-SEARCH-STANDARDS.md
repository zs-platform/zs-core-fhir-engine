# FHIR Search Standards

> **Version:** 2.0.0 | **Updated:** 2026-03-24
> **Standard:** FHIR R5 — https://hl7.org/fhir/R5/search.html

---

## Overview

This document defines the required FHIR search parameters that every ZarishSphere service must support. Search parameters not listed here are optional.

---

## Universal Search Parameters (All Resources)

All FHIR resources in ZarishSphere must support:

| Parameter | Type | Description |
|-----------|------|-------------|
| `_id` | token | Search by FHIR resource ID |
| `_lastUpdated` | date | Filter by last modification date |
| `_tag` | token | Filter by resource tag |
| `_profile` | uri | Filter by declared profile |
| `_count` | number | Limit results per page (max: 100) |
| `_offset` | number | Pagination offset |
| `_sort` | string | Sort results (e.g., `_sort=-_lastUpdated`) |
| `_summary` | token | Return summary form (true/false/count) |

---

## Required Search Parameters by Resource Type

### Patient

| Parameter | Type | Example |
|-----------|------|---------|
| `family` | string | `?family=Ahmed` |
| `given` | string | `?given=Mohammed` |
| `birthdate` | date | `?birthdate=1990-01-15` |
| `gender` | token | `?gender=male` |
| `identifier` | token | `?identifier=NID\|ABC123` |
| `name` | string | `?name=Ahmed` (full-text) |
| `address` | string | `?address=Dhaka` |
| `address-country` | token | `?address-country=BD` |
| `_tenant` | token | `?_tenant=bgd-cxb-camp-1w` (ZS extension) |

### Observation

| Parameter | Type | Example |
|-----------|------|---------|
| `subject` | reference | `?subject=Patient/uuid` |
| `code` | token | `?code=http://loinc.org\|8310-5` |
| `date` | date | `?date=ge2026-01-01` |
| `category` | token | `?category=vital-signs` |
| `status` | token | `?status=final` |
| `value-quantity` | quantity | `?value-quantity=gt37\|Cel` |

### Encounter

| Parameter | Type | Example |
|-----------|------|---------|
| `subject` | reference | `?subject=Patient/uuid` |
| `status` | token | `?status=in-progress` |
| `date` | date | `?date=ge2026-01-01` |
| `class` | token | `?class=AMB` |
| `location` | reference | `?location=Location/bgd-camp1w` |

### MedicationRequest

| Parameter | Type | Example |
|-----------|------|---------|
| `subject` | reference | `?subject=Patient/uuid` |
| `status` | token | `?status=active` |
| `medication` | token | `?medication=rxcui\|1049502` |
| `authored` | date | `?authored=ge2026-01-01` |

---

## Search Chaining

ZarishSphere supports chained search parameters:

```
# Find all observations for patients named Ahmed
GET /fhir/R5/Observation?subject.family=Ahmed

# Find all encounters at a specific location in a district
GET /fhir/R5/Encounter?location.address=Teknaf
```

---

## Multi-Tenant Search Extension

All ZarishSphere search requests must include the tenant context. This is enforced at the JWT level — the `tenant_id` is extracted from the SMART scope.

The `_tenant` parameter is a ZarishSphere extension:

```
GET /fhir/R5/Patient?_tenant=bgd-cxb-camp-1w&family=Ahmed
```

If `_tenant` is not provided, it defaults to the tenant in the JWT claim. Cross-tenant search returns 403 Forbidden.

---

## Pagination Rules

- Default page size: 20
- Maximum page size: 100
- Use `_count` to set page size
- Use Bundle.link with `relation: next` for pagination
- Never return all resources without a limit

```json
{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 1247,
  "link": [
    { "relation": "self", "url": "/fhir/R5/Patient?family=Ahmed&_count=20" },
    { "relation": "next", "url": "/fhir/R5/Patient?family=Ahmed&_count=20&_offset=20" }
  ],
  "entry": [...]
}
```

---

## Performance Requirements

| Search Type | Max Response Time | Notes |
|-------------|------------------|-------|
| Single resource by ID | < 10ms | Cache first |
| Simple search (indexed field) | < 100ms | PostgreSQL GIN |
| Complex chained search | < 500ms | Optimise with explain |
| Full-text search (Typesense) | < 50ms | Patient name search |

---

*Implementation guidance: See ADR-0003 (PostgreSQL GIN indexing) and zs-pkg-go-db.*
