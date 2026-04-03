---
trigger: always_on
---
# ADR-0009 — IBM Carbon Design System 11.x

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **Authors:** @code-and-brain

---

## Context and Problem Statement

ZarishSphere's clinical web interfaces need a design system that is:
- Accessible (WCAG 2.2 AA) — health data is used by people with disabilities
- Healthcare-appropriate — serious, clear, clinical
- Free and open source
- React 19 compatible

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **Carbon DS 11.x** | WCAG 2.2 AA, IBM healthcare focus, React 19, Apache 2.0 | IBM-flavoured, verbose |
| Material UI | Popular, large ecosystem | Google-flavoured, not healthcare-focused |
| Shadcn/ui | Flexible, accessible | Requires significant customisation |
| Ant Design | Feature-rich | Chinese tech company, not healthcare-focused |
| Custom | Full control | Months of work, accessibility risk |

---

## Decision Outcome

**Chosen option: Carbon Design System 11.x**

IBM Carbon is designed for enterprise data-heavy applications, making it ideal for clinical dashboards. It is WCAG 2.2 AA compliant out of the box and is Apache 2.0 licensed. IBM uses it for Watson Health products.

---

## Links

- Carbon DS: https://carbondesignsystem.com/
- GitHub: https://github.com/carbon-design-system/carbon
