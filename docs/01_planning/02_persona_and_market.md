# 02. 타겟 페르소나 및 시장 분석 (Persona & Market Analysis)

## 1. 타겟 유저 페르소나 (Target Persona)

### 페르소나 A: Vibe Coding의 기술 부채에 직면한 시니어 개발자 / 아키텍트 (민우, 36세)
* **배경**: AI 에이전트(Cursor, Roo Code 등)를 도입하여 빠르게 코딩하고 있으나, 에이전트가 만든 코드가 아키텍처 규칙을 위반하고 장애 발생 시 디버깅이 불가능한 'Vibe Slop' 현상으로 스트레스를 받음.
* **니즈**: 별도의 무거운 프레임워크나 MCP를 배우느라 시간을 쏟지 않고, 팀의 기존 TDD, OpenAPI, Makefile 규칙을 에이전트가 즉시 준수하도록 통제 틀(Harness)을 갖추고 싶어함.

### 페르소나 B: AI 파이프라인 체계를 도입하려는 기술 이사 / B2B 리더 (지은, 42세)
* **배경**: 개발팀에 AI 도입을 장려하고 있으나, 개발자들이 감(Vibe)으로 작성한 코드의 프로덕션 장애와 오염에 대한 거버넌스 공포를 지님.
* **니즈**: 기획부터 배포까지 SDLC 표준 산출물이 디렉토리별로 축적되어, AI가 만든 스펙과 코드가 100% 추적 가능한 무결성 파이프라인 체계를 원함.

---

## 2. 시장 분석 및 핵심 차별성 (Market & Competitive Advantage)

1. **기존 시장의 문제 (Harness Over-engineering Purgatory)**:
   * '하네스/가드레일'을 새로운 마법 솔루션처럼 마케팅하며, 거대한 룰셋 작성이나 MCP 튜닝에 과도하게 리소스를 쏟게 만드는 오버엔지니어링 현상 만연.
2. **`create-harness-app`의 핵심 차별성 (Competitive Edges)**:
   * **CLI 기반 극상의 실용성**: `npx create-harness-app` 단 한 줄로 친숙하게 시작.
   * **소프트웨어 공학의 하네스화**: 기존 TDD, Linter, OpenAPI, Docker, Makefile을 에이전트 피드백 루프로 직접 연결.
   * **SDLC Compression**: 문서/스펙 기반 디렉토리 체인(Chain)을 통해 원시 요구사항부터 테스트/배포까지 AI로 속도를 폭발적으로 가속화.
