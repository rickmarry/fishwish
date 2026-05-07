# AI Assistant Workflow Rules

This file contains workflow rules and conventions that apply to **all AI assistants** working on this project (Claude, Gemini, Junie, etc.). These rules ensure consistency across tools and maintainers.

---

## Git Workflow — Audit Trail Required

When making code changes that need to be pushed:

1. **Create a GitHub issue** describing the work
2. **Create a feature branch** from main (use conventional naming: `type/short-description`)
3. **Make the code changes** and commit (conventional commits format)
4. **Push the branch** to origin
5. **Create a PR** that references the issue (use `Closes #N` in the PR body)
6. **Merge the PR** — use `gh pr merge <number> --squash --delete-branch` (auto-merge, no manual gate)
7. **Checkout main and pull** — `git checkout main && git pull`

**Why this workflow:** Every change has a traceable issue → branch → PR → merge path. This serves as a complete audit trail for the project's evolution and ensures all work is documented and reviewable.

---

## Commit Message Conventions

All commits follow **conventional commits** format:

```
type(scope): description
```

**Examples:**
- `feat(spots): add spot creation endpoint`
- `fix(auth): resolve JWT context key injection`
- `docs(adr): add decision record for PostGIS choice`

**Valid types:** `feat`, `fix`, `refactor`, `docs`, `infra`, `ci`, `test`, `chore`

**IMPORTANT:** Never add a `Co-Authored-By` line to commits. Rick owns this IP entirely.

---

## Design Discipline — SDLC

This project follows a strict spec-driven development process to prevent intent drift:

```
HLD → LLD → Specs/Contracts → Implementation → Tests enforce contracts
```

### The Design Gate — Hard Stop

**Design work happens in the planning session, not here.**

The only trigger for implementation in this session is a backlog item with **status `Design Complete`** AND the following files already committed to `main`:
- `docs/bXXX-design.md` — LLD
- `docs/bXXX-spec.md` — EARS-format requirements
- `docs/adr/XXXX-bXXX-*.md` — ADR for non-obvious decisions

**A detailed backlog entry is NOT a design.** The backlog entries in `docs/backlog.md` contain pre-design notes — user stories, open questions, suggested architecture. These exist to inform the design session, not replace it. No matter how detailed the backlog entry looks, it is not a green light to build.

**If the status is `Backlog` → stop. Do not design here. Do not build. Switch back to the planning session.**

**Before implementing any new service or significant feature:**
1. Update the High-Level Design if the feature adds a new service or changes the architecture (`docs/architecture.md` or a new ADR)
2. Confirm the following are committed to `main`:
   - `docs/bXXX-design.md` — LLD
   - `docs/bXXX-spec.md` — **EARS-format** requirements (WHEN / SHALL / IF / WHERE)
   - `docs/adr/XXXX-bXXX-*.md` — ADR for non-obvious decisions
3. Define or update API contracts (request/response structs in `internal/model/`)
4. Implement against the spec
5. Write tests that enforce contract compliance

**Do not skip the gate.** If the design docs are missing, stop and surface it rather than proceeding.

---

## Backlog Management

**Build order is part of item creation — not a separate step.**
Every new backlog item must be placed in the build order the moment it's created. Review the full order, slot the item in the correct position based on dependencies and strategic value, and update any items whose position shifts as a result. Never add an item and leave it unsequenced — that intention evaporates.

"Parking lot" means genuinely unsequenced — no dependencies, no clear strategic position. It is not a holding area for items that haven't been ordered yet.

**Backlog entries must be self-contained.**
Every entry must include: user story, core loop, key design decisions, open questions, suggested architecture, and a Caveats section. A one-liner is not a backlog entry. The entry must be complete enough that a build session can start without re-deriving the thinking.

**Caveats are permanent and must be repeated.**
Every time a backlog item is discussed, its caveats must be surfaced explicitly — not assumed to be remembered. Read the Caveats section before discussing next steps.

**V1/V2 scoping — always surface V2 explicitly.**
If a feature has V1 (MVP) and V2 (future) scopes, both must appear in the backlog entry with a decision trigger for V2. Never bury V2 in a design doc. Update the backlog entry immediately after V1 ships: mark V1 complete, make V2 visible with its trigger condition.

---

## Design Session Opener — Plain Language First

Every design session must open with a plain-language summary of what we're building and why — the user problem, the proposed solution in one sentence, and why this approach over alternatives. No architecture, no schema, no proto definitions until this is stated and agreed. Jumping straight to implementation detail before the problem is clearly framed produces over-engineered solutions to the wrong problem.

---

## Architecture Decision Records (ADRs)

System-level design decisions belong in `docs/adr/`. Each ADR is immutable — never edit an ADR; write a new one that supersedes it.

**Template:** `docs/adr/template.md`
**Numbering:** `NNNN-short-title.md` (zero-padded, e.g., `0001-use-postgis-for-spatial-queries.md`)

**What goes in an ADR:**
- Context: what problem are we solving?
- Options considered
- Decision made
- Consequences (trade-offs, future implications)

**Service-level decisions** (internal to one service) go in `services/<name>/docs/service.md` under "Design Decisions."

---

## Service Documentation

Every service must have:
- `services/<name>/docs/service.md` — purpose, design decisions, API surface, config, local dev instructions
- See `services/spot-service/docs/service.md` as the reference example

---

## Repository Conventions

### Monorepo structure
All services live in `services/`. Each service has its own:
- Build tooling (`go.mod`, `package.json`, etc.)
- Dockerfile
- Documentation (`docs/service.md`)

### Local development
`docker-compose.yml` is the entry point for all infrastructure (Postgres, Redis, MinIO). Run services via `make dev` or individually from `services/<name>/`.

### Shared code
Common code lives in `pkg/`. Services import it via `replace fishwish => ../../` in their go.mod.

### Secrets
Never commit secrets. Use `.env.local` files locally (gitignored) and AWS Secrets Manager in production.

---

## Testing

- **Unit tests:** per service, run via `make test`
- **Integration tests:** verify service-to-database contracts
- **Contract tests:** enforce API request/response struct compliance (not yet implemented)

---

## AI Assistant Etiquette

- Be concise. Match verbosity to task complexity.
- Do not add preamble or postamble unless asked.
- Use the TodoWrite tool for multi-step tasks.
- Follow the git workflow for every code change, no exceptions.
- If unclear, ask. Do not guess at requirements.
