# Google Antigravity Code Refactoring Harness Specification

Google Antigravity 에이전트를 위한 레거시 코드 자동화 분석, 클린 아키텍처 가이드라인 적용 및 TDD 기반 검증 하네스 템플릿 명세서입니다.

---

## 1. 하네스 개요 (Overview)

본 하네스 템플릿은 기존 Go 애플리케이션의 높은 코드 복잡도와 패키지 간의 스파게티 의존성을 체계적으로 리팩토링하기 위해 설계되었습니다. 에이전트는 무작위로 코드를 수정하는 대신, 3단계의 결정적(Deterministic) 노드 파이프라인에 따라 분석, 디자인 패치 작성, 그리고 유닛 테스트 검증을 순차적으로 수행합니다.

---

## 2. 핵심 가드레일 및 아키텍처 규칙 (Architecture Rules)

### 순환 복잡도(Cyclomatic Complexity) 제한
- 순환 복잡도 10 이상의 메서드는 필수 리팩토링 대상으로 분류됩니다.
- 단일 함수 또는 메서드의 길이는 50줄 이내로 분할하여 단일 책임 원칙(SRP)을 준수해야 합니다.

### 클린 아키텍처 단방향 의존성 강제
- 도메인 엔티티 및 유즈케이스 계층은 외부 프레임워크나 DB 계층에 직접 의존할 수 없습니다.
- 외부 모듈과의 통신은 인터페이스 기반의 의존성 역전 원칙(DIP)을 적용하여 구현해야 합니다.

---

## 3. SDLC 3단계 파이프라인 노드 명세 (Pipeline Nodes)

### Phase 1: code-analysis (코드 정적 분석)
- **입력 산출물**: `./src/**/*.go`
- **출력 산출물**: `docs/01_planning/analysis_report.md`
- **주요 작업**: 전체 소스코드의 구조적 문제점, 복잡도 지수, 냄새나는 코드(Code Smells)를 탐지하여 정형화된 보고서로 작성합니다.

### Phase 2: clean-arch-refactor (클린 아키텍처 리팩토링)
- **입력 산출물**: `docs/01_planning/analysis_report.md`
- **출력 산출물**: `docs/02_design/refactor_diff.patch`
- **주요 작업**: 분석 보고서를 기반으로 클린 아키텍처 및 SOLID 원칙을 적용한 수정 사항을 `refactor_diff.patch` 형태의 차분 파일로 생성합니다.

### Phase 3: tdd-verification (TDD 단위 테스트 검증)
- **입력 산출물**: `docs/02_design/refactor_diff.patch`
- **출력 산출물**: `docs/03_development/test_summary.log`
- **주요 작업**: 패치를 적용한 후 `go test ./...`를 수행하여 기존 기능의 회귀(Regression)가 없음을 실측 검증하고 종합 로그를 기록합니다.

---

## 4. 사용 방법 (Usage)

본 하네스는 create-harness-app CLI 및 HarnessHub 레지스트리를 통해 다운로드 및 전사 공유가 가능합니다:

- **템플릿 다운로드**: `create-harness-app hub pull my-refactor-app v1.0.0`
- **프로젝트 적용**: `create-harness-app my-app --template my-refactor-app`
- **템플릿 업로드**: `create-harness-app hub push my-refactor-app v1.0.0`
