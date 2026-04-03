# SNOMED CT Usage Standards

> **Document:** SNOMED-USAGE.md | **Version:** 2.0.0

---

## When to Use SNOMED CT

Use SNOMED CT codes in:
- `Condition.code` — clinical findings and diagnoses (alongside ICD-11)
- `Procedure.code` — clinical procedures
- `AllergyIntolerance.code` — allergen substance codes
- `CDS Hooks` — clinical decision support triggers

---

## FHIR Coding Pattern

```json
{
  "code": {
    "coding": [
      {
        "system": "http://snomed.info/sct",
        "code": "271737000",
        "display": "Anaemia (disorder)"
      },
      {
        "system": "http://id.who.int/icd/release/11/2026-01/mms",
        "code": "3A00",
        "display": "Iron deficiency anaemia"
      }
    ]
  }
}
```

---

## Common SNOMED Codes for ZarishSphere

| Clinical Concept | SNOMED Code |
|----------------|-------------|
| Malnutrition | 248325000 |
| Severe acute malnutrition | 238131007 |
| Moderate acute malnutrition | 302872004 |
| Antenatal care | 424525001 |
| Postnatal care | 133906008 |
| Immunization (procedure) | 127785005 |
| Blood pressure taking | 75367002 |
| Weighing patient | 39857003 |
| Penicillin (allergy) | 372687004 |
| Sulfa drug (allergy) | 387406002 |
