package handlers

import (
	"CHUNKFLOW/processors"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createTestMultipartForm(t *testing.T) (*bytes.Buffer, string) {
	lBody := &bytes.Buffer{}
	lWriter := multipart.NewWriter(lBody)

	_ = lWriter.WriteField("user_id", "u123")
	_ = lWriter.WriteField("session_id", "s123")
	_ = lWriter.WriteField("timestamp", "2025-08-09T10:00:00Z")

	lFileWriter, lErr := lWriter.CreateFormFile("audio", "test.wav")
	if lErr != nil {
		t.Fatalf("Error creating form file: %v", lErr)
	}
	lFileWriter.Write([]byte("fake audio content"))

	lWriter.Close()
	return lBody, lWriter.FormDataContentType()
}

func TestUploadAudio_Success(t *testing.T) {
	processors.StoreAudioChunk = func(pAudioBytes []byte, pMetaData processors.AudioMetadata) (int64, error) {
		return 1, nil
	}
	lBody, lContentType := createTestMultipartForm(t)
	lReq := httptest.NewRequest("POST", "/upload", lBody)
	lReq.Header.Set("Content-Type", lContentType)

	lRR := httptest.NewRecorder()

	UploadAudio(lRR, lReq)

	if lRR.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", lRR.Code)
	}

	if !strings.Contains(lRR.Body.String(), "Audio uploaded successfully") {
		t.Errorf("Unexpected body: %s", lRR.Body.String())
	}
}

func TestUploadAudio_WrongMethod(t *testing.T) {
	lReq := httptest.NewRequest("GET", "/upload", nil)
	lRR := httptest.NewRecorder()

	UploadAudio(lRR, lReq)

	if lRR.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", lRR.Code)
	}
}
