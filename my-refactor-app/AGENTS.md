# 🛡️ Agent Project Constitution (`my-refactor-app`)

## 1. Project-Scoped Operational Boundaries
- ALL operations (reading, writing, modifying code) MUST be confined strictly within this project repository root (`my-refactor-app`).
- NEVER execute destructive commands (e.g. `rm -rf /`, `sudo`, `chmod 777`).

## 2. SDLC Harness Execution Protocol
- Always inspect `.harness/state.json` via `create-harness-app` CLI / MCP to determine the active SDLC phase and pending node.
- Follow the corresponding checklist in `docs/<phase>/SKILL.md` before marking any node as completed.
- Validate deliverables against Definition of Done (DoD) before triggering state updates.

## 3. Allowed Project Commands
- `create-harness-app`
- `go test ./...`
- `make test`
- `make e2e`
