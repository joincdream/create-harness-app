# Antigravity 2.0 Global Harness Guardrails

이 프로젝트는 `create-harness-app` SDLC 런타임 하네스 통제 하에 개발됩니다.

## 🚨 필수 행동 지침 (Mandatory Rules)

1. **상태 관리 툴 호출 의무**:
   각 단계를 수행하고 산출물 파일 생성이 완료되면, 반드시 `create-harness-app state update --node <id> --status completed` 명령을 실행하여 단계를 완료 처리하십시오.

2. **근거 기반 개발 (Data-driven Development)**:
   선행 Phase의 산출물 문서(Input Specs)를 읽지 않고 추측하여 후속 문서를 작성하는 행위를 엄격히 금합니다.

3. **스킬 가이드 준수**:
   새로운 Phase 진행 시 해당 스킬 문서(`.agents/skills/phase_N/SKILL.md`)를 읽고 그에 정의된 입출력 경로 및 검증 규칙을 따라야 합니다.
