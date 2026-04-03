---
trigger: always_on
---
# ADR-0007 — Argo CD 3.3.x for GitOps

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @DevOps-Ariful-Islam

---

## Context and Problem Statement

ZarishSphere is governed via GitHub. Every deployment must flow from a Git commit — no manual `kubectl apply`, no manual Helm installs. Deployments must be automated, auditable, and self-healing.

---

## Decision Drivers

- **GitOps model** — merging a PR to `main` automatically deploys
- **No-coder operation** — platform owners don't need terminal access to deploy
- **CNCF graduated** — production-grade, community-governed
- **App-of-apps pattern** — one ApplicationSet manages all 200+ services
- **Rollback** — automatic rollback on failed health checks

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **Argo CD 3.3.x** | CNCF graduated, App-of-apps, Helm 3.17, auto-sync, free | Learning curve for complex ApplicationSets |
| Flux v2 | CNCF graduated, Git-native | Less visual feedback, weaker UI |
| Jenkins X | CI+CD combined | Complex, heavyweight |
| Spinnaker | Enterprise features | Complex, Java-heavy |
| manual Helm | Simple | Not GitOps, drift possible |

---

## Decision Outcome

**Chosen option: Argo CD 3.3.x**

Argo CD's ApplicationSet controller enables one configuration to manage all 200+ ZarishSphere services. The web UI provides visual deployment status without requiring terminal access. Auto-sync with self-heal ensures the cluster always matches what's in Git.

---

## Links

- Argo CD: https://argo-cd.readthedocs.io/
- GitHub: https://github.com/argoproj/argo-cd
