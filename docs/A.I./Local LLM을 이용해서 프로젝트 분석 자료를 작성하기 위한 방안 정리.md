# Local LLM을 이용해서 프로젝트 분석 자료를 작성하기 위한 방안 정리

## React.JS 프로젝트 분석 방안

1. 폴더 및 파일 구성
2. 클래스 및 모듈 구성
  - API 관련
  - Store 관련
  - 유틸리티 모듈 관련
  - 기타 UI를 제외한 기능 모듈 등
3. Navigation Diagram
4. 호출 관계 그래프 (Job Flow Diagram)

## Navigation Diagram

프로젝트의 페이지 이동 및 API 호출 관계를 아래와 같이 분석하는 다이어그램이다.

다이어그램 구성 요소는 아래와 같다.

* Page: 주소(route)를 가지는 페이지
* (/api): API 호출
* (functionCall): 함수 호출

### 예시

```navigation
Home --> Main : 로그인 성공
Home --> Login : 미로그인 상태
Login --> Terms : 신규 회원
Terms --> Signup : 약관 동의 완료
Signup --> Signup : 오류 발생
Signup --> Main
Login --> (signin) : 로그인 시도
(signin) --> Login : 오류 발생
(signin) --> Main
Main --> Content
Main --> ShuffleTest
Main --> Settings
```

## Job Flow Diagram

클래스나 모듈이 호출하는 관계를 아래와 같이 분석하는 다이어그램이다.
클래스나 모듈의 내부의 흐름은 무시한다.
외부 클래스나 외부 모듈에게 어떤 요청을 하고 받는지에 대한 협력 관계에 집중한다.
각 페이지마다 별도로 작성한다.

### 예시

```jobflow
master: Main
Object: Main, LoginAPI, UserStore
Main.OnMount --> LoginAPI.signin
LoginAPI.signin --> LoginAPI.result.error
LoginAPI.signin --> LoginAPI.result.Ok
LoginAPI.result.Ok --> UserStore.setUserInfo
```
