# Clinical Form Validation Rules

> **Version:** 2.0.0 | **Updated:** 2026-03-24
> **Enforced by:** zs-agent-content-validator (CI on every PR)
> **Schema:** [zs-form-schema-v1.json](../schemas/zs-form-schema-v1.json)

---

## Overview

Every clinical form in `zs-content-forms-*` is automatically validated by the CI pipeline before it can be merged. This document describes all validation rules and how to fix common failures.

---

## Validation Rules

### Rule V-01: Valid JSON Syntax

**Rule:** The form file must be valid JSON.
**CI check:** `python3 -m json.tool form.json`
**Fix:** Use a JSON validator at jsonlint.com or VS Code's JSON support.

---

### Rule V-02: Validates Against ZS Form Schema v1

**Rule:** The form must validate against `schemas/zs-form-schema-v1.json`.
**CI check:** `jsonschema -i form.json schemas/zs-form-schema-v1.json`
**Required fields at root:**
- `$schema` — must be `https://zarishsphere.com/schema/form/v1`
- `id` — must follow pattern `zs-form-{domain}-{number}`
- `title` — must be an i18n key pattern `{{i18n:...}}`
- `version` — semantic version string
- `fhirResource` — valid FHIR R5 resource type
- `sections` — array with at least 1 section

---

### Rule V-03: Every Clinical Field Has a FHIR Mapping

**Rule:** Every field of type `text`, `number`, `date`, `select`, `multiselect`, `boolean`, `textarea` must have a `fhirPath` property.
**CI check:** Python script parses all fields, checks for `fhirPath`.

```json
// ✅ Valid
{
  "id": "field-001",
  "type": "number",
  "label": "{{i18n:forms.vitals.weight_label}}",
  "fhirPath": "Observation.valueQuantity.value",
  "loincCode": "29463-7"
}

// ❌ Invalid — missing fhirPath
{
  "id": "field-001",
  "type": "number",
  "label": "{{i18n:forms.vitals.weight_label}}"
}
```

---

### Rule V-04: Coded Fields Must Have a Terminology Code

**Rule:** Fields with type `select` or `multiselect`, and all measurement fields, must have at least one of: `loincCode`, `snomedCode`, `icd11Code`, `cvxCode`.

```json
// ✅ Valid
{ "type": "number", "loincCode": "8310-5" }

// ❌ Invalid — number field with no LOINC code
{ "type": "number", "fhirPath": "Observation.valueQuantity.value" }
```

---

### Rule V-05: All Labels Are i18n Keys

**Rule:** `label`, `placeholder`, `hint` fields must use the pattern `{{i18n:forms.{domain}.{key}}}`. No inline English text allowed.

```json
// ✅ Valid
{ "label": "{{i18n:forms.vitals.weight_label}}" }

// ❌ Invalid — inline text
{ "label": "Patient Weight (kg)" }
```

---

### Rule V-06: Translation Keys Exist

**Rule:** Every i18n key used in the form must exist in `translations/en.json` AND `translations/bn.json` (Bengali — Bangladesh is the primary target country).
**CI check:** Extracts all `{{i18n:...}}` keys, checks against translation files.

---

### Rule V-07: Form ID Follows Naming Convention

**Rule:** Form ID must match: `^zs-form-[a-z][a-z0-9-]+-[0-9]{2}$`

```
✅ zs-form-anc-01
✅ zs-form-phq9-01
❌ zs_form_anc_01  (underscores)
❌ ZS-FORM-ANC-01  (uppercase)
❌ zs-form-anc     (missing number)
```

---

### Rule V-08: Version is Semantic

**Rule:** `version` must match semver: `^[0-9]+\.[0-9]+\.[0-9]+$`

---

### Rule V-09: No PHI as Default Values

**Rule:** No field's `default` property may contain real patient data.
**CI check:** Regex scan for patterns matching: name formats, dates of birth, ID numbers.

---

### Rule V-10: fhirPath is Valid R5 Path

**Rule:** The `fhirPath` value must be a valid FHIR R5 path for the resource declared in `fhirResource`.
**CI check:** Validated against FHIR R5 resource definitions.

---

## How to Run Validation Locally

```bash
cd zs-content-forms-core  # or any forms repo

# Install validator
pip3 install jsonschema

# Validate a single form
python3 tests/validate_forms.py forms/vitals/vitals-entry-form.json

# Validate all forms
python3 tests/validate_forms.py

# Expected output:
# ✅ forms/vitals/vitals-entry-form.json — VALID
# ❌ forms/lab/lab-order.json — INVALID: field-003 missing fhirPath
```

---

## Common Validation Errors and Fixes

| Error | Cause | Fix |
|-------|-------|-----|
| `fhirPath missing on field-XXX` | Field has no FHIR mapping | Add `"fhirPath": "Resource.path"` |
| `i18n key not in bn.json` | Bengali translation missing | Add key to `translations/bn.json` |
| `Invalid form ID format` | Wrong naming | Use `zs-form-{domain}-{NN}` |
| `loincCode required for measurement` | Missing LOINC | Look up code at loinc.org |
| `Inline label text` | Using text not i18n key | Replace with `{{i18n:forms.x.y}}` |
