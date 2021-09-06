package goBinance

import (
	"encoding/json"
)

func ToString(msg interface{}) string {
	js, _ := json.MarshalIndent(msg, "", "\t")
	s := string(js)
	return s
}

type SubscribeMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     uint     `json:"id"`
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
	}
}

type BookTickerStreamMessage struct {
	UpdateID     int    `json:"u"` // order book updateId
	Symbol       string `json:"s"` // symbol
	BestBidPrice string `json:"b"` // best bid price
	BestBidQty   string `json:"B"` // best bid qty
	BestAskPrice string `json:"a"` // best ask price
	BestAskQty   string `json:"A"` // best ask qty
}

type TickerStream struct {
	EventType              string `json:"e"` // Event type
	EventTyme              int64  `json:"E"` // Event time
	Symbol                 string `json:"s"` // Symbol
	PriceChange            string `json:"p"` // Price change
	PriceChangePercentage  string `json:"P"` // Price change percent
	WeightedAveragePrice   string `json:"w"` // Weighted average price
	FirstTradePrice        string `json:"x"` // First trade(F)-1 price (first trade before the 24hr rolling window)
	LastPrice              string `json:"c"` // Last price
	LastQuantity           string `json:"Q"` // Last quantity
	BestBidPrice           string `json:"b"` // Best bid price
	BestBidQuantity        string `json:"B"` // Best bid quantity
	BestAskPrice           string `json:"a"` // Best ask price
	BestAskQuantity        string `json:"A"` // Best ask quantity
	OpenPrice              string `json:"o"` // Open price
	HighPrice              string `json:"h"` // High price
	LowPrice               string `json:"l"` // Low price
	TotalBaseAssetVolumen  string `json:"v"` // Total traded base asset volume
	TotalQuoteAssetVolumen string `json:"q"` // Total traded quote asset volume
	StatsOpenTime          int    `json:"O"` // Statistics open time
	StatsClosetime         int    `json:"C"` // Statistics close time
	FirstTradeID           int    `json:"F"` // First trade ID
	LastTradeID            int    `json:"L"` // Last trade Id
	TotalNumberOfTrades    int    `json:"n"` // Total number of trades
}
