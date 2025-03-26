package config

import (
	"encoding/json"
	"os"
)

// Config는 애플리케이션의 설정을 담는 구조체입니다.
type Config struct {
	Folders       []string          `json:"folders"`
	IgnoreFolders []string          `json:"ignore-folders"`
	Parsers       map[string]string `json:"parsers"`
	Embedding     EmbeddingConfig   `json:"embedding"`
}

// EmbeddingConfig는 임베딩 관련 설정을 담는 구조체입니다.
type EmbeddingConfig struct {
	Enabled    bool   `json:"enabled"`
	APIKey     string `json:"api-key"`
	ModelName  string `json:"model-name"`
	VectorDim  int    `json:"vector-dim"`
	TestMode   bool   `json:"test-mode"`
	MaxTextSize int   `json:"max-text-size"`
}

// LoadConfig는 지정된 경로의 JSON 파일에서 설정을 로드합니다.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 임베딩 설정 기본값 설정
	if config.Embedding.ModelName == "" {
		config.Embedding.ModelName = "text-embedding-3-large"
	}

	if config.Embedding.VectorDim == 0 {
		config.Embedding.VectorDim = 3072
	}

	if config.Embedding.MaxTextSize == 0 {
		config.Embedding.MaxTextSize = 24 * 1024 // 24KB 제한
	}

	return &config, nil
} 