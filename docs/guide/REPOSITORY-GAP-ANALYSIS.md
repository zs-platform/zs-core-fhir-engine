# Repository Gap Analysis (Plain Language)

This page explains where the repository is strong today and what must be added to become a **complete production-grade FHIR R5 server** for real EMR/EHR/HIS usage.

## 1) What is already good

- You already have a large generated FHIR model library (R4 and R5 resource structs).
- You already have a command-line app (`zs-core-fhir`) that can run a basic server.
- You already have a documentation website pipeline (VitePress + GitHub Pages).
- You already include terminology assets such as FHIR R5 JSON packages and ICD-11 related starter code.

## 2) Critical gaps before hospital production use

Today, the server behaves like an in-memory demo API. For hospital-grade production, these are mandatory upgrades:

1. **Persistent Database**
   - Current data disappears when process restarts.
   - Need PostgreSQL (or equivalent) persistence and migrations.

2. **FHIR Conformance Features**
   - Need `CapabilityStatement`, conformance endpoints, paging, history, conditional operations, `_include`, `_revinclude`, and robust search parameter behavior.

3. **Security and Identity**
   - Need OAuth2/OIDC (SMART on FHIR), role-based access, audit logs, consent-aware access, and tenant isolation.

4. **Terminology Completeness**
   - Current `$expand` support is basic.
   - Need full terminology stack (`$validate-code`, `$lookup`, versioned code systems, import jobs, caching, licensing compliance).

5. **Validation Pipeline**
   - Need profile validation against StructureDefinition and Implementation Guides.
   - Need clear error responses via FHIR `OperationOutcome`.

6. **Integration Layer**
   - Need HL7 v2 ingestion pipeline (ADT/ORM/ORU), mapping rules to FHIR R5 resources, queueing, dead-letter handling.

7. **Operations and Reliability**
   - Need metrics, tracing, backups, disaster recovery, zero-downtime deployments, and SLO monitoring.

## 3) Reality check on version requirement (Go 1.26.0)

Go `1.26.0` is not currently available in this environment yet. The best practical strategy is:

- Use latest stable Go now in CI.
- Keep one compatibility matrix page in docs.
- Upgrade immediately once Go 1.26 is officially released.

## 4) What this repository can become

With staged implementation, this repository can evolve into:

- A robust FHIR R5 core server,
- A terminology service supporting ICD-11 + local code systems,
- An HL7 v2 gateway for legacy hospital systems,
- A published, free documentation hub on GitHub Pages,
- A reusable national template for EMR/EHR/HIS projects.

## 5) External source-of-truth integration

You mentioned another source repo: `https://github.com/zs-docs/zarish-sphere-ssot`.

Recommended pattern:

- Treat that repo as policy/business source-of-truth.
- Mirror only approved machine-readable artifacts into this repo (`CodeSystem`, `ValueSet`, profiles, implementation policies).
- Add a scheduled sync workflow (nightly) with validation gates before merge.

Proceed to:

- [Complete Server Blueprint](/guide/complete-server-blueprint)
- [Execution Roadmap](/guide/execution-roadmap)
