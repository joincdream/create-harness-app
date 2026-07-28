# 01. 상세 기능 및 도메인 엔티티 설계 (Detailed Feature & Entity Spec)

## 1. 모듈별 기능 상세 명세 (SOLID & SoC 적용)

### 1.1 CLI Argument Parser & Option Evaluator (SRP)
- **역할**: 사용자로부터 전달된 CLI 인자(`process.argv`) 및 플래그를 파싱하고 대상 타겟 디렉토리를 산출함.
- **단일 책임**: CLI 입력을 정제된 `Config` 객체로 변환하는 역할만 수행.

### 1.2 Harness Scaffolder Engine (Open-Closed Principle)
- **역할**: 정제된 `Config` 객체를 받아 SDLC 계층 트리 및 가이드라인 마크다운 템플릿을 생성함.
- **확장성**: 새로운 Phase 템플릿이 추가되어도 기존 생성 로직 수정 없이 템플릿 레지스트리 확장을 지원.

### 1.3 Template Registry & File Writer (Dependency Inversion)
- **역할**: 메모리 내 마크다운/YAML 템플릿과 파일 시스템 I/O 레이어를 디커플링하여 안전한 샌드박스 쓰기를 지원.

### 1.4 Harness Blueprint Registry & Modular Template Engine (Terraform-style)
- **역할**: Terraform Module처럼 도메인별(RESTful Web API, CLI App, Batch 등) 표준 워크플로우 템플릿을 선언적 블루프린트(Blueprint)로 관리하고 커뮤니티 및 팀 내 공유를 지원.
- **--template 플래그**: 사용자 지정을 통해 도메인에 특화된 하네스 블루프린트 템플릿 모듈을 동적 주입.

### 1.5 Template Resolution & Override Engine (Go Embed & Custom Overrides)
- **역할**: 템플릿 탐색 시 3단계 우선순위 체인(Fallback Chain)에 따라 블루프린트와 마크다운 파일들을 결정론적으로 로드함.
- **우선순위 체인**:
  1. **Flag/Path Override**: `--template [path]`로 지정된 로직/로컬 디렉토리의 블루프린트 템플릿.
  2. **User Custom Override**: `~/.config/create-harness-app/templates/[name]/`에 위치한 사용자 정의 커스텀 템플릿.
  3. **Default Embedded**: Go 바이너리 내부에 컴파일된 기본 블루프린트 (`//go:embed templates/*`).

---

## 2. 핵심 도메인 엔티티: JSON 선언적 블루프린트 스키마 (Declarative JSON Blueprint)

하네스 스루풋은 단순하고 명확한 **`blueprint.json`** 선언서로 모델링됩니다.
- **Workflow**: SDLC Phase 단계 (물리적 디렉토리)
- **Node**: Workflow 내의 세부 작업 및 산출물 (파일)
- **Guardrail**: Node에 적용되는 AI 통제 규칙 (`AGENTS.md` / 가이드라인)

```json
{
  "name": "web-api-blueprint",
  "version": "1.0.0",
  "description": "RESTful Web API 하네스 블루프린트 모듈",
  "workflows": [
    {
      "phase": 1,
      "dir": "01_planning",
      "description": "기획 & 요구사항 명세 Stage",
      "nodes": [
        { "file": "00_raw_inputs/README.md", "type": "input", "description": "원시 요구사항 입력 창구" },
        { "file": "01_user_requirements.md", "type": "spec", "description": "유저 요구사항 명세" },
        { "file": "03_functional_requirements.md", "type": "spec", "description": "기능 요구사항 & 수락 조건" }
      ]
    },
    {
      "phase": 2,
      "dir": "02_architecture_and_design",
      "description": "아키텍처 & 인터페이스 설계 Stage",
      "nodes": [
        { "file": "01_detailed_feature_spec.md", "type": "spec", "description": "SOLID 도메인 엔티티 설계" },
        { "file": "02_standard_interface_dto.md", "type": "spec", "description": "표준 DTO 스키마" },
        { "file": "05_tech_stack_and_skills.md", "type": "guardrail", "description": "AGENTS.md 불변 아키텍처 가드레일" }
      ]
    }
  ]
}
```

---

## 3. Go 패키지 아키텍처 레이아웃 (Go Project Package Architecture)

실제 Go 구현 시 무결한 관심사 분리(SoC)를 보장하는 표준 프로젝트 구조입니다.

```text
cmd/
└── create-harness-app/
    └── main.go                 # 엔트리포인트 (Flag 파싱 및 CLI 실행)
internal/
├── config/
│   └── config.go               # ScaffoldConfig DTO 및 Flag 파서
├── engine/
│   └── scaffolder.go           # 디렉토리 트리 및 파일 생성 엔진 (Scaffolding Engine)
└── template/
    ├── resolver.go             # 3단계 템플릿 탐색 및 fallback 엔진
    └── blueprint.go            # Declarative JSON Blueprint Struct 정의
templates/                      # //go:embed 전용 기본 블루프린트 패키지 디렉토리
├── default/
│   ├── blueprint.json          # 기본 블루프린트 매니페스트
│   └── files/                  # 기본 md/yaml 템플릿 산출물 파일들
```

---

## 4. 템플릿 탐색 알고리즘 (Template Resolution Algorithm Pseudo-code)

```go
// 3단계 Fallback 탐색 알고리즘
func ResolveBlueprint(templateName string) (fs.FS, *Blueprint, error) {
    // Phase 1: --template 인자가 절대/상대 파일 경로로 주어진 경우
    if isLocalPath(templateName) && pathExists(templateName) {
        return loadFromDisk(templateName)
    }

    // Phase 2: ~/.config/create-harness-app/templates/[templateName] 커스텀 탐색
    customPath := filepath.Join(userHomeDir, ".config/create-harness-app/templates", templateName)
    if pathExists(customPath) {
        return loadFromDisk(customPath)
    }

    // Phase 3: 바이너리 내장 //go:embed templates/ 기본 블루프린트 Fallback
    embeddedFS, err := fs.Sub(embeddedTemplatesFS, "templates/"+templateName)
    if err == nil {
        return loadFromFS(embeddedFS)
    }

    return nil, nil, fmt.Errorf("template '%s' not found in any resolution chain", templateName)
}
```
