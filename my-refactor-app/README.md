# Google Antigravity Code Refactoring Harness ⚡

Google Antigravity 에이전트를 위한 **자동화 코드 리팩토링 및 클린 아키텍처 가이드라인 적용** 하네스 템플릿입니다.

### 🌟 주요 기능 및 혜택
- **복잡도 분석 센서**: cyclomatic complexity 기준 10 이상의 레거시 메서드를 자동 감지
- **클린 아키텍처 가드레일**: SOLID 원칙 및 레이어드 패키지 의존성 단방향 규칙 강제
- **자동 Diff 패치 생성**: 리팩토링 전/후 안전한 `refactor_diff.patch` 생성 및 `go test` 검증

### 💡 사용 방법
```bash
create-harness-app pull antigravity-code-refactor
create-harness-app my-app --template antigravity-code-refactor
```