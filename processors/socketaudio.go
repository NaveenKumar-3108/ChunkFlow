package processors

import (
	"CHUNKFLOW/common"
	"fmt"
)

type WSRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Timestamp string `json:"timestamp"`
	FileName  string `json:"file_name"`
	FileBytes []byte `json:"file_bytes"`
}

type KnowWSResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg,omitempty"`
}

type WSResponse struct {
	Data    AudioMetadata
	ChunkID int64
}
type WsAudioJob struct {
	Req  WSRequest
	Send chan interface{}
}

var JobQueue = make(chan WsAudioJob, 100)

func StartWorkerPool(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go Worker(i)
	}
}

func Worker(id int) {
	for lJob := range JobQueue {

		if len(lJob.Req.FileBytes) == 0 {
			lJob.Send <- KnowWSResponse{
				Status: "not_ok",
				Msg:    "No binaries found",
			}
		}

		lFile, lFileHeader, lErr := common.BytesToMultipartFile(lJob.Req.FileBytes, lJob.Req.FileName)
		if lErr != nil {
			lJob.Send <- KnowWSResponse{
				Status: "not_ok",
				Msg:    lErr.Error(),
			}

		}
		lErrMsg, lValid := ValidateSocketAudioUpload(lJob.Req.UserID, lJob.Req.SessionID, lJob.Req.Timestamp, lFileHeader)
		if !lValid {
			lJob.Send <- KnowWSResponse{
				Status: "not_ok",
				Msg:    lErrMsg,
			}
			continue
		}

		lChecksum, lTranscript := TransformAudio(lJob.Req.FileBytes)

		lContentType := lFileHeader.Header.Get("Content-Type")
		lDuration := 0.0
		switch lContentType {
		case "audio/wav":
			lDuration, lErr = GetWAVDuration(lFile)
			if lErr != nil {
				lJob.Send <- KnowWSResponse{
					Status: "not_ok",
					Msg:    lErr.Error(),
				}
				continue
			}

		case "audio/mpeg", "audio/mp3":
			lDuration, lErr = EstimateMP3Duration(lFile)
			if lErr != nil {
				lJob.Send <- KnowWSResponse{
					Status: "not_ok",
					Msg:    lErr.Error(),
				}
				continue
			}

		}

		lMeta := AudioMetadata{
			UserID:      lJob.Req.UserID,
			SessionID:   lJob.Req.SessionID,
			Timestamp:   lJob.Req.Timestamp,
			FileName:    lJob.Req.FileName,
			Duration:    fmt.Sprintf("%.2f", lDuration),
			Transcript:  lTranscript,
			ContentType: lContentType,
			Checksum:    lChecksum,
			Size:        fmt.Sprintf("%d", len(lJob.Req.FileBytes)),
		}

		lChunkID, lErr := StoreAudioChunk(lJob.Req.FileBytes, lMeta)
		if lErr != nil {
			lJob.Send <- KnowWSResponse{
				Status: "not_ok",
				Msg:    lErr.Error(),
			}
			continue
		}

		lJob.Send <- WSResponse{
			Data:    lMeta,
			ChunkID: lChunkID,
		}

	}
}
