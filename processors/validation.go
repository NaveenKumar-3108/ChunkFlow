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

func ValidateSocketAudioUpload(pUserID, pSessionID, pTimestamp string, pHandler *multipart.FileHeader) (string, bool) {

	if pUserID == "" || pSessionID == "" || pTimestamp == "" {
		log.Println("Error:PVA01 Missing required metadata")
		return "Missing required metadata", false
	}

	return "", true
}
