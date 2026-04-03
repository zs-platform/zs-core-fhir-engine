# Bangladesh Profiles

zs-core-fhir-engine includes specialized FHIR profiles for Bangladesh healthcare requirements.

## BDPatient Profile

Extended Patient profile with Bangladesh-specific identifiers.

### Identifiers

| Identifier Type | Extension URL | Description |
|-----------------|---------------|-------------|
| NID | `http://zdhhs.gov.bd/fhir/StructureDefinition/NID` | National Identity Card Number |
| BRN | `http://zdhhs.gov.bd/fhir/StructureDefinition/BRN` | Birth Registration Number |
| UHID | `http://zdhhs.gov.bd/fhir/StructureDefinition/UHID` | Unique Health ID |

### Example

```json
{
  "resourceType": "Patient",
  "id": "bd-patient-001",
  "identifier": [
    {
      "type": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v2-0203",
            "code": "NID"
          }
        ]
      },
      "system": "http://zdhhs.gov.bd/identifier/nid",
      "value": "1234567890"
    },
    {
      "type": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v2-0203",
            "code": "UHID"
          }
        ]
      },
      "system": "http://zdhhs.gov.bd/identifier/uhid",
      "value": "UH-2024-001234"
    }
  ],
  "name": [
    {
      "use": "official",
      "family": "Chowdhury",
      "given": ["Rahima"]
    }
  ]
}
```

## BDAddress Profile

Extended Address profile for Bangladesh administrative divisions.

### Fields

| Field | Description | Example |
|-------|-------------|---------|
| Division | Division | "Dhaka", "Chattogram", "Sylhet" |
| District | District | "Dhaka", "Cox's Bazar" |
| Upazila | Upazila/Sub-district | "Savar", "Teknaf" |
| Union | Union Council | "Savar" |
| Ward | Ward Number | "01" |
| Village | Village Name | "Kashem Market" |

### Example

```json
{
  "use": "home",
  "type": "physical",
  "line": ["House #123", "Road #5"],
  "district": "Dhaka",
  "state": "Dhaka",
  "postalCode": "1205",
  "country": "Bangladesh"
}
```

## Rohingya Extensions

Specialized extensions for refugee camp operations.

### Identifiers

| Identifier | Description | Example |
|------------|-------------|---------|
| FCN | Family Counting Number | FCN-2023-001234 |
| ProgressID | Progress Card ID | P-12345 |
| MRN | Medical Record Number | MRN-CX-2024-0001 |

### Location Extensions

| Extension | Description |
|-----------|-------------|
| Camp | Camp name (e.g., "Kutupalong", "Balukhali") |
| Block | Block identifier (e.g., "B1", "C2") |
| SubBlock | Sub-block (e.g., "B1-S3") |
| ShelterNumber | Shelter number (e.g., "S-1234") |

### Example

```json
{
  "resourceType": "Patient",
  "id": "rohingya-patient-001",
  "identifier": [
    {
      "system": "http://unhcr.org/fcn",
      "value": "FCN-2023-001234"
    },
    {
      "system": "http://health.cxb.gov.bd/mrn",
      "value": "MRN-CX-2024-0001"
    }
  ],
  "address": [
    {
      "extension": [
        {
          "url": "http://zdhhs.gov.bd/fhir/StructureDefinition/CampName",
          "valueString": "Kutupalong"
        },
        {
          "url": "http://zdhhs.gov.bd/fhir/StructureDefinition/CampBlock",
          "valueString": "B1"
        }
      ],
      "country": "Bangladesh"
    }
  ]
}
```

## Bangladesh ValueSets

### Administrative Divisions

| Level | Example Values |
|-------|----------------|
| Division | Dhaka, Chattogram, Sylhet, Khulna, Rajshahi, Rangpur, Barisal, Mymensingh |
| District | Dhaka, Gazipur, Narayanganj, Cox's Bazar, Bandarban |
| Upazila | Savar, Keraniganj, Teknaf, Ukhiya |

### Facilities

| Type | Example |
|------|---------|
| Hospital | Dhaka Medical College Hospital |
| Clinic | Kutupalong Refugee Health Centre |
| Diagnostic Lab | National Institute of Cancer Research |
