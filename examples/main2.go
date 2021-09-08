package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/Visn0/goBinance"
)

var l = log.New(os.Stdout, "", log.Ltime)

func Log(msg ...interface{}) {
	l.Println(msg)
}

const fee float64 = 0.04 / 100

func createSocketReader(wsId uint, u url.URL, topic string, c *chan goBinance.BookTickerStreamMessage, wg *sync.WaitGroup) {

	defer wg.Done()
	Log(wsId)
	ws := goBinance.WebSocket{}
	defer ws.Close()
	ws.Connect(u.String())

	topicSlice := []string{topic}
	ss := goBinance.SubscribeMessage{Method: "SUBSCRIBE", Params: topicSlice, ID: wsId}
	ws.Subscribe(ss)

	for {
		msg := goBinance.BookTickerStreamMessage{}
		err := ws.Conn.ReadJSON(&msg)
		if err != nil {
			log.Fatal("[READ ERROR]: \n", err)
		}
		select {
		case *c <- msg: // in case channel is free
		default: // if not, next time will try again with newest msg
		}
	}
}

// USDT -> BTC -> ETH -> USDT
func cashCycle(usdt float64, ethusdt float64, ethbtc float64, btcusdt float64) float64 {

	btc := (usdt / btcusdt)      // buying btc->ask
	eth := (btc / ethbtc)        // buying eth->ask
	usdtFinal := (eth * ethusdt) // selling eth->bid
	usdtFinal -= 3 * fee * usdt
	// usdtFinal -= 2*usdt*(0.0360/100) - usdt*(0.0180/100) // 3 trades -> 3 fees
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
		default: // doesnt have new info
			continue
		}

		bu, _ := strconv.ParseFloat(btcusdt.BestAskPrice, 32)
		eb, _ := strconv.ParseFloat(ethbtc.BestAskPrice, 32)
		eu, _ := strconv.ParseFloat(ethusdt.BestBidPrice, 32)
		estimation := cashCycle(usdt, eu, eb, bu)

		Log("----- USDT:", usdt, "ESTIMATION:", estimation)
		if estimation > usdt {
			Log("##### USDT:", usdt, "ESTIMATION:", estimation)
		}
	}
}

func main() {
	topics := []string{
		"btcusdt@bookTicker",
		"ethbtc@bookTicker",
		"ethusdt@bookTicker",
	}
	host := "stream.binance.com:9443"
	// host := "fstream.binance.com"
	u := url.URL{Scheme: "wss", Host: host, Path: "/ws"}
	var ethusdt = make(chan goBinance.BookTickerStreamMessage, 1)
	var ethbtc = make(chan goBinance.BookTickerStreamMessage, 1)
	var btcusdt = make(chan goBinance.BookTickerStreamMessage, 1)

	var wg sync.WaitGroup
	wg.Add(3)
	go createSocketReader(1, u, topics[0], &btcusdt, &wg)
	go createSocketReader(2, u, topics[1], &ethbtc, &wg)
	go createSocketReader(3, u, topics[2], &ethusdt, &wg)

	// // Sender
	go createOrderSender("https", "stream.binance.com:9443", "api", &ethusdt, &ethbtc, &btcusdt)

	wg.Wait()
}
