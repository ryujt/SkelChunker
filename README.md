# SkelChunker

SkelChunker는 소스 코드를 분석하여 스켈레톤(skeleton)과 청크(chunk)로 분리하는 도구입니다. 각 파일의 구조 정보와 세부 구현 내용을 추출하여 JSON 형식으로 저장합니다.

## 주요 기능

- 다양한 프로그래밍 언어 지원 (Java, C, C++, C#, Python, JavaScript, TypeScript, Go, Kotlin, PHP, HTML)
- 자동 언어 감지 및 맞춤형 파서 적용
- 스켈레톤(구조)과 청크(구현) 분리
- MD5 기반 중복 파일 처리
- 설정 가능한 무시 폴더 목록

## 요구사항

- Go 1.16 이상
- Git

## 설치 방법

1. 저장소 클론
```bash
git clone https://github.com/yourusername/SkelChunker.git
cd SkelChunker
```

2. 의존성 설치
```bash
go mod download
```

3. 빌드
```bash
go build -o skelchunker.exe src/main.go
# or 
go build -o skelchunker src/main.go
```

## 설정

`config.json` 파일을 통해 다음 설정을 할 수 있습니다:

```json
{
    "folders": [
        "/projects/service-a"
    ],
    "ignore-folders": [
        "node_modules",
        "dist",
        "build",
        "target",
        "bin",
        "obj",
        "venv",
        "docs",
        "test",
        ".git",
        ".idea",
        ".vscode",
        ".DS_Store",
        ".env",
        ".env.local"
    ],
    "parsers": {
        ".java": "java_parser",
        ".c": "c_parser",
        ".cpp": "cpp_parser",
        ".cs": "csharp_parser",
        ".py": "python_parser",
        ".js": "javascript_parser",
        ".ts": "typescript_parser",
        ".go": "go_parser",
        ".kt": "kotlin_parser",
        ".php": "php_parser",
        ".html": "html_parser"
    }
}
```

### 설정 항목 설명

- `folders`: 분석할 소스 코드 폴더 경로 목록
- `ignore-folders`: 분석에서 제외할 폴더 이름 목록
- `parsers`: 파일 확장자별 파서 매핑 정보

## 실행 방법

기본 설정 파일 사용:
```bash
./skelchunker
```

커스텀 설정 파일 지정:
```bash
./skelchunker -config /path/to/config.json
```

## 출력 형식

분석 결과는 각 소스 파일과 동일한 위치에 `.SkelChunker` 확장자로 저장됩니다.

```json
{
  "path": "/projects/service-a",
  "filename": "main.cpp",
  "md5": "파일내용의 MD5",
  "summary": "코드 요약 결과",
  "embeddings": [임베딩 결과],
  "skeleton": [
    {
      "type": "class",
      "name": "ClassName",
      "members": [
        {
          "md5": "청크 내용의 MD5",
          "type": "method",
          "name": "methodName"
        }
      ]
    }
  ],
  "chunks": [
    {
      "md5": "청크 내용의 MD5",
      "text": "청크 원문",
      "embeddings": [임베딩 결과]
    }
  ]
}
```

## 프로젝트 구조

```
src/
├── main.go                 # 메인 진입점
├── config/
│   └── config.go           # 설정 관련 구조체와 함수
├── parser/
│   ├── parser.go          # 파서 인터페이스
│   └── factory.go         # 파서 팩토리
├── model/
│   ├── skeleton.go        # 스켈레톤 관련 구조체
│   └── chunk.go           # 청크 관련 구조체
├── analyzer/
│   └── analyzer.go        # 코드 분석 로직
└── utils/
    └── hash.go            # MD5 해시 등 유틸리티 함수
```

