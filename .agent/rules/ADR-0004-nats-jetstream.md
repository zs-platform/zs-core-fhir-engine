---
trigger: always_on
---
# ADR-0004 — NATS 2.12.5 JetStream for Messaging

> **Status:** ✅ Accepted
> **Date:** 2026-03-16
> **RFC:** RFC-0002
> **Authors:** @code-and-brain, @DevOps-Ariful-Islam

---

## Context and Problem Statement

ZarishSphere services need to communicate asynchronously. When a patient is created, the search indexer must be updated, an AuditEvent must be written, and subscribed FHIR clients must be notified. This communication must work even when some services are temporarily unavailable.

---

## Decision Drivers

- **Zero cost** — must be self-hosted with no licensing fees
- **Lightweight** — runs alongside other services on Raspberry Pi 5
- **FHIR R5 subscriptions** — FHIR R5 topic-based subscriptions via JetStream
- **Durable consumers** — messages must survive service restarts
- **At-least-once delivery** — no clinical events may be silently dropped
- **JetStream persistence** — messages persisted to disk for audit trail

---

## Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **NATS 2.12.5 JetStream** | 20MB binary, JetStream persistence, FHIR subscription support, free | Smaller ecosystem than Kafka |
| Apache Kafka | Industry standard, huge ecosystem | 500MB+ deployment, ZooKeeper/KRaft complexity, overkill for our volume |
| RabbitMQ | AMQP standard, mature | Erlang runtime, complex clustering, 150MB |
| Google PubSub | Managed, scalable | Cost at scale, vendor lock-in |
| Amazon SQS/SNS | Managed | Cost, AWS vendor lock-in |
| Redis Pub/Sub | Simple | No persistence, messages lost on restart |

---

## Decision Outcome

**Chosen option: NATS 2.12.5 JetStream**

NATS JetStream satisfies every constraint:
1. Single 20MB binary — runs on Raspberry Pi 5
2. JetStream persistence — messages survive service restarts
3. At-least-once delivery — configurable acknowledgement
4. FHIR R5 topic subscriptions — maps directly to NATS subjects
5. Dead letter queue — failed messages to `zs.dlq.*` subjects
6. Zero licence cost, Apache 2.0

## NATS Subject Taxonomy

```
# FHIR resource events
zs.fhir.{ResourceType}.created
zs.fhir.{ResourceType}.updated
zs.fhir.{ResourceType}.deleted

# Clinical alerts
zs.alert.ewars.{alertType}
zs.alert.lab.critical

# Platform events
zs.platform.deploy.{service}
zs.platform.health.{service}

# Dead letter queue
zs.dlq.{originalSubject}
```

---

## Links

- NATS documentation: https://docs.nats.io/
- JetStream guide: https://docs.nats.io/nats-concepts/jetstream
- nats.go client: https://github.com/nats-io/nats.go
