---
trigger: always_on
---
# ADR-0010 — Cloudflare as CDN/DNS/Edge Layer

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **Authors:** @DevOps-Ariful-Islam

---

## Context and Problem Statement

ZarishSphere needs DNS, CDN, DDoS protection, edge compute, email routing, and object storage — all for zero cost.

---

## Decision Outcome

**Chosen option: Cloudflare Free Tier**

Cloudflare's free tier includes:
- DNS (unlimited, instant propagation)
- CDN (unlimited bandwidth)
- DDoS protection (included)
- SSL/TLS (auto-renew, free)
- Workers (100,000 req/day free)
- Pages (unlimited deployments free for OSS)
- R2 storage (10 GB free, zero egress cost)
- Email Routing (free)

No other provider offers this breadth at zero cost.

---

## Links

- Cloudflare: https://cloudflare.com/
