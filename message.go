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
	UpdateID int    `json:"u"` // order book updateId
	Symbol   string `json:"s"` // symbol
	BidPrice string `json:"b"` // best bid price
	BidQty   string `json:"B"` // best bid qty
	AskPrice string `json:"a"` // best ask price
	AskQty   string `json:"A"` // best ask qty
}
