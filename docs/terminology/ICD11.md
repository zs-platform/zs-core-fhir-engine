# ICD-11 Terminology

zs-core-fhir-engine includes support for WHO ICD-11 (International Classification of Diseases, 11th Revision) codes.

## System URL

```
http://id.who.int/icd/release/11/mms
```

## Usage

```bash
# Get all ICD-11 concepts
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms"

# Filter by text
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms&filter=diabetes"
```

## Included Concepts

The terminology server includes the following ICD-11 concepts:

| Code | Display |
|------|---------|
| BA00 | Essential hypertension |
| BA01 | Malignant essential hypertension |
| 1B10 | Tuberculosis of the lung |

## Adding Custom ICD-11 Codes

You can add custom ICD-11 concepts to the terminology server by modifying the `cmd/zs-core-fhir/terminology.go` file:

```go
func StartTerminologyServer(port int) {
    server := NewTerminologyServer()
    
    // Add ICD-11 concepts
    server.AddConcept("http://id.who.int/icd/release/11/mms", "BA00", "Essential hypertension")
    server.AddConcept("http://id.who.int/icd/release/11/mms", "1B10", "Tuberculosis of the lung")
    
    // ... rest of the code
}
```

## ICD-11 Structure

ICD-11 uses a chapter-based structure:

| Chapter | Block | Description |
|---------|-------|-------------|
| I | 1A00-1B9Z | Certain infectious and parasitic diseases |
| II | 2A00-2D90 | Neoplasms |
| ... | ... | ... |

## More Information

- [ICD-11 WHO](https://icd.who.int/browse11)
- [ICD-11 MMS](https://icd.who.int/browse11/i-en)
