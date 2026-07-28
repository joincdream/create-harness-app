# 02. 표준 요청/응답 DTO 및 스키마 명세 (Standard Interface Spec)

## 1. CLI Execution DTO Schema

```typescript
// CLI 호출 입력 DTO
interface CLIInputDTO {
  appName: string;             // 예: 'my-project'
  flags?: {
    force?: boolean;           // 덮어쓰기 강제 여부
    template?: 'default';      // 스캐폴딩 템플릿 종류
  };
}

// CLI 실행 결과 DTO (Result Format)
interface CLIExecutionResultDTO {
  success: boolean;            // 성공 여부
  targetDirectory: string;     // 생성된 절대 경로
  filesCreatedCount: number;   // 생성된 파일 개수
  phasesCreated: string[];     // 생성된 SDLC Phase 목록
  errorMessage?: string;       // 실패 시 에러 메시지
}
```

---

## 2. 공통 JSON API / CLI 로그 포맷 (Standard Log Output)

```json
{
  "timestamp": "2026-07-28T21:42:00.000Z",
  "level": "INFO",
  "action": "SCAFFOLD_PHASE_COMPLETE",
  "details": {
    "phase": "01_planning",
    "files": ["00_raw_inputs/README.md", "01_user_requirements.md"]
  }
}
```
