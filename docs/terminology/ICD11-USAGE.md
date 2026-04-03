# ICD-11 Usage Standards

> **Document:** ICD11-USAGE.md | **Version:** 2.0.0
> **Current Version:** ICD-11 2026-01

---

## When to Use ICD-11

Use ICD-11 codes in:
- `Condition.code` — diagnoses and problem list
- `Observation.code` — when observation represents a condition (e.g., screening result)
- `DiagnosticReport.conclusionCode` — final diagnostic conclusion
- DHIS2 reporting — ICD-11 maps to national reporting requirements

---

## Code Format

ICD-11 uses alphanumeric codes with dot separators:

```
BA00     — Pulmonary tuberculosis
BA00.0   — Pulmonary tuberculosis, confirmed
BA00.Z   — Pulmonary tuberculosis, unspecified
```

---

## FHIR Coding Pattern

```json
{
  "code": {
    "coding": [
      {
        "system": "http://id.who.int/icd/release/11/2026-01/mms",
        "code": "BA00.0",
        "display": "Pulmonary tuberculosis, confirmed"
      }
    ],
    "text": "Pulmonary tuberculosis, confirmed"
  }
}
```

The `system` URI MUST include the ICD-11 year version: `http://id.who.int/icd/release/11/{YEAR}/mms`

---

## Cross-mapping to ICD-10

Many national systems still report in ICD-10. ZarishSphere provides cross-maps in `zs-data-concept-maps`:

```
ICD-11 BA00.0 → ICD-10 A15.0 (Tuberculosis of lung, confirmed by sputum microscopy)
```

Use the concept maps for DHIS2 reporting to legacy systems.

---

## Common ICD-11 Codes for ZarishSphere Programs

| Condition | ICD-11 Code |
|-----------|------------|
| Pulmonary tuberculosis | BA00 |
| Malaria (unspecified) | 1F40 |
| Severe acute malnutrition | 5B53 |
| Moderate acute malnutrition | 5B52 |
| Cholera | 1A00 |
| COVID-19 | RA01.0 |
| Hypertension | BA00 |
| Type 2 diabetes | 5A11 |
| Depressive episode | 6A70 |
| Generalized anxiety disorder | 6B00 |
