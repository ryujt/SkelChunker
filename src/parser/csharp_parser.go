package parser

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"SkelChunker/src/model"
	"strings"
	"fmt"
)

// CSharpParser는 C# 소스 코드를 분석하는 파서입니다.
type CSharpParser struct {
	content []byte
	pos     int
	line    int
	col     int
	tokens  []Token
}

// Token은 C# 코드의 토큰을 나타냅니다.
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Col     int
	Start   int
	End     int
}

// TokenType은 토큰의 종류를 나타냅니다.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdentifier
	TokenKeyword
	TokenString
	TokenNumber
	TokenOperator
	TokenPunctuation
	TokenComment
	TokenWhitespace
)

// 키워드 목록
var keywords = map[string]bool{
	"class":      true,
	"interface":  true,
	"struct":     true,
	"record":     true,
	"public":     true,
	"private":    true,
	"protected":  true,
	"internal":   true,
	"static":     true,
	"virtual":    true,
	"override":   true,
	"abstract":   true,
	"sealed":     true,
	"partial":    true,
	"get":        true,
	"set":        true,
	"operator":   true,
	"namespace":  true,
	"using":      true,
	"readonly":   true,
	"const":      true,
	"new":        true,
	"return":     true,
	"void":       true,
	"if":         true,
	"else":       true,
	"for":        true,
	"foreach":    true,
	"while":      true,
	"do":         true,
	"switch":     true,
	"case":       true,
	"default":    true,
	"break":      true,
	"continue":   true,
	"goto":       true,
	"try":        true,
	"catch":      true,
	"finally":    true,
	"throw":      true,
}

// NewCSharpParser는 새로운 C# 파서를 생성합니다.
func NewCSharpParser() *CSharpParser {
	return &CSharpParser{}
}

// Parse는 소스 코드를 분석하여 스켈레톤과 청크를 반환합니다.
func (p *CSharpParser) Parse(sourceCode string) ([]model.SkeletonNode, []model.Chunk, error) {
	p.content = []byte(sourceCode)
	p.pos = 0
	p.line = 1
	p.col = 1
	p.tokens = []Token{}

	// 토큰화
	if err := p.tokenize(); err != nil {
		return nil, nil, fmt.Errorf("tokenization failed: %w", err)
	}

	// 구문 분석
	nodes, chunks, err := p.parseTokens()
	if err != nil {
		return nil, nil, fmt.Errorf("parsing failed: %w", err)
	}

	// 클래스/메서드가 없는 파일인 경우 전체 파일을 청크로 추가
	if len(nodes) == 0 {
		// 파일 전체 MD5 계산
		fileMD5 := calculateMD5(sourceCode)
		
		// 파일 전체를 하나의 청크로 추가
		chunks = []model.Chunk{
			{
				MD5:  fileMD5,
				Text: sourceCode,
			},
		}
	}

	return nodes, chunks, nil
}

