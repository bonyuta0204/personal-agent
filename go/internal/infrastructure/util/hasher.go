package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// CalculateSHA256 calculates SHA-256 hash of the given content
func CalculateSHA256(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
