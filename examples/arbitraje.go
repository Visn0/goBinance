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
	log.Println(wsId)
	ws := goBinance.WebSocket{}
	defer ws.Close()
	ws.Connect(u.String())

	log.Println("ok")
	topicSlice := []string{topic}
	ss := goBinance.SubscribeMessage{Method: "SUBSCRIBE", Params: topicSlice, ID: wsId}
	log.Println("Trying to subscribe to:", ss)
	ws.Subscribe(ss)

	for {
		msg := goBinance.BookTickerStreamMessage{}
		err := ws.Conn.ReadJSON(&msg)
		if err != nil {
			log.Fatal("[READ ERROR]: \n", err)
		}
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

	var ethusdt goBinance.BookTickerStreamMessage = <-*ethusdtChannel
	var ethbtc goBinance.BookTickerStreamMessage = <-*ethbtcChannel
	var btcusdt goBinance.BookTickerStreamMessage = <-*btcusdtChannel

	for {
		select {
		case ethusdt = <-*ethusdtChannel:
		case ethbtc = <-*ethbtcChannel:
		case btcusdt = <-*btcusdtChannel:
		default:
		}

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
	spothost := "stream.binance.com:9443"
	// futurehost := "fstream.binance.com"
	u := url.URL{Scheme: "wss", Host: spothost, Path: "/ws"} // "wss://stream.binance.com:9443/ws"
	// Readers
	go createSocketReader(1, u, topics[0], &btcusdt)
	go createSocketReader(2, u, topics[1], &ethbtc)
	go createSocketReader(3, u, topics[2], &ethusdt)

	// // Sender
	go createOrderSender("https", "stream.binance.com:9443", "api", &ethusdt, &ethbtc, &btcusdt)
	for {
	}
}
