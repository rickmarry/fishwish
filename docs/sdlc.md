# Software Development Lifecycle

This document defines how work gets designed, built, and validated in FishFinder. The goal is to prevent intent drift — the gap between what was designed and what gets built.

---

## The Problem: Intent Drift

Intent drift happens when:
- Implementation diverges quietly from the design
- Small "while I'm in here" changes accumulate without being reflected back in specs
- Tests pass but don't actually assert the contracts that matter
- The spec gets written and never consulted again

The practices below exist specifically to prevent this.

---

## The Flow

Every non-trivial piece of work follows this sequence. Do not skip layers.

```
HLD → LLD → SPEC → Implementation → Tests enforce contracts
```

### When to apply this

Apply the full flow for:
- Any new service
- Any new API endpoint or data model
- Any significant new feature within an existing service
- Any change to a cross-service data flow

Do **not** apply for:
- Bug fixes
- Infrastructure/ops work (getting the stack running, fixing compatibility issues)
- Developer tooling (smoke tests, dev scripts)

---

## Each Layer

### HLD — High-Level Design

Captured in `docs/architecture.md` for system-level decisions, or `services/<name>/docs/service.md` for service-level design.

Answers: What does this do? How does it fit into the overall system? What are the key design decisions?

### LLD — Low-Level Design

Captured in `services/<name>/DESIGN.md`.

Answers: How exactly does this work internally? What are the data models, the state machines, the failure modes?

### SPEC — Contracts and Acceptance Criteria

Captured in `services/<name>/SPEC.md`.

Answers: What are the exact inputs and outputs? What are the API request/response shapes? What does "done" look like?

API contracts (request/response structs) are defined in `internal/model/` and are the source of truth for service interfaces. Never define API shapes only in handler code.

**SPEC sections are the backlog.** Each section in a SPEC.md carries a status tag that is kept current as work progresses:

- `[PENDING]` — not yet started
- `[IN PROGRESS]` — actively being implemented
- `[COMPLETE]` — implemented and verified against the spec

When picking up work in the code session, scan the SPEC for the next `[PENDING]` section. When done, update the tag. The SPEC is always the source of truth for what's built and what isn't — not a separate task tracker.

**Requirements use EARS format.** Write feature requirements in EARS (Easy Approach to Requirements Syntax) before implementation begins. EARS templates make requirements unambiguous:

- **Ubiquitous:** "The \<system\> shall \<action\>."
- **Event-driven:** "WHEN \<trigger\> the \<system\> shall \<action\>."
- **State-driven:** "WHILE \<state\> the \<system\> shall \<action\>."
- **Conditional:** "WHERE \<condition\> the \<system\> shall \<action\>."
- **Option:** "WHERE \<feature\> is included the \<system\> shall \<action\>."

EARS requirements belong in `SPEC.md` before any implementation starts.

### Implementation

Build to the spec. If you discover during implementation that the spec needs to change, update the spec first, then implement. Never let the implementation quietly diverge.

### Tests Enforce Contracts

Tests are the mechanical enforcement layer. Without them, the discipline is only as strong as human attention.

**Layers:**
- **Smoke test** — coarse end-to-end happy path
- **Contract tests** — assert that API interfaces match the model definitions exactly
- **Unit tests** — internal logic within a service

**The acceptance criterion for a feature is the SPEC.md, not "it works."** When something is done, verify it against the spec explicitly.

---

## Architecture Decision Records

Non-obvious system-level decisions go in `docs/adr/`. Use the template at `docs/adr/template.md`.

Rules:
- Write an ADR before making a significant architectural decision, not after
- Never edit an ADR — write a new one that supersedes it
- **Always reference related ADRs** in the "Context" or a "References" section to ensure architectural traceability
- Service-level decisions go in `services/<name>/docs/service.md`, not in ADRs

---

## GitHub Flow

All work follows this flow regardless of solo/team context. The PR history is the project journal.

1. Create a GitHub issue describing the work
2. Create a feature branch (`feat/`, `fix/`, `docs/`, etc.)
3. Commit with conventional commit messages: `type(scope): description`
4. Open a PR referencing the issue
5. Merge with `gh pr merge --squash`

**Never commit directly to main.**

Commit types: `feat`, `fix`, `refactor`, `docs`, `infra`, `ci`, `test`, `chore`

---

## What Claude Enforces

- If work moves toward implementation without a DESIGN.md and SPEC.md for the service, flag it
- If a new service is built without a `services/<name>/docs/service.md`, flag it
- If a commit is about to go directly to main, flag it
