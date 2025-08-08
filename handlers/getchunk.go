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

type ChunkResponse struct {
	MetaData processors.AudioMetadata `json:"MetaData"`
	Status   string                   `json:"status"`
	Errmsg   string                   `json:"errmsg"`
}

func GetChunkMetadata(w http.ResponseWriter, r *http.Request) {
	log.Println("GetChunkMetadata(+)")
	if strings.EqualFold("GET", r.Method) {
		var lChunkResp ChunkResponse
		lVars := mux.Vars(r)
		lChunkID := lVars["id"]

		if lChunkID == "" {
			log.Println("Error:HGC01 chunk id is missing")
			lChunkResp.Status = "E"
			lChunkResp.Errmsg = "chunk id is missing"
			goto Marshal
		} else {
			lAudioData, lErr := processors.GetChunkData(lChunkID)

			if lErr != nil {
				log.Println("Error:HGC02 ", lErr)
				lChunkResp.Status = "E"
				lChunkResp.Errmsg = lErr.Error()
				goto Marshal
			} else {
				lChunkResp.Status = "S"
				lChunkResp.MetaData = lAudioData
			}
		}

	Marshal:

		lBody, lErr := json.Marshal(lChunkResp)
		if lErr != nil {
			log.Println("Error:HUA07 ", lErr)

		} else {
			fmt.Fprintf(w, string(lBody))
		}
	}
	log.Println("GetChunkMetadata(-)")
}
