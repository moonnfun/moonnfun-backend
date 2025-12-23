package monitor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"meme3/global"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
)

const (
	c_buyTopic    = "0x3d3216984ffd81fb77b1b3c6900c77d4bff1b316b16119b535dfb024f579e02e"
	c_sellTopic   = "0xd87e5212bde853c0f8dcc176803c9e8379c81165d75ac6f7d051a532a1cedea7"
	c_liquidTopic = "0x1955b825797bcf41aa6cbd867e8656367face6df9cc743426835a79dfcc629c0"
	c_launchTopic = "0xe4b8baaa1c927a55ac28b7693040a55b0b48044cf0572d1fc3beb0f418054a20"
)

var (
	v_chan_critical   chan struct{}
	v_criticals_cache *sync.Map
)

type TokenTrading struct {
	Pair        common.Address  `json:"pair"`
	Token       common.Address  `json:"token"`
	Wallet      common.Address  `json:"wallet"`
	Buy         bool            `json:"buy"`
	Liquid      bool            `json:"liquid"`
	Amount      decimal.Decimal `json:"amount"`
	Price       decimal.Decimal `json:"price"`
	TxHash      string          `json:"txhash"`
	Source      decimal.Decimal `json:"source"`
	Number      int64           `json:"number"`
	Timestamp   time.Time       `json:"timestamp"`
	QuoteAmount decimal.Decimal `json:"quoteAmount"`
}

type TradeEvent struct {
	Sender     common.Address  `json:"sender"`
	To         common.Address  `json:"to"`
	Pair       common.Address  `json:"pair"`
	Token      common.Address  `json:"token"`
	Source     decimal.Decimal `json:"source"`
	AmountIn0  decimal.Decimal `json:"amount0In"`
	AmountIn1  decimal.Decimal `json:"amount1In"`
	AmountOut0 decimal.Decimal `json:"amount0Out"`
	AmountOut1 decimal.Decimal `json:"amount1Out"`
}

type LiquidEvent struct {
	Token    common.Address  `json:"token"`
	Time     decimal.Decimal `json:"time"`
	Reverse0 decimal.Decimal `json:"reserve0"`
	Reverse1 decimal.Decimal `json:"reserve1"`
}

var TokenTradingNotify func(tokenTrading any) bool

type WaitTradingList []*TokenTrading

func (t WaitTradingList) Len() int      { return len(t) }
func (t WaitTradingList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t WaitTradingList) Less(i, j int) bool {
	return t[i].Timestamp.UnixNano() < t[j].Timestamp.UnixNano()
}

func HandleTokenTradingTx(height *big.Int, tx *TxData, receiptBuf []byte) {
	isUserContract := strings.Contains(v_user_contract_list, strings.ToUpper(tx.To))
	if !isUserContract && common.HexToAddress(tx.To).Cmp(common.HexToAddress(global.Config.ContractAddress)) != 0 {
		return
	}

	buyTopic := common.HexToHash(c_buyTopic)
	sellTopic := common.HexToHash(c_sellTopic)
	blockHash := common.HexToHash(tx.BlockHash)
	launchTopic := common.HexToHash(c_launchTopic)
	contractAddress := common.HexToAddress(global.Config.ContractAddress)
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{buyTopic, sellTopic, launchTopic}},
	}
	logs, err := GetLogs(v_wss_client, tx, query, receiptBuf)
	if err != nil {
		slog.Error("get TokenTrading failed", "txHash", tx.Hash, "error", err.Error())
		return
	}

	var launchLog *types.Log
	for _, log := range logs {
		if log.Topics[0].String() == c_launchTopic {
			launchLog = &log
		} else {
			slog.Info("HandleTokenTradingTx successed", slog.Any("log", log))
			tokenTrading, err := getTokenTrading(tx, log)
			if err != nil {
				slog.Error("get TokenTrading failed", "log", log, "error", err.Error())
			} else {
				tokenTrading.Number = height.Int64()
				if TokenTradingNotify != nil {
					TokenTradingNotify(tokenTrading)
					break
				}
			}
		}

		if launchLog != nil {
			slog.Info("HandleLaunchLiquidTx successed", slog.Any("log", *launchLog))
			if err := getLaunchLiquid(*launchLog, tx, tx.Time); err != nil {
				slog.Error("get LaunchLiquid failed", "log", *launchLog, "error", err.Error())
			}
		}
	}
}

