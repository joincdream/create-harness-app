# 03. CLI 명령어 및 Flag 스펙 명세서 (CLI Command Spec)

## 1. CLI 명령어 사용법 (Usage)

```bash
# 기본 실행 (현재 디렉토리 또는 지정 디렉토리에 SDLC 하네스 스캐폴딩)
create-harness-app [target_directory] [flags]

# 예시 1: 대화형 위자드 실행 (Interactive Prompt)
create-harness-app

# 예시 2: 원라인 직접 실행 (Non-interactive One-liner)
create-harness-app my-app --template web-api --force
create-harness-app --version
```

### 1.1 실행 모드 명세 (Execution Modes)

1. **Interactive Wizard Mode (인자가 미지정된 경우)**:
   - CLI 실행 시 대상 디렉토리나 템플릿 인자가 없으면 대화형 터미널 위자드가 시작됩니다.
   - **Step 1**: `? 프로젝트 디렉토리명을 입력하세요:` (기본값: `my-harness-app`)
   - **Step 2**: `? 사용할 하네스 블루프린트(템플릿)를 선택하세요:`
     - `> 🌐 web-api  (RESTful Web API: OpenAPI + DTO + TDD)`
     - `  🛠️ cli-app   (CLI Application: Go + Cobra + Unit Test)`
     - `  ⚡ batch     (Event Batch Pipeline)`
     - `  📁 custom    (로컬/외부 커스텀 블루프린트 경로 지정)`
   - **Step 3**: `? 00_raw_inputs 디렉토리를 포함하시겠습니까? (Y/n)`

2. **Non-interactive One-liner Mode (인자가 지정된 경우)**:
   - `create-harness-app my-app --template web-api`처럼 플래그가 지정되면 위자드 질의응답 없이 즉시 스캐폴딩을 완수합니다. (CI/CD 파이프라인 및 파워유저용)

---

## 2. CLI 플래그 및 옵션 명세 (Flags & Options)

| Flag Short | Flag Long   | Type   | Default | Description |
|------------|-------------|--------|---------|-------------|
| `-f`       | `--force`   | bool   | `false` | 대상 디렉토리에 기존 파일이 존재할 경우 강제 덮어쓰기 |
| `-v`       | `--version` | bool   | `false` | `create-harness-app` CLI 버전 정보 출력 |
| `-h`       | `--help`    | bool   | `false` | 도움말 및 명령어 사용법 출력 |

---

## 3. Exit Code (프로세스 종료 코드)

* **Exit Code 0**: 성공적으로 모든 SDLC Phase 디렉토리 및 템플릿 생성을 완료함.
* **Exit Code 1**: 파라미터 인자 오류 또는 잘못된 플래그 입력.
* **Exit Code 2**: 디렉토리 생성 또는 `go:embed` 파일 쓰기 권한/I/O 실패.

---

## 4. Standard I/O (표준 입출력 스펙)

* **Stdout**: 스캐폴딩 진행 상태(`✓ 생성: 01_planning/01_user_requirements.md`) 및 완료 안내 로그 출력.
* **Stderr**: 파일 시스템 권한 오류 또는 유효하지 않은 디렉토리명 에러 출력.
