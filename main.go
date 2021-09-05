package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocket struct {
	Conn *websocket.Conn
	Out  chan []byte
	In   chan []byte
	// Events map[string]EventHandler
}

func (ws *WebSocket) Reader() {
	defer ws.Conn.Close()

	for {
		_, msg, err := ws.Conn.ReadMessage()
		if err != nil {
			log.Println("error", msg)
			return
		}
	}
}

func (ws *WebSocket) Writer() {
	defer ws.Conn.Close()

	for {
		select {
		// case <-done:
		// 	return
		case msg := <-ws.Out:
			jsonMsg, err := json.Marshal(msg)
			err = ws.Conn.WriteMessage(1, []byte(jsonMsg))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

type KlineWebsocketRequest struct {
	Symbol   string
	Interval string
}

func getKline(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_%s", strings.ToLower(kwr.Symbol), string(kwr.Interval))
}

func main() {
	fmt.Println("Hello world")
}
