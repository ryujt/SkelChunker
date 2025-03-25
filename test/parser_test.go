package parser_test

import (
	"fmt"
	"SkelChunker/src/parser"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCSharpParser(t *testing.T) {
	// 현재 디렉토리 가져오기
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	
	// 테스트할 파일 경로
	testFilePath := filepath.Join(dir, "TestClass.cs")

	// 파일 읽기
	content, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("파일을 읽는 중 오류 발생: %v", err)
	}

	// 파서 생성 및 실행
	p := parser.NewCSharpParser()
	nodes, chunks, err := p.Parse(string(content))
	
	if err != nil {
		t.Fatalf("파싱 중 오류 발생: %v", err)
	}

	// 결과 확인
	if len(nodes) == 0 {
		t.Error("노드가 추출되지 않았습니다.")
	}

	if len(chunks) == 0 {
		t.Error("청크가 추출되지 않았습니다.")
	}

	// 결과 출력
	fmt.Println("추출된 노드 수:", len(nodes))
	fmt.Println("추출된 청크 수:", len(chunks))
	
	// 청크 내용 출력
	for i, chunk := range chunks {
		fmt.Printf("청크 %d - MD5: %s\n텍스트 (길이 %d):\n%s\n\n", 
			i, chunk.MD5, len(chunk.Text), chunk.Text)
	}
} 