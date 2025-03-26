package embeddings

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// EmbeddingService는 텍스트 임베딩을 생성하는 서비스 인터페이스입니다.
type EmbeddingService interface {
	CreateEmbedding(text string) ([]float32, error)
	ChunkText(text string, maxSize int) ([]string, error)
}

// OpenAIEmbedding은 OpenAI API를 사용하여 텍스트 임베딩을 생성하는 서비스 구현체입니다.
type OpenAIEmbedding struct {
	apiKey      string
	client      *openai.Client
	modelName   openai.EmbeddingModel
	testMode    bool
	vectorDim   int
}

// NewOpenAIEmbedding은 OpenAIEmbedding의 새 인스턴스를 생성합니다.
func NewOpenAIEmbedding(apiKey string, modelName string, testMode bool, vectorDim int) *OpenAIEmbedding {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	var embeddingModel openai.EmbeddingModel
	if modelName == "" || modelName == "text-embedding-3-large" {
		embeddingModel = openai.LargeEmbedding3
	} else if modelName == "text-embedding-3-small" {
		embeddingModel = openai.SmallEmbedding3
	} else if modelName == "text-embedding-ada-002" {
		embeddingModel = openai.AdaEmbeddingV2
	} else {
		// 기본값으로 설정
		embeddingModel = openai.LargeEmbedding3
	}

	client := openai.NewClient(apiKey)

	return &OpenAIEmbedding{
		apiKey:    apiKey,
		client:    client,
		modelName: embeddingModel,
		testMode:  testMode,
		vectorDim: vectorDim,
	}
}

// CreateEmbedding은 주어진 텍스트에 대한 임베딩을 생성합니다.
func (e *OpenAIEmbedding) CreateEmbedding(text string) ([]float32, error) {
	// 테스트 모드인 경우 더미 임베딩 반환
	if e.testMode {
		dummyEmbedding := make([]float32, e.vectorDim)
		for i := 0; i < e.vectorDim; i++ {
			dummyEmbedding[i] = float32(i) / float32(e.vectorDim)
		}
		return dummyEmbedding, nil
	}

	// 임베딩 요청 생성
	embeddingReq := openai.EmbeddingRequest{
		Input: []string{text},
		Model: e.modelName,
	}

	// OpenAI API 호출
	ctx := context.Background()
	resp, err := e.client.CreateEmbeddings(ctx, embeddingReq)
	if err != nil {
		return nil, fmt.Errorf("임베딩 생성 오류: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("OpenAI API에서 임베딩을 반환하지 않았습니다")
	}

	return resp.Data[0].Embedding, nil
}

// ChunkText는 텍스트를 최대 크기 제한에 맞게 여러 청크로 분할합니다.
func (e *OpenAIEmbedding) ChunkText(text string, maxSize int) ([]string, error) {
	if len(text) <= maxSize {
		return []string{text}, nil
	}

	// 텍스트를 의미 있는 단위로 분할
	chunks := []string{}
	lines := strings.Split(text, "\n")
	currentChunk := strings.Builder{}
	currentSize := 0
	
	// 중괄호 레벨 추적
	braceLevel := 0
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		lineSize := len(line) + 1 // +1 for newline
		
		// 중괄호 레벨 업데이트
		braceLevel += strings.Count(line, "{") - strings.Count(line, "}")
		
		// 현재 라인이 메서드/클래스 시작인지 확인
		isBlockStart := isCodeBlockStart(line)
		
		// 이전 청크가 있고, 새로운 블록이 시작되거나 크기 제한을 초과하는 경우
		if currentSize > 0 && (isBlockStart || currentSize+lineSize > maxSize) {
			// 현재 청크가 비어있지 않으면 저장
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
				currentSize = 0
			}
		}
		
		// 현재 라인이 단독으로 최대 크기를 초과하는 경우
		if lineSize > maxSize {
			// 라인을 더 작은 청크로 분할
			subChunks := splitLongLine(line, maxSize)
			chunks = append(chunks, subChunks...)
			continue
		}
		
		// 현재 라인 추가
		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n")
		}
		currentChunk.WriteString(line)
		currentSize += lineSize
		
		// 블록이 닫히고 청크 크기가 충분히 큰 경우 청크 완료
		if braceLevel == 0 && currentSize >= maxSize/2 {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
			currentSize = 0
		}
	}
	
	// 마지막 청크가 남아있으면 추가
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks, nil
}

// isCodeBlockStart는 주어진 라인이 코드 블록의 시작인지 확인합니다.
func isCodeBlockStart(line string) bool {
	line = strings.TrimSpace(line)
	
	// 클래스 정의
	if strings.HasPrefix(line, "class ") || strings.HasPrefix(line, "interface ") ||
	   strings.HasPrefix(line, "struct ") || strings.HasPrefix(line, "enum ") {
		return true
	}
	
	// 메서드 정의
	if strings.Contains(line, "(") && strings.Contains(line, ")") &&
	   (strings.Contains(line, "public ") || strings.Contains(line, "private ") ||
		strings.Contains(line, "protected ") || strings.Contains(line, "internal ") ||
		strings.Contains(line, "void ") || strings.HasSuffix(line, ")")) {
		return true
	}
	
	// 프로퍼티 정의
	if strings.Contains(line, "get") && strings.Contains(line, "set") {
		return true
	}
	
	return false
}

// splitLongLine은 긴 라인을 여러 청크로 분할합니다.
func splitLongLine(line string, maxSize int) []string {
	var chunks []string
	
	// 주석 라인은 별도 처리
	if strings.TrimSpace(line)[:2] == "//" {
		return []string{line}
	}
	
	// 문자열 리터럴이나 주석이 포함된 경우 전체를 하나의 청크로
	if strings.Contains(line, "\"") || strings.Contains(line, "/*") {
		return []string{line}
	}
	
	// 일반적인 긴 라인 분할
	for len(line) > maxSize {
		// 적절한 분할 지점 찾기
		splitPoint := maxSize
		for splitPoint > 0 && !isSplitPoint(line[splitPoint-1]) {
			splitPoint--
		}
		
		if splitPoint == 0 {
			splitPoint = maxSize
		}
		
		chunks = append(chunks, line[:splitPoint])
		line = line[splitPoint:]
	}
	
	if len(line) > 0 {
		chunks = append(chunks, line)
	}
	
	return chunks
}

// isSplitPoint는 해당 문자에서 라인을 분할할 수 있는지 확인합니다.
func isSplitPoint(ch byte) bool {
	return ch == ' ' || ch == ',' || ch == ';' || ch == '.' ||
		   ch == '(' || ch == ')' || ch == '{' || ch == '}'
} 