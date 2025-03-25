package analyzer

import (
	"SkelChunker/src/model"
	"SkelChunker/src/parser"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Analyzer는 소스 코드 분석을 수행하는 구조체입니다.
type Analyzer struct {
	parserFactory *parser.ParserFactory
}

// NewAnalyzer는 새로운 Analyzer 인스턴스를 생성합니다.
func NewAnalyzer(parserFactory *parser.ParserFactory) *Analyzer {
	return &Analyzer{
		parserFactory: parserFactory,
	}
}

// AnalyzeFile은 단일 파일을 분석하여 결과를 반환합니다.
func (a *Analyzer) AnalyzeFile(filePath string) (*model.AnalysisResult, error) {
	// 파일 읽기
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 파일 확장자 확인
	ext := filepath.Ext(filePath)
	parser, err := a.parserFactory.GetParser(ext)
	if err != nil {
		return nil, fmt.Errorf("no parser available for extension %s: %w", ext, err)
	}

	// MD5 해시 계산
	hash := md5.Sum(content)
	md5Hash := hex.EncodeToString(hash[:])

	// 파일 분석
	skeleton, chunks, err := parser.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// 결과 생성
	result := &model.AnalysisResult{
		Path:     filepath.Dir(filePath),
		Filename: filepath.Base(filePath),
		MD5:      md5Hash,
		Skeleton: skeleton,
		Chunks:   chunks,
	}

	return result, nil
}

// SaveResult는 분석 결과를 파일로 저장합니다.
func (a *Analyzer) SaveResult(result *model.AnalysisResult) error {
	// 파일명 생성
	baseFileName := filepath.Base(result.Filename)
	ext := filepath.Ext(baseFileName)
	outputFileName := baseFileName[:len(baseFileName)-len(ext)] + ".SkelChunker"
	outputPath := filepath.Join(result.Path, outputFileName)

	// JSON 변환
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	// 파일 저장
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write result file: %w", err)
	}

	return nil
} 