package handlers

import (
	"CHUNKFLOW/processors"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func AudioWebSocket(w http.ResponseWriter, r *http.Request) {
	lConn, lErr := upgrader.Upgrade(w, r, nil)
	if lErr != nil {
		log.Println("Upgrade error:", lErr)
		return
	}
	defer lConn.Close()

	lSendChan := make(chan interface{}, 10)

	go func() {
		for lMsg := range lSendChan {
			if lErr := lConn.WriteJSON(lMsg); lErr != nil {
				log.Println("Write error:", lErr)
				return
			}
		}
	}()

	for {
		lMsgType, lMeta, lErr := lConn.ReadMessage()
		if lErr != nil {
			log.Println("Read error :", lErr)
			break
		}

		var lReq processors.WSRequest
		if lErr := json.Unmarshal(lMeta, &lReq); lErr != nil {
			log.Println("Metadata unmarshal error:", lErr)
			break
		}

		lMsgType, lBin, lErr := lConn.ReadMessage()
		if lErr != nil {
			log.Println("Binary read error:", lErr)
			break
		}
		if lMsgType != websocket.BinaryMessage {
			log.Println("Expected binary message but got:", lMsgType)
			continue
		}

		lReq.FileBytes = lBin

		processors.JobQueue <- processors.WsAudioJob{
			Req:  lReq,
			Send: lSendChan,
		}

	}

}
