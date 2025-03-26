package embeddings

// Config는 임베딩 서비스 설정을 담고 있는 구조체입니다.
type Config struct {
	APIKey     string
	ModelName  string
	VectorDim  int
	MaxTextSize int
}

// DefaultConfig는 기본 설정값으로 Config 인스턴스를 반환합니다.
func DefaultConfig() *Config {
	return &Config{
		APIKey:     "",
		ModelName:  "text-embedding-3-large",
		VectorDim:  3072,
		MaxTextSize: 24 * 1024, // 24KB 제한
	}
} 