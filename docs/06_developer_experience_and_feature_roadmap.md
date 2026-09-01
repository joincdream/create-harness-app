# 🚀 create-harness-app Developer Experience (DX) & Feature Roadmap

- **Document ID**: `docs/06_developer_experience_and_feature_roadmap`
- **Target Component**: `create-harness-app` CLI
- **Status**: Proposed Feature Roadmap
- **Last Updated**: `2026-08-03T14:55:00Z`

---

## 1. 개요 (Overview)

본 문서는 `create-harness-app` CLI를 사용하는 실제 **AI 에이전트 개발자(Developer Persona)** 관점에서 개발 생산성과 DX(Developer Experience)를 극대화하기 위해 발굴된 5대 추가 기능과 단계별 구현 로드맵을 정의합니다.

---

## 2. Top 5 기능 상세 명세 (Feature Specifications)

### 2.1. `create-harness-app validate` (로컬 하네스 무결성 검증기) ⭐⭐⭐⭐⭐
- **개념**: `hub push`로 업로드하기 전, 로컬의 `.harness/state.json`, `${file:.harness/description.md}`, `.agents/mcp_config.json`, `docs/` 산출물 경로가 깨지거나 누락되지 않았는지 사전 검증(Lint & Dry-run)합니다.
- **개발자 혜택**: CI/CD 파이프라인이나 사전 푸시 단계에서 문법 오류 및 경로 누락을 즉시 감지하여 업로드 실패를 예방합니다.
```bash
create-harness-app validate
# 출력:
# ✅ .harness/state.json is valid
# ✅ Referenced file .harness/description.md exists
# ✅ All 3 workflow phase nodes are correctly mapped
```

### 2.2. `create-harness-app init` (기존 프로젝트에 하네스 레이어 주입) ⭐⭐⭐⭐⭐
- **개념**: 신규 스캐폴딩(`create-harness-app my-app`)뿐만 아니라, **이미 코드가 존재하는 레거시 프로젝트 디렉터리**에서 `.harness/`, `docs/`, `.agents/mcp_config.json` 하네스 구조만 주입받는 서브커맨드입니다.
- **개발자 혜택**: 이미 운영 중인 수년 된 사내 서비스 프로젝트에 1초 만에 AI 에이전트 가드레일/하네스를 적용할 수 있습니다.
```bash
cd my-existing-legacy-service
create-harness-app init --template my-refactor-app
```

### 2.3. `create-harness-app status` (SDLC 파이프라인 진척도 시각화) ⭐⭐⭐⭐
- **개념**: 현 워크스페이스가 SDLC 파이프라인 3단계 중 어느 Phase/Node에 위치해 있고, 어떤 산출물(`analysis_report.md`, `refactor_diff.patch` 등)이 작성 완료되었는지 터미널 트리 형태로 시각화합니다.
- **개발자 혜택**: 에이전트가 어디까지 작업을 수행했는지 한눈에 실시간 확인 가능합니다.
```bash
create-harness-app status
# 출력: 
# 📦 my-refactor-app Progress (Phase 2/3)
#   [✓] Phase 1: code-analysis ➔ docs/01_planning/analysis_report.md (Created)
#   [➔] Phase 2: clean-arch-refactor ➔ docs/02_design/refactor_diff.patch (In Progress)
#   [ ] Phase 3: tdd-verification ➔ docs/03_development/test_summary.log (Pending)
```

### 2.4. `create-harness-app hub search <query>` (키워드/태그 템플릿 검색) ⭐⭐⭐
- **개념**: Hub에 템플릿이 수십~수백 개로 늘어났을 때, 단순 `list`만으로는 찾기 어려우므로 태그(`refactor`, `golang`, `react`), 에이전트 종류 기반 검색을 지원합니다.
- **개발자 혜택**: 필요한 하네스를 1초 만에 검색 및 조망합니다.
```bash
create-harness-app hub search refactor --agent antigravity
```

### 2.5. `create-harness-app export` / `import` (오프라인 폐쇄망 번들링) ⭐⭐⭐
- **개념**: 외부 인터넷 접속이 차단된 사내 오프라인/폐쇄망(Air-Gapped Enterprise) 환경에서 템플릿을 `.tar.gz` 로컬 바이너리로 상호 오프라인 이관/로드하는 기능입니다.
- **개발자 혜택**: 금융/공공/엔터프라이즈 오프라인 환경을 완벽 지원합니다.
```bash
create-harness-app export my-refactor-app -o my-refactor-app-v1.0.0.harness
create-harness-app import my-refactor-app-v1.0.0.harness
```

---

## 3. 단계별 구현 로드맵 (Implementation Roadmap)

| 단계 | 기능 명칭 | 목표 사양 | 비고 |
| :--- | :--- | :--- | :--- |
| **Phase 1 (1순위)** | `validate` & `init` | 푸시 전 검증 및 기존 프로젝트 하네스 주입 | DX 향상 핵심 |
| **Phase 2 (2순위)** | `status` & `hub search` | 워크스페이스 상태 시각화 및 레지스트리 검색 | 대규모 하네스 관리 |
| **Phase 3 (3순위)** | `export` / `import` | Air-Gapped 폐쇄망 번들 파일 오프라인 지원 | Enterprise 대응 |
