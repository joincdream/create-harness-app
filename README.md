# create-harness-app 🚀

> **Declarative SDLC Harness & Specification Scaffolding CLI for AI Pair-Programming**

`create-harness-app` is an ultra-fast, zero-dependency Go CLI tool designed to control the non-determinism and technical debt of AI coding (Vibe Coding). It automatically scaffolds a **4-phase Software Development Life Cycle (SDLC) specification and deterministic guardrail harness** into your project environment with a single command.

---

## 🌟 Key Features

* **⚡ Go 1.22 Native & Zero Dependency**: Compiles into a single static binary with `<0.001s` startup time and minimal footprint.
* **📦 `go:embed` 3-Stage Fallback Template Resolver**:
  1. `--template [path]` CLI flag or explicit local path override
  2. `~/.config/create-harness-app/templates/` user custom override
  3. Binary-embedded default templates fallback (`//go:embed`)
* **🧩 Declarative JSON Blueprints**:
  * Modeled via simple `blueprint.json` manifest: **Directory = Workflow (Stage)**, **File = Node (Job)**, **AGENTS.md = Guardrail (Rule)**.
* **🧙‍♂️ Vite-style Interactive Wizard**:
  * Runs a friendly interactive prompt when executed without arguments, while fully supporting non-interactive one-liners for CI/CD pipelines.

---

## 🚀 Quick Start

### 1. Build and Run
```bash
# Build binary
go build -o create-harness-app ./cmd/create-harness-app

# 1) Interactive Wizard Mode
./create-harness-app

# 2) Non-interactive One-liner Mode
./create-harness-app my-app --template default
```

### 2. Run TDD Unit Tests
```bash
go test ./... -v
```

---

## 📁 Generated SDLC Directory Architecture

```text
my-app/
├── docs/                               # [SDLC Documents & Specification Directory]
│   ├── 01_planning/                    #   - [Phase 1: Planning & Spec Stage]
│   │   ├── 00_raw_inputs/              #     * Raw Inputs (PDF, Meeting Notes, Emails)
│   │   └── 01_user_requirements.md     #     * User Requirements Spec
│   └── 02_architecture_and_design/     #   - [Phase 2: Architecture & Design Stage]
│       └── 05_tech_stack_and_skills.md #     * AGENTS.md Immutable Guardrail Rules
```

---

## 📜 License

MIT License
