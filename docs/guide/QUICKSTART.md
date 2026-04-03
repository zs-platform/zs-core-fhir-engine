# Quick Start

> **Version:** 2.0.0 | **Standards:** FHIR R5 (ADR-0002), Go 1.26.1 (ADR-0001)

This guide will help you get started with ZarishSphere FHIR Engine in just a few minutes.

---

## 🚀 Start the FHIR Server

The fastest way to get started is to build and run the FHIR server:

```bash
# Build the server
go build -o fhir-engine ./cmd/fhir-engine

# Run with default settings
./fhir-engine serve --port 8080

# Or run with full features
./fhir-engine serve --debug --port 8080
```

You should see output like:

```
2026/04/03 10:00:00 INFO ZarishSphere FHIR Engine v2.0.0
2026/04/03 10:00:00 INFO Loading IG data from ./config...
2026/04/03 10:00:00 INFO Loaded 10 CodeSystems and 25 ValueSets
2026/04/03 10:00:00 INFO FHIR R5 Server starting on :8080...
2026/04/03 10:00:00 INFO SMART on FHIR 2.1 authentication enabled
2026/04/03 10:00:00 INFO NATS JetStream events enabled (ADR-0004)
2026/04/03 10:00:00 INFO Server ready - All endpoints available
```

## 🧪 Test the Server

### 1. Check Server Health

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"ok","timestamp":"2026-04-03T10:00:01Z"}
```

### 2. Get FHIR Server Metadata (CapabilityStatement)

```bash
curl http://localhost:8080/fhir/R5/metadata
```

### 3. Create a Patient (with Bangladesh identifiers)

```bash
curl -X POST http://localhost:8080/fhir/R5/Patient \
  -H "Content-Type: application/fhir+json" \
  -H "X-Tenant-ID: demo-hospital" \
  -d '{
    "resourceType": "Patient",
    "active": true,
    "name": [{
      "family": "Chowdhury",
      "given": ["Rahima"]
    }],
    "gender": "female",
    "birthDate": "1990-05-15",
    "identifier": [{
      "system": "http://bangladesh.gov/nid",
      "value": "1234567890123"
    }]
  }'
```

### 4. Read the Patient

```bash
# Replace {id} with the ID from the previous response
curl http://localhost:8080/fhir/R5/Patient/{id} \
  -H "X-Tenant-ID: demo-hospital"
```

### 5. Search for Patients

```bash
curl http://localhost:8080/fhir/Patient
```

### 6. Update a Patient

```bash
curl -X PUT http://localhost:8080/fhir/Patient/{id} \
  -H "Content-Type: application/fhir+json" \
  -d '{
    "resourceType": "Patient",
    "id": "{id}",
    "active": true,
    "name": [{
      "family": "Chowdhury",
      "given": ["Rahima", "Begum"]
    }],
    "gender": "female",
    "birthDate": "1990-05-15"
  }'
```

### 7. Delete a Patient

```bash
curl -X DELETE http://localhost:8080/fhir/Patient/{id}
```

## 🏥 Use the Terminology Server

Start the standalone terminology server:

```bash
./zs-core-fhir-engine terminology --port 8081
```

### Expand a ValueSet

```bash
# Get all ICD-11 concepts
curl "http://localhost:8081/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms"

# Filter results
curl "http://localhost:8081/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms&filter=hyper"
```

## 📊 Monitor System Health

### Check Metrics

```bash
curl http://localhost:8080/metrics
```

### View System Status

```bash
curl http://localhost:8080/health
```

## 🛠 Use the FHIR Library

### Create a Patient Resource

```go
package main

import (
    "encoding/json"
    "fmt"
    
    "github.com/zarishsphere/zs-core-fhir-engine/fhir/r5"
    "github.com/zarishsphere/zs-core-fhir-engine/fhir/primitives"
)

func main() {
    active := true
    birthDate := primitives.MustDate("1990-05-15")
    
    patient := r5.Patient{
        Active: &active,
        Name: []r5.HumanName{
            {
                Use:    r5.HumanNameUseOfficial,
                Family: "Chowdhury",
                Given:  []string{"Rahima"},
            },
        },
        Gender:    r5.PatientGenderFemale,
        BirthDate: &birthDate,
    }
    
    // Marshal to JSON
    data, _ := json.MarshalIndent(patient, "", "  ")
    fmt.Println(string(data))
}
```

## 🚀 Production Deployment

For production use, use our automated deployment script:

```bash
./scripts/production-deploy.sh
```

This script:
- ✅ Checks dependencies (Go, PostgreSQL)
- ✅ Builds the application
- ✅ Sets up production configuration
- ✅ Starts the server with health checks
- ✅ Provides monitoring endpoints

## 🐳 Docker Quick Start

```bash
# Build Docker image
docker build -t zs-core-fhir-engine .

# Run FHIR server
docker run -p 8080:8080 zs-core-fhir-engine serve --port 8080

# Run with PostgreSQL backend
docker run -p 8080:8080 \
  -e FHIR_DB_MODE=postgresql \
  -e FHIR_DB_HOST=postgres \
  zs-core-fhir-engine serve --port 8080
```

## 🧪 Validate FHIR Resources

```bash
# Validate a patient resource
./zs-core-fhir-engine validate examples/patient.json

# Validate with detailed output
./zs-core-fhir-engine validate --verbose examples/patient.json
```

## 📚 Next Steps

### For Healthcare Professionals
- [FHIR Overview](/fhir/overview) - Understanding FHIR resources
- [Patient Resources](/fhir/patient) - Patient data management
- [Bangladesh Profiles](/fhir/profiles) - Local healthcare standards

### For Developers
- [API Reference](/api/overview) - Complete REST API documentation
- [Installation Guide](/guide/installation) - Detailed setup instructions
- [Development Guide](/guide/introduction) - Contributing to the project

### For System Administrators
- [Production Deployment](/guide/complete-server-blueprint) - Production setup
- [Configuration](/guide/configuration) - System configuration
- [Monitoring](/guide/monitoring) - Observability and metrics

## 🔧 Troubleshooting

### Server Won't Start
```bash
# Check if port is available
netstat -tlnp | grep :8080

# Try a different port
./zs-core-fhir-engine serve --port 8081
```

### Database Connection Issues
```bash
# Use in-memory mode for testing
./zs-core-fhir-engine serve --port 8080

# Check PostgreSQL connection
psql -h localhost -U fhir_user -d fhir_db
```

### Validation Errors
```bash
# Validate with detailed output
./zs-core-fhir-engine validate --verbose your-resource.json

# Check FHIR specification
curl https://hl7.org/fhir/r5/patient.html
```

---

**Need help?** Check our [API Reference](/api/overview) or open an issue on GitHub.
