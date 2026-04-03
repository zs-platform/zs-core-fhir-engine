---
trigger: always_on
---
# ADR-0001 — Go 1.26.1 as Primary Backend Language

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @code-and-brain

---

## Context and Problem Statement

ZarishSphere requires a backend language for:
- The FHIR R5 REST engine (primary constraint)
- 51 microservices (clinical, public health, ERP)
- ARM64 binaries that run on Raspberry Pi 5 (last-mile offline)
- Single-binary deployment without external runtime

The language must compile to self-contained binaries, have excellent concurrency, and support ARM64 cross-compilation from Linux.

---

## Decision Drivers

- **Zero cost** — no runtime licence fees
- **Single binary** — for Raspberry Pi and offline deployment
- **ARM64 native** — compile from x86 to ARM64 for Raspberry Pi 5
- **FHIR library availability** — fhir-toolbox-go (damedic), gofhir-models (fastenhealth)
- **Concurrency model** — goroutines handle thousands of concurrent FHIR requests
- **Performance** — sub-10ms FHIR resource reads on commodity hardware
- **CNCF ecosystem** — Docker, Kubernetes, Prometheus all have native Go SDKs

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **Go 1.26.1** | Single binary, ARM64, FHIR libs available, CNCF native, fast | Smaller talent pool than Java/Python |
| Java 21 (Spring) | Large HAPI FHIR library, huge talent pool | JVM runtime required, 500MB+ container, slow cold start |
| Python 3.13 (FastAPI) | Easy to write, huge ecosystem | GIL limits concurrency, no single-binary, slow for FHIR validation |
| Rust | Fastest possible, zero runtime | FHIR libraries immature, steep learning curve for contributors |
| Node.js 22 | Large ecosystem, easy deployment | Single-threaded (worker_threads complexity), memory usage |

---

## Decision Outcome

**Chosen option: Go 1.26.1**

Go 1.26.1 is the only option that satisfies ALL constraints simultaneously:
1. Compiles to a single self-contained binary
2. Native ARM64 cross-compilation (for Raspberry Pi 5)
3. Green Tea GC (default in 1.26) eliminates GC pause concerns
4. Mature FHIR R5 libraries available (fhir-toolbox-go)
5. OpenTelemetry, Prometheus, Kubernetes all have first-class Go SDKs
6. Zero runtime licence cost
7. Excellent concurrency via goroutines

---

## Positive Consequences

- Single binary deploys anywhere (Raspberry Pi 5, Kubernetes, Docker)
- ~5 MB container images (distroless) vs 500MB JVM
- Native ARM64 cross-compilation supports offline edge deployments
- Green Tea GC eliminates pause-the-world GC events
- NATS, PostgreSQL (pgx), OpenTelemetry all have first-class Go clients

## Negative Consequences / Trade-offs

- Smaller available talent pool than Java or Python
- Error handling is verbose compared to exceptions
- No generics for all patterns (though Go 1.21+ generics cover most cases)
- FHIR library ecosystem smaller than HAPI Java

---

## Implementation Notes

All Go services must use:
- Go 1.26.1 (specified in `go.mod`)
- chi v5 for HTTP routing
- pgx v5 for PostgreSQL
- zerolog for structured logging
- testcontainers-go for integration tests

---

## Links

- RFC: RFC-0002 (Repository Structure Standards)
- FHIR library: https://github.com/damedic/fhir-toolbox-go
- Go download: https://go.dev/dl/
