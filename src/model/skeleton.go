package model

// Member는 클래스의 멤버(메서드)를 나타내는 구조체입니다.
type Member struct {
	MD5  string `json:"md5"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// SkeletonNode는 코드의 구조적 요소(클래스, 함수)를 나타내는 구조체입니다.
type SkeletonNode struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Members []Member `json:"members,omitempty"`
	MD5     string   `json:"md5,omitempty"`
}

// Chunk는 코드의 실제 구현 내용을 담는 구조체입니다.
type Chunk struct {
	MD5        string    `json:"md5"`
	Text       string    `json:"text"`
	Embeddings []float64 `json:"embeddings"`
}

// AnalysisResult는 파일 분석 결과를 나타내는 구조체입니다.
type AnalysisResult struct {
	Path       string         `json:"path"`
	Filename   string        `json:"filename"`
	MD5        string        `json:"md5"`
	Embeddings []float64     `json:"embeddings"`
	Skeleton   []SkeletonNode `json:"skeleton"`
	Chunks     []Chunk       `json:"chunks"`
} 