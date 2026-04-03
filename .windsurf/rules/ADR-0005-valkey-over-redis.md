---
trigger: always_on
---
# ADR-0005 — Valkey 9.0.3 over Redis

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @code-and-brain

---

## Context and Problem Statement

ZarishSphere requires a fast in-memory key-value store for:
- FHIR resource hot cache (reduce PostgreSQL reads)
- Terminology lookup cache (ICD-11, SNOMED codes — 24h TTL)
- User session tokens
- Background job queue (Asynq)
- Rate limiting counters

---

## Decision Drivers

- **Zero cost forever** — Redis changed to SSPL in 2024, restricting use
- **Linux Foundation governance** — community-owned, no proprietary control
- **API compatibility** — Valkey is a Redis 7.x drop-in replacement
- **Atomic slot migration** — Valkey 9.0.3 adds atomic slot migration (improves cluster reliability)

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **Valkey 9.0.3** | BSD-3, Linux Foundation, Redis-compatible, atomic slot migration | Newer, smaller community than Redis |
| Redis 7.x | Widely known, large ecosystem | SSPL licence restricts use in hosted services |
| Redis 8.x | New features | SSPL licence |
| Memcached | Simple, fast | No persistence, no data structures |
| Apache Ignite | Distributed, persistent | Heavy, complex |
| DragonflyDB | Redis-compatible, fast | BSL licence |

---

## Decision Outcome

**Chosen option: Valkey 9.0.3**

Redis's change to SSPL in 2024 means it cannot be used freely in a platform that enables third parties to build on it. Valkey (Linux Foundation, BSD-3 licence) is fully API-compatible, free forever, and in 9.0.3 has atomic slot migration for cluster reliability.

---

## Links

- Valkey: https://valkey.io/
- valkey-go client: https://github.com/valkey-io/valkey-go
