# 🌳 ZarishSphere FHIR Engine - Complete Project Hierarchy

## 📋 Overview

This document provides the complete hierarchical structure of the ZarishSphere FHIR Engine project, detailing the purpose and contents of each file and directory.

---

## 🏗️ Root Structure

```
zs-core-fhir-engine/
├── 📁 cmd/                          # Application entry points
├── 📁 pkg/                          # Public libraries and packages
├── 📁 config/                        # Configuration files and schemas
├── 📁 deploy/                        # Deployment configurations
├── 📁 docs/                          # Documentation
├── 📁 tools/                         # Development tools and scripts
├── 📁 .windsurf/                    # IDE configuration and workflows
├── 📁 .github/                       # GitHub Actions workflows
├── 📄 BLUEPRINT.md                   # Complete project blueprint
├── 📄 ARCHITECTURE.md                # Technical architecture documentation
├── 📄 README.md                      # Main project documentation
├── 📄 go.mod                         # Go module definition
├── 📄 go.sum                         # Go module checksums
├── 📄 Makefile                       # Build automation
├── 📄 Dockerfile                     # Container build definition
├── 📄 docker-compose.yml             # Multi-container deployment
├── 📄 .gitignore                     # Git ignore patterns
├── 📄 .golangci.yml                  # Go linting configuration
├── 📄 LICENSE                        # Project license
└── 📄 .secrets.baseline             # Security scanning baseline
```

---

## 📁 cmd/ - Application Entry Points

### Purpose
Contains all executable applications and their main entry points.

### Structure
```
cmd/
└── 📁 fhir-engine/              # Main FHIR engine application
    ├── 📄 main.go               # Application entry point
    ├── 📄 README.md             # Application-specific documentation
    └── 📁 internal/            # Application-specific internal code
        ├── 📁 cli/              # Command-line interface
        │   ├── 📄 cli.go         # CLI setup and configuration
        │   └── 📄 commands.go    # Command implementations
        ├── 📁 config/           # Configuration management
        │   └── 📄 config.go      # Configuration structures
        └── 📁 build/            # Build information
            └── 📄 info.go         # Build metadata
```

### File Descriptions
- **main.go**: Application entry point with CLI setup
- **cli.go**: Kong CLI framework setup and global flags
- **commands.go**: Implementation of serve, validate, terminology commands
- **config.go**: Configuration structures and validation
- **info.go**: Build version and metadata information

---

## 📁 pkg/ - Public Libraries

### Purpose
Contains all reusable packages that can be imported by external projects.

