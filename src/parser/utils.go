package parser

import (
	"crypto/md5"
	"encoding/hex"
)

// calculateMD5는 주어진 문자열의 MD5 해시를 계산하여 반환합니다.
func calculateMD5(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
} 