# 05. 기술 스택 및 불변 아키텍처 규칙 (AGENTS.md Guardrails)

## 1. 기술 스택 (Tech Stack Constraints)

* **Language/Runtime**: Go (Golang 1.22+) - 단일 정적 바이너리 배포
* **CLI Engine**: Go Standard Library (`embed`, `flag`, `os`) - Zero External Dependency
* **Testing Framework**: Go Native `testing` package (go test)
* **Interface Schema**: OpenAPI 3.0.3 YAML Spec

---

## 2. 불변 아키텍처 규칙 (Architectural Guardrails)

1. **Rule 1: Spec First Constraint**:
   * 어떤 경우에도 `02_architecture_and_design` 스펙 문서(OpenAPI, DTO, SOLID 설계) 검증이 완료되기 전에 구현 코드 작성으로 이행할 수 없다.
2. **Rule 2: Zero External Bloat**:
   * CLI 스캐폴더 엔진은 불필요한 무거운 외부 무거운 프레임워크에 의존하지 않고 경량화된 표준 API를 유지한다.
3. **Rule 3: Deterministic Test Sensors**:
   * 모든 생성물은 `go test ./...` 검증을 통과해야 배포/완료로 인정한다.
4. **Rule 4: Template Priority & Override Chain Rule**:
   * 템플릿 탐색 시 `--template` 지정 경로 ➔ `~/.config/create-harness-app/templates/` ➔ `go:embed` 기본 바이너리 순서의 3단계 오버라이드 체인을 항상 보장해야 한다.
