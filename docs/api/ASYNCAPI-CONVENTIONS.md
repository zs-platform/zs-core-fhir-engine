# AsyncAPI 3.0 Event Schema Standards

> **Version:** 2.0.0 | **Updated:** 2026-03-24
> **Spec:** AsyncAPI 3.0 — https://www.asyncapi.com/docs/reference/specification/v3.0.0

---

## Overview

ZarishSphere uses NATS 2.12.5 JetStream for all asynchronous service-to-service communication. Every published event must have an AsyncAPI 3.0 schema definition.

---

## Subject (Topic) Naming Convention

```
zs.{domain}.{resource}.{action}

Domain examples: fhir, alert, platform, clinical
Resource: the entity type (patient, encounter, observation)
Action: created, updated, deleted, flagged, synced

Examples:
  zs.fhir.Patient.created
  zs.fhir.Observation.created
  zs.alert.ewars.threshold_exceeded
  zs.platform.service.health_check
  zs.clinical.nutrition.muac_critical
```

**Rules:**
- All lowercase
- Dots as separators (NATS hierarchy)
- Wildcards: `zs.fhir.*.created` (all FHIR creates)
- Wildcards: `zs.>` (all ZS events)

---

## Event Envelope Schema

Every event published to NATS must use this envelope:

```json
{
  "$schema": "https://zarishsphere.com/schema/event/v1",
  "id": "uuid-v7-here",
  "specVersion": "1.0",
  "type": "zs.fhir.Patient.created",
  "source": "zs-svc-patient",
  "subject": "Patient/uuid",
  "time": "2026-03-24T10:00:00Z",
  "dataContentType": "application/fhir+json",
  "tenantId": "bgd-cxb-camp-1w",
  "data": {
    "resourceType": "Patient",
    "id": "uuid",
    "...": "full FHIR R5 resource"
  }
}
```

---

## AsyncAPI 3.0 Schema File

Each service that publishes events must include an `asyncapi.yaml` in its `docs/` directory:

```yaml
asyncapi: 3.0.0
info:
  title: zs-svc-patient Events
  version: 1.0.0
  description: Events published by the patient registration service

channels:
  patientCreated:
    address: zs.fhir.Patient.created
    messages:
      patientCreated:
        $ref: '#/components/messages/PatientCreated'

operations:
  publishPatientCreated:
    action: send
    channel:
      $ref: '#/channels/patientCreated'

components:
  messages:
    PatientCreated:
      name: PatientCreated
      title: Patient Created
      payload:
        $ref: '#/components/schemas/FHIRPatient'

  schemas:
    FHIRPatient:
      type: object
      required: [resourceType, id, tenantId]
      properties:
        resourceType:
          type: string
          enum: [Patient]
        id:
          type: string
          format: uuid
        tenantId:
          type: string
```

---

## Dead Letter Queue (DLQ)

Failed message processing routes to: `zs.dlq.{original.subject}`

Example: Failed processing of `zs.fhir.Patient.created` → `zs.dlq.zs.fhir.Patient.created`

DLQ messages include the original message plus:
```json
{
  "dlq": {
    "originalSubject": "zs.fhir.Patient.created",
    "failedAt": "2026-03-24T10:05:00Z",
    "reason": "terminology lookup timeout",
    "attemptCount": 5
  }
}
```

---

## Event Retention (JetStream)

| Subject Pattern | Retention | Max Messages | Storage |
|----------------|-----------|-------------|---------|
| `zs.fhir.>` | 30 days | 10M | File (persisted) |
| `zs.alert.>` | 90 days | 1M | File |
| `zs.platform.>` | 7 days | 100k | Memory |
| `zs.dlq.>` | 60 days | 1M | File |
