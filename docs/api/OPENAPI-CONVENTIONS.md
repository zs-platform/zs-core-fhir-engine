# OpenAPI 3.1 Conventions

> **Document:** OPENAPI-CONVENTIONS.md | **Version:** 2.0.0

---

## File Location

Every `zs-svc-*` and `zs-core-*` repository must have:

```
{repo}/docs/openapi.yaml
```

---

## Required OpenAPI Info Block

```yaml
openapi: "3.1.0"
info:
  title: "ZarishSphere {Service Name} API"
  description: |
    {Service description from PRD}

    ## Authentication
    All endpoints require a valid SMART on FHIR 2.1 JWT Bearer token.
    See https://docs.zarishsphere.com/auth for details.

  version: "1.0.0"
  contact:
    name: "ZarishSphere Platform"
    email: "platform@zarishsphere.com"
    url: "https://zarishsphere.com"
  license:
    name: "Apache 2.0"
    url: "https://www.apache.org/licenses/LICENSE-2.0"

servers:
  - url: "https://api.zarishsphere.com"
    description: "Production"
  - url: "http://localhost:{port}"
    description: "Local development"
    variables:
      port:
        default: "8001"

security:
  - bearerAuth: []

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

---

## FHIR Endpoint Pattern

```yaml
paths:
  /fhir/R5/Patient:
    post:
      summary: Create a FHIR Patient resource
      operationId: createPatient
      tags: [Patient]
      requestBody:
        required: true
        content:
          application/fhir+json:
            schema:
              $ref: "#/components/schemas/Patient"
      responses:
        "201":
          description: Patient created successfully
          headers:
            Location:
              schema:
                type: string
              description: URL of the created resource
          content:
            application/fhir+json:
              schema:
                $ref: "#/components/schemas/Patient"
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "422":
          $ref: "#/components/responses/UnprocessableEntity"
```

---

## Standard Error Responses

```yaml
components:
  responses:
    BadRequest:
      description: Bad request — invalid input
      content:
        application/fhir+json:
          schema:
            $ref: "#/components/schemas/OperationOutcome"
    Unauthorized:
      description: Unauthorized — missing or invalid JWT
    Forbidden:
      description: Forbidden — insufficient SMART scope
    NotFound:
      description: Resource not found
    UnprocessableEntity:
      description: FHIR validation failure
      content:
        application/fhir+json:
          schema:
            $ref: "#/components/schemas/OperationOutcome"
```
