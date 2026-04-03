# ZarishSphere FHIR Engine CLI

A comprehensive FHIR R5 implementation with Bangladesh and Rohingya localization, terminology services, and a full RESTful FHIR server.

## Features

### Core Library (`github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir`)

- **Complete FHIR R5 Support**: All resources and types
- **Type-Safe Primitives**: Custom Date, DateTime, Time, and Instant types with validation
- **Standards Compliant**: Generated from official FHIR StructureDefinitions
- **JSON Support**: Full marshaling/unmarshaling with `encoding/json`
- **Bangladesh Profiles**: Localized BDPatient, BDAddress with NID, BRN, UHID support
- **Rohingya Support**: FCN, Progress ID, MRN identifiers and Camp location extensions

### CLI Tool (`cmd/fhir-engine`)

A command-line tool for FHIR operations:

| Command | Description |
|---------|-------------|
| `serve` | Start full FHIR REST server |
| `terminology` | Start terminology server |

### FHIR Server

- **RESTful API**: Full CRUD operations (Create, Read, Update, Delete)
- **Search**: Basic search capabilities
- **Terminology Service**: ValueSet/$expand operation
- **ICD-11 Support**: WHO ICD-11 codes included
- **Bangladesh Geography**: Division, District, Upazila codes

## Installation

```bash
# Install from source
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine

# Build
go build -o fhir-engine ./cmd/fhir-engine

# Or install
go install ./cmd/fhir-engine
```

## Quick Start

### Start FHIR Server

```bash
# Start full FHIR server on port 8080
./fhir-engine serve --port 8080

# Start terminology server only
./fhir-engine terminology --port 8080
```

### Start with Custom IG

```bash
./fhir-engine serve --port 8080 --ig ./config
```

## FHIR Server Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/fhir` | Server metadata (CapabilityStatement) |
| POST | `/fhir/{resourceType}` | Create resource |
| GET | `/fhir/{resourceType}` | Search resources |
| GET | `/fhir/{resourceType}/{id}` | Read resource |
| PUT | `/fhir/{resourceType}/{id}` | Update resource |
| DELETE | `/fhir/{resourceType}/{id}` | Delete resource |
| GET | `/fhir/ValueSet/$expand?url={system}` | Expand ValueSet |

## Example Usage

### Create a Patient

```bash
curl -X POST http://localhost:8080/fhir/Patient \
  -H "Content-Type: application/fhir+json" \
  -d '{
    "resourceType": "Patient",
    "id": "example",
    "active": true,
    "name": [{
      "family": "Mia",
      "given": ["Chowdhury"]
    }],
    "gender": "female"
  }'
```

### Search for Patients

```bash
curl http://localhost:8080/fhir/Patient
```

### Expand a ValueSet

```bash
curl "http://localhost:8080/fhir/ValueSet/$expand?url=http://id.who.int/icd/release/11/mms"
```

## Configuration

| Flag | Default | Description |
|-------|---------|-------------|
| `serve` | false | Start full FHIR server |
| `terminology` | false | Start terminology server only |
| `port` | 8080 | Server port |
| `ig` | `./config` | Path to FHIR Implementation Guide |

## Bangladesh-Specific Features

### Supported Identifiers

- **NID**: National ID
- **BRN**: Birth Registration Number
- **UHID**: Unique Health ID
- **FCN**: Family Counting Number (Rohingya)
- **Progress ID**: Refugee Progress ID
- **MRN**: Medical Record Number

### Administrative Divisions

- Division (e.g., Dhaka, Chattogram)
- District
- Upazila
- Union
- Ward
- Village

## Documentation

- [FHIR Library Documentation](../../pkg/fhir/README.md)
- [API Reference](../../docs/api/OVERVIEW.md)
- [Terminology Server](../../docs/terminology/OVERVIEW.md)
- [Bangladesh Profiles](../../docs/fhir/PROFILES.md)

## License

MIT License - see LICENSE file for details.

## Support

- GitHub Issues: https://github.com/zarishsphere/zs-core-fhir-engine/issues
