# LOINC Usage Standards

> **Document:** LOINC-USAGE.md | **Version:** 2.0.0
> **Current Version:** LOINC 2.80+

---

## When to Use LOINC

Use LOINC codes in:
- `Observation.code` — ALL clinical observations and measurements
- `DiagnosticReport.code` — lab panels and report types
- `Questionnaire.item.code` — questionnaire item codes

---

## FHIR Coding Pattern

```json
{
  "code": {
    "coding": [
      {
        "system": "http://loinc.org",
        "code": "8480-6",
        "display": "Systolic blood pressure"
      }
    ],
    "text": "Systolic blood pressure"
  },
  "valueQuantity": {
    "value": 120,
    "unit": "mmHg",
    "system": "http://unitsofmeasure.org",
    "code": "mm[Hg]"
  }
}
```

---

## Mandatory LOINC Codes for ZarishSphere Programs

| Observation | LOINC Code | Unit (UCUM) |
|------------|-----------|------------|
| Weight | 29463-7 | kg |
| Height | 8302-2 | cm |
| MUAC (Mid-upper arm circumference) | 56072-2 | cm |
| Systolic blood pressure | 8480-6 | mm[Hg] |
| Diastolic blood pressure | 8462-4 | mm[Hg] |
| Heart rate | 8867-4 | /min |
| Respiratory rate | 9279-1 | /min |
| Body temperature | 8310-5 | Cel |
| Oxygen saturation (SpO2) | 59408-5 | % |
| Blood glucose | 15074-8 | mmol/L |
| Hemoglobin | 718-7 | g/dL |
| PHQ-9 total score | 44261-6 | {score} |
| GAD-7 total score | 69737-5 | {score} |
| Gestational age | 11884-4 | wk |
| Fundal height | 11882-8 | cm |

---

## Panel Codes

When ordering lab panels, use the panel LOINC code at the DiagnosticReport level and individual component codes at the Observation level:

| Panel | Panel LOINC | Components |
|-------|------------|-----------|
| CBC | 58410-2 | Hemoglobin 718-7, Hematocrit 20570-8, WBC 6690-2 |
| Malaria RDT | 51587-4 | Result 91371-5 |
| Blood glucose (fasting) | 1558-6 | Result 1558-6 |
