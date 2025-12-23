package monitor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/store"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	c_swapTopic     = "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"
	c_swapBuyTopic  = "0xeedf30835a639f25df5f21fa4231d87519341f2e1bbd617cfd7a82681e5ebc4e"
	c_swapSellTopic = "0xe3a8bbe0834808f3b0fbf54934ba50ffd1627bb3dd6bd8904e5e8df6bca58334"

	// C_w_address_main      = "0xE30feDd158A2e3b13e9badaeABaFc5516e95e8C7"
	// C_w_address_test      = "0xF8EB55EC97B59d91fe9E91A1d61147e0d2A7b6F7"
	// C_w_address_dev       = "0x712516e61C8B383dF4A63CFe83d7701Bce54B03e"
	// C_router_address_main = "0xa4cF2F53D1195aDDdE9e4D3aCa54f556895712f2"
	// C_router_address_test = "0x527b42CA5e11370259EcaE68561C14dA415477C8"
	// C_router_address_dev  = "0x59F2f1fCfE2474fD5F0b9BA1E73ca90b143Eb8d0"
)

type SwapTrade struct {
	Wallet common.Address  `json:"caller"`
	Token  common.Address  `json:"token"`
	Source decimal.Decimal `json:"source"`
}

var w_address = ""
var v_monitor_pairs = sync.Map{}

func InitPairs() error {
	// if global.Config.Testnet {
	// 	if global.Config.Debug {
	// 		w_address = C_w_address_dev
	// 	} else {
	// 		w_address = C_w_address_test
	// 	}
	// }
	w_address = global.Config.RPC.WrapperAddress

	tokens, err := store.DBList[model.Token](model.C_Token, bson.M{"launch": true}, false)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		if token.Pair != "" {
			RegisterPair(token.Pair)
		}
	}
	return nil
}

func RegisterPair(pairAddress string) {
	v_monitor_pairs.Store(pairAddress, byte('0'))
}

func IsLaunchedTokenSwap(pairAddress string) bool {
	_, ok := v_monitor_pairs.Load(pairAddress)
	if !ok {
		return false
	}
	return true
}

func GetRouter() string {
	// if global.Config.Testnet {
	// 	if global.Config.Debug {
	// 		return C_router_address_dev
	// 	} else {
	// 		return C_router_address_test
	// 	}
	// } else {
	// 	return C_router_address_main
	// }
	return global.Config.RPC.SwapRouterAddress
}

func HandleTokenSwapTx(height *big.Int, tx *TxData, receiptBuf []byte) {
	if common.HexToAddress(tx.To).Cmp(common.HexToAddress(global.Config.ContractAddress)) != 0 && common.HexToAddress(tx.To).Cmp(common.HexToAddress(GetRouter())) != 0 {
		return
	}

	blockHash := common.HexToHash(tx.BlockHash)
	swapTopic := common.HexToHash(c_swapTopic)
	routerAddressTopic := common.HexToHash(GetRouter())
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{{swapTopic}, {routerAddressTopic}},
	}
	logs, err := GetLogs(v_wss_client, tx, query, receiptBuf)
	if err != nil {
		slog.Error("get TokenSwap failed", "txHash", tx.Hash, "error", err.Error())
		return
	}
	// global.DebugForce("get logs successed", "logs", logs, "routerAddress", GetRouter())

	var swapTrade *SwapTrade
	if st, err := GetSwapTrade(height, tx, receiptBuf); err == nil && st != nil {
		swapTrade = st
	}

	for _, log := range logs {
		if IsLaunchedTokenSwap(log.Address.String()) {
			slog.Info("HandleTokenSwapTx successed", slog.Any("log", log))
			tokenTrading, err := getTokenSwap(log, tx, tx.Time, swapTrade)
			if err != nil {
				slog.Error("get TokenSwap failed", "log", log, "error", err.Error())
			} else {
				tokenTrading.Number = height.Int64()
				// tokenTrading.Timestamp = time.UnixMilli(int64(header.Time))
				if TokenTradingNotify != nil {
					TokenTradingNotify(tokenTrading)
					return
				}
			}
		}
	}
}

