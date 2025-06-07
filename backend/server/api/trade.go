package api

import (
	"meme3/service/logic"
	"meme3/service/model"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func initTrade(api *swag.API) {
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/trade/hot",
			endpoint.Tags("Trade"),
			endpoint.Handler(webTradeHot),
			endpoint.Summary("hot trades"),
			endpoint.Description("get hot trades"),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.Trade{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/trade/list",
			endpoint.Handler(webTrades),
			endpoint.Tags("Trade"),
			endpoint.Summary("trade list"),
			endpoint.Description("get trade list"),
			endpoint.Query("tag", "string", "token tag", false),
			endpoint.Query("creator", "string", "token creator", false),
			endpoint.Query("address", "string", "token creator", false),
			endpoint.Query("order", "string", "pages order", false),
			endpoint.Query("orderField", "string", "order field", false),
			endpoint.Query("offset", "string", "pages offset", false),
			endpoint.Query("limit", "string", "pages limit", false),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.Trade{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
}

func webTradeHot(w http.ResponseWriter, r *http.Request) {
	// tokenHolders, err := logic.GetHotTrades()
	// if err != nil {
	// 	WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
	// 	return
	// }
	// WebResponseJson(w, r, ApiResponse(tokenHolders, true), http.StatusOK)
}

func webTrades(w http.ResponseWriter, r *http.Request) {
	tag := WebParams(r).Get("tag")
	order := WebParams(r).Get("order")
	limit := WebParams(r).Get("limit")
	offset := WebParams(r).Get("offset")
	creator := WebParams(r).Get("creator")
	address := WebParams(r).Get("address")
	orderField := WebParams(r).Get("orderField")

	if tag != "" {
		tag = tags[tag]
	}
	if orderField != "" {
		orderField = orderFields[orderField]
	}

	iLimit, _ := strconv.Atoi(limit)
	iOffset, _ := strconv.Atoi(offset)

	tradeList, total, err := logic.GetModelListPage[model.Trade](model.C_Trade, address, creator, tag, order, orderField, iOffset, iLimit)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}

	if iLimit == 10 {
		ret := make([]any, 0)
		for _, trade := range tradeList {
			t, _ := logic.GetToken(trade.Address)
			retTrade := struct {
				model.Trade
				model.Token
			}{}
			tokenPrice := logic.GetTokenPrice(t.Address)
			retTrade.Trade = *trade
			retTrade.Token = *t
			retTrade.Token.MarketCap = decimal.NewFromInt(retTrade.Token.Supply).Mul(tokenPrice)
			ret = append(ret, retTrade)
		}
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else {
		WebResponseJson(w, r, ApiResponseList(tradeList, total), http.StatusOK)
	}
}
