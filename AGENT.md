# Trustica – Agent Context

This file defines the permanent context for the Trustica project.
All contributions and changes must respect the principles and constraints defined here.

---

## Project Identity

Trustica is a certificate trust orchestration tool for developer environments.
It configures CA trust for tools and runtimes without modifying the operating system trust store.

Trustica focuses on:
- Explicit trust configuration
- Environment-aware behavior
- Idempotent operations
- Verifiable results

---

## Non-Goals

Trustica must NOT:
- Modify OS trust stores
- Guess behavior based on binaries or PATH
- Intercept or proxy network traffic
- Auto-fetch certificates
- Act as a system security tool

---

## Architecture

Trustica uses Hexagonal Architecture (Ports & Adapters).

- The Core contains pure domain logic
- Environment-specific logic lives in adapters
- The CLI is an adapter
- Adapters never talk to each other
- The Core never touches the OS directly

---

## Core Concepts

- Environment: a logical trust domain
- Trust Target: an environment with detect/apply/verify
- Registry: catalog of known environments

---

## MVP Scope

Initial supported environments:
- OpenSSL (system-level usage only)
- Python (certifi-based)
- Node.js
- Git

---

## Design Principles

- Explicit over implicit
- Detection over guessing
- Idempotency is mandatory
- Failures must be isolated
- Trust must be verifiable

---

## AI Collaboration Rules

When generating code or design:
- Never violate non-goals
- Never introduce OS-level mutations
- Prefer clarity over cleverness
- Ask before expanding scope