func getTokenSwap(log types.Log, tx *TxData, txTime time.Time, swapTrade *SwapTrade) (*TokenTrading, error) {
	bBuy := false
	topics := ""
	tradeEvent := &TradeEvent{}
	for _, topic := range log.Topics {
		topics += topic.String() + ","
	}
	topics = strings.TrimSuffix(topics, ",")
	if err := doGetTokenSwap(topics, "0x"+common.Bytes2Hex(log.Data), tradeEvent); err != nil {
		return nil, err
	}
	// global.DebugForce("get SwapTopic successed", "swapTopic", topics)

	txHash := log.TxHash.String()
	tokenTrading := &TokenTrading{
		TxHash: txHash,
		Pair:   log.Address,
	}
	tokenTrading.Buy = bBuy
	tokenTrading.Token = tradeEvent.Token
	tokenTrading.Wallet = tradeEvent.To
	tokenTrading.Timestamp = txTime
	if tradeEvent.To.String() == GetRouter() {
		tokenTrading.Wallet = common.HexToAddress(GetTxSender(tx.Hash))
	}
	if swapTrade != nil {
		tokenTrading.Token = swapTrade.Token
		tokenTrading.Source = swapTrade.Source
		tokenTrading.Wallet = swapTrade.Wallet
	}

	buyQuoteAmount := tradeEvent.AmountIn0
	sellQuoteAmount := tradeEvent.AmountOut0
	buyTokenAmount := tradeEvent.AmountOut1
	sellTokenAmount := tradeEvent.AmountIn1
	token0Address := GetToken0(log.Address.String())
	if token0Address != "" && common.HexToAddress(token0Address).Cmp(common.HexToAddress(w_address)) != 0 {
		slog.Info("get sell trading successed", "tradeEvent", tradeEvent, "w_address", w_address, "token0Address", token0Address)
		buyQuoteAmount = tradeEvent.AmountIn1
		sellQuoteAmount = tradeEvent.AmountOut1
		buyTokenAmount = tradeEvent.AmountOut0
		sellTokenAmount = tradeEvent.AmountIn0
	} else if token0Address == "" {
		slog.Warn("get token0 failed", "pairAddress", log.Address.String(), "w_address", w_address)
	}

	// constPriceValue := decimal.NewFromInt(1).Div(decimal.NewFromInt(1000000000).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))))
	// global.Debug(buyQuoteAmount, buyTokenAmount, sellQuoteAmount, sellTokenAmount, priceQuote)
	if sellQuoteAmount.IsZero() {
		tokenTrading.Buy = true
		tokenTrading.Price = buyQuoteAmount.Div(buyTokenAmount) //.Mul(priceQuote)
		tokenTrading.Amount = buyTokenAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		tokenTrading.QuoteAmount = buyQuoteAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))

		// swap trading fees: 1/%
		if swapTrade != nil && tokenTrading.Buy {
			quoteAmount := buyQuoteAmount.Div(decimal.NewFromFloat(0.99))
			tokenTrading.QuoteAmount = quoteAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		}
	} else {
		tokenTrading.Buy = false
		tokenTrading.Price = sellQuoteAmount.Div(sellTokenAmount) //.Mul(priceQuote)
		tokenTrading.Amount = sellTokenAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		tokenTrading.QuoteAmount = sellQuoteAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		// tokenTrading.Price = tokenTrading.Amount.Mul(tokenTrading.Amount).Mul(constPriceValue).Mul(priceQuote)
	}
	slog.Info("Get TokenTrading successed", "tokenTrading", tokenTrading)

	return tokenTrading, nil
}

func doGetTokenSwap(topic, data string, output any) error {
	parseArgs, err := ParseSwapLog(topic, data)
	if err != nil {
		slog.Error("parse swap log failed", "topic", topic, "log", data, "error", err.Error())
		return err
	}

	if err := json.Unmarshal([]byte(parseArgs), output); err != nil {
		slog.Error("get TokenSwap failed", "parseArgs", parseArgs, "error", err.Error())
		return err
	}
	global.Debug("Get TokenSwap successed", "output", output)
	return nil
}

func GetSwapTrade(height *big.Int, tx *TxData, receiptBuf []byte) (*SwapTrade, error) {
	if tx.To != global.Config.ContractAddress {
		return nil, nil
	}

	blockHash := common.HexToHash(tx.BlockHash)
	swapBuyTopic := common.HexToHash(c_swapBuyTopic)
	swapSellTopic := common.HexToHash(c_swapSellTopic)
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{{swapBuyTopic, swapSellTopic}},
	}
	logs, err := GetLogs(v_wss_client, tx, query, receiptBuf)
	if err != nil {
		// slog.Error("get TokenSwap failed", "txHash", tx.Hash().String(), "error", err.Error())
		return nil, err
	}
	// global.DebugForce("get logs successed", "logs", logs, "routerAddress", GetRouter())

	for _, log := range logs {
		swapTrade := &SwapTrade{}
		if err := doGetTokenTrading(log.Topics[0].String(), "0x"+common.Bytes2Hex(log.Data), swapTrade); err == nil {
			return swapTrade, nil
		}
	}
	return nil, fmt.Errorf("invalid swap trade event, tx: %s", tx.Hash)
}
