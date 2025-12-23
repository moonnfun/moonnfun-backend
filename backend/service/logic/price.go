package logic

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

var tokenPriceMap sync.Map
var earliestPriceMap sync.Map

func pushMsgID(tokenAddress string, timestamp int64) string {
	return fmt.Sprintf("%s_%v", tokenAddress, timestamp)
}

func StartKLinePush(userID, tokenAddress, timePeroid string) (pushStopCh chan struct{}) {
	slog.Info("start init kLine", "userID", userID, "tokenAddress", tokenAddress, "timePeroid", timePeroid)
	tradingList, _, err := GetModelListPage[model.Trade](model.C_Trade, tokenAddress, "", "", "asc", "timestamp", 0, 0)
	if err != nil {
		slog.Error("start kline push failed", "tokenAddress", tokenAddress, "timePeroid", timePeroid, "error", err.Error())
		return
	}
	InitKLine(KlineType_Init, tokenAddress, timePeroid)

	// 第一根k线
	klines := make([]model.TokenOHLCV, 0)
	token, err := GetToken(tokenAddress)
	if err != nil {
		slog.Error("start kline push failed", "tokenAddress", tokenAddress, "timePeroid", timePeroid, "error", err.Error())
		return
	}
	_, kline, err := UpdateKLine(KlineType_Init, tokenAddress, timePeroid, token.StartPrice, time.UnixMilli(token.StartTime), decimal.Zero)
	if err != nil {
		slog.Error("update kline failed", "token", token, "error", err.Error())
		return
	}
	klines = append(klines, kline)
	SaveKLine(KlineType_Init, tokenAddress, timePeroid, kline)

	for _, tokenTrading := range tradingList {
		bUpdate, currentKLine, err := UpdateKLine(KlineType_Init, tokenTrading.Address, timePeroid, tokenTrading.MPrice, time.UnixMilli(tokenTrading.Timestamp), tokenTrading.Amount.Mul(tokenTrading.Price))
		if err != nil {
			slog.Error("update kline failed", "tokenTrading", tokenTrading, "error", err.Error())
			continue
		}
		if currentKLine.T > 0 {
			global.Debug("kline: ", currentKLine)
			if bUpdate && len(klines) > 0 {
				klines[len(klines)-1] = currentKLine
			} else {
				klines = append(klines, currentKLine)
			}
		}

	}
	if len(klines) > 0 {
		slog.Info("before init kLine", "userID", userID, "tokenAddress", tokenAddress, "timePeroid", timePeroid, "klines", len(klines))
		global.WebsocketSend(userID, pushMsgID(tokenAddress, time.Now().UnixMilli()), global.WssMsgType(global.C_Msg_Price_Init, timePeroid), klines)
	}
	return
}

func StopKLinePush(pushStopCh chan struct{}) {
	if pushStopCh != nil {
		close(pushStopCh)
	}
}

func UpdateTokenTrading(data any) (bHandle bool) {
	tokenTrading := data.(*monitor.TokenTrading)

	var token *model.Token
	if global.Config.Mode == global.C_Mode_Sync {
		if token = GetTokenEx(tokenTrading.Token.String(), tokenTrading.Pair.String()); token == nil {
			slog.Error("get token from database failed", "tokenTrading", tokenTrading)
			return
		}
	} else {
		if token = MustGetToken(tokenTrading.Token.String(), tokenTrading.Pair.String(), time.Duration(5)*time.Minute); token == nil {
			slog.Error("get token from database failed with timeout 5 minutes", "tokenTrading", tokenTrading)
			return
		}
	}
	global.Debug("get token successed", "token", token)
	tokenTrading.Token = common.HexToAddress(token.Address)

	// udpate trade
	var trade *model.Trade
	if tr, err := UpdateTrade(token, tokenTrading, token.Launch); err != nil {
		slog.Error("update trade failed", "tokenTrading", tokenTrading, "error", err.Error())
		return
	} else {
		trade = tr
	}

	// update kline
	for _, v := range getTimeList() {
		bInit := GetKLine(KlineType_Update, tokenTrading.Token.String(), v) == nil
		if bInit {
			InitKLine(KlineType_Update, tokenTrading.Token.String(), v)
		}
		bUpdate, kline, err := UpdateKLine(KlineType_Update, tokenTrading.Token.String(), v, trade.MPrice, tokenTrading.Timestamp, tokenTrading.Amount.Mul(tokenTrading.Price))
		if err != nil {
			slog.Error("update kline failed", "tokenTrading", tokenTrading, "error", err.Error())
			continue
		}

		if bInit {
			global.WebsocketSend("", pushMsgID(tokenTrading.Token.String(), int64(kline.T)), global.WssMsgType(global.C_Msg_Price_Init, v), []model.TokenOHLCV{kline})
		} else if bUpdate {
			global.WebsocketSend("", pushMsgID(tokenTrading.Token.String(), int64(kline.T)), global.WssMsgType(global.C_Msg_Price_Update, v), []model.TokenOHLCV{kline})
		} else {
			global.WebsocketSend("", pushMsgID(tokenTrading.Token.String(), int64(kline.T)), global.WssMsgType(global.C_Msg_Price_Append, v), []model.TokenOHLCV{kline})
		}
	}
	return true
}

