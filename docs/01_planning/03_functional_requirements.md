# 03. 기능 요구사항 및 수락 조건 (Functional Requirements & Acceptance Criteria)

## 1. 상세 기능 요구사항 (Detailed Functional Requirements)

* **[FR-01] 대화형 CLI 스캐폴더**:
  * `npx create-harness-app [target-dir]` 실행 지원.
  * 대상 디렉토리가 미지정된 경우 대화형 프로그래밍 방식으로 프로젝트명과 템플릿 옵션 입력 수용.
* **[FR-02] SDLC 4단계 디렉토리 트리 자동 자동 스캐폴딩**:
  * `01_planning/` (원시 입력, 유저 요구사항, 페르소나, 세부 기능 요구사항)
  * `02_architecture_and_design/` (상세 명세, DTO 스키마, OpenAPI 3.0 yaml, 시퀀스 다이어그램, `AGENTS.md`)
  * `03_development_cycle/` (환경/하네스 구축, TDD 자가 교정 가이드, 백엔드/프론트엔드 코드 및 단위 테스트)
  * `04_qa_and_verification/` (스펙 추적 QA Matrix, E2E 시나리오 테스트, 릴리즈 체크리스트)
* **[FR-03] AI 에이전트용 불변 가이드라인 자동 주입**:
  * 각 디렉토리 및 Job마다 에이전트가 읽고 준수해야 할 마크다운 규칙 및 프롬프트 샌드박스 가이드라인 템플릿 포함.
* **[FR-04] 가드레일 문서 기반 Input ➔ Output 연쇄 파이프라인**:
  * 사전 주입된 가드레일 문서(SOLID, 계층 구조, DI, 미들웨어 패턴)를 지시 소스로 삼아, 개발자의 한 줄 지시만으로 Well-Architected 스펙/코드가 연쇄 생성되는 워크플로우 보장.

---

## 2. 세부 수락 조건 (Explicit Acceptance Criteria - Given-When-Then)

### Scenario 1: CLI 초기 실행 및 디렉토리 스캐폴딩 검증
* **Given**: 개발자가 터미널에서 `create-harness-app my-app` 명령을 실행했을 때
* **When**: CLI 엔진이 인자를 파싱하고 템플릿 생성을 시작하면
* **Then**:
  1. 1초 이내에 `my-app/` 아래 SDLC 4단계 계층 트리가 즉시 생성되어야 한다.
  2. `01_planning/00_raw_inputs/` 디렉토리가 생성되어 원시 자료를 수용할 준비가 되어야 한다.
  3. 총 10개 이상의 하네스 가이드라인 마크다운 및 YAML 스펙 파일이 생성되어야 한다.

### Scenario 2: 가드레일 연쇄 생성 워크플로우 검증
* **Given**: `01_planning/00_raw_inputs/`에 회의록 원시 자료가 들어있고 가드레일 문서가 주입되어 있을 때
* **When**: 개발자가 AI 에이전트에 "가드레일에 따라 요구사항 명세를 생성해 줘"라고 한 줄 지시를 내리면
* **Then**: AI 에이전트가 외부로 이탈하지 않고 `01_user_requirements.md` 스펙 및 수락 조건을 정밀하게 완성해야 한다.
