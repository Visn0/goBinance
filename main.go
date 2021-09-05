package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

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

type KlineStreamMessage struct {
	EventType string `json:"e"`
	EventTime int    `json:"E"`
	Symbol    string `json:"s"`
	K         struct {
		KlineStartTime           int    `json:"t"`
		KlineCloseTime           int    `json:"T"`
		Symbol                   string `json:"s"`
		Interval                 string `json:"i"`
		FirstTradeID             int    `json:"f"`
		LastTradeID              int    `json:"L"`
		OpenPrice                string `json:"o"`
		ClosePrice               string `json:"c"`
		HighPrice                string `json:"h"`
		LowPrice                 string `json:"l"`
		BaseAssetVolume          string `json:"v"`
		NumberOfTrades           int    `json:"n"`
		IsKlineClosed            bool   `json:"x"`
		QuoteAssetVolume         string `json:"q"`
		TakerBuyBaseAssetVolume  string `json:"V"`
		TakerBuyQuoteAssetVolume string `json:"Q"`
		Ignore                   string `json:"B"`
	} `json:"k"`
}

func (k KlineStreamMessage) String() string {
	js, _ := json.MarshalIndent(k, "", "\t")
	s := string(js)
	return s
}

func getKline(kwr KlineWebsocketRequest) {
	// url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_%s", strings.ToLower(kwr.Symbol), string(kwr.Interval))
	// request := fmt.Sprintf("{\"method\":\"SUBSCRIBE\",\"params\":[\"btcusdt@kline_1m\"],\"id\":1}\"")
	host := "stream.binance.com"
	port := "9443"
	topic := "btcusdt@kline_1m"
	m := make(map[string]interface{})
	m["method"] = "SUBSCRIBE"
	m["params"] = topic
	m["id"] = 1

	u := url.URL{Scheme: "wss", Host: host + ":" + port, Path: "ws/" + topic}
	log.Println(u.String())
	log.Println(m)

	c, hr, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(hr.Status)

	// err = c.WriteJSON(m)
	if err != nil {
		log.Println(err)
	}
	go func() {
		defer c.Close()
		log.Println()
		for {
			var msg KlineStreamMessage
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Fatal("[READ ERROR]: \n", err)
			}
			log.Println("[MESSAGE]: ", msg)

		}
	}()
}

func main() {
	fmt.Println("Hello world")
	kwr := KlineWebsocketRequest{Symbol: "btcusdt", Interval: "1m"}
	getKline(kwr)

	// fmt.Scanf()
	for {
	}
}
