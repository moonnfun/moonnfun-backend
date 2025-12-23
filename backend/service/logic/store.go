package logic

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"os"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	C_Min_HolderAmount = 0.001
)

func GetModelList[T any](modelName string, where bson.M) ([]*T, error) {
	tokens, err := store.DBList[T](modelName, where, false)
	if err != nil {
		return []*T{}, err
	}
	return tokens, nil
}

func GetModelListPage[T any](modelName, address, creator, tag, order, orderField string, offset, limit int) ([]*T, int, error) {
	whereModel := bson.M{}
	if order == "" {
		order = "desc"
	}

	workOrderField := "createdAt"
	if orderField != "" {
		workOrderField = orderField
	}
	if creator != "" {
		whereModel["creator"] = creator
	}
	if address != "" {
		whereModel["address"] = address
	}

	tokens, total, err := store.DBListPage[T](modelName, whereModel, order, workOrderField, offset, limit, false)
	if err != nil {
		return []*T{}, 0, err
	}
	// total := store.DBCount(modelName, whereModel)
	return tokens, int(total), nil
}

func GetModelListPageEx[T any](modelName string, whereModel bson.M, order, orderField string, offset, limit int) ([]*T, int, error) {
	if order == "" {
		order = "desc"
	}

	workOrderField := "createdAt"
	if orderField != "" {
		workOrderField = orderField
	}

	tokens, total, err := store.DBListPage[T](modelName, whereModel, order, workOrderField, offset, limit, false)
	if err != nil {
		return []*T{}, 0, err
	}
	// total := store.DBCount(model.C_Token, whereModel)
	return tokens, int(total), nil
}

func UpdateTokenForTrading(token *model.Token, tokenTrading *monitor.TokenTrading) error {
	var tokenAmount, totalAmount decimal.Decimal
	if ttrading := monitor.GetReverse(token.Address, global.Config.ContractAddress, token.Launch); ttrading != nil {
		tokenAmount = ttrading.TokenAmount
		totalAmount = ttrading.TotalAmount
	}
	if tokenAmount.GreaterThanOrEqual(decimal.NewFromInt(token.Supply)) {
		tokenAmount = decimal.NewFromInt(token.Supply)
	}
	if tokenAmount.LessThanOrEqual(decimal.Zero) {
		tokenAmount = decimal.Zero
	}

	tokenPrice := GetTokenPrice(token.Address, token.Launch, true)
	launchQuoteAmount, progress := GetProgress(totalAmount)
	if !token.Launch {
		token.Progress = progress
		if progress == 100 {
			token.Launch = true
			token.LaunchTime = time.Now().UnixMilli()
		}
	} else {
		token.Progress = 100
		totalAmount = launchQuoteAmount
	}
	if totalAmount.GreaterThanOrEqual(launchQuoteAmount) {
		totalAmount = launchQuoteAmount
	}
	if totalAmount.LessThanOrEqual(decimal.Zero) {
		totalAmount = decimal.Zero
	}

	quotePrice := monitor.GetQuotePrice()
	token.DBID = primitive.ObjectID{}
	token.VolumeTime = time.Now().UnixMilli()
	token.Volume24 = GetVolume24hForToken(token, true)
	token.Holders = int(store.DBCount(model.C_Token_holder, bson.M{"address": tokenTrading.Token.String()}))
	token.TokenTrading.MarketCap, _ = decimal.NewFromString(tokenPrice.PriceUsd.Mul(decimal.NewFromInt(token.Supply)).StringFixed(20)) // 采用USDT计价
	token.TokenTrading.Volume, _ = decimal.NewFromString(token.TokenTrading.Volume.Add(tokenTrading.QuoteAmount.Mul(quotePrice)).StringFixed(20))
	token.TokenTrading.VirtualLiquidity = totalAmount.Mul(quotePrice).Mul(decimal.NewFromInt(2))
	token.TokenTrading.TokenAmount, _ = decimal.NewFromString(tokenAmount.StringFixed(20))
	token.TokenTrading.TotalAmount, _ = decimal.NewFromString(totalAmount.StringFixed(20))
	if dbToken, dbErr := GetToken(token.Address); dbToken != nil && dbErr == nil {
		if dbToken.Launch && !token.Launch {
			token.Launch = dbToken.Launch
			token.Progress = dbToken.Progress
		}
		if err := store.DBSet(model.C_Token, token, bson.M{"address": token.Address}); err != nil {
			return err
		}
	} else {
		return dbErr
	}
	slog.Info("update token successed", "token", token)
	return nil
}

