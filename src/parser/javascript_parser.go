package parser

import (
	"SkelChunker/src/model"
	"strings"
)

// JavaScriptParser는 JavaScript 소스 코드를 분석하는 파서입니다.
type JavaScriptParser struct{}

// NewJavaScriptParser는 새로운 JavaScript 파서를 생성합니다.
func NewJavaScriptParser() *JavaScriptParser {
	return &JavaScriptParser{}
}

// Parse는 JavaScript 소스 코드를 분석하여 스켈레톤과 청크를 반환합니다.
func (p *JavaScriptParser) Parse(sourceCode string) ([]model.SkeletonNode, []model.Chunk, error) {
	// 이 예제에서는 간단한 구현만 제공합니다.
	// 실제로는 더 복잡한 파싱 로직이 필요합니다.
	
	// 소스 코드 청크를 위한 MD5 계산
	contentMD5 := calculateMD5(sourceCode)
	
	// 기본 청크는 전체 파일 내용
	basicChunk := model.Chunk{
		MD5:  contentMD5,
		Text: sourceCode,
	}
	
	// 스켈레톤 노드 구성 (간단히 함수와 클래스만 검출)
	var nodes []model.SkeletonNode
	
	// 라인 단위로 분리
	lines := strings.Split(sourceCode, "\n")
	
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		
		// 함수 정의 검출
		if strings.HasPrefix(line, "function ") || 
		   strings.Contains(line, "function(") || 
		   strings.Contains(line, "=> {") {
			// 함수명 추출 (간단한 구현)
			functionName := extractFunctionName(line)
			
			// 스켈레톤 노드에 추가
			nodes = append(nodes, model.SkeletonNode{
				Type: "function",
				Name: functionName,
				MD5:  contentMD5,
			})
		}
		
		// 클래스 정의 검출
		if strings.HasPrefix(line, "class ") {
			// 클래스명 추출
			className := extractClassName(line)
			
			// 클래스 노드 생성
			classNode := model.SkeletonNode{
				Type:    "class",
				Name:    className,
				Members: []model.Member{},
			}
			
			// 클래스 내 메서드 검출 (간단한 구현)
			for j := i + 1; j < len(lines); j++ {
				methodLine := strings.TrimSpace(lines[j])
				
				// 클래스 끝 검출
				if methodLine == "}" {
					break
				}
				
				// 메서드 정의 검출
				if strings.Contains(methodLine, "(") && strings.Contains(methodLine, ")") &&
				   !strings.HasPrefix(methodLine, "//") {
					methodName := extractMethodName(methodLine)
					
					// 메서드를 멤버로 추가
					classNode.Members = append(classNode.Members, model.Member{
						Type: "method",
						Name: methodName,
						MD5:  calculateMD5(methodLine),
					})
				}
			}
			
			// 스켈레톤 노드에 클래스 추가
			nodes = append(nodes, classNode)
		}
	}
	
	// 노드가 없으면 파일 전체를 하나의 청크로 처리
	if len(nodes) == 0 {
		return nil, []model.Chunk{basicChunk}, nil
	}
	
	return nodes, []model.Chunk{basicChunk}, nil
}

// 함수명 추출 (간단한 구현)
func extractFunctionName(line string) string {
	if strings.HasPrefix(line, "function ") {
		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			// 괄호 위치로 함수명 추출
			functionName := parts[1]
			parenIndex := strings.Index(functionName, "(")
			if parenIndex > 0 {
				return functionName[:parenIndex]
			}
			return functionName
		}
	} else if strings.Contains(line, "=") {
		// 화살표 함수 또는 함수 표현식
		parts := strings.Split(line, "=")
		if len(parts) >= 1 {
			return strings.TrimSpace(parts[0])
		}
	}
	
	return "anonymous"
}

// 클래스명 추출
func extractClassName(line string) string {
	parts := strings.Split(line, " ")
	if len(parts) >= 2 {
		className := parts[1]
		// extends 부분 제거
		if strings.Contains(className, "extends") {
			className = strings.Split(className, "extends")[0]
		}
		return strings.TrimSpace(className)
	}
	return "UnknownClass"
}

// 메서드명 추출
func extractMethodName(line string) string {
	// 주석, 공백 제거
	line = strings.TrimSpace(line)
	
	// 접근 제한자 등 제거
	keywords := []string{"static", "async", "public", "protected", "private"}
	for _, keyword := range keywords {
		if strings.HasPrefix(line, keyword+" ") {
			line = strings.TrimSpace(line[len(keyword):])
		}
	}
	
	// 괄호 전까지 추출
	parenIndex := strings.Index(line, "(")
	if parenIndex > 0 {
		return strings.TrimSpace(line[:parenIndex])
	}
	
	return "unknown"
}

// GetLanguage는 파서가 처리하는 언어를 반환합니다.
func (p *JavaScriptParser) GetLanguage() string {
	return "JavaScript"
}

// GetFileExtensions는 파서가 처리할 수 있는 파일 확장자를 반환합니다.
func (p *JavaScriptParser) GetFileExtensions() []string {
	return []string{".js"}
}
