package goBinance

import (
	"log"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Conn *websocket.Conn
	// Out  chan []byte
	// In   chan []byte
	// Events map[string]EventHandler
}

func (ws *WebSocket) Connect(url string) {
	log.Printf("Connecting to %v", url)

	wsd, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("[DIAL ERROR]: \n", err)
	}
	ws.Conn = wsd
}

func (ws *WebSocket) Close() {
	ws.Conn.Close()
}

func (ws *WebSocket) SendMessage(request interface{}) []byte {
	err := ws.Conn.WriteJSON(request)
	if err != nil {
		log.Fatal("[ERROR]: \n", err)
	}

	_, msg, err := ws.Conn.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	log.Println("[RESPONSE", string(msg))
	return msg
}