func UpdateTrade(token *model.Token, tokenTrading *monitor.TokenTrading, bLaunch bool) (*model.Trade, error) {
	qamount, _ := decimal.NewFromString(tokenTrading.QuoteAmount.StringFixed(20))
	amount, _ := decimal.NewFromString(tokenTrading.Amount.StringFixed(20))
	price, _ := decimal.NewFromString(tokenTrading.Price.StringFixed(20))
	tprice := GetTokenPrice(tokenTrading.Token.String(), bLaunch, true)
	if tprice.Price.LessThanOrEqual(decimal.Zero) {
		tprice.Price = price
	}
	if tprice.PriceUsd.LessThanOrEqual(decimal.Zero) {
		tprice.PriceUsd = price.Mul(monitor.GetQuotePrice())
	}

	// // udpate quoteAmount
	// totalQuote := token.TotalAmount.Add(tokenTrading.QuoteAmount)
	// launchQuoteAmount, progress := GetProgress(totalQuote)
	// if !token.Launch && progress == 100 && tokenTrading.Buy {
	// 	if totalQuote.GreaterThan(launchQuoteAmount) {
	// 		diff := totalQuote.Sub(launchQuoteAmount)
	// 		qamount, _ = decimal.NewFromString(tokenTrading.QuoteAmount.Sub(diff).StringFixed(20))
	// 		global.Debug("update quote amount successed", "diff", diff, "quoteAmount", tokenTrading.QuoteAmount, "qamount", qamount)
	// 	}
	// }

	trade := &model.Trade{
		Address:     tokenTrading.Token.String(),
		Wallet:      tokenTrading.Wallet.String(),
		Buy:         tokenTrading.Buy,
		Liquid:      tokenTrading.Liquid,
		Amount:      amount,
		Price:       price,
		Number:      tokenTrading.Number,
		MPrice:      tprice.Price, // tprice.PriceUsd,
		TxHash:      tokenTrading.TxHash,
		Source:      tokenTrading.Source,
		Timestamp:   tokenTrading.Timestamp.UnixMilli(),
		QuoteAmount: qamount,
	}
	dbID := primitive.NewObjectID()
	if dbTrade, dbErr := store.DBGet[model.Trade](model.C_Trade, bson.M{"txhash": tokenTrading.TxHash}); dbTrade != nil && dbErr == nil {
		slog.Error("before update trade with exists txhash", "trade", dbTrade)
		return nil, fmt.Errorf("before update trade with exists txhash, trade: %+v", dbTrade)
	} else {
		trade.DBID = dbID
		trade.CreatedAt = time.Now().UnixMilli()
	}
	if err := store.DBSet(model.C_Trade, trade, bson.M{"_id": dbID}); err != nil {
		slog.Error("create trade failed", "trade", trade, "error", err.Error())
		return nil, err
	}
	UdpateTokenPriceByTrade(trade)

	// earliestPriceMap.Store(trade.Address, &model.TokenPrice{
	// 	Price:       trade.MPrice,
	// 	CurrentTime: trade.Timestamp,
	// })

	slog.Info("update trade successed", "trade", trade)
	return trade, nil
}

