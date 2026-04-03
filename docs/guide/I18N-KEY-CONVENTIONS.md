# Internationalisation (i18n) Key Conventions

> **Document:** I18N-KEY-CONVENTIONS.md | **Version:** 2.0.0

---

## Supported Languages (MVP)

| Code | Language | Script | Status |
|------|----------|--------|--------|
| `en` | English | Latin | ✅ Required |
| `bn` | Bengali / Bangla | Bengali | ✅ Required (Bangladesh, West Bengal) |
| `my` | Burmese | Burmese | 🟡 In Progress |
| `ur` | Urdu | Arabic | 📋 Planned |
| `hi` | Hindi | Devanagari | 📋 Planned |
| `th` | Thai | Thai | 📋 Planned |

All forms must have EN + BN translations before merging. Other languages can be added progressively.

---

## Key Naming Convention

```
{namespace}.{key}

Namespaces:
  forms.{domain}       ← domain-specific form keys
  forms.common         ← common across all forms
  forms.options        ← shared option values (yes/no/unknown)
  units.{unit}         ← unit names
  nav.{section}        ← navigation labels
  errors.{code}        ← error messages
  alerts.{type}        ← alert messages

Examples:
  forms.maternity.anc_contact_number_label
  forms.maternity.anc_contact_number_hint
  forms.common.date_of_birth_label
  forms.options.yes
  forms.options.no
  forms.options.unknown
  units.kg
  units.cm
  errors.required_field
  alerts.critical_low_muac
```

---

## Translation File Format

```json
// translations/en.json
{
  "forms": {
    "maternity": {
      "title": "Antenatal Care — First Contact",
      "section_1_title": "Maternal Information",
      "edd_label": "Expected Date of Delivery",
      "edd_hint": "Enter the date from ultrasound or LMP calculation",
      "gravida_label": "Number of pregnancies (Gravida)",
      "gravida_hint": "Include current pregnancy",
      "lmp_label": "Last Menstrual Period (LMP)",
      "gestational_age_weeks_label": "Gestational Age (weeks)"
    },
    "common": {
      "date_of_birth_label": "Date of Birth",
      "date_of_birth_hint": "DD/MM/YYYY",
      "sex_label": "Sex",
      "weight_label": "Weight",
      "height_label": "Height"
    },
    "options": {
      "yes": "Yes",
      "no": "No",
      "unknown": "Unknown",
      "male": "Male",
      "female": "Female",
      "other": "Other"
    }
  },
  "units": {
    "kg": "kg",
    "g": "g",
    "cm": "cm",
    "mmhg": "mmHg",
    "celsius": "°C",
    "weeks": "weeks",
    "percent": "%"
  }
}
```

---

## Adding a New Translation

1. Add all keys to `en.json` first
2. Open a PR to `zs-data-translations` with the English keys
3. Add `bn.json` translations in the same PR (mandatory)
4. Other language PRs can follow separately
5. CI validates that all keys in `en.json` exist in `bn.json` before merge
