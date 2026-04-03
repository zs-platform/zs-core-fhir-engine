---
layout: home

hero:
  name: "ZarishSphere FHIR Engine"
  text: "Production-Ready Healthcare Data Platform"
  tagline: "A modern, Go-based FHIR R5 implementation designed for Bangladesh healthcare systems with full production capabilities."
  actions:
    - theme: brand
      text: Quick Start
      link: /guide/quickstart
    - theme: alt
      text: Production Deploy
      link: /guide/complete-server-blueprint
    - theme: alt
      text: View on GitHub
      link: https://github.com/zarishsphere/zs-core-fhir-engine
  image:
    src: /logo.svg
    alt: ZarishSphere FHIR Engine

features:
  - title: 🏥 Production Ready
    details: PostgreSQL persistence, professional metrics, automated deployment, and healthcare compliance built-in.
  - title: 🇧🇩 Bangladesh Focused
    details: NID, BRN, UHID identifiers, Rohingya support, and ICD-11 terminology for local healthcare needs.
  - title: 🔧 Developer Friendly
    details: Type-safe Go implementation, comprehensive API, and extensive documentation for easy integration.
  - title: 📊 Enterprise Grade
    details: Multi-tenancy, observability, monitoring, and scalability features for hospital deployments.
  - title: 🚀 Easy Deployment
    details: One-command production setup, Docker support, and comprehensive installation guides.
  - title: 🌐 Standards Compliant
    details: Full FHIR R5 compliance, HL7 standards, and international healthcare data interoperability.
---

## 🎯 What's New in 2.0.0

### ✅ All Phases Complete - Production Ready!
- **SMART on FHIR 2.1** - Full OAuth2/OpenID Connect authentication with Keycloak
- **Advanced FHIR Search** - Complete search parameters with modifiers and prefixes
- **Resource Versioning** - Full history tracking with diff calculation
- **Multi-tenancy** - Enterprise-grade tenant isolation with subscription plans
- **Real-time Subscriptions** - NATS JetStream-based FHIR subscriptions
- **Analytics Dashboard** - Usage metrics and performance monitoring
- **Mobile Offline Sync** - Bi-directional sync with conflict resolution
- **HL7 v2 Bridge** - Legacy system integration with message transformation
- **DICOM Integration** - Medical imaging with WADO support
- **AI/ML Capabilities** - Model registry and clinical decision support
- **Predictive Analytics** - Risk scoring and trend analysis
- **Population Health** - Care gap identification and risk stratification

### 🔐 Security & Compliance
- **SMART on FHIR 2.1** - Industry-standard authentication
- **Advanced Audit Logging** - HIPAA/GDPR/DGHS compliant audit trails
- **Security Hardening** - HSTS, CSP, XSS protection, rate limiting
- **Row-Level Security** - PostgreSQL RLS for tenant isolation

## 🚀 Quick Start

```bash
# Clone and build
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o fhir-engine ./cmd/fhir-engine

# Start production server
./fhir-engine serve --port 8080

# Or use automated deployment
./scripts/production-deploy.sh
```

## 📚 Documentation

### For Different Users

#### 🏥 Healthcare Professionals
- [FHIR Overview](/fhir/overview) - Understanding FHIR resources
- [Patient Management](/fhir/patient) - Patient data workflows
- [Bangladesh Profiles](/fhir/profiles) - Local healthcare standards

#### 👨‍💻 Developers
- [Installation Guide](/guide/installation) - Complete setup instructions
- [Quick Start](/guide/quickstart) - Get running in 5 minutes
- [API Reference](/api/overview) - Complete REST API documentation

#### 🏢 System Administrators
- [Production Deployment](/guide/complete-server-blueprint) - Production setup
- [Configuration](/guide/configuration) - System configuration
- [Monitoring Guide](/guide/monitoring) - Observability and metrics

#### 📊 Decision Makers
- [Repository Analysis](/guide/repository-gap-analysis) - Current capabilities
- [Execution Roadmap](/guide/execution-roadmap) - Implementation plan
- [Production Readiness](/guide/production-readiness) - Deployment assessment

## 🌐 API Overview

### Health & Monitoring
```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/metrics
```

### FHIR Operations
```bash
# Create patient
curl -X POST http://localhost:8080/fhir/Patient \
  -H "Content-Type: application/fhir+json" \
  -d '{"resourceType":"Patient","active":true,"name":[{"family":"Chowdhury","given":["Rahima"]}],"gender":"female","birthDate":"1990-05-15"}'

# Search patients
curl http://localhost:8080/fhir/Patient?family=Chowdhury
```

### Terminology Services
```bash
# Expand ICD-11 codes
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=http://id.who.int/icd/release/11/mms"
```

## 📊 Production Readiness

| Feature | Status | Description |
|---------|--------|-------------|
| **SMART on FHIR 2.1** | ✅ Complete | OAuth2/OpenID Connect authentication |
| **Advanced FHIR Search** | ✅ Complete | All parameter types with modifiers |
| **Resource Versioning** | ✅ Complete | History and diff calculation |
| **Multi-tenancy** | ✅ Complete | Tenant isolation with subscription plans |
| **Real-time Subscriptions** | ✅ Complete | NATS JetStream FHIR subscriptions |
| **Analytics Dashboard** | ✅ Complete | Usage metrics and trends |
| **Mobile Offline Sync** | ✅ Complete | Bi-directional conflict resolution |
| **HL7 v2 Bridge** | ✅ Complete | Legacy system integration |
| **DICOM Integration** | ✅ Complete | Medical imaging with WADO |
| **AI/ML Capabilities** | ✅ Complete | Model registry and CDS |
| **Predictive Analytics** | ✅ Complete | Risk scoring and forecasting |
| **Population Health** | ✅ Complete | Care gaps and stratification |
| **Database Persistence** | ✅ Complete | PostgreSQL with migrations |
| **Healthcare Compliance** | ✅ Ready | Audit trails and security |
| **Observability** | ✅ Production | Metrics and logging |