func UpdateTokenHoler(trade *model.Trade, tokenTrading *monitor.TokenTrading) error {
	var createdAt int64
	amountFromDB := decimal.Zero
	var tokenDBID primitive.ObjectID
	if dbHolder, err := store.DBGet[model.TokenHolder](model.C_Token_holder, bson.M{"wallet": tokenTrading.Wallet.String(), "address": tokenTrading.Token.String()}); err == nil && dbHolder != nil {
		tokenDBID = dbHolder.DBID
		createdAt = dbHolder.CreatedAt
		amountFromDB = dbHolder.Amount
	} else {
		tokenDBID = primitive.NewObjectID()
		createdAt = time.Now().UnixMilli()
	}

	// update fileds
	amount := monitor.GetTokenBalance(tokenTrading.Token.String(), trade.Wallet)
	if !tokenTrading.Buy {
		if amount.LessThanOrEqual(decimal.NewFromFloat(C_Min_HolderAmount)) && !amount.Equal(decimal.NewFromInt(-1)) {
			if err := store.DBDelete(model.C_Token_holder, bson.M{"_id": tokenDBID}); err != nil {
				slog.Error("delete tokenHolder failed", "tokenTrading", tokenTrading, "error", err.Error())
				return err
			}
			slog.Info("delete tokenHolder successed", "tokenTrading", tokenTrading)
			return nil
		}
		amountFromDB = amountFromDB.Sub(tokenTrading.Amount)
	} else {
		amountFromDB = amountFromDB.Add(tokenTrading.Amount)
	}
	if amount.Equal(decimal.NewFromInt(-1)) {
		amount = amountFromDB
	}
	percent := amount.Div(decimal.NewFromInt(model.C_total_supply)).Mul(decimal.NewFromInt(100))

	dbAmount, _ := decimal.NewFromString(amount.StringFixed(20))
	dbPercent, _ := decimal.NewFromString(percent.StringFixed(20))
	avgPrice, err := GetAveragePrice(tokenTrading.Token.String(), tokenTrading.Wallet.String(), amount)
	if err != nil || avgPrice.Equal(decimal.Zero) {
		slog.Error("GetAveragePrice failed", "tokenTrading", tokenTrading)
		avgPrice = trade.Price
	}

	tokenHolder := &model.TokenHolder{
		Address:  tokenTrading.Token.String(),
		Wallet:   tokenTrading.Wallet.String(),
		Amount:   dbAmount,
		Percent:  dbPercent,
		AvgPrice: avgPrice,
	}
	if createdAt != 0 {
		tokenHolder.CreatedAt = createdAt
	}
	if err := store.DBSet(model.C_Token_holder, tokenHolder, bson.M{"_id": tokenDBID}); err != nil {
		slog.Error("update tokenHolder failed", "tokenHolder", tokenHolder, "error", err.Error())
		return err
	}
	slog.Info("update tokenHolder successed", "tokenHolder", tokenHolder)
	return nil
}

func GetUserOwnedTokens(tag, address, wallet, order, orderField string, offset, limit int) ([]*model.Token, int, error) {
	iOrder := -1
	if order == "asc" {
		iOrder = 1
	}
	workOrderField := "createdAt"
	if orderField != "" {
		workOrderField = orderField
	}
	if offset > 0 {
		offset -= 1
	}
	if limit == 0 {
		limit = -1
	}

	//	mongo.Pipeline{
	//		{{"$group", bson.D{{"_id", "$state"}, {"totalPop", bson.D{{"$sum", "$pop"}}}}}},
	//		{{"$match", bson.D{{"totalPop", bson.D{{"$gte", 10*1000*1000}}}}}},
	//	}
	pipeline := mongo.Pipeline{
		{{
			"$match", bson.D{{"wallet", wallet}, {"amount", bson.D{{"$gte", decimal.NewFromFloat(C_Min_HolderAmount)}}}},
		}},
		{{
			"$lookup", []bson.E{
				{"from", "token"},
				{"localField", "address"},
				{"foreignField", "address"},
				{"as", "tokens"},
			},
		}},
		{{
			"$sort", []bson.E{
				{"tokens." + workOrderField, iOrder},
			},
		}},
	}
	// if tag != "" {
	// 	tagE := bson.E{"tag", tag}
	// 	pipeline = mongo.Pipeline{
	// 		{{
	// 			"$match", bson.D{{"wallet", wallet}, {"amount", bson.D{{"$gte", decimal.NewFromFloat(C_Min_HolderAmount)}}}},
	// 		}},
	// 		{{
	// 			"$lookup", []bson.E{
	// 				{"from", "token"},
	// 				{"localField", "address"},
	// 				{"foreignField", "address"},
	// 				{"as", "tokens"},
	// 			},
	// 		}},
	// 		{{
	// 			"$project", []bson.E{{
	// 				"tokens", []bson.E{
	// 					{
	// 						"filter", []bson.E{
	// 							{"input", "$tokens"},
	// 							{"as", "token"},
	// 							{"cond", []bson.E{
	// 								tagE,
	// 							}},
	// 						},
	// 					},
	// 					// {
	// 					// 	"sort", []bson.E{
	// 					// 		{"tokens." + workOrderField, iOrder},
	// 					// 	},
	// 					// },
	// 					// {"skip", int64(offset)},
	// 					// {"limit", int64(limit)},
	// 				},
	// 			}},
	// 		}},
	// 	}
	// }
	// if address != "" {
	// 	addressE := bson.E{}
	// 	if strings.HasPrefix(address, "0x") {
	// 		addressE = bson.E{"address", address}
	// 	} else {
	// 		addressE = bson.E{"symbol", address}
	// 	}
	// 	pipeline = mongo.Pipeline{
	// 		{{
	// 			"$match", bson.D{{"wallet", wallet}, {"amount", bson.D{{"$gte", decimal.NewFromFloat(C_Min_HolderAmount)}}}},
	// 		}},
	// 		{{
	// 			"$lookup", []bson.E{
	// 				{"from", "token"},
	// 				{"localField", "address"},
	// 				{"foreignField", "address"},
	// 				{"as", "tokens"},
	// 			},
	// 		}},
	// 		{{
	// 			"$project", []bson.E{{
	// 				"tokens", []bson.E{
	// 					{
	// 						"filter", []bson.E{
	// 							{"input", "$tokens"},
	// 							{"as", "token"},
	// 							{"cond", []bson.E{
	// 								addressE,
	// 							}},
	// 						},
	// 					},
	// 					// {
	// 					// 	"sort", []bson.E{
	// 					// 		{"tokens." + workOrderField, iOrder},
	// 					// 	},
	// 					// },
	// 					// {"skip", int64(offset)},
	// 					// {"limit", int64(limit)},
	// 				},
	// 			}},
	// 		}},
	// 	}
	// }

	tokenHolderList := []struct {
		model.TokenHolder
		Tokens []*model.Token `bson:"tokens"`
	}{}
	if err := store.GetDB().Collection(model.C_Token_holder).Aggregate(context.Background(), pipeline).All(&tokenHolderList); err != nil {
		return nil, 0, err
	}
	global.Debug("GetUserOwnedTokens successed", "tokenHolderList", tokenHolderList)

	tokenList := make([]*model.Token, 0)
	for _, tokenHolder := range tokenHolderList {
		tokenList = append(tokenList, tokenHolder.Tokens...)
	}
	return tokenList, len(tokenList), nil
}

