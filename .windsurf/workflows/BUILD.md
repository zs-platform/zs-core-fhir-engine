---
auto_execution_mode: 2
description: build
---
Build ZarishSphere FHIR Engine and related components
This workflow builds the core FHIR R5 engine and its associated services based on ADR specifications.

## Current Project Structure

**Primary Repository**: zs-core-fhir-engine
- Location: `/home/ariful/Desktop/zarishsphere/_cloned/zs-core-fhir-engine`
- Purpose: FHIR R5 REST engine with validation and terminology services
- Status: ✅ Core implementation complete

## Build Targets

### 1. Core FHIR Engine (zs-core-fhir-engine)
**Status**: ✅ Ready for build
**Components**:
- FHIR R5 REST server
- Resource validation framework  
- Terminology server (ValueSet/$expand)
- In-memory storage with interface for PostgreSQL
- CLI with serve, validate, terminology, version commands

**Build Command**:
```bash
cd /home/ariful/Desktop/zarishsphere/_cloned/zs-core-fhir-engine
go mod tidy
go build ./cmd/zs-core-fhir-engine
```

**Docker Build**:
```bash
docker build -t zs-core-fhir-engine .
```

### 2. Planned Additional Repositories (Not Yet Built)
Based on ADR specifications, the following repositories are planned:

**High Priority**:
- zs-core-fhir-r4-bridge - R4↔R5 translation
- zs-core-fhir-validator - Standalone validator service
- zs-core-fhir-subscriptions - NATS-backed subscriptions
- zs-pkg-go-fhir - Shared Go FHIR library

**Medium Priority**:
- zs-core-fhirpath - FHIRPath 2.0 evaluator
- zs-core-cds-hooks - CDS Hooks 2.0 service
- zs-data-fhir-profiles - Bangladesh FHIR Implementation Guide

## Build Verification Steps

### Context
Verify that the FHIR engine builds correctly and all components function.

### Status Check
1. Check Go module dependencies
2. Verify compilation succeeds
3. Test CLI functionality
4. Validate Docker build

### Execution
```bash
# 1. Clean and build
go clean -cache
go mod tidy
go build ./cmd/zs-core-fhir-engine

# 2. Test CLI
./zs-core-fhir-engine --help
./zs-core-fhir-engine version
./zs-core-fhir-engine validate --help

# 3. Docker build (if Docker available)
docker build -t zs-core-fhir-engine .
```

### Verification
- ✅ Binary compiles without errors
- ✅ CLI commands respond correctly
- ✅ Version information displays properly
- ✅ Help text is complete and accurate
- ✅ Docker image builds successfully

## Production Readiness Checklist

### Completed ✅
- [x] Core FHIR server implementation
- [x] Resource validation framework
- [x] CLI interface with multiple commands
- [x] Docker containerization
- [x] Basic terminology service
- [x] Storage interface abstraction

### Next Steps for Production 🚧
- [ ] PostgreSQL persistence layer
- [ ] FHIR search parameters implementation
- [ ] Authentication/authorization system
- [ ] Resource versioning and history
- [ ] TLS/HTTPS support
- [ ] OpenTelemetry observability
- [ ] Comprehensive test suite
- [ ] CI/CD pipeline setup

## Architecture Decisions Applied

From `.windsurf/rules/` ADRs:
- **Go 1.26.1** backend language (ADR-0001)
- **FHIR R5** data model (ADR-0002) 
- **PostgreSQL 18.3** database (ADR-0003)
- **NATS 2.12.5** messaging (ADR-0004)
- **Valkey 9.0.3** caching (ADR-0005)
- **OpenTofu 1.11.x** infrastructure (ADR-0006)
- **Argo CD 3.3.x** GitOps (ADR-0007)
- **Cilium 1.17+** networking (ADR-0008)

## Integration Points

### External Systems
- **Bangladesh DGHS**: Health data exchange
- **WHO**: ICD-11 terminology integration
- **HL7**: FHIR standards compliance

### Internal Services
- **Authentication**: SMART on FHIR OAuth2
- **Audit**: Comprehensive audit logging
- **Monitoring**: Prometheus metrics + Grafana dashboards