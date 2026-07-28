# create-harness-app 기획 초안 (Planning Draft)

> **문서 출처 및 기반**: [posts/deep-dive/2026-07-28-harness-engineering-vibe-coding.md](file:///home/yundream/myjob/cloit/ai-info/posts/deep-dive/2026-07-28-harness-engineering-vibe-coding.md)

---

## 1. 프로젝트 개요 & 목적 (Overview & Goals)

* **프로젝트명**: `create-harness-app`
* **핵심 비전**: 감(Vibe)에 의존하는 AI 코딩의 비결정성(Non-determinism)과 기술 부채를 통제하기 위해, **기존 소프트웨어 공학(SDLC) 기반의 검증 센서와 명세(Spec) 체계를 명령 한 줄로 프로젝트 초기 환경에 스캐폴딩(Scaffolding)하는 CLI 툴**을 구축합니다.
* **핵심 타겟**:
  * AI 에이전트를 도입하여 생산성을 극대화하고자 하는 시니어/주니어 엔지니어 및 팀.
  * 복잡한 MCP 구축이나 거대한 가드레일 작성(오버엔지니어링) 없이, 검증된 소프트웨어 공학 도구(TDD, Linter, OpenAPI, Makefile)를 AI 피드백 루프에 즉시 물리고자 하는 조직.

---

## 2. 배경 및 문제 의식 (Background & Problem Statement)

1. **Vibe Coding의 위험**: 통제 장치 없는 AI 코딩은 맥락 유실, 할루시네이션, 아키텍처 오염 및 장애(Vibe Slop)를 초래함.
2. **트렌드의 현혹과 오버엔지니어링**: 하네스라는 버즈워드에 얽매여 거대한 마법 프레임워크나 프롬프트 튜닝 연옥(Architecture Astronauts)에 리소스를 낭비하는 경향 발생.
3. **거장들의 본질 정의 (Martin Fowler & Kent Beck)**:
   * **`Agent = Model + Harness`**
   * 하네스의 본질은 거창한 신기술이 아니라, **"이미 검증된 기존의 TDD, OpenAPI 명세, Linter, CI/CD 스크립트 등 결정론적 검증 센서를 에이전트에 가두어 연결하는 것"**임.
4. **해법 (SDLC Compression)**: 개발자가 일일이 수동으로 가드레일을 깎지 않고, CLI 도구를 통해 디렉토리 기반 스펙 체인(Chain)을 자동 구축하여 **원시 요구사항부터 테스트/배포까지 개발 주기를 AI로 압착(Compression)**함.

---

## 3. 기능 요구사항 리스트 (Functional Requirements)

* **[FR-01] 대화형 CLI 실행 스캐폴더**:
  * `npx create-harness-app [project-name]` 명령 단 한 줄로 동작해야 함.
* **[FR-02] SDLC 4단계 디렉토리 자동 구성**:
  * `01_planning/`: 원시 자료(`00_raw_inputs/`), 유저 요구사항, 페르소나, 수락 조건.
  * `02_architecture_and_design/`: 상세 기능 명세, DTO, OpenAPI 3.0 yaml, 시퀀스 다이어그램, `AGENTS.md`.
  * `03_development_cycle/`: 환경/하네스 구축(Boilerplate, Docker, Makefile) 및 백엔드/프론트엔드 TDD 개발.
  * `04_qa_and_verification/`: 스펙 추적 QA Matrix, E2E 시나리오 테스트, 릴리즈 체크리스트.
* **[FR-03] AI 에이전트용 불변 가이드라인 자동 주입**:
  * 각 디렉토리마다 에이전트가 읽고 준수해야 할 마크다운 규칙 및 프롬프트 샌드박스 가이드라인 템플릿 포함.
* **[FR-04] 가드레일 문서 기반 Input ➔ Output 연쇄 파이프라인 지원**:
  * 사전 주입된 가드레일 문서(SOLID, 계층 구조, DI, 미들웨어 패턴)를 지시 소스로 삼아, 개발자의 한 줄 지시만으로 Well-Architected 스펙/코드가 연쇄 생성되는 워크플로우 보장.

---

## 4. 수락 조건 (Acceptance Criteria)

* **Given**: 개발자가 터미널에서 `create-harness-app` 명령을 실행할 때
* **When**: 프로젝트명을 입력하거나 현재 디렉토리를 지정하면
* **Then**:
  1. 1초 이내에 4단계 SDLC 계층 트리 디렉토리가 완전 생성된다.
  2. 10개 이상의 하네스 가이드라인 및 스펙 템플릿 마크다운 파일이 올바르게 주입된다.
  3. `01_planning/00_raw_inputs/` 디렉토리가 즉시 원시 자료(PDF, 회의록)를 받을 준비가 된다.
