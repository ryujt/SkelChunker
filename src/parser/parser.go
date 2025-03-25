package parser

import (
	"SkelChunker/src/model"
)

// Parser는 소스 코드를 분석하여 스켈레톤과 청크를 추출하는 인터페이스입니다.
type Parser interface {
	// Parse는 소스 코드 문자열을 받아 스켈레톤 노드와 청크들을 반환합니다.
	Parse(sourceCode string) ([]model.SkeletonNode, []model.Chunk, error)
	
	// GetLanguage는 파서가 처리할 수 있는 프로그래밍 언어를 반환합니다.
	GetLanguage() string
	
	// GetFileExtensions는 파서가 처리할 수 있는 파일 확장자들을 반환합니다.
	GetFileExtensions() []string
} 