func getTokenTrading(tx *TxData, log types.Log) (*TokenTrading, error) {
	tradeEvent := &TradeEvent{}
	if err := doGetTokenTrading(log.Topics[0].String(), "0x"+common.Bytes2Hex(log.Data), tradeEvent); err != nil {
		return nil, err
	}

	tokenTrading := &TokenTrading{
		TxHash: log.TxHash.String(),
		Pair:   common.HexToAddress(global.Config.ContractAddress),
	}
	tokenTrading.Token = tradeEvent.Token
	tokenTrading.Wallet = tradeEvent.To
	tokenTrading.Timestamp = tx.Time

	buyQuoteAmount := tradeEvent.AmountIn1
	sellQuoteAmount := tradeEvent.AmountOut1
	buyTokenAmount := tradeEvent.AmountOut0
	sellTokenAmount := tradeEvent.AmountIn0
	// constPriceValue := decimal.NewFromInt(1).Div(decimal.NewFromInt(1000000000).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))))
	// global.Debug(buyQuoteAmount, buyTokenAmount, sellQuoteAmount, sellTokenAmount, priceQuote)
	if log.Topics[0].String() == c_buyTopic {
		// if sellQuoteAmount.IsZero() {
		tokenTrading.Buy = true
		tokenTrading.Price = buyQuoteAmount.Div(buyTokenAmount) //.Mul(priceQuote)
		tokenTrading.Amount = buyTokenAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		tokenTrading.QuoteAmount = buyQuoteAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		// tokenTrading.Price = tokenTrading.Amount.Mul(tokenTrading.Amount).Mul(constPriceValue).Mul(priceQuote)
	} else if log.Topics[0].String() == c_sellTopic {
		tokenTrading.Buy = false
		tokenTrading.Price = sellQuoteAmount.Div(sellTokenAmount) //.Mul(priceQuote)
		tokenTrading.Amount = sellTokenAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		tokenTrading.QuoteAmount = sellQuoteAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		// tokenTrading.Price = tokenTrading.Amount.Mul(tokenTrading.Amount).Mul(constPriceValue).Mul(priceQuote)
	} else {
		slog.Warn("ignore tx event", slog.Any("tx", tx), slog.Any("log", log))
		return nil, fmt.Errorf("ignore tx event: %s", log.Topics[0].String())
	}
	slog.Info("Get TokenTrading successed", "tokenTrading", tokenTrading)

	return tokenTrading, nil
}

func doGetTokenTrading(topic, data string, output any) error {
	parseArgs, err := ParseLog(topic, data)
	if err != nil {
		slog.Error("parse log failed", "topic", topic, "log", data, "error", err.Error())
		return err
	}

	if err := json.Unmarshal([]byte(parseArgs), output); err != nil {
		slog.Error("get TokenTrading failed", "parseArgs", parseArgs, "error", err.Error())
		return err
	}
	global.Debug("Get TokenTrading successed", "parseArgs", parseArgs, "output", output)
	return nil
}

// func getTokenTradingForLiquid(txhash string, liquidEvent *LiquidEvent) (*TokenTrading, error) {
// 	tokenTrading := &TokenTrading{
// 		TxHash: txhash,
// 	}
// 	priceQuote := GetQuotePrice()
// 	if priceQuote.LessThan(decimal.Zero) {
// 		return nil, errors.New("invalid Quote price")
// 	}

// 	tokenTrading.Liquid = true
// 	tokenTrading.Timestamp = time.Unix(liquidEvent.Time.BigInt().Int64(), 0)

// 	quoteAmount, err := decimal.NewFromString(liquidEvent.Reverse1.String())
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid quoteAmount: %s", err.Error())
// 	}
// 	tokenAmount, err := decimal.NewFromString(liquidEvent.Reverse0.String())
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid tokenAmount: %s", err.Error())
// 	}

// 	tokenTrading.Price = quoteAmount.Div(tokenAmount).Mul(priceQuote)
// 	tokenTrading.QuoteAmount = quoteAmount.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
// 	return tokenTrading, nil
// }

func pushTradingNotify(tokenTrading *TokenTrading) {
	sec := tokenTrading.Timestamp.Second()
	bPrev := sec >= 0 && sec <= 9

	if bPrev {
		if v_criticals_cache == nil {
			v_criticals_cache = new(sync.Map)
			v_chan_critical = make(chan struct{})
			go runCriticalsTimeout()
		}
		v_criticals_cache.Store(tokenTrading.TxHash, tokenTrading)
		global.Debug("cache tx for 0 second -> 9 second", "second", sec, "time", tokenTrading.Timestamp.Unix())
		return
	} else if v_criticals_cache != nil && v_chan_critical != nil && sec >= 10 {
		global.Debug("before run cache tx for 0 second -> 9", "second", sec, "time", tokenTrading.Timestamp.Unix())
		close(v_chan_critical)
	}
	TokenTradingNotify(tokenTrading)
}

func runCriticalsTimeout() {
	for {
		select {
		case <-v_chan_critical:
			v_chan_critical = nil
			v_criticals_cache.Range(func(key, value any) bool {
				global.Debug("start run cache tx for 0 second -> 9 second", "time", value.(*TokenTrading).Timestamp.Unix())
				TokenTradingNotify(value.(*TokenTrading))
				return true
			})
			v_criticals_cache.Clear()
			v_criticals_cache = nil
			return
		case <-time.After(time.Duration(20) * time.Second):
			v_criticals_cache.Range(func(key, value any) bool {
				global.Debug("start run cache tx for 0 second -> 9 second", "time", value.(*TokenTrading).Timestamp.Unix())
				TokenTradingNotify(value.(*TokenTrading))
				return true
			})
			v_criticals_cache.Clear()
			v_criticals_cache = nil
			return
		}
	}
}
