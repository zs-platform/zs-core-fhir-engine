# Terminology Governance Policy

> **Document:** TERMINOLOGY-GOVERNANCE.md | **Version:** 2.0.0

---

## Overview

Coded clinical data is only interoperable when everyone uses the same codes with the same meaning. This policy defines which terminology systems ZarishSphere uses for each clinical domain and how they are maintained.

---

## Approved Terminology Systems

| Domain | Primary System | Secondary | Source Repository |
|--------|---------------|-----------|------------------|
| Diagnoses | ICD-11 (2026-01) | SNOMED CT | `zs-data-icd11` |
| Laboratory observations | LOINC 2.80+ | — | `zs-data-loinc` |
| Clinical findings / symptoms | SNOMED CT | — | `zs-data-snomed` |
| Vaccines | CVX (CDC) | — | `zs-data-cvx` |
| Medications | RxNorm | Local formulary | `zs-data-rxnorm` |
| OpenMRS concepts | CIEL (OpenMRS) | — | `zs-data-ciel` |
| Units of measure | UCUM | — | Built-in FHIR |
| Administrative | HL7 FHIR value sets | — | Built-in FHIR |

---

## Coding Requirements by Clinical Domain

### Observations (Vitals and Labs)
- **ALL** `Observation` resources must have `Observation.code` with a LOINC code
- LOINC panel codes used for ordered labs
- LOINC component codes for individual measurements
- Example: Blood pressure panel LOINC 55284-4; systolic component 8480-6

### Diagnoses
- All `Condition` resources must have a code from ICD-11 or SNOMED CT
- ICD-11 preferred for administrative/reporting purposes
- SNOMED CT preferred for clinical decision support
- If both are known: include both as coding array entries

### Medications
- All `MedicationRequest` resources must reference an RxNorm code
- Country-specific drug codes may be included as secondary coding
- Generic names always included alongside brand names

### Vaccines
- All `Immunization` resources must have a CVX code
- Country EPI schedule codes may be included as secondary

---

## Terminology Update Policy

| System | Update Frequency | Process |
|--------|-----------------|---------|
| ICD-11 | Annual (January) | `zs-agent-dependency-updater` opens PR with new release |
| LOINC | Bi-annual (Feb/Aug) | Same automated PR process |
| SNOMED CT | Bi-annual (Jan/Jul) | Same automated PR process |
| CVX | As needed | Manual PR when CDC releases update |
| RxNorm | Continuous (API-based) | NLM API called at runtime; local cache refreshed weekly |
| CIEL | Irregular | Manual PR when OpenMRS releases update |

---

## Local Code Extensions

Countries may add local code systems for concepts not covered by standard terminologies:

```json
{
  "system": "https://zarishsphere.com/CodeSystem/bgd-camp-codes",
  "code": "camp-1w",
  "display": "Camp 1 West (Kutupalong)"
}
```

Local code systems must:
1. Use the `https://zarishsphere.com/CodeSystem/{cc}-{name}` URI pattern
2. Be defined as a FHIR `CodeSystem` resource in `zs-data-fhir-profiles`
3. Include a complete concept list with displays in at least EN and the local language
