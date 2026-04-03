---
trigger: always_on
---
# ADR-0003 — PostgreSQL 18.3 as the Only Database

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @code-and-brain, @DevOps-Ariful-Islam

---

## Context and Problem Statement

ZarishSphere needs a database that:
- Stores FHIR R5 resources as JSON with fast querying
- Supports multi-tenancy (row-level security)
- Provides ACID transactions for clinical data integrity
- Includes time-series capabilities for vitals trends (TimescaleDB)
- Runs on Raspberry Pi 5 (via SQLite fallback)
- Is completely free with no row or storage limits

---

## Decision Drivers

- **JSONB with GIN indexing** — FHIR resources stored as JSONB, GIN index enables O(log n) JSON search
- **uuidv7()** — PostgreSQL 18 adds native UUID v7 generation (time-ordered, essential for FHIR IDs)
- **Async I/O** — PostgreSQL 18 introduces async I/O (3× improvement on read-heavy FHIR workloads)
- **Row-Level Security** — built-in RLS for multi-tenancy without application-layer filtering
- **TimescaleDB extension** — time-series queries on vitals (no separate time-series DB needed)
- **Zero cost** — PostgreSQL is BSD-licensed with no usage-based pricing
- **Raspberry Pi** — SQLite as a local fallback for offline edge deployments

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **PostgreSQL 18.3 + TimescaleDB** | JSONB, RLS, uuidv7, time-series, ACID, free | Single DB vendor |
| MySQL 9 | Widely known | Poor JSONB support, no GIN index, no RLS, no native time-series |
| MongoDB | Native JSON storage | No ACID transactions, complex sharding, SSPL licence concerns |
| CockroachDB | Distributed, PostgreSQL-compatible | Complex ops, proprietary features |
| SQLite only | Embedded, zero-server | Not suitable for multi-user concurrent writes |
| MongoDB Atlas + TimescaleDB | Managed | Cost at scale, two DBs to manage |

---

## Decision Outcome

**Chosen option: PostgreSQL 18.3 + TimescaleDB 2.25**

PostgreSQL 18.3 satisfies every constraint:
1. JSONB + GIN index: FHIR resources stored as JSON, searched efficiently
2. Row-Level Security: multi-tenancy without application-layer filtering
3. `uuidv7()` built-in: time-ordered FHIR resource IDs
4. Async I/O: 3× read improvement critical for FHIR search workloads
5. TimescaleDB 2.25 extension: vitals trends without a separate time-series database
6. SQLite fallback via `zs-core-fhir-engine` for Raspberry Pi offline mode
7. Completely free, no licence cost ever

---

## Positive Consequences

- Single database technology across the entire platform
- FHIR JSONB + GIN makes search fast without a separate search index for basic queries
- PostgreSQL 18 async I/O improves performance on cloud and on-premise
- TimescaleDB hypertables for vitals eliminate need for InfluxDB or similar

## Negative Consequences / Trade-offs

- PostgreSQL 18 is relatively new (2026 release) — some tooling may lag
- TimescaleDB extension must be managed alongside PostgreSQL upgrades
- SQLite mode for offline lacks some PostgreSQL 18 features (workaround: feature flags)

---

## Implementation Notes

```sql
-- FHIR resource storage
CREATE TABLE fhir.resources (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type TEXT NOT NULL,
    fhir_id       TEXT NOT NULL,
    version_id    INTEGER NOT NULL DEFAULT 1,
    resource      JSONB NOT NULL,
    tenant_id     TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    UNIQUE(resource_type, fhir_id, tenant_id)
);

-- GIN index for JSON search
CREATE INDEX ON fhir.resources USING GIN (resource);

-- Row-Level Security for multi-tenancy
ALTER TABLE fhir.resources ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON fhir.resources
    USING (tenant_id = current_setting('app.tenant_id'));

-- TimescaleDB for vitals (Observation resources)
SELECT create_hypertable('fhir.observations_ts', 'recorded_at');
```

---

## Links

- PostgreSQL 18 release: https://www.postgresql.org/about/news/
- TimescaleDB 2.25: https://github.com/timescale/timescaledb/releases
- pgx v5: https://github.com/jackc/pgx