func GetEarliestTokenPrice(token *model.Token) decimal.Decimal {
	timestamp := time.Now().Add(-24 * time.Hour).UnixMilli()
	whereModel := bson.D{{"address", token.Address}, {"timestamp", bson.D{{"$gte", timestamp}}}}
	tradeList, _, err := store.DBListPage[model.Trade](model.C_Trade, whereModel, "asc", "timestamp", 0, 1, false)
	if err == nil && len(tradeList) > 0 {
		slog.Info("get earliest token price successed within 24h", "timestamp", timestamp, "token", token, "trade", tradeList[0])
		earliestPriceMap.Store(tradeList[0].Address, &model.TokenPrice{
			Price:       tradeList[0].MPrice,
			CurrentTime: tradeList[0].Timestamp,
		})
		return tradeList[0].MPrice
	} else {
		return decimal.Zero
	}
}

func GetTokenPricePercent(token *model.Token) (decimal.Decimal, float64) {
	price := decimal.Zero
	startTime := token.CreatedAt
	timestamp := time.Now().Add(-24 * time.Hour).UnixMilli()
	if startTime >= timestamp {
		price = token.StartPrice
	} else {
		v, ok := earliestPriceMap.Load(token.Address)
		if ok && v.(*model.TokenPrice).CurrentTime >= timestamp {
			price = v.(*model.TokenPrice).Price
			global.Debug("get earliest price successed from memory", "price", price)
		} else {
			price = GetEarliestTokenPrice(token)
			global.Debug("get earliest price successed from trade", "price", price)
		}
	}

	tokenPrice := GetTokenPrice(token.Address, token.Launch, false)
	if price.GreaterThan(decimal.Zero) {
		global.Debug("====================", "price", price, "tokenPrice", tokenPrice, "percent", tokenPrice.Price.Sub(price).Div(price).Mul(decimal.NewFromInt(100)))
		percentS := tokenPrice.Price.Sub(price).Div(price).Mul(decimal.NewFromInt(100)).String()
		percentF, _ := ParseFloatNoSci(percentS)
		return tokenPrice.PriceUsd, percentF
	}
	// global.Debug("get tokenPrice ok", "tokenPrice", tokenPrice)
	return tokenPrice.PriceUsd, 0
}

func ParseFloatNoSci(s string) (float64, error) {
	ret, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	strRet := fmt.Sprintf("%v", ret)
	if strings.Contains(strRet, "e") || strings.Contains(strRet, "E") {
		return 0, errors.New("scientific notation not allowed")
	}
	return ret, nil
}

func InitTokenPriceTimer(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(global.Config.RefreshTokenPriceTimer) * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				UpdateTokenPrice()
			case <-ctx.Done():
				slog.Info("stop token price update")
				return
			}
		}
	}()
}

func GetTokenPrice(tokenAddress string, bLaunch bool, bForce bool) *model.TokenPrice {
	t, ok := tokenPriceMap.Load(tokenAddress)
	if ok && !bForce {
		if t.(*model.TokenPrice).PriceUsd.GreaterThan(decimal.Zero) {
			return t.(*model.TokenPrice)
		}
	}

	// tokenPrice := GetTokenPriceFromTrade(tokenAddress, bLaunch)
	tokenPrice := GetTokenPriceFromContract(tokenAddress, bLaunch)
	if tokenPrice.PriceUsd.LessThanOrEqual(decimal.Zero) {
		tokenPrice = GetTokenPriceFromTrade(tokenAddress, bLaunch)
	}
	tokenPriceMap.Store(tokenAddress, tokenPrice)
	global.Debug("GetTokenPriceFromContract", "tokenAddress", tokenAddress, "tokenPrice", tokenPrice)
	return tokenPrice
}

