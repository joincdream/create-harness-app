# LLM-인간 협업 워크플로우 런타임 하네스 검토 명세

이 문서는 LLM이 워크플로우를 사람과 협업하여 각 단계를 시작하고 최종 산출물을 만들어낼 수 있도록 필요한 하네스(Harness) 기능 검토 및 아키텍처 명세를 정의합니다.

---

## 1. 현재 작업 상태 확인 (`status`)

LLM이 현재 작업 진행률 및 미완료 산출물을 파악하기 위해 동적 상태 관리가 필요합니다.

* **현황 데이터**: `blueprint.json`에 정의된 정적 블루프린트 구조만 존재하며 런타임 진척도 파일이 없음.
* **아키텍처 설계**:
  * **기계용 원본**: `.harness/state.json` (JSON 포맷)에 원자적 상태 저장 및 스키마 검증 수행.
  * **인간 협업 뷰**: `create-harness-app status` 실행 시 마크다운/터미널 요약 뷰로 자동 변환 렌더링.

---

## 2. 다음 작업 추론 (`next`)

현재 작업 상태에 기초하여 LLM이 다음에 수행해야 할 구체적인 목표와 작업을 자동으로 추론할 수 있어야 합니다.

* **현황 데이터**: Phase 순서만 정의되어 있으며 노드 간 상세 의존성 및 순서 로직이 미구현됨.
* **아키텍처 설계**:
  * `blueprint.json` 노드 명세에 `depends_on`, `required_inputs`, `expected_outputs` 필드 확장.
  * `create-harness-app next --json` CLI 명령을 제공하여 미완료 노드 중 실행 가능한 다음 작업의 Goal, 입력 파일, 출력 파일 명세를 자동 반환.

---

## 3. 단계별 가드레일 및 컨텍스트 로딩 (`context`)

현재 워크플로우 단계에 필요한 가드레일 규칙과 스펙 문서만 선별하여 LLM 런타임 프롬프트에 제공할 수 있어야 합니다.

* **현황 데이터**: `type: "guardrail"` 노드는 스캐폴딩되지만 런타임 시점에 해당 단계 문서만 번들링해주는 인터페이스가 없음.
* **아키텍처 설계**:
  * `create-harness-app context --phase <current>` CLI 명령을 지원하여 현재 Phase의 가드레일(`AGENTS.md` 등) 및 선행 단계의 스펙을 통합(Bundle) 제공.
  * 표준 MCP(Model Context Protocol) Server 인터페이스를 내장하여 LLM Agent가 Tool Call을 통해 가드레일과 상태 조회를 직접 수행할 수 있도록 확장.

---

## 4. LLM 추론과 결정적 Tool의 역할 분담 & 자연어 매핑

* **역할 분담**: LLM은 가드레일 해석 및 툴 호출 판단(Reasoning)을 담당하고, `create-harness-app` CLI Tool은 산출물 파일 실체 검증 및 `state.json` 안전 갱신(Deterministic Execution)을 담당함.
* **자연어-CLI 매핑**: 사용자의 자연어 명령("상태 알려줘", "작업 완료해줘")을 파악하여 `status`, `state update` 커맨드로 적절한 파라미터와 함께 변환 실행.

---

## 5. 프롬프트 비대화 방지 (JIT Context Optimization)

* **Dynamic Context Loading**: 전체 지시어가 아닌 현재 Phase/Node 관련 가드레일만 `context` 명령으로 동적 주입.
* **Progressive Disclosure**: 타 Phase 규칙을 배제하여 LLM 컨텍스트 토큰 최적화.
* **단순화된 인터페이스**: `status`, `next`, `context`, `state update` 4개 핵심 인터페이스로 툴 정의 단순화.

---

## 6. Multi-Agent Runtime Target 지원 (Antigravity, Claude CLI 등)

* **에이전트별 특화 파일 규격**:
  * **Antigravity CLI**: `.agents/AGENTS.md` 가드레일 및 `.agents/skills/<name>/SKILL.md` 지능형 스킬 구조 지원.
  * **Claude Code / CLI**: `CLAUDE.md` 규칙 파일 및 `.mcp.json` 도구 설정 구조 지원.
* **타깃 분기 스캐폴딩**:
  * `create-harness-app --agent <antigravity|claude|generic>` 옵션과 `blueprint.json` 지목을 통해 대상 AI 에이전트에 최적화된 파일/디렉토리 구조를 가변 생성.
