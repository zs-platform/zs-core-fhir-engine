# FHIR Resources Overview

zs-core-fhir-engine provides complete support for FHIR R5 resources.

## Supported Resources

The library includes all FHIR R5 resources including but not limited to:

| Category | Resources |
|----------|-----------|
| **Administrative** | Patient, Practitioner, Organization, Location |
| **Clinical** | Observation, DiagnosticReport, Encounter, Condition |
| **Medications** | Medication, MedicationRequest, MedicationAdministration |
| **Workflow** | Appointment, Schedule, Task |
| **Financial** | Claim, Coverage, Invoice |

## Bangladesh-Specific Profiles

### BDPatient

Extended Patient profile for Bangladesh:

```go
type BDPatient struct {
    Patient
    // Bangladesh-specific identifiers
    NID    string  // National ID
    BRN    string  // Birth Registration Number
    UHID   string  // Unique Health ID
}
```

### BDAddress

Extended Address profile for Bangladesh administrative divisions:

```go
type BDAddress struct {
    Address
    Division   string  // e.g., "Dhaka", "Chattogram"
    District   string  // e.g., "Dhaka", "Cox's Bazar"
    Upazila    string  // Sub-district
    Union      string  // Union council
    Ward       string  // Ward number
    Village    string  // Village name
}
```

## Rohingya Support

Specialized extensions for refugee camp operations:

### Identifiers

| Type | Description | Use Case |
|------|-------------|----------|
| FCN | Family Counting Number | UNHCR registration |
| Progress ID | Progress Card ID | Camp management |
| MRN | Medical Record Number | Health facility records |

### Location Extensions

| Field | Description |
|-------|-------------|
| Camp | Refugee camp name |
| Block | Camp block number |
| SubBlock | Sub-block identifier |
| ShelterNumber | Shelter assignment |

## Usage Example

```go
import "github.com/zarishsphere/zs-core-fhir-engine/fhir/r5"

// Create a Patient
patient := r5.Patient{
    Active: r5.Bool(true),
    Name: []r5.HumanName{
        {
            Use:    r5.HumanNameUseOfficial,
            Family: "Chowdhury",
            Given:  []string{"Rahima"},
        },
    },
    Gender:    r5.PatientGenderFemale,
    BirthDate: primitives.MustDate("1990-05-15"),
}

// Create a Bangladesh Patient
bdPatient := BDPatient{
    Patient: patient,
    NID:     "1234567890",
    UHID:    "UH-2024-001234",
}
```

## Primitive Types

The library provides type-safe primitive types with validation:

| Type | Description | Example |
|------|-------------|---------|
| `Date` | Partial date support | `"2024"`, `"2024-01"`, `"2024-01-15"` |
| `DateTime` | Date with time | `"2024-01-15T10:30:00Z"` |
| `Time` | Time of day | `"10:30:00"` |
| `Instant` | Timestamp | `"2024-01-15T10:30:00.000Z"` |

## Next Steps

- [Patient Resource](/fhir/patient) - Detailed Patient documentation
- [Bangladesh Profiles](/fhir/profiles) - Bangladesh-specific profiles
