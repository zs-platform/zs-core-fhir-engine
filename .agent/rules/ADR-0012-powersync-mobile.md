---
trigger: always_on
---
# ADR-0012 — PowerSync for Mobile Offline Sync

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **Authors:** @code-and-brain, @BGD-Health-Program

---

## Context and Problem Statement

ZarishSphere mobile apps (CHW, clinic) must work offline for hours or days, then sync when connectivity returns. The sync mechanism must handle conflict resolution, partial sync, and security.

---

## Decision Outcome

**Chosen option: PowerSync 1.x (self-hosted)**

PowerSync provides:
- SQLite offline storage on device (Flutter drift)
- Bi-directional sync to PostgreSQL 18.3 (our database)
- Row-level sync rules (per tenant, per user role)
- Self-hosted server (zero cost, Apache 2.0)
- Conflict resolution strategies built-in

Alternatives considered: PocketBase (not PostgreSQL-native), CRDTs (complex), manual sync (unreliable).

---

## Links

- PowerSync: https://www.powersync.com/
