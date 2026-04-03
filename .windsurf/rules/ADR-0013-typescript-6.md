---
trigger: always_on
---
# ADR-0013 — TypeScript 6.0 for All Frontend Code

> **Status:** ✅ Accepted
> **Date:** 2026-03-24
> **Authors:** @code-and-brain

---

## Context and Problem Statement

ZarishSphere's 43 React microfrontends and 8 TypeScript packages need a consistent, well-typed language. TypeScript 6.0 (released 2026-03-06) is the current stable release.

---

## Decision Outcome

**Chosen option: TypeScript 6.0 (strict mode)**

TypeScript 6.0 is the last JavaScript-based TypeScript release before a planned Go/Rust rewrite of the compiler. It is the most stable, most compatible version. Strict mode enabled by default catches an entire class of runtime errors at compile time — essential for clinical software.

Key rules:
- `strict: true` in all tsconfig.json
- Zero `any` types — use `unknown` and narrow
- All React components fully typed — no implicit `children: any`

---

## Links

- TypeScript 6.0: https://devblogs.microsoft.com/typescript/
