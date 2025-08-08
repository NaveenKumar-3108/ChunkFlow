package processors

import (
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func ValidateAudioUpload(pUserID, pSessionID, pTimestamp string, pHandler *multipart.FileHeader) (string, bool) {

	if pUserID == "" || pSessionID == "" || pTimestamp == "" {
		log.Println("Error:PVA01 Missing required metadata")
		return "Missing required metadata", false
	}

	lMimeType := pHandler.Header.Get("Content-Type")
	if !strings.HasPrefix(lMimeType, "audio/") {
		log.Println("Error:PVA02 Unsupported media type:", lMimeType)
		return "Unsupported media type", false
	}

	lAllowedMIMEs := map[string]bool{
		"audio/wav":   true,
		"audio/x-wav": true,
		"audio/mpeg":  true,
		"audio/ogg":   true,
	}

	if !lAllowedMIMEs[lMimeType] {
		log.Println("Error:PVA03 Unsupported audio format:", lMimeType)
		return "Unsupported audio format", false
	}

	const lMaxSize = 10 * 1024 * 1024
	if pHandler.Size > lMaxSize {
		log.Println("Error:PVA04 File too large:", pHandler.Size)
		return "File size exceeds limit", false
	}

	if pHandler.Size == 0 {
		log.Println("Error:PVA05 Empty file")
		return "Empty audio file", false
	}

	lExt := strings.ToLower(filepath.Ext(pHandler.Filename))
	lAllowedExts := map[string]bool{
		".wav": true,
		".mp3": true,
		".ogg": true,
	}
	if !lAllowedExts[lExt] {
		log.Println("Error:PVA06 Invalid file extension:", lExt)
		return "Invalid audio file extension", false
	}

	return "", true
}
