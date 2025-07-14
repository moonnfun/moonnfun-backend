package global

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

// web3 network
const Testnet = "testnet"
const Mainnet = "mainnet"

// websocket message type
const (
	C_Msg_Topic_KLine = "kline"

	C_Msg_Price_Init   = "price_init_"
	C_Msg_Price_Update = "price_update_"
	C_Msg_Price_Append = "price_append_"

	C_PushMode_Sync  = "PushMode_Sync"
	C_PushMode_ASync = "PushMode_ASync"
)

var GetLatestNumber func() int64
var PushTokenToCache func(field string, token any, totalSize int) error
var GetAveragePrice func(tokenAddress, wallet string, totalAmount decimal.Decimal) (decimal.Decimal, error)

var WebsocketSubscribe func(id, address, topic string) chan struct{}
var WebsocketHandler func(w http.ResponseWriter, r *http.Request) error
var WebsocketSend func(fromID, toID, msgType string, payload any) error

func Str(v any) string {
	return fmt.Sprintf("%v", v)
}

func WssMsgType(prefix, topic string) string {
	return fmt.Sprintf("%s%s", prefix, topic)
}
