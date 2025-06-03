package global

import (
	"fmt"
	"net/http"
)

// web3 network
const Testnet = "testnet"
const Mainnet = "mainnet"

// websocket message type
const (
	C_Msg_Topic_KLine = "kline"

	C_Msg_Price_Init   = "price_init"
	C_Msg_Price_Update = "price_update"

	C_PushMode_Sync  = "PushMode_Sync"
	C_PushMode_ASync = "PushMode_ASync"
)

var WebsocketHandler func(w http.ResponseWriter, r *http.Request) error

func Str(v any) string {
	return fmt.Sprintf("%v", v)
}
