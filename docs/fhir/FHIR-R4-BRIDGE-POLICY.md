# FHIR R4 Bridge Policy

> **Version:** 2.0.0 | **Updated:** 2026-03-24
> **Related ADR:** ADR-0002 (FHIR R5 as clinical model)
> **Bridge Repository:** zs-core-fhir-r4-bridge

---

## Background

Approximately 73% of health information systems that ZarishSphere interacts with (OpenMRS, legacy DHIS2, hospital EMRs) use FHIR R4 (4.0.1). ZarishSphere stores all data as FHIR R5 internally but must exchange data with these R4 systems.

The R4 bridge (`zs-core-fhir-r4-bridge`) provides bidirectional R4↔R5 translation.

---

## When to Use the R4 Bridge

**Use the R4 bridge for:**
- Receiving data FROM systems that send FHIR R4
- Sending data TO systems that only accept FHIR R4
- Integration with OpenMRS (R4/R5)
- Integration with legacy HIS with FHIR R4 APIs

**Do NOT use the R4 bridge for:**
- Internal service-to-service communication (always use R5)
- Storing data — storage is always R5
- Generating FHIR subscriptions (always R5 topics)
- SMART on FHIR app launches (use R5)

---

## Translation Rules

### R4 → R5 (Inbound)

When receiving R4 from a partner system:

| R4 Resource | R5 Equivalent | Notes |
|-------------|--------------|-------|
| `Patient` | `Patient` | R5 adds `Patient.link` changes |
| `Observation` | `Observation` | R5 component handling differs |
| `MedicationRequest` | `MedicationRequest` | R5 rename: `medicationCodeableConcept` → `medication` |
| `Condition` | `Condition` | R5 `clinicalStatus` binding updated |
| `Appointment` | `Appointment` | R5 significant restructure — use bridge carefully |

### R5 → R4 (Outbound)

When sending R5 to an R4 system:

1. Translate R5-only features to closest R4 equivalent
2. Add `meta.tag` with `system: zs-r4-translated, code: true`
3. Log translation in AuditEvent
4. If translation is lossy, document which fields are dropped

---

## Lossy Translation Register

Some R5 features have no R4 equivalent. When translating R5→R4, these are dropped:

| R5 Field | R4 Equivalent | Lossiness |
|----------|--------------|-----------|
| `Subscription.topic` | None | Lost (R4 has no topic subscriptions) |
| `Appointment.previousAppointment` | None | Lost |
| `Patient.genderIdentity` | None (extension) | Downgraded to extension |

All lossy translations must be logged in the ZarishSphere AuditEvent trail.

---

## Testing Requirements

The bridge must have integration tests covering:
- Round-trip translation (R5 → R4 → R5) with no data loss for non-lossy fields
- All listed R4 resource types
- Error handling for malformed R4 input
- AuditEvent generation on every translation

---

## Versioning

The R4 bridge version is independent of the platform version. Bridge version schema:
```
{major}.{r4-spec-version}.{r5-spec-version}
Example: 1.4.0-5.0.0 (bridge v1, R4 4.0.x, R5 5.0.0)
```
