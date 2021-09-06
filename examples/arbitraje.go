package main

import (
	"log"
	"net/url"
	"strconv"

	"github.com/Visn0/goBinance"
)

var ethusdt = make(chan goBinance.BookTickerStreamMessage, 1)
var ethbtc = make(chan goBinance.BookTickerStreamMessage, 1)
var btcusdt = make(chan goBinance.BookTickerStreamMessage, 1)

const fee float64 = 0.04 / 100

func createSocketReader(wsId uint, u url.URL, topic string, c *chan goBinance.BookTickerStreamMessage) {
	// u := url.URL{Scheme: scheme, Host: host, Path: path} // "wss://stream.binance.com:9443/ws"

	ws := goBinance.WebSocket{}
	defer ws.Close()
	ws.Connect(u.String())

	topicSlice := []string{topic}
	ss := goBinance.SubscribeMessage{Method: "SUBSCRIBE", Params: topicSlice, ID: wsId}
	log.Println("Trying to subscribe to:", ss)
	ws.Subscribe(ss)

	log.Println("ok")

	for {
		// message := goBinance.KlineStreamMessage{}
		// msg := make(map[string]interface{})
		msg := goBinance.BookTickerStreamMessage{}
		err := ws.Conn.ReadJSON(&msg)
		if err != nil {
			log.Fatal("[READ ERROR]: \n", err)
		}
		// // *c <- msg
		select {
		case *c <- msg:
			// log.Printf("[MESSAGE-%s]: %v\n", path, msg)
		default:
			// log.Printf("[MESSAGE-%s]: %v\n", path, msg)
		}
	}
}

// USDT -> ETH -> BTC -> USDT
func cashCycle(usdt float64, ethusdt float64, ethbtc float64, btcusdt float64) float64 {
	// res := ((usdt * fee) / ethusdt) * ethbtc * btcusdt * fee * fee

	btc := (usdt / btcusdt)
	eth := (btc / ethbtc)
	usdtFinal := (eth * ethusdt)
	usdtFinal -= 2*usdt*(0.0360/100) - usdt*(0.0180/100) // 3 trades -> 3 fees
	return usdtFinal
}

func createOrderSender(scheme, host, path string,
	ethusdtChannel, ethbtcChannel, btcusdtChannel *chan goBinance.BookTickerStreamMessage) {
	usdt := 100.0
	for {
		// USDT -> ETH -> BTC -> USDT
		// fmt.Println("PREVIOUS")
		ethusdt := <-*ethusdtChannel
		// fmt.Printf("\n###### ethusdt\n")
		ethbtc := <-*ethbtcChannel
		// fmt.Printf("\n###### ethubtc\n")
		btcusdt := <-*btcusdtChannel
		// fmt.Printf("\n###### btcusdt\n")

		eu, _ := strconv.ParseFloat(ethusdt.BestAskPrice, 32)
		eb, _ := strconv.ParseFloat(ethbtc.BestAskPrice, 32)
		bu, _ := strconv.ParseFloat(btcusdt.BestBidPrice, 32)
		estimation := cashCycle(usdt, eu, eb, bu)

		log.Printf("----- USDT: %f, ESTIMATION: %f\n", usdt, estimation)
		if estimation > usdt {
			log.Printf(" ######## USDT: %f, ESTIMATION: %f\n", usdt, estimation)
		}
	}
}

func main() {
	topics := []string{
		"btcusdt@bookTicker",
		"ethbtc@bookTicker",
		"ethusdt@bookTicker",
	}
	// spothost := "stream.binance.com:9443"
	futurehost := "fstream.binance.com"
	url := url.URL{Scheme: "wss", Host: futurehost, Path: "/ws"} // "wss://stream.binance.com:9443/ws"
	// Readers
	go createSocketReader(1, url, topics[0], &btcusdt)
	go createSocketReader(2, url, topics[1], &ethbtc)
	go createSocketReader(3, url, topics[2], &ethusdt)

	// Sender
	go createOrderSender("https", "stream.binance.com:9443", "api", &ethusdt, &ethbtc, &btcusdt)

	for {
	}
}
