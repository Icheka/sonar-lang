package handlers

import (
	"fmt"
	"language-server/utils"
	"log"
	"net/http"
	"strings"

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
	if err != nil {
		log.Println(err)
	}

	// defer ws.Close()
	log.Println("Connected")
	EmitJson(ws, Message{
		Type: "connected",
	})

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println(err)
			break
		}

		log.Printf("%+v", msg)

		switch msg.Type {
		case "evaluate":
			source := fmt.Sprint(msg.Data)
			if len(strings.Trim(source, " ")) == 0 {
				break
			}

			out, err := utils.Evaluate(source)
			if len(strings.Trim(err, " ")) != 0 {
				log.Println("StdErr:", err)
				EmitJson(ws, Message{
					Type: "stderr",
					Data: err,
				})
				break
			}
			log.Println("StdOut:", out)
			EmitJson(ws, Message{
				Type: "stdout",
				Data: out,
			})
		}
	}
}
