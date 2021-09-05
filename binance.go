package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

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
	}
}

type Subscribe struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     uint     `json:"id"`
}

func NewSuscribe() *Subscribe {
	return &Subscribe{
		Method: "SUBSCRIBE",
		Params: []string{"btcusdt@depth"},
		ID:     1,
	}
}

func (k KlineStreamMessage) String() string {
	js, _ := json.MarshalIndent(k, "", "\t")
	s := string(js)
	return s
}

func websocketConnection(scheme string, host string, path string, requestHeader http.Header) {
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	log.Printf("Connecting to %v", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), requestHeader)
	if err != nil {
		log.Fatal("[DIAL ERROR]: \n", err)
	}
	defer ws.Close()

	// SEND SUSCRIPTION REQUEST
	ss := *NewSuscribe()
	err = ws.WriteJSON(ss)
	if err != nil {
		log.Fatal("[SUBSCRIBE ERROR]: \n", err)
	}
	// RECEIVE RESPONSE
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	log.Println("[SUSCRIPTION RESPONSE", string(msg))

	func() {
		for {
			message := KlineStreamMessage{}
			err := ws.ReadJSON(&message)
			if err != nil {
				log.Fatal("[READ ERROR]: \n", err)
			}
			log.Println("[MESSAGE]: ", message)
		}
	}()

}

func main() {
	websocketConnection("wss", "stream.binance.com:9443", "ws/btcusdt@kline_1m", nil)
}
