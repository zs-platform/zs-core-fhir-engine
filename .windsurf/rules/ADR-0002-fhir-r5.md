---
trigger: always_on
---
# ADR-0002 — FHIR R5 as the Clinical Data Model

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @code-and-brain, @BGD-Health-Program

---

## Context and Problem Statement

ZarishSphere needs a clinical data model that:
- Enables interoperability with national HIS systems (DHIS2, OpenMRS, legacy)
- Supports the full EMR/EHR clinical domain (patient, encounter, observation, medication, etc.)
- Has stable, internationally recognised semantics
- Can represent structured clinical data with coded terminology
- Supports audit and provenance tracking

---

## Decision Drivers

- **Interoperability** — partner systems speak FHIR (or HL7v2 via bridge)
- **Completeness** — FHIR R5 covers all required clinical domains
- **Standards body** — HL7 International (not a proprietary schema)
- **SMART on FHIR** — app launch protocol requires FHIR
- **Topic subscriptions** — FHIR R5 introduces topic-based subscriptions (not in R4)
- **No proprietary lock** — any FHIR-compliant system can read our data

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **FHIR R5 (5.0.0)** | Topic subscriptions, R4 bridge available, current standard | Newer, fewer R5-native tools than R4 |
| FHIR R4 (4.0.1) | More mature tooling, wider adoption | No topic subscriptions, not current standard |
| OpenMRS data model | Strong EMR domain model | Proprietary, Java-only, poor interoperability |
| Custom schema | Full control | Zero interoperability, impossible to integrate |
| FHIR R4B | Transitional | Being phased out |

---

## Decision Outcome

**Chosen option: FHIR R5 (5.0.0)**

FHIR R5 is chosen over R4 because:
1. Topic-based subscriptions (required for real-time NATS integration)
2. Improved Appointment and Schedule resources
3. Enhanced CarePlan and CareTeam resources
4. Better provenance and audit support
5. It is the current standard — R4 will be in maintenance mode

An **R4↔R5 translation bridge** (`zs-core-fhir-r4-bridge`) is provided for partner systems that are still on R4. Internal storage is always R5.

---

## Positive Consequences

- Full interoperability with any FHIR R5-compliant system
- SMART on FHIR 2.1 app ecosystem available
- Topic subscriptions enable real-time event streaming via NATS
- Standard terminology bindings (LOINC, SNOMED, ICD-11) defined by HL7

## Negative Consequences / Trade-offs

- Must maintain R4 bridge for the ~73% of partner systems on R4
- FHIR R5 tooling is still maturing (fewer R5-native validators than R4)
- Validation requires HAPI FHIR validator (Java) in CI

---

## Links

- FHIR R5 specification: https://hl7.org/fhir/R5/
- fhir-toolbox-go: https://github.com/damedic/fhir-toolbox-go
