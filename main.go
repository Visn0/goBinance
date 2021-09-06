package main

import (
	"log"
	"net/url"

	btypes "github.com/Visn0/goBinance/binanceTypes"
	stream "github.com/Visn0/goBinance/stream"
)

func main() {
	topic := "btcusdt@kline_1m"
	topics := []string{topic}
	spoturl := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws"} // "wss://stream.binance.com:9443/ws"

	ws := stream.WebSocket{}
	defer ws.Close()
	ws.Connect(spoturl.String())

	ss := btypes.SubscribeMessage{Method: "SUBSCRIBE", Params: topics, ID: 1}
	log.Println(ss)
	ws.Subscribe(ss)

	func() {
		for {
			message := btypes.KlineStreamMessage{}
			err := ws.Conn.ReadJSON(&message)
			if err != nil {
				log.Fatal("[READ ERROR]: \n", err)
			}
			log.Println("[MESSAGE]: ", message)
		}
	}()

}
