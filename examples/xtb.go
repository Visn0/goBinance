package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/Visn0/goBinance"
)

type XTBMessage struct {
	Command string `json:"command"`
}

type LoginMessage struct {
	XTBMessage
	Arguments struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
		AppId    string `json:"appId"`
		AppName  string `json:"appName"`
	} `json:"arguments"`
}

func (lm LoginMessage) String() string {
	js, _ := json.MarshalIndent(lm, "", "\t")
	s := string(js)
	return s
}

func NewLoginMessage(id, password *string) *LoginMessage {
	msg := &LoginMessage{
		XTBMessage: XTBMessage{
			Command: "login",
		},
	}

	msg.Arguments.UserId = *id
	msg.Arguments.Password = *password
	// msg.Arguments.AppId = *appId
	// msg.Arguments.AppName = *appName

	return msg
}

type LogoutMessage struct {
	XTBMessage
}

func NewLogoutMessage() *LogoutMessage {
	return &LogoutMessage{
		XTBMessage: XTBMessage{
			Command: "logout",
		},
	}
}

type GetCommissionMessage struct {
	XTBMessage
	Arguments struct {
		Symbol string  `json:"symbol"`
		Volume float32 `json:"volume"`
	} `json:"arguments"`
}

func NewGetCommisionMessage(symbol string, volume float32) *GetCommissionMessage {
	msg := &GetCommissionMessage{
		XTBMessage: XTBMessage{
			Command: "getCommissionDef",
		},
	}

	msg.Arguments.Symbol = symbol
	msg.Arguments.Volume = volume

	return msg
}

type GetTickPricesMessage struct {
	XTBMessage
	// Arguments struct {
	Symbol    string `json:"symbol"`
	SessionID string `json:"streamSessionId"`
	// } `json:"arguments"`
	// MinArrivalTime int    `json:"minArrivalTime"`
	// MaxLevel       int    `json:"maxLevel"`
}

func NewGetTickPricesMessage(symbol string, sessionId string) *GetTickPricesMessage {
	msg := &GetTickPricesMessage{
		XTBMessage: XTBMessage{
			Command: "getTickPrices",
		},
		// Symbol:    symbol,
		// SessionID: sessionId,
	}

	// msg.Arguments.Symbol = symbol
	// msg.Arguments.SessionID = sessionId

	return msg
}

func (lm GetTickPricesMessage) String() string {
	js, _ := json.MarshalIndent(lm, "", "\t")
	s := string(js)
	return s
}

var (
	id       = flag.String("id", "", "User account id.")
	password = flag.String("password", "", "User account id.")
)

func Login(id *string, password *string) (*goBinance.WebSocket, string) {
	log.Println(">> Loging in")

	ws := new(goBinance.WebSocket)
	ws.Connect("wss://ws.xtb.com/demo")

	response := ws.SendMessage(*NewLoginMessage(id, password))
	m := make(map[string]interface{})
	json.Unmarshal(response, &m)
	return ws, fmt.Sprintf("%v", m["streamSessionId"])
}

func Logout(ws *goBinance.WebSocket) {
	log.Println("<< Loging out")
	ws.SendMessage(*NewLogoutMessage())
}

func Trade(ws *goBinance.WebSocket, sessionId string) {
	// ws.SendMessage(*NewGetCommisionMessage("EURUSD", 1.0))
	msg := *NewGetTickPricesMessage("EURUSD", sessionId)
	log.Println(msg)
	ws.SendMessage(msg)
}

func main() {
	flag.Parse()

	fmt.Printf("id: %s, pass: %s\n", *id, *password)
	ws, sessionId := Login(id, password)
	Trade(ws, sessionId)
	Logout(ws)
}
