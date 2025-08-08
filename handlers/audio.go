package handlers

import (
	"CHUNKFLOW/common"
	"CHUNKFLOW/processors"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func UploadAudio(w http.ResponseWriter, r *http.Request) {
	log.Println("UploadAudio(+)")

	if !strings.EqualFold("POST", r.Method) {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	lParentCtx := r.Context()
	lCtx, lCancel := context.WithTimeout(lParentCtx, 4*time.Second)
	defer lCancel()

	var lResp processors.Response

	lUserID := r.FormValue("user_id")
	lSessionID := r.FormValue("session_id")
	lTimestamp := r.FormValue("timestamp")

	lFile, lHandler, lErr := r.FormFile("audio")
	if lErr != nil {
		log.Println("Error:HUA01", lErr)
		lResp.Status = common.ERROR
		lResp.ErrMsg = lErr.Error()
		fmt.Fprintf(w, lResp.ErrMsg)
		return

	}
	defer lFile.Close()

	lFileByteArr, lErr := io.ReadAll(lFile)
	if lErr != nil {
		log.Println("Error:HUA02", lErr)
		lResp.Status = common.ERROR
		lResp.ErrMsg = lErr.Error()
		fmt.Fprintf(w, lResp.ErrMsg)
		return

	}

	lContentType := lHandler.Header.Get("Content-Type")

	lJob := &processors.AudioJob{
		Ctx:         lCtx,
		UserID:      lUserID,
		SessionID:   lSessionID,
		Timestamp:   lTimestamp,
		FileName:    lHandler.Filename,
		ContentType: lContentType,
		FileArr:     lFileByteArr,
		Handler:     lHandler,
		File:        lFile,
	}

	lJobs := make(chan *processors.AudioJob, 1)
	lResults := make(chan processors.Response, 1)

	processors.StartAudioWorkerPool(lCtx, lJobs, lResults)

	lJobs <- lJob
	close(lJobs)

	select {
	case <-lCtx.Done():
		log.Println("Context Done:", lCtx.Err())
		lResp.Status = common.ERROR
		lResp.ErrMsg = "Request cancelled or timed out"

	case lResult := <-lResults:
		lResp = lResult
	}

	lBody, lErr := json.Marshal(lResp)
	if lErr != nil {
		log.Println("Error:HUA07 ", lErr)

	} else {
		fmt.Fprintf(w, string(lBody))
	}
	log.Println("UploadAudio(-)")
}