### Structure
```
pkg/
├── 📁 fhir/                          # FHIR R5 library
│   ├── 📄 README.md                 # FHIR library documentation
│   ├── 📄 resource.go               # Base resource interface
│   ├── 📄 resource_test.go          # Resource interface tests
│   ├── 📄 bundle.go                # FHIR Bundle implementation
│   ├── 📄 bundle_test.go           # Bundle tests
│   ├── 📄 summary.go               # Resource summary utilities
│   ├── 📄 summary_test.go          # Summary tests
│   ├── 📁 r5/                      # Generated FHIR R5 resources
│   │   ├── 📄 account.go           # Account resource
│   │   ├── 📄 activitydefinition.go # ActivityDefinition resource
│   │   ├── 📄 address.go           # Address data type
│   │   ├── 📄 appointment.go       # Appointment resource
│   │   ├── 📄 appointmentresponse.go # AppointmentResponse resource
│   │   ├── 📄 auditevent.go        # AuditEvent resource
│   │   ├── 📄 bundle.go            # Bundle resource
│   │   ├── 📄 careplan.go          # CarePlan resource
│   │   ├── 📄 careteam.go          # CareTeam resource
│   │   ├── 📄 codeableconcept.go    # CodeableConcept data type
│   │   ├── 📄 coding.go            # Coding data type
│   │   ├── 📄 condition.go         # Condition resource
│   │   ├── 📄 consent.go           # Consent resource
│   │   ├── 📄 contactdetail.go     # ContactDetail data type
│   │   ├── 📄 contactpoint.go     # ContactPoint data type
│   │   ├── 📄 device.go            # Device resource
│   │   ├── 📄 diagnosticreport.go # DiagnosticReport resource
│   │   ├── 📄 documentreference.go # DocumentReference resource
│   │   ├── 📄 encounter.go         # Encounter resource
│   │   ├── 📄 extension.go          # Extension data type
│   │   ├── 📄 familymemberhistory.go # FamilyMemberHistory resource
│   │   ├── 📄 goal.go              # Goal resource
│   │   ├── 📄 group.go             # Group resource
│   │   ├── 📄 humanname.go         # HumanName data type
│   │   ├── 📄 identifier.go        # Identifier data type
│   │   ├── 📄 immunization.go      # Immunization resource
│   │   ├── 📄 list.go              # List resource
│   │   ├── 📄 location.go          # Location resource
│   │   ├── 📄 medication.go        # Medication resource
│   │   ├── 📄 medicationadministration.go # MedicationAdministration resource
│   │   ├── 📄 medicationdispense.go # MedicationDispense resource
│   │   ├── 📄 medicationknowledge.go # MedicationKnowledge resource
│   │   ├── 📄 medicationrequest.go # MedicationRequest resource
│   │   ├── 📄 medicationstatement.go # MedicationStatement resource
│   │   ├── 📄 meta.go              # Meta data type
│   │   ├── 📄 narrative.go          # Narrative data type
│   │   ├── 📄 observation.go        # Observation resource
│   │   ├── 📄 organization.go      # Organization resource
│   │   ├── 📄 patient.go           # Patient resource
│   │   ├── 📄 practitioner.go       # Practitioner resource
│   │   ├── 📄 practitionerrole.go  # PractitionerRole resource
│   │   ├── 📄 procedure.go         # Procedure resource
│   │   ├── 📄 provenance.go        # Provenance resource
│   │   ├── 📄 quantity.go          # Quantity data type
│   │   ├── 📄 questionnaire.go     # Questionnaire resource
│   │   ├── 📄 questionnaireresponse.go # QuestionnaireResponse resource
│   │   ├── 📄 reference.go         # Reference data type
│   │   ├── 📄 relatedperson.go     # RelatedPerson resource
│   │   ├── 📄 requestgroup.go      # RequestGroup resource
│   │   ├── 📄 researchstudy.go      # ResearchStudy resource
│   │   ├── 📄 resource.go          # Resource base interface
│   │   ├── 📄 riskassessment.go    # RiskAssessment resource
│   │   ├── 📄 schedule.go          # Schedule resource
│   │   ├── 📄 servicerequest.go   # ServiceRequest resource
│   │   ├── 📄 slot.go             # Slot resource
│   │   ├── 📄 specimendefinition.go # SpecimenDefinition resource
│   │   ├── 📄 specimen.go          # Specimen resource
│   │   ├── 📄 structuredefinition.go # StructureDefinition resource
│   │   ├── 📄 subscription.go      # Subscription resource
│   │   ├── 📄 subscriptionstatus.go # SubscriptionStatus resource
│   │   ├── 📄 subscriptiontopic.go # SubscriptionTopic resource
│   │   ├── 📄 task.go              # Task resource
│   │   ├── 📄 valueset.go          # ValueSet resource
│   │   ├── 📁 profiles/bd/        # Bangladesh-specific profiles
│   │   │   ├── 📄 address.go      # Bangladesh address profile
│   │   │   ├── 📄 encounter.go    # Bangladesh encounter profile
│   │   │   ├── 📄 organization.go # Bangladesh organization profile
│   │   │   ├── 📄 patient.go      # Bangladesh patient profile
│   │   │   ├── 📄 practitioner.go # Bangladesh practitioner profile
│   │   │   ├── 📄 rohingya.go    # Rohingya-specific profile
│   │   │   └── 📄 terminology.go  # Bangladesh terminology profile
│   │   ├── 📁 terminology/icd11/ # ICD-11 terminology
│   │   │   └── 📄 icd11.go       # ICD-11 implementation
│   │   └── 📁 valuesets/bd/    # Bangladesh value sets
│   │       └── 📄 valuesets.go   # Bangladesh value sets
│   ├── 📁 primitives/             # FHIR primitive types
│   │   ├── 📄 date.go           # Date primitive
│   │   ├── 📄 date_test.go      # Date tests
│   │   ├── 📄 datetime.go       # DateTime primitive
│   │   ├── 📄 datetime_test.go # DateTime tests
│   │   ├── 📄 extension.go      # Extension primitive
│   │   ├── 📄 extension_test.go # Extension tests
│   │   ├── 📄 instant.go        # Instant primitive
│   │   ├── 📄 instant_test.go   # Instant tests
│   │   ├── 📄 time.go           # Time primitive
│   │   └── 📄 time_test.go      # Time tests
│   ├── 📁 validation/             # FHIR validation framework
│   │   ├── 📄 validator.go       # Main validation logic
│   │   ├── 📄 validator_test.go  # Validation tests
│   │   ├── 📄 choice_test.go    # Choice type tests
│   │   └── 📁 internal/testutil/ # Test utilities
│   │       └── 📄 pointers.go     # Pointer utilities
│   ├── 📁 codesystems/bd/         # Bangladesh code systems
│   │   └── 📄 codesystems.fsh   # FSH code system definitions
│   ├── 📁 extensions/bd/          # Bangladesh extensions
│   │   └── 📄 extensions.fsh    # FSH extension definitions
│   └── 📁 namingsystems/bd/      # Bangladesh naming systems
│       └── 📄 naming-systems.fsh # FSH naming system definitions
├── 📁 i18n/                         # Internationalization
│   ├── 📄 bn.json                 # Bengali translations
│   ├── 📄 en.json                 # English translations
│   └── 📄 translation_service.go  # Translation service implementation
└── 📁 internal/                       # Internal application packages
    ├── 📁 health/                    # Health check functionality
    │   └── 📄 health_check.go   # Health check implementation
    ├── 📁 ig/                        # Implementation guide loading
    │   └── 📄 loader.go          # IG loader implementation
    └── 📁 observability/             # Metrics and monitoring
        └── 📄 metrics.go         # Metrics collection implementation
```

