package processors

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateChecksum(pData []byte) string {
	h := sha256.New()
	h.Write(pData)
	return hex.EncodeToString(h.Sum(nil))
}

func FakeTranscript(pData []byte) string {
	return "Welcome "
}

func TransformAudio(pData []byte) (string, string) {
	return GenerateChecksum(pData), FakeTranscript(pData)
}
