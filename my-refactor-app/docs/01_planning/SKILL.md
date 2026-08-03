# 📝 Phase 1: Cyclomatic Complexity & Static Analysis Skill Manual

> **Phase**: 1 (Planning & Sensing)  
> **Node ID**: `code-analysis`  
> **Skill File**: `docs/01_planning/SKILL.md`

---

## 1. Role & Objective
Go 모듈 내 소스 코드를 정적 분석하여 순환 복잡도(Cyclomatic Complexity) 10 이상의 레거시 메서드를 탐지한다.

## 2. Guardrails & Principles (불변 제약사항)
- `src/` 및 기존 Go 패키지 레이아웃을 임의로 변경하거나 손대지 말 것.
- 정적 분석 단계에서는 기존 인터페이스(Struct Interface) 수정을 금지함.

## 3. Inputs & Outputs
- **Inputs**: `./src/**/*.go` (또는 `./internal/service/*.go`)
- **Outputs**: `docs/01_planning/analysis_report.md`

## 4. Definition of Done (DoD) & Status Transition
1. 감지된 고복잡도 메서드 분석 리포트를 `docs/01_planning/analysis_report.md` Markdown 파일로 작성한다.
2. 완료 후 MCP 도구를 호출하여 노드를 완료 처리한다:
   `harness_update_state(node_id="code-analysis", status="completed")`
