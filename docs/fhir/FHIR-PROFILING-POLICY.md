# FHIR Profiling Policy

> **Document:** FHIR-PROFILING-POLICY.md | **Version:** 2.0.0
> **FHIR Version:** R5 (5.0.0)

---

## When to Create a Profile

Create a ZarishSphere FHIR profile (StructureDefinition) when:
- A FHIR resource needs country-specific constraints
- A FHIR resource needs mandatory ZarishSphere extensions (tenant_id, program_code)
- A particular combination of terminology bindings is required for all uses

Do NOT create a profile just to document how you use a resource — use implementation guides instead.

---

## Profile Authoring Tool

All ZarishSphere profiles are authored using **FHIR Shorthand (FSH)** and compiled with **SUSHI**.

```bash
# Install SUSHI
npm install -g fsh-sushi

# Compile profiles
sushi .
```

---

## Mandatory ZarishSphere Extensions

All ZarishSphere-profiled resources MUST include:

```fsh
// ZarishSphere mandatory extensions
Extension: ZSTenantExtension
Id: zs-tenant
Title: "ZarishSphere Tenant ID"
Description: "The organization/facility tenant identifier"
* value[x] only string

Extension: ZSProgramExtension
Id: zs-program
Title: "ZarishSphere Health Program"
Description: "The health program context (e.g., bgd-refugee-response)"
* value[x] only string
```

---

## Profile Naming Convention

```
ZS{ResourceType}Profile   ← Base ZarishSphere profile
ZS{CC}{ResourceType}Profile  ← Country-specific profile

Examples:
ZSPatientProfile          ← Base Patient profile
ZSBGDPatientProfile       ← Bangladesh-specific Patient profile
ZSCXBPatientProfile       ← Cox's Bazar-specific Patient profile
```

---

## Terminology Binding Strength

Use these binding strengths for coded fields:

| Field Type | Binding Strength |
|-----------|-----------------|
| Diagnosis codes | `required` (ICD-11 or SNOMED) |
| Observation codes | `required` (LOINC) |
| Medication codes | `preferred` (RxNorm) |
| Vaccine codes | `required` (CVX) |
| Units of measure | `required` (UCUM) |
| Status codes | `required` (FHIR ValueSet) |
| Administrative codes | `preferred` |

---

## Profile Validation in CI

All profiles in `zs-data-fhir-profiles` are validated in CI using HAPI FHIR validator:

```yaml
# .github/workflows/validate.yml
- name: Validate FHIR profiles
  run: |
    java -jar validator_cli.jar \
      -version 5.0.0 \
      -ig output/ \
      -recurse \
      resources/**/*.json
```
