# HarnessHub 백엔드 연동 및 통합 개요

> **문서 이동 안내**: 상세 백엔드 API 연동 개발 계획 및 Docker 스타일 API 연동 설계서 전문은 `harness-hub` 프로젝트 문서군으로 통합 이관되었습니다.

---

## 1. 개요 및 참고 문서 링크

`create-harness-app` CLI는 중앙 `HarnessHub 백엔드`(`harness-hub/backend`) 서버와 연동하여 템플릿 검색, 조회, 다운로드/설치(`pull`), 로컬 삭제(`templates remove`) 기능을 제공합니다.

### 📌 이관된 상세 명세 문서 링크
* **CLI 연동 개발 계획서**: [`harness-hub/docs/07_cli_integration_plan.md`](file:///home/yundream/myjob/cloit/poc/harness-hub/docs/07_cli_integration_plan.md)
* **REST API & Docker 스타일 CLI 설계서**: [`harness-hub/docs/08_cli_integration_design.md`](file:///home/yundream/myjob/cloit/poc/harness-hub/docs/08_cli_integration_design.md)

---

## 2. CLI 주요 커맨드 요약

* `create-harness-app hub list`: 원격 레지스트리 템플릿 목록 조회
* `create-harness-app hub info <name> [version]`: 특정 템플릿 상세 및 Blueprint JSON 조회
* `create-harness-app hub pull <name> [version]`: 원격 템플릿 다운로드 및 로컬 캐시(`~/.config/create-harness-app/templates/<name>`) 압축해제/저장
* `create-harness-app templates remove <name>`: 로컬 캐시 템플릿 삭제
