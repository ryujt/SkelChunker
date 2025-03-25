package parser

import "fmt"

// ParserFactory는 파일 확장자에 따라 적절한 파서를 생성하는 팩토리입니다.
type ParserFactory struct {
	extensionToParser map[string]Parser
}

// NewParserFactory는 새로운 ParserFactory 인스턴스를 생성합니다.
func NewParserFactory() *ParserFactory {
	return &ParserFactory{
		extensionToParser: make(map[string]Parser),
	}
}

// RegisterParser는 파서를 등록합니다.
func (f *ParserFactory) RegisterParser(parser Parser) {
	for _, ext := range parser.GetFileExtensions() {
		f.extensionToParser[ext] = parser
	}
}

// GetParser는 파일 확장자에 맞는 파서를 반환합니다.
func (f *ParserFactory) GetParser(fileExtension string) (Parser, error) {
	parser, exists := f.extensionToParser[fileExtension]
	if !exists {
		return nil, fmt.Errorf("no parser registered for extension: %s", fileExtension)
	}
	return parser, nil
} 