func WriteFile(content []byte, outfile string) error {

	file, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.Write(content)
	writer.Flush()
	return err
}

func GetPairAddress(tokenAddress string) string {
	ticker := time.NewTicker(time.Second)
	routerAddress := global.Config.RPC.SwapRouterAddress
	for {
		select {
		case <-ticker.C:
			pairAddress, _ := monitor.GetPair(tokenAddress, routerAddress)
			if pairAddress != "" {
				return pairAddress
			}
		}
	}
}

func GetLatestNumber() int64 {
	tradeList, _, err := GetModelListPageEx[model.Trade](model.C_Trade, bson.M{}, "desc", "timestamp", 0, 1)
	if err == nil && len(tradeList) > 0 {
		tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "blocknumber", 0, 1)
		if err == nil && len(tradeList) > 0 {
			if tokenList[0].BlockNumber > tradeList[0].Number {
				return tokenList[0].BlockNumber
			}
		}
		return tradeList[0].Number
	}
	return 0
}

func GetTxRecordListForSync() []*model.TxRecord {
	timestamp := time.Now().Add(-24 * time.Hour).UnixMilli()
	whereModel := bson.D{{"createdat", bson.D{{"$gte", timestamp}}}}
	txRecordList, _, err := store.DBListPage[model.TxRecord](model.C_TxRecord, whereModel, "asc", "createdat", 0, 1, false)
	if err == nil && len(txRecordList) > 0 {
		return txRecordList
	}
	return []*model.TxRecord{}
}

func SaveTxRecord(txRecord *model.TxRecord) error {
	dbID := primitive.NewObjectID()
	if dbRecord, dbErr := store.DBGet[model.TxRecord](model.C_TxRecord, bson.M{"number": txRecord.Number, "iscreate": txRecord.IsCreate}); dbRecord != nil && dbErr == nil {
		dbID = dbRecord.DBID
		txRecord.DBID = primitive.ObjectID{}
	} else {
		txRecord.DBID = dbID
		txRecord.CreatedAt = time.Now().UnixMilli()
	}

	if err := store.DBSet(model.C_TxRecord, txRecord, bson.M{"_id": dbID}); err != nil {
		slog.Error("create tx record failed", "txRecord", txRecord, "error", err.Error())
		return err
	}
	return nil
}
