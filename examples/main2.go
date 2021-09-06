// package main

// import (
// 	"log"
// 	"net/url"
// 	"strconv"

// 	"github.com/Visn0/goBinance"
// )

// type TickerStream struct {
// 	EventType              string `json:"e"` // Event type
// 	EventTyme              int64  `json:"E"` // Event time
// 	Symbol                 string `json:"s"` // Symbol
// 	PriceChange            string `json:"p"` // Price change
// 	PriceChangePercentage  string `json:"P"` // Price change percent
// 	WeightedAveragePrice   string `json:"w"` // Weighted average price
// 	FirstTradePrice        string `json:"x"` // First trade(F)-1 price (first trade before the 24hr rolling window)
// 	LastPrice              string `json:"c"` // Last price
// 	LastQuantity           string `json:"Q"` // Last quantity
// 	BestBidPrice           string `json:"b"` // Best bid price
// 	BestBidQuantity        string `json:"B"` // Best bid quantity
// 	BestAskPrice           string `json:"a"` // Best ask price
// 	BestAskQuantity        string `json:"A"` // Best ask quantity
// 	OpenPrice              string `json:"o"` // Open price
// 	HighPrice              string `json:"h"` // High price
// 	LowPrice               string `json:"l"` // Low price
// 	TotalBaseAssetVolumen  string `json:"v"` // Total traded base asset volume
// 	TotalQuoteAssetVolumen string `json:"q"` // Total traded quote asset volume
// 	StatsOpenTime          int    `json:"O"` // Statistics open time
// 	StatsClosetime         int    `json:"C"` // Statistics close time
// 	FirstTradeID           int    `json:"F"` // First trade ID
// 	LastTradeID            int    `json:"L"` // Last trade Id
// 	TotalNumberOfTrades    int    `json:"n"` // Total number of trades
// }

// var ethusdt = make(chan TickerStream, 1)
// var ethbtc = make(chan TickerStream, 1)
// var btcusdt = make(chan TickerStream, 1)

// const fee float64 = 0.04 / 100

// func createSocketReader(wsId uint, scheme string, host string, path string, topic string, c *chan TickerStream) {
// 	u := url.URL{Scheme: scheme, Host: host, Path: path} // "wss://stream.binance.com:9443/ws"

// 	ws := goBinance.WebSocket{}
// 	defer ws.Close()
// 	ws.Connect(u.String())

// 	topicSlice := []string{topic}
// 	ss := goBinance.SubscribeMessage{Method: "SUBSCRIBE", Params: topicSlice, ID: wsId}
// 	log.Println(ss)
// 	ws.Subscribe(ss)

// 	for {
// 		// message := goBinance.KlineStreamMessage{}
// 		// msg := make(map[string]interface{})
// 		msg := TickerStream{}
// 		err := ws.Conn.ReadJSON(&msg)
// 		if err != nil {
// 			log.Fatal("[READ ERROR]: \n", err)
// 		}
// 		// // *c <- msg
// 		select {
// 		case *c <- msg:
// 			// log.Printf("[MESSAGE-%s]: %v\n", path, msg)
// 		default:
// 			// log.Printf("[MESSAGE-%s]: %v\n", path, msg)
// 		}
// 	}
// }

// // USDT -> ETH -> BTC -> USDT
// func cashCycle(usdt float64, ethusdt float64, ethbtc float64, btcusdt float64) float64 {
// 	// res := ((usdt * fee) / ethusdt) * ethbtc * btcusdt * fee * fee

// 	btc := (usdt / btcusdt)
// 	eth := (btc / ethbtc)
// 	usdtFinal := (eth * ethusdt)
// 	usdtFinal -= 2*usdt*(0.0360/100) - usdt*(0.0180/100) // 3 trades -> 3 fees
// 	return usdtFinal
// }

// func createOrderSender(scheme string, host string, path string, ethusdtChannel *chan TickerStream, ethbtcChannel *chan TickerStream, btcusdtChannel *chan TickerStream) {
// 	usdt := 100.0
// 	for {
// 		// USDT -> ETH -> BTC -> USDT
// 		// fmt.Println("PREVIOUS")
// 		ethusdt := <-*ethusdtChannel
// 		// fmt.Printf("\n###### ethusdt\n")
// 		ethbtc := <-*ethbtcChannel
// 		// fmt.Printf("\n###### ethubtc\n")
// 		btcusdt := <-*btcusdtChannel
// 		// fmt.Printf("\n###### btcusdt\n")

// 		eu, _ := strconv.ParseFloat(ethusdt.LastPrice, 32)
// 		eb, _ := strconv.ParseFloat(ethbtc.LastPrice, 32)
// 		bu, _ := strconv.ParseFloat(btcusdt.LastPrice, 32)
// 		estimation := cashCycle(usdt, eu, eb, bu)

// 		log.Printf("----- USDT: %f, ESTIMATION: %f\n", usdt, estimation)
// 		if estimation > usdt {
// 			log.Printf(" ######## USDT: %f, ESTIMATION: %f\n", usdt, estimation)
// 		}
// 	}
// }

// func main() {
// 	// Readers
// 	go createSocketReader(1, "wss", "stream.binance.com:9443", "ws", "btcusdt@ticker", &btcusdt)
// 	go createSocketReader(2, "wss", "stream.binance.com:9443", "ws", "ethbtc@ticker", &ethbtc)
// 	go createSocketReader(3, "wss", "stream.binance.com:9443", "ws", "ethusdt@ticker", &ethusdt)

// 	// Sender
// 	go createOrderSender("https", "stream.binance.com:9443", "api", &ethusdt, &ethbtc, &btcusdt)

// 	for {
// 	}
// }
