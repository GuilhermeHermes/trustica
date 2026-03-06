# Trustica Architecture

This document describes the internal architecture of the Trustica project.
It explains how responsibilities are divided, how components interact, and which constraints must be respected when extending the system.

---

## Architectural Style

Trustica follows **Hexagonal Architecture (Ports & Adapters)**.

This choice is driven by the nature of the problem:
- Multiple external systems with different trust models
- A stable core domain that must not depend on tools or OS details
- The need to add new environments without rewriting orchestration logic

Hexagonal architecture allows Trustica to:
- Isolate environment-specific behavior
- Keep domain logic pure and testable
- Prevent accidental coupling to the operating system or tooling

---

## High-Level Overview
+------------------+
| CLI | (Inbound Adapter)
+--------+---------+
|
+--------v---------+
| Core | (Domain / Orchestrator)
+--------+---------+
|
+--------v--------------------------------+
| Adapters                                   |        |      |     |        |
| ------------------------------------------ | ------ | ---- | --- | ------ |
| OpenSSL                                    | Python | Node | Git | JVM... |
| +----------------------------------------+ |        |      |     |        |

## Core Domain

### Responsibilities

The Core is responsible for:
- Orchestrating the trust workflow
- Managing execution state
- Aggregating results
- Enforcing domain rules

The Core does **not**:
- Execute shell commands
- Access the filesystem
- Read or write environment variables
- Know how any specific tool works

---

### Core Concepts

#### Environment

An **Environment** is a logical trust domain that defines:
- How certificate trust is configured
- Where certificates are installed
- How trust can be verified

An environment is not a binary or executable.
It is a conceptual runtime or tooling context.

---

#### Trust Target

A **Trust Target** is an environment that provides:
1. Detection logic
2. Application logic
3. Verification logic

Each trust target encapsulates all knowledge required to configure trust for that environment.

---

#### Environment State

Each environment transitions through explicit states:

- `NotDetected`
- `Detected`
- `Applied`
- `AlreadyTrusted`
- `Failed`
- `Skipped`

State transitions are controlled by the Core.

---

## Ports (Contracts)

Ports define **what the Core expects**, not how it is implemented.

### Environment Port

Every environment adapter must implement the following conceptual contract:

- **Detect**  
  Determines whether the environment exists and is applicable.

- **Apply**  
  Applies the provided CA certificate to the environment.

- **Verify**  
  Confirms that TLS trust works correctly after application.

The Core treats all environments uniformly via this port.

---

### System Port

The System Port abstracts interactions with the host system, such as:
- Filesystem access
- Environment variables
- Process execution

Adapters may depend on the System Port.
The Core must never depend on it.

---

## Adapters

### Environment Adapters

Each supported environment is implemented as a dedicated adapter.

An adapter:
- Knows how its environment works
- Implements Detect / Apply / Verify
- Operates independently from other adapters

Adapters must:
- Be self-contained
- Be idempotent
- Fail gracefully

Adapters must not:
- Call other adapters
- Modify OS trust stores
- Maintain global state

---

### CLI Adapter

The CLI is an inbound adapter.

Responsibilities:
- Parse user input
- Translate flags into Core commands
- Render Core results for humans

The CLI:
- Contains no business logic
- Makes no trust decisions
- Never applies certificates directly

---

## Registry

The Registry is a catalog of known environment adapters.

Responsibilities:
- Declare which environments Trustica supports
- Provide adapters to the Core during discovery

The Registry:
- Is static and explicit
- Does not perform detection or execution
- Does not contain environment logic

---

## Execution Flow

1. CLI receives user input (certificate path, options)
2. CLI invokes the Core
3. Core retrieves adapters from the Registry
4. Core runs `Detect` on each adapter
5. Core builds an execution plan
6. Core runs `Apply` per applicable environment
7. Core runs `Verify` per applied environment
8. Core aggregates results
9. CLI renders output

Failures in one environment do not affect others.

---

## Idempotency Model

Trustica enforces idempotency at the environment level.

Running the same operation multiple times must:
- Produce the same result
- Never duplicate certificates
- Never corrupt configuration

Adapters must explicitly detect already-trusted states.

---

## Error Handling

- Errors are scoped to the environment where they occur
- The Core never aborts execution due to a single failure
- Errors are reported as structured results

---

## MVP Environments

The initial implementation includes:

- OpenSSL (system-level usage only)
- Python (certifi-based)
- Node.js
- Git

Each environment is implemented as an independent adapter.

---

## Architectural Constraints

The following constraints are non-negotiable:

- No OS trust store modification
- No PATH scanning or binary inference
- No cross-adapter communication
- No hidden side effects
- No implicit trust assumptions

---

## Extensibility

New environments can be added by:
- Implementing a new adapter
- Registering it in the Registry

No Core changes should be required.

This ensures long-term scalability without architectural drift.

---

## Summary

Trustica’s architecture enforces a strict separation between:
- **What** trust means (Core)
- **How** trust is applied (Adapters)
- **How** users interact (CLI)

This separation is essential to maintain correctness, safety, and extensibility as the number of supported environments grows.