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

	return &config, nil
} 