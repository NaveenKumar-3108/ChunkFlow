package handlers

import (
	"CHUNKFLOW/processors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type UserChunks struct {
	UserChunks []processors.AudioMetadata `json:"userchunks"`
	Status     string                     `json:"status"`
	ErrMsg     string                     `json:"errmsg"`
}

func GetUserChunksdata(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold("GET", r.Method) {
		log.Println("GetUserChunksdata(+)")
		var lUserChunkResp UserChunks
		lVars := mux.Vars(r)
		lUserID := lVars["id"]

		if lUserID == "" {
			log.Println("Error:HGU01 user id is missing")
			lUserChunkResp.Status = "E"
			lUserChunkResp.ErrMsg = "user id is missing"
			goto Marshal
		} else {
			lAudioDataArr, lErr := processors.GetUserChunks(lUserID)

			if lErr != nil {
				log.Println("Error:HGU02 ", lErr)
				lUserChunkResp.Status = "E"
				lUserChunkResp.ErrMsg = lErr.Error()
				goto Marshal
			} else {
				if len(lAudioDataArr) == 0 {
					lUserChunkResp.UserChunks = []processors.AudioMetadata{}
					lUserChunkResp.Status = "E"
					lUserChunkResp.ErrMsg = "No Records found"
				} else {
					lUserChunkResp.Status = "S"
					lUserChunkResp.UserChunks = lAudioDataArr
				}

			}
		}

	Marshal:

		lBody, lErr := json.Marshal(lUserChunkResp)
		if lErr != nil {
			log.Println("Error:HGU03 ", lErr)

		} else {
			fmt.Fprintf(w, string(lBody))
		}
	}
	log.Println("GetUserChunksdata(-)")
}
