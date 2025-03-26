# 소스코드 스켈레톤 및 청크 추출 프로젝트 요구사항 명세서

## 개요  
본 프로젝트는 지정된 소스코드 디렉토리(폴더) 내의 파일들을 분석하여, 각 파일의 구조 정보를 담은 ‘스켈레톤(skeleton)’과 세부 구현 내용을 포함하는 ‘청크(chunk)’로 분리하고 저장하는 기능을 수행한다.  
대상 언어를 자동으로 분류하고, 언어별 맞춤형 파서를 사용하여 구조적 분석을 수행한다.

---

## 사용 기술 및 모델
- 구현 언어: Go (Golang)
- 요약 모델: `CHAT_MODEL = "gpt-4o-mini"`
- 임베딩 모델: `EMBEDDINGS_MODEL = "text-embedding-3-large"`

---

## 상세 요구사항

### 1. 기본 동작 흐름
- 설정 파일(`config.json`)을 통해 다음 정보를 입력받는다:
  - 분석 대상 폴더 목록
  - 확장자별 파서 매핑 정보
- 지정된 폴더를 재귀적으로 탐색하여 해당 확장자의 파일을 모두 수집하고 분석을 수행한다.
- 분석 결과는 각 파일마다 JSON 형식으로 저장된다.

---

### 2. 설정 파일 (`config.json`)
```json
{
    "folders": [
      "/projects/service-a",
      "/projects/service-b"
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

- `folders`: 분석 대상 루트 디렉토리 배열  
- `ignore-folders`: 분석 제외 대상 루트 디렉토리 배열  
- `parsers`: 파일 확장자별 파서 지정

---

### 3. 분석 대상 파일 처리
- 확장자가 `parsers` 항목에 명시된 경우에만 분석 대상이 된다.
- 디렉토리는 재귀적으로 탐색한다.
- 같은 파일 내용(`md5`)이 이미 분석된 경우에는 재처리하지 않는다.

---

### 4. 언어별 파서 요구사항
- 각 파서는 다음 구조적 요소만 추출한다:
  - 클래스(class), 메서드(method), 함수(function)
- 주석, import, 상세 구현 등은 무시한다.
- 구조 요소별로 고유 청크를 생성하여 분석한다.

---

### 5. 출력 형식

```json
{
  "path": "/projects/service-a",
  "filename": "main.cpp",
  "md5": "파일내용의 MD5",
  "embeddings": [임베딩결과(파일내용)],
  "skeleton": [
    {
      "type": "class",
      "name": "ClassName",
      "members": [
        {
          "md5": "청크내용의 MD5",
          "type": "method",
          "name": "methodName"
        }
      ]
    },
    {
      "md5": "청크내용의 MD5",
      "type": "function",
      "name": "functionName"
    },
  ],
  "chunks": [
    {
      "md5": "청크내용의 MD5",
      "text": "청크내용",
      "embeddings": [임베딩결과(청크내용)]
    }
  ]
}
```

- 결과 파일은 다음 위치에 저장:
  - `{원본파일경로}/{파일명(확장자 제거)}.SkelChunker`

---

### 6. 임베딩 용량 제한 처리 방법

* 대상 텍스트 크기가 큰 경우에는 여러 조각으로 나누서 저장한다.
  * "embeddings": [임베딩결과1, 임베딩결과2, ...]
  * 임베딩결과는 임베딩 모델에 의해서 리턴된 배열이다.
* 대상 텍스트를 라인 단위로 붙여가면서 24KB가 넘지 않도록 한다.

---

## 비고
- 분석 제외 언어(지정되지 않은 확장자)는 무시하며 로그만 출력한다.
- 파일 단위 결과 외에 전체 프로젝트 단위 메타 정보 집계 파일 생성은 별도 기능으로 지원 가능.
