package processors

import (
	"CHUNKFLOW/internal/db"
	"context"
	"log"
	"mime/multipart"
	"sync"
)

type AudioMetadata struct {
	ChunkID     string `json:"chunkid,omitempty"`
	UserID      string `json:"user_id"`
	SessionID   string `json:"session_id"`
	Timestamp   string `json:"timestamp"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Duration    string `json:"duration"`
	Transcript  string `json:"transcript"`
	Checksum    string `json:"checksum"`
	Size        string `json:"size"`
}
type AudioJob struct {
	Ctx         context.Context
	UserID      string
	SessionID   string
	Timestamp   string
	FileName    string
	ContentType string
	FileArr     []byte
	Handler     *multipart.FileHeader
	File        multipart.File
	Error       error
	Duration    string
	Transcript  string
	Checksum    string
	Size        string
}
type Response struct {
	ChunkID int64  `json:"chunkid"`
	Status  string `json:"status"`
	ErrMsg  string `json:"errmsg"`
	Msg     string `json:"msg"`
}

var storeAudioMutex sync.Mutex

var StoreAudioChunk = func(pAudioBytes []byte, pMetaData AudioMetadata) (int64, error) {
	log.Println("StoreAudioChunk(+)")
	var lChunkID int64

	storeAudioMutex.Lock()
	defer storeAudioMutex.Unlock()

	query := `
    INSERT INTO audio_chunks (user_id, session_id, timestamp,
        file_name, content_type, audio_data,
        duration, size, transcript, checksum
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	lResult, lErr := db.GDB.Exec(query,
		pMetaData.UserID,
		pMetaData.SessionID,
		pMetaData.Timestamp,
		pMetaData.FileName,
		pMetaData.ContentType,
		pAudioBytes,
		pMetaData.Duration,
		pMetaData.Size,
		pMetaData.Transcript,
		pMetaData.Checksum,
	)

	if lErr != nil {
		log.Println("ERROR:PSA01", lErr)
		return lChunkID, lErr
	}

	lChunkID, lErr = lResult.LastInsertId()
	if lErr != nil {
		log.Println("ERROR:PSA02 failed to get last insert id", lErr)
		return lChunkID, lErr
	}

	log.Println("StoreAudioChunk(-)")
	return lChunkID, nil
}
