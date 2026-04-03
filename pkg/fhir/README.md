# ZarishSphere FHIR Library for Go

A comprehensive Go library for working with FHIR (Fast Healthcare Interoperability Resources) R5 specifications with Bangladesh-specific extensions and profiles.

## Features

- **Complete FHIR R5 Support**: All 158 resources and 44 complex types
- **Type-Safe Primitives**: Custom Date, DateTime, Time, and Instant types with validation
- **Standards Compliant**: Generated directly from official FHIR StructureDefinitions
- **JSON Support**: Full marshaling/unmarshaling with `encoding/json`
- **Bangladesh Profiles**: Localized BDPatient, BDAddress with NID, BRN, UHID support
- **Rohingya Support**: FCN, Progress ID, MRN identifiers and Camp location extensions
- **Validation Framework**: Built-in validation for primitive types with proper error handling

## Installation

```bash
go get github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir
```

## Quick Start

### Creating a FHIR Resource

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5/resources"
    "github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/primitives"
)

func main() {
    // Create a Patient resource
    active := true
    birthDate := primitives.MustDate("1974-12-25")
    
    patient := resources.Patient{
        ID:     stringPtr("example"),
        Active: &active,
        Name: []resources.HumanName{
            {
                Use:    stringPtr("official"),
                Family: stringPtr("Chowdhury"),
                Given:  []string{"Rahima"},
            },
        },
        Gender: stringPtr("female"),
        BirthDate: &birthDate,
    }
    
    // Marshal to JSON
    data, err := json.MarshalIndent(patient, "", "  ")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(data))
}

func stringPtr(s string) *string {
    return &s
}

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5/resources"
)

func main() {
    // Read FHIR JSON file
    data, err := os.ReadFile("patient.json")
    if err != nil {
        panic(err)
    }
    
    // Unmarshal into Patient struct
    var patient resources.Patient
    if err := json.Unmarshal(data, &patient); err != nil {
        panic(err)
    }
    
    // Access fields
    fmt.Printf("Patient ID: %s\n", *patient.ID)
    fmt.Printf("Name: %s %s\n", *patient.Name[0].Given[0], *patient.Name[0].Family)
    fmt.Printf("Birth Date: %s\n", patient.BirthDate.String())
}
```

## Working with Primitive Types

### Date Type

```go
import "github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/primitives"

// Create dates with different precision
yearOnly := primitives.MustDate("2024")                    // Year precision
monthOnly := primitives.MustDate("2024-01")                  // Year-month precision
fullDate := primitives.MustDate("2024-01-15")               // Full date precision

// Convert to time.Time
t, err := fullDate.Time()
if err != nil {
    panic(err)
}
fmt.Printf("As time.Time: %v\n", t)

// Check precision
fmt.Printf("Precision: %s\n", fullDate.Precision()) // "day"
```

### DateTime Type

```go
// Create with timezone
dt := primitives.MustDateTime("2024-01-15T10:30:00+06:00")
fmt.Printf("DateTime: %s\n", dt.String())

// Convert to time.Time
t, err := dt.Time()
if err != nil {
    panic(err)
}
```

### Time Type

```go
// Create time values
time1 := primitives.MustTime("10:30:00")
time2 := primitives.MustTime("10:30:00.123")  // With fractional seconds

// Calculate duration
duration, _ := time1.Duration(time2.Time())
fmt.Printf("Duration: %v\n", duration) // 0.123s
```

### Instant Type

```go
// Create instant (always includes timezone)
instant := primitives.MustInstant("2024-01-15T10:30:00Z")
fmt.Printf("Instant: %s\n", instant.String())

// Convert to time.Time
t, err := instant.Time()
if err != nil {
    panic(err)
}
```

## Bangladesh-Specific Features

### Using Bangladesh Profiles

```go
import (
    "github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5/resources"
    "github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/profiles/bd"
)

func createBangladeshPatient() *resources.Patient {
    return &resources.Patient{
        ID: stringPtr("bd-patient-001"),
        Extension: []resources.Extension{
            {
                Url: stringPtr("http://zarishsphere.org/fhir/StructureDefinition/bangladesh-national-id"),
                Value: &bd.BangladeshNationalId{
                    Value: stringPtr("1234567890123"),
                },
            },
            {
                Url: stringPtr("http://zarishsphere.org/fhir/StructureDefinition/birth-registration-number"),
                Value: &bd.BirthRegistrationNumber{
                    Value: stringPtr("BRN20234567890"),
                },
            },
        },
        Name: []resources.HumanName{
            {
                Family: stringPtr("Islam"),
                Given:  []string{"Mohammad"},
            },
        },
    }
}
```

### Rohingya Support

```go
func createRohingyaPatient() *resources.Patient {
    return &resources.Patient{
        ID: stringPtr("rohingya-patient-001"),
        Extension: []resources.Extension{
            {
                Url: stringPtr("http://zarishsphere.org/fhir/StructureDefinition/family-counting-number"),
                Value: &bd.FamilyCountingNumber{
                    Value: stringPtr("FCN123456"),
                },
            },
            {
                Url: stringPtr("http://zarishsphere.org/fhir/StructureDefinition/progress-id"),
                Value: &bd.ProgressId{
                    Value: stringPtr("PID2023456789"),
                },
            },
            {
                Url: stringPtr("http://zarishsphere.org/fhir/StructureDefinition/camp-location"),
                Value: &bd.CampLocation{
                    Value: stringPtr("Camp-7-Kutupalong"),
                },
            },
        },
    }
}
```

## Validation

### Primitive Type Validation

```go
// Valid dates
validDates := []string{
    "2024",           // Year precision - OK
    "2024-01",         // Year-month - OK
    "2024-01-15",      // Full date - OK
}

