package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/Visn0/goBinance"
)

func main() {
	// topic :=
	topics := []string{"btcusdt@kline_1m", "ethusdt@bookTicker"}
	// spothost := "stream.binance.com:9443"
	futurehost := "fstream.binance.com"
	url := url.URL{Scheme: "wss", Host: futurehost, Path: "/ws"} // "wss://stream.binance.com:9443/ws"

	ws := goBinance.WebSocket{}
	defer ws.Close()
	ws.Connect(url.String())

	ss := goBinance.SubscribeMessage{Method: "SUBSCRIBE", Params: topics, ID: 1}
	log.Println("Trying to subscribe to:", ss)
	ws.Subscribe(ss)

	func() {
		message := goBinance.KlineStreamMessage{}

		for {
			start := time.Now()

			m := make(map[string]interface{})
			_, p, err := ws.Conn.ReadMessage()
			if err != nil {
				log.Fatal("[READ ERROR]: \n", err)
			}
			_ = json.Unmarshal(p, &m)

			if m["e"] == "kline" {
				err = json.Unmarshal(p, &message)
				if err != nil {
					log.Println(err)
				}
				log.Println(message)
			} else {
				log.Println(m)
			}

			log.Println("Elapsed:", time.Since(start))
		}
	}()

}
