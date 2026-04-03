---
trigger: always_on
---
# ADR-0014 — Flutter 3.41.2 for Mobile Applications

> **Status:** ✅ Accepted
> **Date:** 2026-03-24
> **Authors:** @code-and-brain, @BGD-Health-Program

---

## Context and Problem Statement

ZarishSphere needs mobile apps for:
- Community Health Workers (offline household visits)
- Clinicians (offline patient consultations)
- Supervisors (field monitoring)
- Patients (appointment and record access)

These must run on Android (dominant in target countries) and iOS from a single codebase.

---

## Decision Outcome

**Chosen option: Flutter 3.41.2 + Dart 3.11**

Flutter:
- Single codebase for Android + iOS
- Dart 3.11 supports dot shorthands and is null-safe
- PowerSync 1.x Flutter SDK for offline sync to PostgreSQL 18.3
- Flutter 3.41.2 quarterly stable release cycle — predictable upgrades
- Riverpod 2.6 for testable, null-safe state management
- 90fps performance on low-end Android devices (target: Android 10+, 2GB RAM)

---

## Links

- Flutter: https://flutter.dev/
- Flutter 3.41.2 release notes: https://docs.flutter.dev/release/release-notes