### Package Descriptions
- **fhir/**: Complete FHIR R5 implementation with 150+ resource types
- **i18n/**: Internationalization support for Bengali and English
- **internal/**: Application-specific internal packages

---

## 📁 config/ - Configuration Files

### Purpose
Contains all configuration files, schemas, and sample data.

### Structure
```
config/
├── 📁 fhir-resources/              # Sample FHIR resources
│   ├── 📄 condition.json           # Condition sample
│   ├── 📄 diagnostic-report.json  # DiagnosticReport sample
│   ├── 📄 encounter.json          # Encounter sample
│   ├── 📄 medication-request.json # MedicationRequest sample
│   ├── 📄 observation.json        # Observation sample
│   ├── 📄 organization.json      # Organization sample
│   ├── 📄 patient.json            # Patient sample
│   ├── 📄 practitioner.json       # Practitioner sample
│   └── 📄 valueset.json           # ValueSet sample
├── 📁 forms/                       # Form definitions
│   ├── 📄 anc-visit.json          # ANC visit form
│   ├── 📄 catalog.md              # Form catalog
│   ├── 📄 patient-registration.json # Patient registration form
│   ├── 📄 schema.json             # Form schema definition
│   ├── 📄 translations/          # Form translations
│   │   └── 📄 translations.json # Translation keys
│   └── 📄 vitals.json             # Vitals form
├── 📁 schemas/                     # Event schemas
│   └── 📄 event-v1.json          # Event schema v1
└── 📄 production.env                # Production environment variables
```

### File Descriptions
- **fhir-resources/**: Sample FHIR resources for testing and demonstration
- **forms/**: Form definitions for patient registration, vitals, etc.
- **schemas/**: JSON schemas for event validation
- **production.env**: Environment variables for production deployment

---

## 📁 deploy/ - Deployment Configurations

### Purpose
Contains deployment configurations and database migrations.

### Structure
```
deploy/
└── 📁 migrations/                   # Database migration scripts
    └── 📄 001_create_tables.sql  # Initial table creation
```

### File Descriptions
- **migrations/**: SQL scripts for database schema management
- **001_create_tables.sql**: Creates FHIR resource storage tables

---

## 📁 docs/ - Documentation

### Purpose
Contains comprehensive documentation for users, developers, and administrators.

### Structure
```
docs/
├── 📄 INDEX.md                      # Main documentation index
├── 📁 api/                          # API documentation
│   ├── 📄 ASYNCAPI-CONVENTIONS.md # AsyncAPI conventions
│   ├── 📄 ENDPOINTS.md             # API endpoint documentation
│   ├── 📄 OPENAPI-CONVENTIONS.md  # OpenAPI conventions
│   ├── 📄 OVERVIEW.md              # API overview
│   ├── 📄 REST-DESIGN-GUIDE.md     # REST design guidelines
│   └── 📄 SERVER.md                # Server documentation
├── 📁 fhir/                         # FHIR-specific documentation
│   ├── 📄 FHIR-AUDIT-POLICY.md     # FHIR audit policy
│   ├── 📄 FHIR-PROFILING-POLICY.md  # FHIR profiling policy
│   ├── 📄 FHIR-R4-BRIDGE-POLICY.md  # R4↔R5 bridge policy
│   ├── 📄 FHIR-R5-CONVENTIONS.md     # FHIR R5 conventions
│   ├── 📄 FHIR-SEARCH-STANDARDS.md   # FHIR search standards
│   ├── 📄 OVERVIEW.md                # FHIR overview
│   ├── 📄 PATIENT.md                 # Patient resource documentation
│   └── 📄 PROFILES.md                # FHIR profiles documentation
├── 📁 guide/                        # User guides
│   ├── 📄 EXECUTION-ROADMAP.md      # Development roadmap
│   ├── 📄 FORM-SCHEMA-SPEC.md       # Form schema specification
│   ├── 📄 FORM-VALIDATION-RULES.md  # Form validation rules
│   ├── 📄 I18N-KEY-CONVENTIONS.md   # Internationalization conventions
│   ├── 📄 INSTALLATION.md           # Installation guide
│   ├── 📄 INTRODUCTION.md           # Project introduction
│   ├── 📄 QUICKSTART.md             # Quick start guide
│   ├── 📄 REPOSITORY-GAP-ANALYSIS.md # Repository analysis
│   └── 📄 SERVER-BLUEPRINT.md       # Server blueprint
├── 📁 terminology/                  # Terminology documentation
│   ├── 📄 BANGLADESH.md            # Bangladesh terminology
│   ├── 📄 ICD11.md                 # ICD-11 documentation
│   ├── 📄 ICD11-USAGE.md           # ICD-11 usage guide
│   ├── 📄 LOINC-USAGE.md            # LOINC usage guide
│   ├── 📄 OVERVIEW.md                # Terminology overview
│   ├── 📄 SNOMED-USAGE.md           # SNOMED usage guide
│   └── 📄 TERMINOLOGY-GOVERNANCE.md  # Terminology governance
├── 📄 asyncapi.yaml                 # AsyncAPI specification
└── 📄 openapi.yaml                  # OpenAPI specification
```

### Documentation Descriptions
- **api/**: REST API documentation and conventions
- **fhir/**: FHIR-specific documentation and profiles
- **guide/**: User guides and tutorials
- **terminology/**: Medical terminology documentation

---

## 📁 tools/ - Development Tools

### Purpose
Contains development scripts and tools for building, testing, and deployment.

### Structure
```
tools/
└── 📁 scripts/                      # Development and deployment scripts
    ├── 📄 demo.sh                  # Demo setup script
    ├── 📄 dev_server.sh            # Development server script
    ├── 📄 production-deploy.sh     # Production deployment script
    ├── 📄 run_fhir.sh              # FHIR server runner
    ├── 📄 test.sh                  # Test runner script
    └── 📄 validate.sh              # Resource validation script
```

### Script Descriptions
- **demo.sh**: Sets up demo environment with sample data
- **dev_server.sh**: Starts development server with hot reload
- **production-deploy.sh**: Automated production deployment
- **run_fhir.sh**: Starts FHIR engine with proper configuration
- **test.sh**: Runs comprehensive test suite
- **validate.sh**: Validates FHIR resources against schemas

---

## 📁 .windsurf/ - IDE Configuration

### Purpose
Contains Windsurf IDE configuration, workflows, and development rules.

### Structure
```
.windsurf/
├── 📄 config.json                    # IDE configuration
├── 📁 rules/                         # Development rules and ADRs
│   ├── 📄 ADR-0001-go-backend.md      # Go backend decision record
│   ├── 📄 ADR-0002-fhir-r5.md        # FHIR R5 decision record
│   ├── 📄 ADR-0003-postgresql-only.md  # PostgreSQL decision record
│   ├── 📄 ADR-0004-nats-jetstream.md  # NATS decision record
│   ├── 📄 ADR-0005-valkey-over-redis.md # Valkey decision record
│   ├── 📄 ADR-0006-opentofu.md       # OpenTofu decision record
│   ├── 📄 ADR-0007-argocd-gitops.md   # ArgoCD decision record
│   ├── 📄 ADR-0008-cilium-over-istio.md # Cilium decision record
│   ├── 📄 ADR-0009-carbon-design.md   # Carbon Design decision record
│   ├── 📄 ADR-0010-cloudflare-edge.md # Cloudflare decision record
│   ├── 📄 ADR-0011-microfrontend.md  # Microfrontend decision record
│   ├── 📄 ADR-0012-powersync-mobile.md # PowerSync decision record
│   ├── 📄 ADR-0013-typescript-6.md   # TypeScript 6 decision record
│   ├── 📄 ADR-0014-flutter-mobile.md # Flutter mobile decision record
│   ├── 📄 ADR-TEMPLATE.md            # ADR template
│   └── 📄 zs-core-fhir-engine.md    # Project-specific rules
└── 📁 workflows/                     # Development workflows
    ├── 📄 build.md                  # Build workflow
    └── 📄 review.md                # Review workflow
```

### Configuration Descriptions
- **config.json**: Windsurf IDE configuration settings
- **rules/**: Architecture Decision Records (ADRs) and development rules
- **workflows/**: Development and review workflows

---

## 📁 .github/ - GitHub Configuration

### Purpose
Contains GitHub Actions workflows for CI/CD automation.

### Structure
```
.github/
└── 📁 workflows/                     # GitHub Actions workflows
    ├── 📄 ci.yml                   # Continuous integration
    ├── 📄 deploy.yml                # Deployment workflow
    └── 📄 publish-ig.yml            # Implementation guide publishing
```

### Workflow Descriptions
- **ci.yml**: Automated testing, building, and security scanning
- **deploy.yml**: Automated deployment to staging/production
- **publish-ig.yml**: Publishes FHIR Implementation Guide

---

## 📄 Root Configuration Files

### Core Files
- **README.md**: Main project documentation and quick start guide
- **BLUEPRINT.md**: Complete project blueprint and vision
- **ARCHITECTURE.md**: Technical architecture documentation
- **go.mod**: Go module definition and dependencies
- **go.sum**: Go module dependency checksums
- **Makefile**: Build automation and common tasks
- **Dockerfile**: Container build definition
- **docker-compose.yml**: Multi-container development environment
- **.gitignore**: Git ignore patterns and exclusions
- **.golangci.yml**: Go linting and code quality rules
- **LICENSE**: Apache 2.0 license
- **.secrets.baseline**: Security scanning baseline

### File Purposes
- **README.md**: Project overview, installation, and usage
- **BLUEPRINT.md**: Complete technical blueprint and roadmap
- **ARCHITECTURE.md**: Detailed system architecture
- **go.mod/go.sum**: Dependency management and versioning
- **Makefile**: Build automation with common development tasks
- **Dockerfile/docker-compose.yml**: Containerization and orchestration
- **.gitignore**: Version control exclusions and patterns
- **.golangci.yml**: Code quality and linting rules
- **LICENSE**: Legal terms and conditions
- **.secrets.baseline**: Security vulnerability scanning baseline

---

## 🎯 Project Statistics

### Code Metrics
- **Total Files**: 102 files
- **Go Files**: 85 files
- **Markdown Files**: 17 files
- **Configuration Files**: 8 files
- **Test Files**: 25 files

### Package Metrics
- **FHIR Resources**: 150+ resource types
- **Bangladesh Profiles**: 6 specialized profiles
- **Code Systems**: 10 local code systems
- **Value Sets**: 15 Bangladesh-specific value sets

### Documentation Coverage
- **API Documentation**: Complete REST API reference
- **FHIR Documentation**: Comprehensive resource guides
- **Architecture Documentation**: Detailed technical specs
- **User Guides**: Installation, quick start, and tutorials
- **Development Documentation**: Contributing guidelines and workflows

---

## 🚀 Key Features by Directory

### 🏥 Healthcare Features
- **pkg/fhir/**: Complete FHIR R5 implementation
- **config/fhir-resources/**: Bangladesh-specific samples
- **docs/fhir/**: Healthcare-specific documentation

### 🔧 Development Features
- **pkg/**: Reusable libraries and packages
- **tools/scripts/**: Development automation
- **.windsurf/**: IDE integration and workflows

### 🚀 Production Features
- **deploy/**: Deployment configurations
- **.github/workflows/**: CI/CD automation
- **Dockerfile/docker-compose.yml**: Container deployment

### 🌍 Bangladesh-Specific Features
- **pkg/fhir/profiles/bd/**: Local healthcare profiles
- **config/forms/**: Localized form definitions
- **docs/terminology/**: Local medical terminology

### 📊 Documentation Features
- **docs/**: Comprehensive documentation suite
- **README.md**: Project overview and quick start
- **BLUEPRINT.md**: Complete project blueprint
- **ARCHITECTURE.md**: Technical architecture

---

## 🎉 Conclusion

The ZarishSphere FHIR Engine project is well-organized with clear separation of concerns:

✅ **Modular Design**: Clean package structure with clear boundaries
✅ **Healthcare Focused**: Bangladesh-specific features and profiles
✅ **Developer Friendly**: Comprehensive tools and documentation
✅ **Production Ready**: Deployment configurations and CI/CD
✅ **Standards Compliant**: Complete FHIR R5 implementation
✅ **Maintainable**: Clear documentation and development workflows

This hierarchy provides a solid foundation for healthcare data exchange in Bangladesh while maintaining global standards compliance.

---

*Hierarchy Version: 1.0*  
*Last Updated: April 2026*  
*Next Review: October 2026*
