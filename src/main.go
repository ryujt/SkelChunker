package main

import (
	"SkelChunker/src/analyzer"
	"SkelChunker/src/config"
	"SkelChunker/src/embeddings"
	"SkelChunker/src/parser"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// shouldIgnoreFolder는 주어진 경로가 무시해야 할 폴더인지 확인합니다.
func shouldIgnoreFolder(path string, ignoreFolders []string) bool {
	// 경로를 슬래시로 분리하여 각 부분을 확인
	parts := strings.Split(path, string(os.PathSeparator))
	
	// 각 부분이 무시 폴더 목록에 있는지 확인
	for _, part := range parts {
		for _, ignoreFolder := range ignoreFolders {
			if part == ignoreFolder {
				return true
			}
		}
	}
	return false
}

func main() {
	// 커맨드 라인 인자 파싱
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// 설정 로드
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 파서 팩토리 생성
	parserFactory := parser.NewParserFactory()

	// 파서 등록
	parserFactory.RegisterParser(parser.NewCSharpParser())
	parserFactory.RegisterParser(parser.NewJavaScriptParser())

	// 임베딩 서비스 설정
	var embeddingService embeddings.EmbeddingService
	var embeddingConfig *embeddings.Config
	
	if cfg.Embedding.Enabled {
		embeddingConfig = &embeddings.Config{
			APIKey:     cfg.Embedding.APIKey,
			ModelName:  cfg.Embedding.ModelName,
			VectorDim:  cfg.Embedding.VectorDim,
			MaxTextSize: cfg.Embedding.MaxTextSize,
		}
		
		embeddingService = embeddings.NewOpenAIEmbedding(
			cfg.Embedding.APIKey,
			cfg.Embedding.ModelName,
			cfg.Embedding.VectorDim,
		)
		
		fmt.Println("Embedding service initialized with model:", cfg.Embedding.ModelName)
	} else {
		fmt.Println("Embedding service is disabled")
	}

	// 분석기 초기화
	analyzer := analyzer.NewAnalyzer(parserFactory, embeddingService, embeddingConfig)

	// 각 폴더 처리
	for _, folder := range cfg.Folders {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 디렉토리인 경우 무시 폴더 체크
			if info.IsDir() {
				if shouldIgnoreFolder(path, cfg.IgnoreFolders) {
					fmt.Printf("Skipping ignored folder: %s\n", path)
					return filepath.SkipDir
				}
				return nil
			}

			// 파일 확장자 확인
			ext := filepath.Ext(path)
			if _, exists := cfg.Parsers[ext]; !exists {
				return nil
			}

			// 파일 분석
			result, err := analyzer.AnalyzeFile(path)
			if err != nil {
				fmt.Printf("Error analyzing file %s: %v\n", path, err)
				return nil
			}

			// 결과 저장
			if err := analyzer.SaveResult(result); err != nil {
				fmt.Printf("Error saving result for file %s: %v\n", path, err)
				return nil
			}

			fmt.Printf("Successfully processed: %s\n", path)
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing folder %s: %v\n", folder, err)
		}
	}
}
