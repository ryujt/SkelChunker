# C# 파서 개발 요구사항

## 1. 개요
본 문서는 **C# 언어 파서**를 개발하기 위한 상세 요구사항을 정의한다.  
해당 파서는 코드의 구조(스켈레톤)와 내용(청크)을 분리하여 **RAG(Retrieval-Augmented Generation)** 기반 시스템에서 활용될 수 있도록 구성된다.

기본 파서 작성 원칙은 `parser-requirements.md` 문서를 따른다.

---

## 2. 목적
- 코드의 **구조적 맥락**을 표현하는 *스켈레톤(skeleton)*과  
- 실제 **세부 코드 내용**을 담은 *청크(chunk)*를 분리 추출한다.
- 모든 결과는 **클래스/메서드/함수 단위**로 정의되어야 한다.

---

## 3. 구현 원칙

### 3.1. 단일 파일 구성
- **스캐너(Tokenizer)**와 **파서(Parser)**를 **하나의 파일**로 작성한다.
- 외부 라이브러리에 의존하지 않는 방식으로 구현한다 (표준 기능만 사용).

---

## 4. 구조 분석 규칙

### 4.1. 클래스
- `class`, `interface`, `struct`, `record` 키워드를 모두 **클래스 정의로 처리**한다.
- 클래스 이름, 범위(시작~끝 위치), 내부 구성요소(메서드 등)를 추출한다.
- 클래스 안에 또 다른 클래스가 선언된 경우:
  - **중첩 클래스를 평면화하여 처리**한다.
  - 즉, 중첩되어 있어도 **독립된 클래스**처럼 별도로 추출한다.
  - 중첩의 깊이에 상관없이 모두 같은 방식으로 선형화한다.

### 4.2. 메서드
- 클래스 내부에 정의된 모든 `method`, `constructor`, `property getter/setter`, `operator` 등을 **메서드로 처리**한다.
- 시그니처, 이름, 시작~끝 위치를 정확히 추출한다.

### 4.3. 클래스 외부 코드
- 클래스 외부에 존재하는 코드는 다음 기준에 따라 처리한다:
  - 함수 형태(`returnType name(args)`)로 정의된 경우: **function**
  - 함수가 아닌 경우 (예: 변수 선언, using 문, 상수 등): **etc**

---

## 5. 출력 구조

파서의 출력은 다음과 같은 구조 정보를 포함해야 한다:

```json
{
  "skeleton": [
    {
      "type": "class",
      "name": "MyClass",
      "members": [
        {
          "type": "method",
          "name": "MyMethod",
          "range": [시작라인, 종료라인]
        }
      ]
    },
    {
      "type": "function",
      "name": "GlobalFunc",
      "range": [시작라인, 종료라인]
    },
    {
      "type": "etc",
      "name": null,
      "range": [시작라인, 종료라인]
    }
  ],
  "chunks": [
    {
      "type": "method",
      "name": "MyMethod",
      "text": "청크 전체 텍스트"
    }
  ]
}
```

---

## 6. 기타 유의사항

- 범위(range)는 라인 기준이 아닌 **문자 인덱스 기준**으로 추출해도 무방하나, 구현에 따라 유연하게 선택 가능.
- XML 주석이나 `#region`, `#pragma` 등의 메타정보는 **무시**한다.
- 조건부 컴파일 코드(`#if`, `#endif`)도 단순 코드로 간주한다.
