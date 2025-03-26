package analyzer

import (
	"SkelChunker/src/embeddings"
	"SkelChunker/src/model"
	"SkelChunker/src/parser"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"bytes"
)

// CompactResult는 임베딩 배열을 한 줄로 저장하기 위한 구조체입니다.
type CompactResult struct {
	Path       string              `json:"path"`
	Filename   string              `json:"filename"`
	MD5        string              `json:"md5"`
	Embeddings json.RawMessage     `json:"embeddings,omitempty"`
	Skeleton   []model.SkeletonNode `json:"skeleton"`
	Chunks     []CompactChunk       `json:"chunks"`
}

// CompactChunk는 청크의 임베딩을 한 줄로 저장하기 위한 구조체입니다.
type CompactChunk struct {
	MD5        string          `json:"md5"`
	Text       string          `json:"text"`
	Embeddings json.RawMessage `json:"embeddings,omitempty"`
}

// Analyzer는 소스 코드 분석을 수행하는 구조체입니다.
type Analyzer struct {
	parserFactory *parser.ParserFactory
	embeddingService embeddings.EmbeddingService
	embeddingConfig *embeddings.Config
}

// NewAnalyzer는 새로운 Analyzer 인스턴스를 생성합니다.
func NewAnalyzer(parserFactory *parser.ParserFactory, embeddingService embeddings.EmbeddingService, embeddingConfig *embeddings.Config) *Analyzer {
	return &Analyzer{
		parserFactory: parserFactory,
		embeddingService: embeddingService,
		embeddingConfig: embeddingConfig,
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

	// SkelChunker 파일 경로 생성
	baseFileName := filepath.Base(filePath)
	ext = filepath.Ext(baseFileName)
	skelChunkerPath := filepath.Join(filepath.Dir(filePath), baseFileName[:len(baseFileName)-len(ext)]+".SkelChunker")

	// 기존 SkelChunker 파일이 있는지 확인
	var existingResult *model.AnalysisResult
	if existingData, err := os.ReadFile(skelChunkerPath); err == nil {
		existingResult = &model.AnalysisResult{}
		if err := json.Unmarshal(existingData, existingResult); err == nil {
			// 파일 전체 MD5 비교
			if existingResult.MD5 == md5Hash {
				// 파일이 변경되지 않았으므로 기존 결과 반환
				return existingResult, nil
			}
		}
	}

	// 파일 분석
	skeleton, chunks, err := parser.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// 구조화된 코드가 없는 경우 (skeleton이 nil이거나 비어있는 경우)
	if skeleton == nil || len(skeleton) == 0 {
		// 파일 전체를 하나의 청크로 생성
		chunks = []model.Chunk{
			{
				MD5:  md5Hash,
				Text: string(content),
			},
		}
	}

	// 결과 생성
	result := &model.AnalysisResult{
		Path:     filepath.Dir(filePath),
		Filename: filepath.Base(filePath),
		MD5:      md5Hash,
		Skeleton: skeleton,
		Chunks:   chunks,
	}

	// 파일 전체를 처리한 후 최종 검증
	// 청크가 비어있으면 파일 전체를 청크로 추가
	if result.Chunks == nil || len(result.Chunks) == 0 {
		result.Chunks = []model.Chunk{
			{
				MD5:  md5Hash,
				Text: string(content),
			},
		}
	}

	// 임베딩 서비스가 있는 경우 임베딩 생성 수행
	if a.embeddingService != nil {
		// 파일 전체 임베딩 생성
		fileEmbeddings, err := a.createEmbeddingsForText(string(content))
		if err != nil {
			fmt.Printf("Warning: Failed to create embeddings for file %s: %v\n", filePath, err)
		} else {
			result.Embeddings = fileEmbeddings
		}

		// 각 청크에 대한 임베딩 생성
		for i := range result.Chunks {
			// 청크 전처리
			chunkText := preprocessCodeForEmbedding(result.Chunks[i].Text)
			
			// 청크가 너무 큰 경우 분할하여 임베딩 생성
			chunkEmbeddings, err := a.createEmbeddingsForText(chunkText)
			if err != nil {
				fmt.Printf("Warning: Failed to create embeddings for chunk in file %s: %v\n", filePath, err)
				continue
			}
			
			// 청크가 분할된 경우 첫 번째 임베딩만 사용
			if len(chunkEmbeddings) > 0 {
				result.Chunks[i].Embeddings = chunkEmbeddings[0]
			}
		}
	}

	return result, nil
}

// createEmbeddingsForText는 주어진 텍스트에 대한 임베딩을 생성합니다.
func (a *Analyzer) createEmbeddingsForText(text string) ([][]float32, error) {
	// 임베딩 서비스가 없는 경우 빈 배열 반환
	if a.embeddingService == nil {
		return nil, nil
	}

	// 텍스트 전처리
	text = preprocessCodeForEmbedding(text)

	// 텍스트 크기가 제한을 초과하는 경우 청크로 분할
	chunks, err := a.embeddingService.ChunkText(text, a.embeddingConfig.MaxTextSize)
	if err != nil {
		return nil, fmt.Errorf("failed to chunk text: %w", err)
	}

	// 각 청크에 대해 임베딩 생성
	var embeddings [][]float32
	for _, chunk := range chunks {
		embedding, err := a.embeddingService.CreateEmbedding(chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding: %w", err)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// preprocessCodeForEmbedding는 코드 텍스트를 임베딩에 적합하게 전처리합니다.
func preprocessCodeForEmbedding(text string) string {
	// 불필요한 공백 제거
	text = strings.TrimSpace(text)
	
	// 연속된 빈 줄을 하나로 통합
	text = regexp.MustCompile(`\n\s*\n\s*\n`).ReplaceAllString(text, "\n\n")
	
	// 주석에서 불필요한 공백 제거
	text = regexp.MustCompile(`//\s*(.*)`).ReplaceAllString(text, "// $1")
	text = regexp.MustCompile(`/\*\s*(.*?)\s*\*/`).ReplaceAllString(text, "/* $1 */")
	
	// 중괄호 스타일 통일
	text = regexp.MustCompile(`\s*{\s*\n`).ReplaceAllString(text, " {\n")
	text = regexp.MustCompile(`\n\s*}`).ReplaceAllString(text, "\n}")
	
	return text
}

// SaveResult는 분석 결과를 파일로 저장합니다.
func (a *Analyzer) SaveResult(result *model.AnalysisResult) error {
	// 파일명 생성
	baseFileName := filepath.Base(result.Filename)
	ext := filepath.Ext(baseFileName)
	outputFileName := baseFileName[:len(baseFileName)-len(ext)] + ".SkelChunker"
	outputPath := filepath.Join(result.Path, outputFileName)

	// 임베딩 배열을 한 줄로 마샬링
	compactResult := CompactResult{
		Path:     result.Path,
		Filename: result.Filename,
		MD5:      result.MD5,
		Skeleton: result.Skeleton,
	}

	// 파일 전체 임베딩 처리
	if result.Embeddings != nil {
		// 임베딩 배열을 한 줄로 마샬링
		embeddings, err := json.Marshal(result.Embeddings)
		if err != nil {
			return fmt.Errorf("failed to marshal embeddings: %w", err)
		}
		// 줄바꿈과 공백을 제거하여 한 줄로 만들기
		compactResult.Embeddings = bytes.ReplaceAll(
			bytes.ReplaceAll(embeddings, []byte("\n"), []byte("")),
			[]byte(" "), []byte(""),
		)
	}

	// 청크 임베딩 처리
	compactResult.Chunks = make([]CompactChunk, len(result.Chunks))
	for i, chunk := range result.Chunks {
		compactChunk := CompactChunk{
			MD5:  chunk.MD5,
			Text: chunk.Text,
		}
		if chunk.Embeddings != nil {
			// 청크 임베딩을 한 줄로 마샬링
			embeddings, err := json.Marshal(chunk.Embeddings)
			if err != nil {
				return fmt.Errorf("failed to marshal chunk embedding: %w", err)
			}
			// 줄바꿈과 공백을 제거하여 한 줄로 만들기
			compactChunk.Embeddings = bytes.ReplaceAll(
				bytes.ReplaceAll(embeddings, []byte("\n"), []byte("")),
				[]byte(" "), []byte(""),
			)
		}
		compactResult.Chunks[i] = compactChunk
	}

	// 최종 JSON 변환 (임베딩 배열은 한 줄로, 나머지는 들여쓰기)
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(compactResult); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	// 파일 저장
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write result file: %w", err)
	}

	return nil
} 