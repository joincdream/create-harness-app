# 엔터프라이즈 고도화 로드맵 및 핵심 기능 명세

이 문서는 `create-harness-app`을 대규모 팀 및 기업 환경(Enterprise)에 도입하고 조직 차원에서 SDLC 하네스를 성공적으로 운용하기 위해 필요한 4대 핵심 고도화 로드맵을 정의합니다.

---

## 1. CI/CD 파이프라인 Quality Gate (`make verify` / `verify`)

* **개념**: 특정 CI/CD 벤더에 의존하지 않고, **스캐폴딩 시점에 프로젝트 루트에 `Makefile` 타깃(`make verify`)으로 내장**하여 로컬과 CI 서버에서 100% 동일한 검증 파이프라인을 구동합니다.
* **명세 및 실행**:
  * `make verify` (또는 `create-harness-app verify [--strict]`)
  * **검증 내용**:
    1. 산출물 문서 파일 존재 여부 및 마크다운 표준 구조(섹션, 테이블) 유효성 검사
    2. OpenAPI 스키마 syntax 검사
    3. TDD 단위 테스트 성공 여부 및 커버리지 조건 검증
  * 검증 실패 시 CI/CD PR 파이프라인 자동 블로킹(Exit Code 1).
  * **장점**: 로컬-CI 검증 환경 100% 동기화, CI 벤더 종속성 완전 제거, LLM의 자율 `make verify` 호출 검증 지원.

---

## 2. 사내 중앙 템플릿 레지스트리 (`registry push / pull`)

* **개념**: 전사 조직이 검증된 표준 아키텍처, 가드레일, 스킬 패키지를 중앙에서 일관되게 관리하고 공유합니다.
* **명세 및 실행**:
  * `create-harness-app registry push <template_name> --registry <url>`: 현장에서 개발/최적화된 하네스를 사내 레지스트리로 업로드
  * `create-harness-app pull <template_name> --registry <url>`: 사내 레지스트리의 최신 하네스 템플릿을 로컬로 동기화받아 스캐폴딩

---

## 3. Human-in-the-Loop 승인 게이트 (`approval_required`)

* **개념**: 보안, 핵심 아키텍처 설계 등 리스크가 높은 Phase/Node에 대해 LLM의 자율 완료 처리를 제한하고 테크 리드의 서명(Sign-off)을 강제합니다.
* **명세 및 연동**:
  * `blueprint.json` 노드 명세에 `"approval_required": true` 옵션 제공.
  * LLM이 `state update` 호출 시, `approval_required` 필드가 참이면 즉시 완료되지 않고 `status: "awaiting_approval"` 상태로 전환.
  * 테크 리드가 `create-harness-app approve --node <id>` 커맨드를 실행한 후에만 `completed` 상태로 전환되어 다음 Phase 개시 허용.

---

## 4. 역할별 멀티 에이전트 오케스트레이션 (Multi-Agent Subagent Coordination)

* **개념**: 단일 LLM의 컨텍스트 비대화 및 환각을 방지하기 위해, Main Coordinator가 `create-harness-app next`로 지정된 Phase별 전담 Subagent를 스폰하고 지목된 `SKILL.md`만 로드하여 작업하는 아키텍처입니다.
* **`create-harness-app next` 파이프라인 흐름**:
  1. **Next Target 분석**: Main Agent가 `create-harness-app next` 실행 ➔ 미완료 노드 및 지정 `skill_path` (`SKILL.md`) 반환.
  2. **명시적 스킬 주입 & Subagent 스폰**: Main Agent가 반환된 `skill_path`의 `SKILL.md` 지침만 주입하여 깨끗한 컨텍스트의 전담 Subagent 스폰.
  3. **산출물 작성 & State Update**: Subagent가 `SKILL.md`에 정의된 입출력 문서 작성 후 `create-harness-app state update` 툴을 호출하여 완수 처리.
* **핵심 이점**: 컨텍스트 격리(Context Isolation)를 통한 환각 최소화, 레이저 포커스 지시어 로딩, 독립 노드의 병렬 실행 지원.

---

## 5. 웹 애플리케이션 시각화 및 대시보드 (Workflow & State Visualization)

* **개념**: 템플릿 관리를 넘어, SDLC 워크플로우 그래프와 런타임 진척도를 웹 대시보드로 시각화하여 인간과 AI 에이전트 간의 협업 관제 탑(Visual Control Tower)을 제공합니다.
* **3대 핵심 시각화 기능**:
  1. **DAG 워크플로우 청사진 시각화**: `blueprint.json`의 Phase, Node, `depends_on` 의존성 및 입출력 문서 구조를 인터랙티브 DAG 그래프로 렌더링.
  2. **실시간 런타임 진척도 대시보드**: `.harness/state.json`과 동기화하여 `pending`, `in_progress`, `completed`, `blocked` 상태를 실시간 상태 맵으로 모니터링.
  3. **원클릭 HITL 승인 웹 UI**: `approval_required` 노드의 산출물 문서를 웹 UI에서 검토 후 버튼 클릭으로 `Approve` 처리.
