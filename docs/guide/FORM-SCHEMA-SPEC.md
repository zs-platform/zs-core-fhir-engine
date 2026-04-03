# ZS Form Schema v1 — Complete Specification

> **Document:** FORM-SCHEMA-SPEC.md | **Version:** 1.0.0
> **Schema URI:** https://zarishsphere.com/schema/form/v1
> **Schema File:** [schemas/zs-form-schema-v1.json](../schemas/zs-form-schema-v1.json)

---

## Overview

ZS Form Schema v1 is the standard for all clinical forms in ZarishSphere. Every form in `zs-content-forms-*` must validate against this schema.

The schema ensures:
- Every field has a FHIR R5 mapping
- Every coded field has a LOINC/SNOMED/ICD-11 code
- Every label is an i18n key (not inline text)
- Forms are machine-processable by `zs-pkg-ui-form-engine`

---

## Complete Form Structure

```json
{
  "$schema": "https://zarishsphere.com/schema/form/v1",
  "id": "zs-form-{domain}-{number}",
  "title": "{{i18n:forms.{domain}.title}}",
  "description": "{{i18n:forms.{domain}.description}}",
  "version": "1.0.0",
  "fhirResource": "Observation",
  "status": "active",
  "tags": ["maternity", "anc", "core"],
  "programs": ["bgd-refugee-response"],
  "sections": [
    {
      "id": "section-{n}",
      "title": "{{i18n:forms.{domain}.section_{n}_title}}",
      "description": "{{i18n:forms.{domain}.section_{n}_description}}",
      "repeating": false,
      "fields": [
        {
          "id": "field-{nnn}",
          "type": "number",
          "label": "{{i18n:forms.{domain}.field_{nnn}_label}}",
          "hint": "{{i18n:forms.{domain}.field_{nnn}_hint}}",
          "placeholder": "{{i18n:forms.{domain}.field_{nnn}_placeholder}}",
          "fhirPath": "Observation.valueQuantity.value",
          "fhirResource": "Observation",
          "loincCode": "29463-7",
          "loincDisplay": "Body weight",
          "unit": "kg",
          "ucumUnit": "kg",
          "required": true,
          "readOnly": false,
          "hidden": false,
          "validation": {
            "min": 0,
            "max": 300,
            "decimalPlaces": 1
          },
          "displayCondition": null
        }
      ]
    }
  ],
  "logic": [],
  "calculatedFields": []
}
```

---

## Field Types

| Type | Description | FHIR Output | Validation |
|------|-------------|-------------|-----------|
| `text` | Short free-text input | `valueString` | maxLength, pattern |
| `textarea` | Multi-line text | `valueString` | maxLength |
| `number` | Numeric input | `valueQuantity` | min, max, decimalPlaces |
| `integer` | Whole number | `valueInteger` | min, max |
| `date` | Date picker | `valueDateTime` | minDate, maxDate |
| `datetime` | Date + time | `valueDateTime` | minDate, maxDate |
| `select` | Single choice (dropdown) | `valueCoding` | options (coded list) |
| `multiselect` | Multiple choices | `valueCodeableConcept[]` | options (coded list) |
| `boolean` | Yes/No toggle | `valueBoolean` | — |
| `scale` | Likert scale | `valueInteger` | min, max, labels |
| `signature` | Digital signature | `Attachment` | — |
| `photo` | Camera capture | `Attachment` | maxSizeMB |
| `gps` | GPS coordinates | `Extension(geolocation)` | — |
| `barcode` | Barcode scanner | `valueString` | pattern |

---

## Select/Multiselect Option Format

Options for coded fields must reference standard terminology:

```json
{
  "type": "select",
  "options": [
    {
      "value": "LA19263-5",
      "display": "{{i18n:forms.options.yes}}",
      "system": "http://loinc.org"
    },
    {
      "value": "LA32-8",
      "display": "{{i18n:forms.options.no}}",
      "system": "http://loinc.org"
    }
  ]
}
```

---

## Display Conditions (Show/Hide Logic)

Fields can be conditionally shown based on other field values:

```json
{
  "id": "field-005",
  "label": "{{i18n:forms.maternity.delivery_complications}}",
  "displayCondition": {
    "fieldId": "field-004",
    "operator": "equals",
    "value": "complicated"
  }
}
```

Operators: `equals`, `notEquals`, `greaterThan`, `lessThan`, `contains`, `notEmpty`

---

## i18n Key Convention

All labels, hints, and option displays must be i18n keys:

```
Format: {{i18n:{namespace}.{key}}}

Namespaces:
  forms.{domain}    ← form-specific keys
  forms.common      ← shared across forms
  forms.options     ← shared option labels (yes/no/unknown)
  units.{unit}      ← unit labels

Examples:
  {{i18n:forms.maternity.anc_contact_number}}
  {{i18n:forms.common.date_of_birth}}
  {{i18n:forms.options.yes}}
  {{i18n:units.kg}}
```

Keys must exist in ALL supported language files before the form can be merged.
