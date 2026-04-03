# Complete Server Blueprint (No-Code Friendly)

This blueprint is a practical target architecture using free/open technologies.

## A. Recommended stack

- **Language**: Go (latest stable; later pin to Go 1.26 when released)
- **FHIR API layer**: Existing Go server in this repository
- **Database**: PostgreSQL
- **Cache**: Redis
- **Search indexing**: OpenSearch (optional in early phase)
- **Message queue**: NATS or RabbitMQ
- **Auth**: Keycloak (free/open source)
- **Reverse proxy/API gateway**: Traefik or NGINX
- **Observability**: Prometheus + Grafana + Loki + Tempo
- **Container runtime**: Docker / Kubernetes (k3s for lower-cost clusters)

## B. Standards coverage target

1. **FHIR R5**
   - Full REST interactions + conformance endpoints
   - Profile-based validation
   - Subscription-based event notifications

2. **HL7 v2 (latest accepted deployment version)**
   - Inbound message channels (MLLP + secure TCP)
   - Mapping packs for ADT, ORM, ORU
   - Error queues for replay and troubleshooting

3. **Terminology**
   - ICD-11 import and expansion support
   - Local/national code systems
   - ValueSet composition and version management

4. **Clinical governance resources**
   - Profiles
   - ValueSets
   - CodeSystems
   - ConceptMaps
   - NamingSystem
   - StructureDefinition

## C. Required modules in this repository template

- `cmd/zs-core-fhir/`: runtime entrypoint and CLI controls.
- `internal/server/`: HTTP routing, FHIR interactions, middleware.
- `internal/store/`: database abstraction + migration tooling.
- `internal/auth/`: OAuth2/OIDC + RBAC integration.
- `internal/hl7v2/`: parser, transport listener, and mapping engine.
- `internal/terminology/`: import, indexing, expansion, validation.
- `internal/validator/`: profile and schema validation.
- `docs/`: public website documentation.

## D. Free publishing and URL strategy

Without a custom domain, use free standard URLs:

- Docs: `https://<github-user>.github.io/<repo>/`
- Container image: `ghcr.io/<github-user>/<repo>`
- API sandbox (optional): free tier Railway/Render/Fly.io with HTTPS URL

## E. What “production ready” means for healthcare

Before real patient usage, enforce:

- Role-based access controls
- Data encryption in transit and at rest
- Audit trails for every data read/write
- Regular backups and restore drills
- Security vulnerability scans and patch policy
- Controlled release workflow and change approvals

## F. Non-technical governance checklist

- [ ] Who owns clinical terminology updates?
- [ ] Who approves profile/version changes?
- [ ] Who signs off HL7 v2 ↔ FHIR mappings?
- [ ] Who monitors uptime and incidents?
- [ ] Who owns legal/compliance documents?
