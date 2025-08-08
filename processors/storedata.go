package processors

import (
	"CHUNKFLOW/internal/db"
	"context"
	"log"
	"mime/multipart"
)

type AudioMetadata struct {
	ChunkID     string
	UserID      string
	SessionID   string
	Timestamp   string
	FileName    string
	ContentType string
	Duration    string
	Transcript  string
	Checksum    string
	Size        string
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

func StoreAudioChunk(pAudioBytes []byte, pMetaData AudioMetadata) (int64, error) {
	var lChunkID int64
	query := `
	INSERT INTO audio_chunks ( user_id, session_id, timestamp,
		file_name, content_type, audio_data,
		duration,  size, transcript, checksum
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
	} else {

		lChunkID, lErr = lResult.LastInsertId()
		if lErr != nil {
			log.Println("ERROR:PSA02 failed to get last insert id", lErr)
			return lChunkID, lErr
		}
	}

	return lChunkID, lErr
}