// tokenize는 소스 코드를 토큰으로 분리합니다.
func (p *CSharpParser) tokenize() error {
	reader := bufio.NewReader(bytes.NewReader(p.content))
	var buffer strings.Builder
	inString := false
	inChar := false
	inLineComment := false
	inBlockComment := false
	
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		p.pos++
		p.col++

		if ch == '\n' {
			p.line++
			p.col = 1
			
			if inLineComment {
				p.tokens = append(p.tokens, Token{
					Type:  TokenComment,
					Value: buffer.String(),
					Line:  p.line - 1,
					Col:   p.col,
					Start: p.pos - buffer.Len() - 1,
					End:   p.pos - 1,
				})
				buffer.Reset()
				inLineComment = false
			}
		}

		// 문자열 처리
		if inString {
			buffer.WriteRune(ch)
			if ch == '"' && p.pos > 1 && p.content[p.pos-2] != '\\' {
				p.tokens = append(p.tokens, Token{
					Type:  TokenString,
					Value: buffer.String(),
					Line:  p.line,
					Col:   p.col - buffer.Len(),
					Start: p.pos - buffer.Len(),
					End:   p.pos,
				})
				buffer.Reset()
				inString = false
			}
			continue
		}

		// 문자 처리
		if inChar {
			buffer.WriteRune(ch)
			if ch == '\'' && p.pos > 1 && p.content[p.pos-2] != '\\' {
				p.tokens = append(p.tokens, Token{
					Type:  TokenString,
					Value: buffer.String(),
					Line:  p.line,
					Col:   p.col - buffer.Len(),
					Start: p.pos - buffer.Len(),
					End:   p.pos,
				})
				buffer.Reset()
				inChar = false
			}
			continue
		}

		// 한 줄 주석 처리
		if inLineComment {
			buffer.WriteRune(ch)
			continue
		}

		// 블록 주석 처리
		if inBlockComment {
			if ch == '/' && p.pos > 1 && p.content[p.pos-2] == '*' {
				buffer.WriteRune(ch)
				p.tokens = append(p.tokens, Token{
					Type:  TokenComment,
					Value: buffer.String(),
					Line:  p.line,
					Col:   p.col - buffer.Len(),
					Start: p.pos - buffer.Len(),
					End:   p.pos,
				})
				buffer.Reset()
				inBlockComment = false
			} else {
				buffer.WriteRune(ch)
			}
			continue
		}

		// 새로운 토큰 시작
		switch {
		case ch == '"':
			inString = true
			buffer.WriteRune(ch)

		case ch == '\'':
			inChar = true
			buffer.WriteRune(ch)

		case ch == '/' && p.pos < len(p.content) && p.content[p.pos] == '/':
			inLineComment = true
			buffer.WriteString("//")
			p.pos++
			reader.ReadRune()
			p.col++

		case ch == '/' && p.pos < len(p.content) && p.content[p.pos] == '*':
			inBlockComment = true
			buffer.WriteString("/*")
			p.pos++
			reader.ReadRune()
			p.col++

		case isLetter(ch) || ch == '_':
			buffer.WriteRune(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || (!isLetter(nextCh) && !isDigit(nextCh) && nextCh != '_') {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				buffer.WriteRune(nextCh)
				p.pos++
				p.col++
			}

			word := buffer.String()
			tokenType := TokenIdentifier
			if keywords[word] {
				tokenType = TokenKeyword
			}

			p.tokens = append(p.tokens, Token{
				Type:  tokenType,
				Value: word,
				Line:  p.line,
				Col:   p.col - buffer.Len(),
				Start: p.pos - buffer.Len(),
				End:   p.pos,
			})
			buffer.Reset()

		case isDigit(ch):
			buffer.WriteRune(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || (!isDigit(nextCh) && nextCh != '.') {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				buffer.WriteRune(nextCh)
				p.pos++
				p.col++
			}

			p.tokens = append(p.tokens, Token{
				Type:  TokenNumber,
				Value: buffer.String(),
				Line:  p.line,
				Col:   p.col - buffer.Len(),
				Start: p.pos - buffer.Len(),
				End:   p.pos,
			})
			buffer.Reset()

		case isOperator(ch):
			buffer.WriteRune(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || !isOperator(nextCh) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				buffer.WriteRune(nextCh)
				p.pos++
				p.col++
			}

			p.tokens = append(p.tokens, Token{
				Type:  TokenOperator,
				Value: buffer.String(),
				Line:  p.line,
				Col:   p.col - buffer.Len(),
				Start: p.pos - buffer.Len(),
				End:   p.pos,
			})
			buffer.Reset()

		case isPunctuation(ch):
			p.tokens = append(p.tokens, Token{
				Type:  TokenPunctuation,
				Value: string(ch),
				Line:  p.line,
				Col:   p.col,
				Start: p.pos - 1,
				End:   p.pos,
			})

		case isWhitespace(ch):
			p.tokens = append(p.tokens, Token{
				Type:  TokenWhitespace,
				Value: string(ch),
				Line:  p.line,
				Col:   p.col,
				Start: p.pos - 1,
				End:   p.pos,
			})
		}
	}

	return nil
}

// cleanupText는 텍스트에서 불필요한 문자나 흰색 공간을 정리합니다.
func cleanupText(text string) string {
	// 원본 텍스트를 그대로 반환
	return text
}

// parseTokens는 토큰을 분석하여 스켈레톤과 청크를 생성합니다.
func (p *CSharpParser) parseTokens() ([]model.SkeletonNode, []model.Chunk, error) {
	var nodes []model.SkeletonNode
	var chunks []model.Chunk

	// 전체 소스 코드 텍스트
	originalSource := string(p.content)

	// 클래스/인터페이스/구조체/레코드가 있는지 확인
	hasStructuredContent := false

	for i := 0; i < len(p.tokens); i++ {
		token := p.tokens[i]

		// 클래스/인터페이스/구조체/레코드 처리
		if token.Type == TokenKeyword && (token.Value == "class" || token.Value == "interface" || token.Value == "struct" || token.Value == "record") {
			hasStructuredContent = true
			// 클래스 이름 찾기
			className := ""
			classStart := i
			classNamePos := i + 1
			
			for j := i + 1; j < len(p.tokens) && className == ""; j++ {
				if p.tokens[j].Type == TokenIdentifier {
					className = p.tokens[j].Value
					classNamePos = j
					break
				}
			}
			
			if className == "" {
				continue // 클래스 이름을 찾지 못함
			}

			// 클래스 끝 위치 찾기
			classEnd := p.findBlockEnd(classNamePos)
			if classEnd == -1 {
				continue
			}

			// 클래스 영역 추출
			classStartOffset := p.tokens[classStart].Start
			classEndOffset := p.tokens[classEnd].End
			
			if classStartOffset >= 0 && classEndOffset <= len(originalSource) && classStartOffset < classEndOffset {
				classContent := originalSource[classStartOffset:classEndOffset]
				classContent = cleanupText(classContent)

				// 스켈레톤 노드 생성
				classNode := &model.SkeletonNode{
					Type:    "class",
					Name:    className,
					MD5:     "",
					Members: []model.Member{},
				}

				// 메서드 찾기
				methods := p.findMethodsInRange(classNamePos, classEnd)
				for _, method := range methods {
					// 메서드 코드 전체를 직접 추출한 내용 사용
					methodContent := method.content
					// 전체 메서드 코드를 정리
					methodContent = cleanupText(methodContent)
					// MD5 해시 계산
					methodMD5 := calculateMD5(methodContent)

					// 멤버 추가
					classNode.Members = append(classNode.Members, model.Member{
						Type: "method",
						Name: method.name,
						MD5:  methodMD5,
					})

					// 청크 추가 (변경된 메서드만 새로 생성)
					chunks = append(chunks, model.Chunk{
						MD5:  methodMD5,
						Text: methodContent,
					})
				}

				nodes = append(nodes, *classNode)
			}

			// 다음 토큰으로 건너뛰기
			i = classEnd
		}
	}

	// 구조화된 내용이 없거나 nodes가 비어있는 경우 파일 전체를 청크로 처리
	if !hasStructuredContent || len(nodes) == 0 {
		// 파일 전체 내용을 청크로 추가
		fileContent := cleanupText(originalSource)
		fileMD5 := calculateMD5(fileContent)
		chunks = append(chunks, model.Chunk{
			MD5:  fileMD5,
			Text: fileContent,
		})
	}

	return nodes, chunks, nil
}

// extractMethodContent는 메서드의 전체 내용을 추출합니다. 
// 원본 소스에서 직접 추출하므로 토큰 범위 문제를 해결합니다.
func extractMethodContent(source string, startLine, endLine int) string {
	lines := strings.Split(source, "\n")
	
	if startLine < 0 {
		startLine = 0
	}
	
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	
	// 시작 라인 조정 (주석 포함)
	for i := startLine - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//") || line == "" {
			startLine = i
		} else {
			break
		}
	}
	
	var result strings.Builder
	for i := startLine; i <= endLine; i++ {
		result.WriteString(lines[i])
		if i < endLine {
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// findMethodsInRange는 주어진 범위 내에서 모든 메서드를 찾습니다.
func (p *CSharpParser) findMethodsInRange(start, end int) []methodInfo {
	var methods []methodInfo
	
	// 원본 소스
	originalSource := string(p.content)
	
	for i := start; i < end; i++ {
		if isMethodStart(p.tokens, i) {
			methodName := getMethodName(p.tokens, i)
			
			// 메서드 시작 위치 찾기 - 라인 번호 기준
			startLine := p.tokens[i].Line - 1  // 0-based 라인 번호로 변환
			
			// 메서드 본문 블록 찾기
			methodBodyStart := i
			for j := i; j < end; j++ {
				if p.tokens[j].Type == TokenPunctuation && p.tokens[j].Value == "{" {
					methodBodyStart = j
					break
				}
			}
			
			// 메서드 끝 찾기
			methodEnd := p.findBlockEnd(methodBodyStart)
			if methodEnd == -1 || methodEnd > end {
				continue
			}
			
			// 메서드 끝 라인 번호
			endLine := p.tokens[methodEnd].Line - 1  // 0-based 라인 번호
			
			// 메서드 내용 추출
			methodContent := extractMethodContent(originalSource, startLine, endLine)
			
			methods = append(methods, methodInfo{
				name:     methodName,
				startPos: p.tokens[i].Start,
				endPos:   p.tokens[methodEnd].End,
				content:  methodContent,  // 직접 추출한 메서드 내용
			})
			
			// 다음 토큰으로 건너뛰기
			i = methodEnd
		}
	}
	
	return methods
}

// findBlockEnd는 중괄호 블록의 끝을 찾습니다.
func (p *CSharpParser) findBlockEnd(start int) int {
	braceLevel := 0
	openBraceFound := false
	
	for i := start; i < len(p.tokens); i++ {
		if p.tokens[i].Type == TokenPunctuation {
			if p.tokens[i].Value == "{" {
				openBraceFound = true
				braceLevel++
			} else if p.tokens[i].Value == "}" {
				braceLevel--
				if openBraceFound && braceLevel == 0 {
					return i
				}
			}
		}
	}
	
	return -1
}

// extractTextBetween은 두 위치 사이의 원본 텍스트를 추출합니다.
func (p *CSharpParser) extractTextBetween(start, end int) string {
	if start >= len(p.tokens) || end >= len(p.tokens) || start < 0 || end < 0 {
		return ""
	}
	
	startPos := p.tokens[start].Start
	endPos := p.tokens[end].End
	
	if startPos < 0 {
		startPos = 0
	}
	
	if endPos > len(p.content) {
		endPos = len(p.content)
	}
	
	if startPos >= endPos || startPos >= len(p.content) || endPos <= 0 {
		return ""
	}
	
	return string(p.content[startPos:endPos])
}

// isMethodStart는 현재 위치가 메서드의 시작인지 확인합니다.
func isMethodStart(tokens []Token, pos int) bool {
	if pos+2 >= len(tokens) {
		return false
	}

	// 제어자 건너뛰기 (public, private, static 등)
	i := pos
	for i < len(tokens) && tokens[i].Type == TokenKeyword && isModifier(tokens[i].Value) {
		i++
	}
	
	if i+2 >= len(tokens) {
		return false
	}

	// 반환 타입 + 이름 + ( 패턴
	if ((tokens[i].Type == TokenIdentifier || tokens[i].Type == TokenKeyword) && 
		tokens[i+1].Type == TokenIdentifier && 
		i+2 < len(tokens) && tokens[i+2].Type == TokenPunctuation && tokens[i+2].Value == "(") {
		return true
	}

	// 속성 패턴 (Type Name { get; set; })
	if ((tokens[i].Type == TokenIdentifier || tokens[i].Type == TokenKeyword) && 
		tokens[i+1].Type == TokenIdentifier && 
		i+2 < len(tokens) && tokens[i+2].Type == TokenPunctuation && tokens[i+2].Value == "{") {
		for j := i+3; j < len(tokens) && j < i+10; j++ {
			if tokens[j].Type == TokenKeyword && (tokens[j].Value == "get" || tokens[j].Value == "set") {
				return true
			}
		}
	}

	// 생성자 패턴 (ClassName(params))
	if tokens[i].Type == TokenIdentifier && 
		i+1 < len(tokens) && tokens[i+1].Type == TokenPunctuation && tokens[i+1].Value == "(" {
		return true
	}

	// 연산자 오버로드 패턴 (operator Type())
	if tokens[i].Type == TokenKeyword && tokens[i].Value == "operator" {
		return true
	}

	return false
}

// getMethodName은 메서드 이름을 추출합니다.
func getMethodName(tokens []Token, pos int) string {
	// 제어자 건너뛰기 (public, private, static 등)
	i := pos
	for i < len(tokens) && tokens[i].Type == TokenKeyword && isModifier(tokens[i].Value) {
		i++
	}
	
	if i+1 >= len(tokens) {
		return ""
	}
	
	// 기본 메서드 패턴 (ReturnType Name(...))
	if ((tokens[i].Type == TokenIdentifier || tokens[i].Type == TokenKeyword) && 
		tokens[i+1].Type == TokenIdentifier) {
		return tokens[i+1].Value
	}
	
	// 생성자 패턴 (ClassName(...))
	if tokens[i].Type == TokenIdentifier && 
		i+1 < len(tokens) && tokens[i+1].Type == TokenPunctuation && tokens[i+1].Value == "(" {
		return tokens[i].Value
	}
	
	// 연산자 오버로드 패턴 (operator Type(...))
	if tokens[i].Type == TokenKeyword && tokens[i].Value == "operator" && i+1 < len(tokens) {
		return "operator" + tokens[i+1].Value
	}
	
	return ""
}

// isModifier는 주어진 키워드가 접근 제한자인지 확인합니다.
func isModifier(word string) bool {
	modifiers := map[string]bool{
		"public":    true,
		"private":   true,
		"protected": true,
		"internal":  true,
		"static":    true,
		"virtual":   true,
		"override":  true,
		"abstract":  true,
		"sealed":    true,
		"readonly":  true,
		"async":     true,
		"extern":    true,
		"unsafe":    true,
		"new":       true,
		"partial":   true,
	}
	
	return modifiers[word]
}

// GetLanguage는 파서가 처리하는 언어를 반환합니다.
func (p *CSharpParser) GetLanguage() string {
	return "C#"
}

// GetFileExtensions는 파서가 처리할 수 있는 파일 확장자를 반환합니다.
func (p *CSharpParser) GetFileExtensions() []string {
	return []string{".cs"}
}

// 유틸리티 함수들
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isOperator(ch rune) bool {
	return strings.ContainsRune("+-*/%=<>!&|^~", ch)
}

func isPunctuation(ch rune) bool {
	return strings.ContainsRune("(){}[];,.", ch)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func needsSpace(current, next Token) bool {
	if current.Type == TokenPunctuation || next.Type == TokenPunctuation {
		return false
	}
	if current.Type == TokenWhitespace || next.Type == TokenWhitespace {
		return false
	}
	return true
}

func calculateMD5(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// methodInfo는 메서드 정보를 저장하는 구조체입니다.
type methodInfo struct {
	name     string
	startPos int
	endPos   int
	content  string  // 직접 추출한 메서드 내용
} 