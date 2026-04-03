# Terminology Service Overview

The zs-core-fhir-engine terminology service provides FHIR terminology operations including ValueSet expansion.

## Features

- **ValueSet/$expand**: Expand ValueSets to get all concepts
- **ICD-11 Support**: WHO ICD-11 codes
- **Local Codes**: Bangladesh administrative divisions
- **Filter Support**: Filter concepts by text search

## Starting the Terminology Server

```bash
# Standalone terminology server
./zs-core-fhir --term-server --port 8080

# Or use the FHIR server which includes terminology
./zs-core-fhir --server --port 8080
```

## ValueSet/$expand Operation

The `$expand` operation returns all concepts in a ValueSet.

### Endpoint

```
GET /fhir/ValueSet/$expand?url={system}[&filter={text}]
```

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `url` | string | **Required**. The ValueSet system URL |
| `filter` | string | Optional text filter |

### Examples

```bash
# Get all ICD-11 codes
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms"

# Filter ICD-11 codes
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms&filter=diabetes"

# Get Bangladesh divisions
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=https://health.zarishsphere.com/fhir/ValueSet/bd-divisions"
```

## Response Format

```json
[
  {
    "code": "BA00",
    "display": "Essential hypertension",
    "system": "http://id.who.int/icd/release/11/mms"
  },
  {
    "code": "BA01",
    "display": "Malignant essential hypertension",
    "system": "http://id.who.int/icd/release/11/mms"
  }
]
```

## Built-in ValueSets

| System | Description |
|--------|-------------|
| `http://id.who.int/icd/release/11/mms` | ICD-11 (WHO) |
| `https://health.zarishsphere.com/fhir/ValueSet/bd-divisions` | Bangladesh Divisions |
| `https://health.zarishsphere.com/fhir/ValueSet/bd-districts` | Bangladesh Districts |

## Next Steps

- [ICD-11 Reference](/terminology/icd11) - ICD-11 code details
- [Bangladesh Divisions](/terminology/bangladesh) - Administrative divisions