func GetTokenPriceFromTrade(tokenAddress string, bLaunch bool) *model.TokenPrice {
	method := "GetTokenPrice"
	tokenPrice := model.TokenPrice{}
	if bLaunch {
		return GetTokenPriceFromContract(tokenAddress, bLaunch)
	} else {
		trades, _, _ := GetModelListPageEx[model.Trade](model.C_Trade, bson.M{"address": tokenAddress}, "", "", 0, 1)
		if trades == nil || len(trades) == 0 {
			return GetTokenPriceFromContract(tokenAddress, bLaunch)
		} else {
			quotePrice := monitor.GetQuotePrice()
			tokenPrice.Price = trades[0].MPrice
			tokenPrice.PriceUsd = trades[0].MPrice.Mul(quotePrice)
			// tokenPrice.PriceUsd = trades[0].QuoteAmount.Div(trades[0].Amount).Mul(quotePrice)
		}
	}
	if tokenPrice.Price.LessThanOrEqual(decimal.Zero) || tokenPrice.PriceUsd.LessThanOrEqual(decimal.Zero) {
		return GetTokenPriceFromContract(tokenAddress, bLaunch)
	}

	tokenPrice.Method = method
	return &tokenPrice
}

func GetTokenPriceFromContract(tokenAddress string, bLaunch bool) *model.TokenPrice {
	method := "GetTokenPrice"
	if bLaunch {
		method = "GetTokenPriceEx"
	}
	tokenPrice := monitor.GetTokenPrice(method, tokenAddress, global.Config.ContractAddress)
	tokenPrice.Method = method
	return &tokenPrice
}

func UdpateTokenPriceByTrade(trade *model.Trade) {
	t, ok := tokenPriceMap.Load(trade.Address)
	if ok {
		if t.(*model.TokenPrice).PriceUsd.GreaterThan(decimal.Zero) && trade.Amount.GreaterThan(decimal.Zero) {
			tokenPrice := t.(*model.TokenPrice)
			quotePrice := monitor.GetQuotePrice()
			tokenPrice.PriceUsd = trade.QuoteAmount.Div(trade.Amount).Mul(quotePrice)
			tokenPriceMap.Store(trade.Address, tokenPrice)
		}
	}
}

func UpdateTokenPrice() {
	tokenPriceMap.Range(func(key, value any) bool {
		tokenPrice := monitor.GetTokenPrice(value.(*model.TokenPrice).Method, fmt.Sprintf("%v", key), global.Config.ContractAddress) //+"1")
		tokenPrice.Method = value.(*model.TokenPrice).Method
		tokenPriceMap.Store(key, &tokenPrice)
		return true
	})
}

func GetProgress(quoteAmount decimal.Decimal) (decimal.Decimal, float64) {
	launchQuoteAmount := decimal.NewFromFloat(global.Config.LaunchMC)
	progress := quoteAmount.Div(launchQuoteAmount).Mul(decimal.NewFromInt(100)).InexactFloat64()
	if progress >= 100 {
		progress = 100
	} else if progress <= 0 {
		progress = 0
	}
	global.Debug("get progress successed", quoteAmount, launchQuoteAmount, progress)
	return launchQuoteAmount, ToFixed(progress, 2)
}

func ToFixed(num float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(num*pow) / pow
}

func GetAveragePrice(tokenAddress, wallet string, totalAmount decimal.Decimal) (decimal.Decimal, error) {
	trades, _, err := GetModelListPageEx[model.Trade](model.C_Trade, bson.M{"address": tokenAddress, "wallet": wallet}, "asc", "createdAt", 0, 0)
	if err != nil {
		return decimal.Zero, err
	}

	totalCost := decimal.NewFromInt(0) // 总成本, Quote计价
	for _, t := range trades {
		if t.Buy {
			totalCost = totalCost.Add(t.QuoteAmount)
		} else {
			totalCost = totalCost.Sub(t.QuoteAmount)
		}
	}

	if totalAmount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, nil
	}
	result := totalCost.Div(totalAmount)
	if result.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, nil
	}
	return result, nil
}
