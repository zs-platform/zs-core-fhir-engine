---
trigger: always_on
---
# ADR-0011 — Microfrontend Architecture (Next.js 16.2)

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **Authors:** @code-and-brain

---

## Context and Problem Statement

ZarishSphere's frontend spans: clinical workflows, public health dashboards, operational tools, and patient portal. Different teams (clinical, public health, ERP) must be able to deploy independently without coordinating releases.

---

## Decision Outcome

**Chosen option: Microfrontend Architecture via Vite Module Federation + Next.js 16.2**

Each functional area (patient registration, encounters, pharmacy, etc.) is an independently deployable React application. Shell apps compose them via Module Federation. This enables:
- Independent deployment per feature area
- Independent team ownership
- Independent testing
- Shared component library (`zs-pkg-ui-design-system`)

---

## Links

- Next.js 16.2: https://nextjs.org/
- Module Federation: https://module-federation.io/
