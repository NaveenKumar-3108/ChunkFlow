package processors

import (
	"CHUNKFLOW/common"
	"context"
	"errors"
	"fmt"
	"strconv"
)

const WorkerCount = 5

func StartAudioWorkerPool(pCtx context.Context, pJobs <-chan *AudioJob, pResults chan<- Response) {
	for i := 0; i < WorkerCount; i++ {
		go func() {
			for {
				select {
				case <-pCtx.Done():
					return
				case lJob, ok := <-pJobs:
					if !ok {
						return
					}
					if lJob.Error != nil || pCtx.Err() != nil {
						continue
					}

					ProcessAudioPipeline(pCtx, lJob)

					if lJob.Error != nil {
						pResults <- Response{Status: common.ERROR, ErrMsg: lJob.Error.Error()}

						continue
					}

					lMeta := AudioMetadata{
						UserID:      lJob.UserID,
						SessionID:   lJob.SessionID,
						Timestamp:   lJob.Timestamp,
						FileName:    lJob.FileName,
						ContentType: lJob.ContentType,
						Duration:    lJob.Duration,
						Transcript:  lJob.Transcript,
						Checksum:    lJob.Checksum,
						Size:        strconv.Itoa(len(lJob.FileArr)),
					}

					lChunkID, lErr := StoreAudioChunk(lJob.FileArr, lMeta)
					if lErr != nil {
						pResults <- Response{Status: common.ERROR, ErrMsg: lJob.Error.Error()}
					} else {
						pResults <- Response{Status: "S", ChunkID: lChunkID, Msg: "Audio uploaded successfully"}
					}
				}
			}
		}()
	}
}

func ProcessAudioPipeline(pCtx context.Context, pJob *AudioJob) {
	if pCtx.Err() != nil {
		pJob.Error = pCtx.Err()
		return
	}

	lErrMsg, lValid := ValidateAudioUpload(pJob.UserID, pJob.SessionID, pJob.Timestamp, pJob.Handler)
	if !lValid {
		pJob.Error = errors.New(lErrMsg)
		return
	}

	if pCtx.Err() != nil {
		pJob.Error = pCtx.Err()
		return
	}

	pJob.Checksum, pJob.Transcript = TransformAudio(pJob.FileArr)

	if pCtx.Err() != nil {
		pJob.Error = pCtx.Err()
		return
	}

	switch pJob.ContentType {
	case "audio/wav":
		lDuration, lErr := GetWAVDuration(pJob.File)
		if lErr != nil {
			pJob.Error = lErr
			return
		}

		pJob.Duration = fmt.Sprintf("%.2f", lDuration)

	case "audio/mpeg", "audio/mp3":
		lDuration, lErr := EstimateMP3Duration(pJob.File)
		if lErr != nil {
			pJob.Error = lErr
			return
		}

		pJob.Duration = fmt.Sprintf("%.2f", lDuration)

	}
}
