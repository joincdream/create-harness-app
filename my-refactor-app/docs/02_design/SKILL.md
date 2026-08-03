# 🛡️ Phase 2: Clean Architecture & SOLID Refactoring Skill Manual

> **Phase**: 2 (Design & Refactoring)  
> **Node ID**: `clean-arch-refactor`  
> **Skill File**: `docs/02_design/SKILL.md`

---

## 1. Role & Objective
1단계에서 생성된 `docs/01_planning/analysis_report.md` 리포트를 읽고, 클린 아키텍처 및 SOLID 원칙을 적용한 Diff 패치를 생성한다.

## 2. Guardrails & Principles (불변 제약사항)
- **단일 책임 원칙 (SRP)**: 거대 서비스 함수를 전용 UseCase 인터페이스로 분리할 것.
- **의존성 역전 원칙 (DIP)**: 비즈니스 로직이 구체 구조체가 아닌 인터페이스에 의존하게 만들 것.
- **행위 불변 가드레일**: 기존 외부 API 시그니처 및 입출력 규약을 100% 보존할 것.

## 3. Inputs & Outputs
- **Inputs**: `docs/01_planning/analysis_report.md`
- **Outputs**: `docs/02_design/refactor_diff.patch`

## 4. Definition of Done (DoD) & Status Transition
1. 리팩토링 변경사항을 `docs/02_design/refactor_diff.patch` 파일로 생성한다.
2. 완료 후 MCP 도구를 호출하여 노드를 완료 처리한다:
   `harness_update_state(node_id="clean-arch-refactor", status="completed")`
