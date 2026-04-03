---
trigger: always_on
---
# ADR-0006 — OpenTofu 1.11.x over Terraform

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @DevOps-Ariful-Islam

---

## Context and Problem Statement

ZarishSphere's infrastructure (Kubernetes clusters, PostgreSQL, NATS, Cloudflare DNS, Keycloak) must be defined as code and reproducible. The IaC tool must be free for all users including organisations that build services on ZarishSphere.

---

## Decision Drivers

- **Zero cost** — HashiCorp changed Terraform to BSL in 2023
- **Linux Foundation** — OpenTofu is CNCF/Linux Foundation, MPL-2.0 forever
- **Terraform compatibility** — OpenTofu is a 1:1 fork, all modules work
- **Active development** — OpenTofu 1.11.x adds features beyond Terraform OSS

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **OpenTofu 1.11.x** | MPL-2.0, Linux Foundation, Terraform-compatible | Smaller community than Terraform |
| Terraform BSL | Largest community, most modules | BSL licence restricts use in managed services |
| Pulumi | Multi-language IaC | Complex, not HCL-compatible |
| Ansible | Agentless, simple | Not declarative state management |
| raw kubectl/Helm | Simple | No state management, drift detection |

---

## Decision Outcome

**Chosen option: OpenTofu 1.11.x**

OpenTofu is the Linux Foundation's community-maintained fork of Terraform under MPL-2.0. All existing Terraform modules are directly compatible. HashiCorp's BSL makes Terraform unusable for a platform like ZarishSphere that enables third-party deployments.

---

## Links

- OpenTofu: https://opentofu.org/
- GitHub: https://github.com/opentofu/opentofu
