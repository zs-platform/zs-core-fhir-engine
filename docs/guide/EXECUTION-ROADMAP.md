# Execution Roadmap (Program Plan)

This is a staged plan from the current repository state to a complete deployable healthcare platform.

## Phase 0 — Stabilize the foundation (2-4 weeks)

- Correct documentation claims to match actual features.
- Add clear capability matrix (implemented vs planned).
- Add CI quality gates for tests, lint, and security scan.
- Define formal version policy (FHIR, Go, API).

## Phase 1 — Core FHIR persistence (4-8 weeks)

- Add PostgreSQL-backed resource storage.
- Add migration framework.
- Implement read/search pagination and `_count` controls.
- Implement transaction bundle handling.

## Phase 2 — Conformance and validation (4-8 weeks)

- Add `CapabilityStatement` endpoint.
- Add profile validation pipeline with detailed `OperationOutcome`.
- Add conditional create/update and resource history endpoints.

## Phase 3 — Terminology expansion (4-8 weeks)

- Terminology index and import pipeline.
- `$expand`, `$validate-code`, `$lookup` hardening.
- ICD-11 operational package management.
- National/local terminology publication workflow.

## Phase 4 — HL7 v2 interoperability (6-10 weeks)

- Add HL7 v2 listener service.
- Build ADT/ORM/ORU mapping templates.
- Add queue + retry + dead-letter architecture.
- Add reconciliation dashboards.

## Phase 5 — Security, privacy, and trust (4-8 weeks)

- Integrate Keycloak for OAuth2/OIDC.
- Add SMART on FHIR scopes.
- Add detailed audit event storage.
- Add consent enforcement hooks.

## Phase 6 — Operations and public launch (3-6 weeks)

- Add observability stack and alerting.
- Create runbooks and incident playbooks.
- Harden deployment pipelines.
- Publish release documentation website with stakeholder-friendly guides.

## Parallel documentation deliverables

- Executive overview (non-technical).
- System operations manual.
- API and integration handbook.
- Terminology governance handbook.
- Compliance and data protection checklist.

## Suggested first milestone for your team

If you want quickest real-world progress, prioritize:

1. Persistent storage.
2. Security (authentication/authorization).
3. Terminology operations.
4. HL7 v2 gateway.

This order gives practical EMR/HIS integration value early while keeping future scale possible.
