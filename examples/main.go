package main

import (
	"log"
	"net/url"

	"github.com/Visn0/goBinance"
)

func main() {
	topic := "btcusdt@kline_1m"
	topics := []string{topic}
	spoturl := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws"} // "wss://stream.binance.com:9443/ws"

	ws := goBinance.WebSocket{}
	defer ws.Close()
	ws.Connect(spoturl.String())

	ss := goBinance.SubscribeMessage{Method: "SUBSCRIBE", Params: topics, ID: 1}
	log.Println(ss)
	ws.Subscribe(ss)

	func() {
		for {
			message := goBinance.KlineStreamMessage{}
			err := ws.Conn.ReadJSON(&message)
			if err != nil {
				log.Fatal("[READ ERROR]: \n", err)
			}
			log.Println("[MESSAGE]: ", message)
		}
	}()

}
