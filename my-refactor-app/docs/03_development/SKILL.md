# 🧪 Phase 3: Regression Test & Quality Gate Skill Manual

> **Phase**: 3 (Testing & Verification)  
> **Node ID**: `tdd-verification`  
> **Skill File**: `docs/03_development/SKILL.md`

---

## 1. Role & Objective
2단계에서 생성된 `docs/02_design/refactor_diff.patch`를 적용하고 단위 테스트 및 회귀 테스트를 실행하여 결함이 없음을 검증한다.

## 2. Guardrails & Principles (불변 제약사항)
- 모든 단위 테스트가 성공(PASS)해야만 릴리즈 승인 게이트를 통과시킨다.
- 테스트 커버리지가 떨어지거나 빌드 에러 발생 시 상태를 `blocked` 처리한다.

## 3. Inputs & Outputs
- **Inputs**: `docs/02_design/refactor_diff.patch`
- **Outputs**: `docs/03_development/test_summary.log`

## 4. Definition of Done (DoD) & Status Transition
1. `go test ./...` 실행 결과 및 로그를 `docs/03_development/test_summary.log`로 저장한다.
2. 완료 후 MCP 도구를 호출하여 노드를 완료 처리한다:
   `harness_update_state(node_id="tdd-verification", status="completed")`
