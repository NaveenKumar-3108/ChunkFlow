package processors

import (
	"CHUNKFLOW/internal/db"
	"log"
)

func GetChunkData(pChunkID string) (AudioMetadata, error) {

	var lAudioData AudioMetadata
	lQuery := ` 
	SELECT chunk_id, user_id, session_id, timestamp,
	       file_name, content_type, duration, size,
	       transcript, checksum
	FROM audio_chunks
	WHERE chunk_id = ?
	`

	lErr := db.GDB.QueryRow(lQuery, pChunkID).Scan(
		&lAudioData.ChunkID, &lAudioData.UserID, &lAudioData.SessionID, &lAudioData.Timestamp,
		&lAudioData.FileName, &lAudioData.ContentType, &lAudioData.Duration,
		&lAudioData.Size, &lAudioData.Transcript, &lAudioData.Checksum,
	)
	if lErr != nil {
		log.Println("Error:PGC01", lErr)
		return lAudioData, lErr
	}

	return lAudioData, lErr

}

func GetUserChunks(pUserID string) ([]AudioMetadata, error) {
	var lUserAudioData []AudioMetadata
	lQuery := ` 
	SELECT chunk_id, user_id, session_id, timestamp,
	       file_name, content_type, duration, size,
	       transcript, checksum
	FROM audio_chunks
	WHERE user_id = ?
	`

	lRows, lErr := db.GDB.Query(lQuery, pUserID)
	if lErr != nil {
		log.Println("Error:PGC01", lErr)
		return lUserAudioData, lErr
	} else {
		defer lRows.Close()

		for lRows.Next() {
			var lAudioData AudioMetadata
			lErr := lRows.Scan(
				&lAudioData.ChunkID, &lAudioData.UserID, &lAudioData.SessionID, &lAudioData.Timestamp,
				&lAudioData.FileName, &lAudioData.ContentType, &lAudioData.Duration,
				&lAudioData.Size, &lAudioData.Transcript, &lAudioData.Checksum,
			)
			if lErr != nil {
				log.Println("Error:PGC02", lErr)
				return lUserAudioData, lErr
			} else {
				lUserAudioData = append(lUserAudioData, lAudioData)
			}
		}

	}

	return lUserAudioData, lErr
}
