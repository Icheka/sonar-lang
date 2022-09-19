package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func EmitJson(ws *websocket.Conn, msg Message) {
	if err := ws.WriteJSON(msg); err != nil {
		log.Printf("An error occurred while emitting JSON (%+v)", msg)
	}
}

func HandleWebSocket(w http.ResponseWriter, req *http.Request) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := wsUpgrader.Upgrade(w, req, nil)
	if err == nil {
		log.Println(err)
	}

	// defer ws.Close()
	log.Println("Connected")

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println(err)
			break
		}

		log.Printf("%+v", msg)
	}
}