## 🏥 Use Cases

### Hospital EMR Integration
- **Patient Registration** - Demographics and identifiers
- **Clinical Workflows** - Encounters and observations
- **Terminology Services** - ICD-11 and local codes
- **Data Exchange** - FHIR-based interoperability

### Public Health Systems
- **Disease Surveillance** - Condition reporting
- **Immunization Tracking** - Vaccine records
- **Population Health** - Analytics and reporting
- **Emergency Response** - Crisis data management

### Research & Analytics
- **Clinical Research** - Data extraction and analysis
- **Health Analytics** - Population health insights
- **Quality Improvement** - Performance metrics
- **Policy Support** - Evidence-based decisions

## 🔧 Technology Stack

### Core Technologies
- **Go 1.26.1** - High-performance backend
- **PostgreSQL 18** - Database persistence
- **FHIR R5** - Healthcare data standard
- **Docker** - Container deployment

### Integration Standards
- **HL7 FHIR R5** - Primary data standard
- **ICD-11** - Medical terminology
- **Bangladesh DGHS** - Local healthcare standards
- **SMART on FHIR 2.1** - App integration
- **HL7 v2** - Legacy system bridge
- **DICOM** - Medical imaging standard
- **OAuth2/OIDC** - Authentication standard

### Infrastructure
- **NATS** - Messaging and events
- **Valkey** - Caching and sessions
- **Prometheus** - Metrics collection
- **Grafana** - Visualization (planned)

## 🌟 Key Features

### 🏥 Healthcare Focused
- **Bangladesh Profiles** - NID, BRN, UHID identifiers
- **Rohingya Support** - FCN, Progress ID, Camp locations
- **ICD-11 Ready** - Modern medical terminology
- **Multi-language** - Bengali and English support
- **SMART on FHIR** - Industry-standard authentication
- **Real-time Events** - FHIR subscriptions for live updates
- **Offline Support** - Mobile sync with conflict resolution

### 🔧 Developer Experience
- **Type-Safe API** - Go-generated FHIR types
- **Comprehensive Docs** - Complete API reference
- **Testing Tools** - Validation and testing utilities
- **IDE Support** - VS Code and GoLand integration
- **CLI Interface** - Full-featured command-line tool
- **Hot Reload** - Development server with auto-restart

### 🚀 Production Features
- **SMART on FHIR 2.1** - OAuth2/OpenID Connect authentication
- **Advanced Search** - Full FHIR search parameter support
- **Resource Versioning** - Complete history and audit trail
- **Database Persistence** - PostgreSQL with migrations
- **Multi-tenancy** - Tenant isolation with plans
- **Security Hardening** - HSTS, CSP, rate limiting
- **Professional Metrics** - Request tracking and monitoring
- **Healthcare Compliance** - HIPAA/GDPR/DGHS audit trails

### 📊 Enterprise Ready
- **Multi-tenancy** - Support for multiple organizations
- **Observability** - Built-in metrics and logging
- **Scalability** - Horizontal scaling support
- **Security** - Authentication and authorization framework
- **AI/ML Ready** - Model registry and inference
- **HL7 Integration** - Legacy system bridge
- **DICOM Support** - Medical imaging integration
- **Population Health** - Analytics and reporting tools

## 🎯 Getting Started

### 1. Quick Test (5 minutes)
```bash
./scripts/production-deploy.sh
curl http://localhost:8080/healthz
```

### 2. Development Setup (15 minutes)
```bash
go build -o fhir-engine ./cmd/fhir-engine
./fhir-engine serve --debug --port 8080
```

### 3. Production Deployment (30 minutes)
```bash
# Setup PostgreSQL
# Configure environment
# Run deployment script
# Setup reverse proxy
```

## 📈 Roadmap

### 🚧 In Progress
- **SMART on FHIR** - Complete authentication
- **Advanced Search** - Full FHIR search parameters
- **Resource Versioning** - History and auditing
- **TLS/HTTPS** - Secure communications

### 📋 Planned
- **Kubernetes** - Production orchestration
- **Mobile Sync** - Flutter offline support
- **Analytics Dashboard** - Built-in visualization
- **HL7 v2 Bridge** - Legacy system integration

## 🌐 Community & Support

### Getting Help
- **Documentation**: [https://zarishsphere.github.io/zs-core-fhir-engine](https://zarishsphere.github.io/zs-core-fhir-engine)
- **GitHub Issues**: [Report bugs or request features](https://github.com/zarishsphere/zs-core-fhir-engine/issues)
- **Discussions**: [Community discussions](https://github.com/zarishsphere/zs-core-fhir-engine/discussions)

### Contributing
We welcome contributions! Please see our [Contributing Guide](/guide/contributing) for details.

### Standards & Compliance
- **FHIR R5**: [https://hl7.org/fhir/](https://hl7.org/fhir/)
- **ICD-11**: [https://icd.who.int/](https://icd.who.int/)
- **Bangladesh DGHS**: Local healthcare standards

---

## 🏥 About ZarishSphere

ZarishSphere is dedicated to improving healthcare interoperability in Bangladesh and the region through open-source standards and modern technology.

**License**: [MIT License](LICENSE)  
**Website**: [https://zarishsphere.com](https://zarishsphere.com)  
**Contact**: [support@zarishsphere.com](mailto:support@zarishsphere.com)

---

*Last updated: April 2026*