for _, dateStr := range validDates {
    date, err := primitives.NewDate(dateStr)
    if err != nil {
        fmt.Printf("Error creating date '%s': %v\n", dateStr, err)
    } else {
        fmt.Printf("Valid date '%s': %s (precision: %s)\n", 
            dateStr, date.String(), date.Precision())
    }
}

// Invalid dates
invalidDates := []string{
    "2024-13-15",      // Invalid month
    "invalid",           // Invalid format
    "2024-02-30",      // Invalid day for February
}

for _, dateStr := range invalidDates {
    _, err := primitives.NewDate(dateStr)
    if err != nil {
        fmt.Printf("Expected error for '%s': %v\n", dateStr, err)
    }
}
```

## Resource Examples

### Creating an Observation

```go
func createVitalsObservation() *resources.Observation {
    systolic := 120.0
    diastolic := 80.0
    
    return &resources.Observation{
        ID:     stringPtr("vitals-001"),
        Status: stringPtr("final"),
        Code: &resources.CodeableConcept{
            Coding: []resources.Coding{
                {
                    System: stringPtr("http://loinc.org"),
                    Code:    stringPtr("85354-9"),
                    Display: stringPtr("Blood pressure panel"),
                },
            },
        },
        Subject: &resources.Reference{
            Reference: stringPtr("Patient/example-patient"),
        },
        Component: []resources.ObservationComponent{
            {
                Code: &resources.CodeableConcept{
                    Coding: []resources.Coding{
                        {
                            System: stringPtr("http://loinc.org"),
                            Code:    stringPtr("8480-6"),
                            Display: stringPtr("Systolic blood pressure"),
                        },
                    },
                },
                ValueQuantity: &resources.Quantity{
                    Value:  &systolic,
                    Unit:   stringPtr("mm[Hg]"),
                    System: stringPtr("http://unitsofmeasure.org"),
                    Code:   stringPtr("mm[Hg]"),
                },
            },
            {
                Code: &resources.CodeableConcept{
                    Coding: []resources.Coding{
                        {
                            System: stringPtr("http://loinc.org"),
                            Code:    stringPtr("8462-4"),
                            Display: stringPtr("Diastolic blood pressure"),
                        },
                    },
                },
                ValueQuantity: &resources.Quantity{
                    Value:  &diastolic,
                    Unit:   stringPtr("mm[Hg]"),
                    System: stringPtr("http://unitsofmeasure.org"),
                    Code:   stringPtr("mm[Hg]"),
                },
            },
        },
    }
}
```

### Creating a Bundle

```go
func createSearchSetBundle() *resources.Bundle {
    return &resources.Bundle{
        ID:     stringPtr("search-bundle-001"),
        Type:    stringPtr("searchset"),
        Total:   uintPtr(2),
        Entry: []resources.BundleEntry{
            {
                FullUrl: stringPtr("https://example.com/Patient/1"),
                Resource: &patient1, // Assume patient1 is defined
            },
            {
                FullUrl: stringPtr("https://example.com/Patient/2"),
                Resource: &patient2, // Assume patient2 is defined
            },
        },
    }
}
```

## Testing

The library includes comprehensive test coverage:

```bash
# Test primitive types
go test ./pkg/fhir/primitives/...

# Test FHIR resources
go test ./pkg/fhir/r5/...

# Test validation
go test ./pkg/fhir/validation/...

# Run all tests
go test ./pkg/fhir/...
```

## Performance Considerations

- **Zero Allocations**: Primitive types use value types where possible
- **Lazy Parsing**: Time conversions only happen when needed
- **Memory Efficiency**: Pointer types for optional fields save memory
- **JSON Efficiency**: Uses standard `encoding/json` with no reflection overhead

## Comparison with Other Libraries

### vs google/fhir

- **go-zs-core-fhir/fhir**: Custom primitives with validation, Bangladesh profiles
- **google/fhir**: Protocol buffers, more complex setup

### vs samply/golang-fhir

- **go-zs-core-fhir/fhir**: FHIR R5 support, generated from specs
- **samply/golang-fhir**: FHIR R4 only, hand-written types

## Contributing

Contributions are welcome! Please:

1. Run tests: `go test ./...`
2. Format code: `go fmt ./...`
3. Check linting: `golangci-lint run`
4. Update documentation

## License

MIT License - see LICENSE file for details.

## Resources

- [FHIR R5 Specification](https://hl7.org/fhir/R5/)
- [Bangladesh DGHS](https://dghs.gov.bd/)
- [FHIR Implementation Guides](https://build.fhir.org/ig/)
- [ZarishSphere Documentation](../../docs/)
