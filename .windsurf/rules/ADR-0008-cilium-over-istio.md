---
trigger: always_on
---
# ADR-0008 — Cilium over Istio for Networking

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @DevOps-Ariful-Islam, @code-and-brain

---

## Context and Problem Statement

ZarishSphere's microservices need network-level security policies (zero-trust) and observability. The solution must work on low-resource clusters (Raspberry Pi 5, ARM64) without excessive overhead.

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **Cilium 1.17+** | eBPF-native, 60% less overhead than Istio, CNCF, no sidecar | Newer, kernel version requirement |
| Istio | Mature, widely adopted | 2× memory overhead per pod, sidecar injection complexity |
| Linkerd | Lightweight, Rust proxy | Less features than Istio |
| Raw NetworkPolicy | Simple | No mTLS, no observability |

---

## Decision Outcome

**Chosen option: Cilium 1.17+**

Cilium uses eBPF instead of sidecar proxies. This reduces per-pod overhead by ~60%, which is critical for Raspberry Pi deployments. Cilium 1.17 supports both NetworkPolicy and Hubble (observability) without Istio's complexity.

---

## Links

- Cilium: https://cilium.io/
