# Patient Resource

The Patient resource represents demographic and administrative information about an individual receiving healthcare services.

## Basic Usage

```go
import (
    "github.com/zarishsphere/zs-core-fhir-engine/fhir/r5"
    "github.com/zarishsphere/zs-core-fhir-engine/fhir/primitives"
)

// Create a simple patient
active := true
patient := r5.Patient{
    ID:     r5.Id("example"),
    Active: &active,
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
```

## Patient Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `resourceType` | string | Must be "Patient" |

### Common Fields

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `*Id` | Logical ID |
| `Active` | `*bool` | Whether patient is active |
| `Name` | `[]HumanName` | Patient names |
| `Gender` | `PatientGender` | Gender identity |
| `BirthDate` | `*Date` | Date of birth |
| `Address` | `[]Address` | Addresses |
| `Telecom` | `[]ContactPoint` | Contact details |

## HumanName

```go
name := r5.HumanName{
    Use:     r5.HumanNameUseOfficial,
    Family:  "Chowdhury",
    Given:   []string{"Rahima", "Begum"},
    Prefix:  []string{"Mrs"},
    Suffix:  []string{"PhD"},
}
```

## Gender Values

| Code | Description |
|------|-------------|
| `male` | Male |
| `female` | Female |
| `other` | Other |
| `unknown` | Unknown |

## Address (Bangladesh)

```go
address := r5.Address{
    Use: r5.AddressUseHome,
    Line: []string{"123", "Main Road"},
    City: "Dhaka",
    District: "Dhaka",
    State: "Dhaka",
    PostalCode: "1205",
    Country: "Bangladesh",
}
```

## ContactPoint

```go
telecom := []r5.ContactPoint{
    {
        System: r5.ContactPointSystemPhone,
        Value:  "+880-2-12345678",
        Use:    r5.ContactPointUseHome,
    },
    {
        System: r5.ContactPointSystemPhone,
        Value:  "+8801712345678",
        Use:    r5.ContactPointUseMobile,
    },
    {
        System: r5.ContactPointSystemEmail,
        Value:  "rahima@example.com",
    },
}
```

## JSON Representation

```json
{
  "resourceType": "Patient",
  "id": "example",
  "active": true,
  "name": [
    {
      "use": "official",
      "family": "Chowdhury",
      "given": ["Rahima"]
    }
  ],
  "gender": "female",
  "birthDate": "1990-05-15",
  "address": [
    {
      "use": "home",
      "city": "Dhaka",
      "district": "Dhaka",
      "state": "Dhaka",
      "country": "Bangladesh"
    }
  ]
}
```
