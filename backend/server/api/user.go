package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"meme3/global"
	"meme3/server/web"
	"meme3/service/logic"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"

	"github.com/shopspring/decimal"
	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"go.mongodb.org/mongo-driver/bson"
)

func initUser(api *swag.API) {
	// API User
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/detail",
			endpoint.Tags("User"),
			endpoint.Handler(webUserDetail),
			endpoint.Summary("User detail"),
			endpoint.Description("User detail"),
			endpoint.Query("address", "string", "wallet address", false),
			endpoint.Response(http.StatusOK, "Successfully get user", endpoint.SchemaResponseOption(model.User{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/user/login",
			endpoint.Tags("User"),
			endpoint.Handler(WebUserLogin),
			endpoint.Summary("User login"),
			endpoint.Description("User login"),
			endpoint.FormData("refer", "string", "user referral", false),
			endpoint.FormData("address", "string", "wallet address", true),
			endpoint.FormData("message", "string", "signature message", true),
			endpoint.FormData("signature", "string", "wallet signature", true),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.User{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/user/bind",
			endpoint.Tags("User"),
			endpoint.Handler(WebUserBind),
			endpoint.Summary("User bind"),
			endpoint.Description("User bind"),
			endpoint.FormData("refer", "string", "user referral", false),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.User{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/user/follow",
			endpoint.Tags("User"),
			endpoint.Handler(WebUserFollow),
			endpoint.Summary("User bind"),
			endpoint.Description("User bind"),
			endpoint.FormData("tokenID", "string", "token ID", true),
			endpoint.FormData("follow", "string", "follow or not", true),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.User{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/verify",
			endpoint.Tags("User"),
			endpoint.Handler(webUserVerify),
			endpoint.Summary("User verify"),
			endpoint.Description("User verify"),
			endpoint.Response(http.StatusOK, "Successfully verify player", endpoint.SchemaResponseOption("successed")),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/logout",
			endpoint.Tags("User"),
			endpoint.Handler(WebUserLogout),
			endpoint.Summary("User logout"),
			endpoint.Description("User logout"),
			endpoint.Response(http.StatusOK, "Successfully added logout status", endpoint.SchemaResponseOption("successed")),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/statistics",
			endpoint.Tags("User"),
			endpoint.Handler(webUserStatistics),
			endpoint.Summary("User statistics"),
			endpoint.Description("User statistics"),
			endpoint.Query("address", "string", "wallet address", false),
			endpoint.Response(http.StatusOK, "Successfully get user", endpoint.SchemaResponseOption(model.UserStatistics{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/referrals",
			endpoint.Tags("User"),
			endpoint.Handler(webUserReferrals),
			endpoint.Summary("User referrals"),
			endpoint.Description("User referrals"),
			endpoint.Query("order", "string", "pages order", false),
			endpoint.Query("orderField", "string", "order field", false),
			endpoint.Query("offset", "string", "pages offset", false),
			endpoint.Query("limit", "string", "pages limit", false),
			endpoint.Response(http.StatusOK, "Successfully get user", endpoint.SchemaResponseOption([]model.Refferal{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/holders",
			endpoint.Handler(webUserHolders),
			endpoint.Tags("User"),
			endpoint.Summary("user token list"),
			endpoint.Description("get user token list"),
			endpoint.Query("tag", "string", "token tag", false),
			endpoint.Query("creator", "string", "token creator", false),
			endpoint.Query("address", "string", "token creator", false),
			endpoint.Query("order", "string", "pages order", false),
			endpoint.Query("orderField", "string", "order field", false),
			endpoint.Query("offset", "string", "pages offset", false),
			endpoint.Query("limit", "string", "pages limit", false),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.TokenHolder{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/user/follows",
			endpoint.Handler(webUserFollowList),
			endpoint.Tags("User"),
			endpoint.Summary("user follow list"),
			endpoint.Description("get user follow list"),
			endpoint.Query("tag", "string", "token tag", false),
			endpoint.Query("creator", "string", "token creator", false),
			endpoint.Query("address", "string", "token address", false),
			endpoint.Query("order", "string", "pages order", false),
			endpoint.Query("orderField", "string", "order field", false),
			endpoint.Query("offset", "string", "pages offset", false),
			endpoint.Query("limit", "string", "pages limit", false),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.Token{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/user/support",
			endpoint.Tags("User"),
			endpoint.Handler(WebUserSupport),
			endpoint.Summary("User support"),
			endpoint.Description("User support"),
			endpoint.Body(model.Project{}, "Help object that needs to be added to the store", true),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.User{})),
		),
	)
}

func WebUserSupport(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	user, err := logic.GetUser(fmt.Sprintf("%v", userID), false)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	_, project, err := WebBody[model.Project](r)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	// make key
	project.Key = store.NewId()

	user.Project = *project

	if err := logic.SaveUserToDB(user, true); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	} else {
		WebResponseJson(w, r, ApiResponse(user, true), http.StatusOK)
	}
}

func WebUserBind(w http.ResponseWriter, r *http.Request) {
	refer := WebParams(r).Get("refer")

	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	user, err := logic.GetUser(fmt.Sprintf("%v", userID), false)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	if err := logic.SaveUser(refer, user, true); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponse(user, true), http.StatusOK)
}

func WebUserFollow(w http.ResponseWriter, r *http.Request) {
	tokenID := WebParams(r).Get("tokenID")
	follow := WebParams(r).Get("follow")

	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	token, err := logic.AddFollowToken(userID, tokenID, follow == "true")
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponse(token, true), http.StatusOK)
}

func WebUserLogin(w http.ResponseWriter, r *http.Request) {
	refer := WebParams(r).Get("refer")
	address := WebParams(r).Get("address")
	message := WebParams(r).Get("message")
	signature := WebParams(r).Get("signature")

	// check
	if userID := web.PopFromSession(r, web.C_Session_User); userID != nil {
		user, _ := logic.GetUser(userID, true)
		WebResponseJson(w, r, ApiResponse(user, true), http.StatusOK)
		return
	}

	// user login
	user, err := logic.UserLogin(address, signature, message, refer)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	// session
	if err := web.PushToSession(w, r, web.C_Session_User, address); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponse(user, true), http.StatusOK)
}

func webUserDetail(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	user, err := logic.GetUser(fmt.Sprintf("%v", userID), false)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponse(user, true), http.StatusOK)
}

func webUserVerify(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}
	WebResponseJson(w, r, ApiResponse("successed", true), http.StatusOK)
}

func WebUserLogout(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID != nil {
		logic.RemoveUser(userID)
	}
	WebResponseJson(w, r, ApiResponse("successed", true), http.StatusOK)
}

func webUserHoldersForTokens(w http.ResponseWriter, r *http.Request) {
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
	slog.Info("before get user's owned tokens", "limit", iLimit, "offset", iOffset)

	tokenList, total, err := logic.GetUserOwnedTokens("", "", creator, order, orderField, iOffset, iLimit)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}
	for i, token := range tokenList {
		tokenPriceUsd, percent := logic.GetTokenPricePercent(token)
		tokenList[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(1000000000))
		tokenList[i].Percent = percent
	}

	ret := make([]*model.Token, 0)
	if tag != "" {
		for _, t := range tokenList {
			if t.Tag == tag {
				ret = append(ret, t)
			}
		}
	}
	if address != "" {
		for _, t := range tokenList {
			if t.Address == address {
				ret = append(ret, t)
			}
		}
	}
	if tag != "" || address != "" {
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else {
		WebResponseJson(w, r, ApiResponseList(tokenList, total), http.StatusOK)
	}
}

func webUserHolders(w http.ResponseWriter, r *http.Request) {
	tag := WebParams(r).Get("tag")
	order := WebParams(r).Get("order")
	limit := WebParams(r).Get("limit")
	offset := WebParams(r).Get("offset")
	creator := WebParams(r).Get("creator")
	address := WebParams(r).Get("address")
	orderField := WebParams(r).Get("orderField")
	creator = model.FormatAddress(creator)

	whereModel := bson.M{}
	if tag != "" {
		tag = tags[tag]
	}
	if orderField != "" {
		orderField = orderFields[orderField]
	}
	if creator != "" {
		whereModel["wallet"] = creator
	}
	if address != "" {
		if strings.HasPrefix(address, "0x") {
			whereModel["address"] = address
		}
	}
	whereModel["amount"] = bson.M{"$gte": decimal.NewFromFloat(logic.C_Min_HolderAmount)}

	iLimit, _ := strconv.Atoi(limit)
	iOffset, _ := strconv.Atoi(offset)
	slog.Info("before get user's owned tokens", "limit", iLimit, "offset", iOffset)
	tokenHolders, total, err := logic.GetModelListPageEx[model.TokenHolder](model.C_Token_holder, whereModel, order, orderField, iOffset, iLimit)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	tokenHolderList := make([]*model.UserTokenHolder, 0)
	for _, holder := range tokenHolders {
		token, err := logic.GetToken(holder.Address)
		if err != nil {
			slog.Error("get token for holder failed", "holder", holder)
			continue
		}
		if tag != "" && token.Tag != tag {
			continue
		}
		if address != "" {
			if !strings.HasPrefix(address, "0x") && !strings.HasPrefix(token.Symbol, address) {
				continue
			}
		}
		token.Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), token.ID)
		userHolder := &model.UserTokenHolder{
			TokenAddress: token.Address,
			TokenHolder:  *holder,
			Token:        *token,
		}
		if token.Launch {
			userHolder.Token.Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), token.ID)
			tokenPriceUsd, percent := logic.GetTokenPricePercent(token)
			userHolder.Token.MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
			userHolder.Token.Volume24 = logic.GetVolume24hForToken(token, false)
			userHolder.Token.Percent = percent
			userHolder.TokenHolder.Amount = monitor.GetTokenBalance(holder.Address, holder.Wallet)
			userHolder.TokenHolder.Percent = userHolder.Amount.Div(decimal.NewFromInt(token.Supply)).Mul(decimal.NewFromInt(100))
		}
		tokenHolderList = append(tokenHolderList, userHolder)
	}
	if orderField == "marketcap" {
		ret := model.UserTokenHolderListMC(tokenHolderList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else if orderField == "volume" {
		ret := model.UserTokenHolderListVolume(tokenHolderList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else if orderField == "progress" {
		ret := model.UserTokenHolderListProgress(tokenHolderList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else {
		ret := model.UserTokenHolderList(tokenHolderList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(tokenHolderList, total), http.StatusOK)
	}
}

func webUserFollowList(w http.ResponseWriter, r *http.Request) {
	tag := WebParams(r).Get("tag")
	address := WebParams(r).Get("address")
	orderField := WebParams(r).Get("orderField")

	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}
	user, err := logic.GetUser(userID, true)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	symbol, tokenAddress := "", ""
	if tag != "" {
		tag = tags[tag]
	}
	if orderField != "" {
		orderField = orderFields[orderField]
	}
	if address != "" {
		if !strings.HasPrefix(address, "0x") {
			symbol = address
		} else {
			tokenAddress = address
		}
	}

	tokenList := make([]*model.Token, 0)
	for _, tokenID := range user.FollowList {
		token, err := store.DBGet[model.Token](model.C_Token, bson.M{"id": tokenID})
		if err != nil {
			slog.Error("get token for followList failed", "err", err.Error())
			continue
		}
		if tag != "" && token.Tag != tag {
			continue
		}
		if symbol != "" && token.Symbol != symbol {
			continue
		}
		if tokenAddress != "" && strings.ToUpper(token.Address) != strings.ToUpper(tokenAddress) {
			continue
		}
		token.Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), token.ID)
		tokenList = append(tokenList, token)
	}

	total := len(tokenList)
	if orderField == "marketcap" {
		ret := model.SortMarketCaps(tokenList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else if orderField == "volume" {
		ret := model.SortVolumes(tokenList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else if orderField == "progress" {
		ret := model.SortProgress(tokenList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
	} else {
		WebResponseJson(w, r, ApiResponseList(tokenList, total), http.StatusOK)
	}
}

func webUserStatistics(w http.ResponseWriter, r *http.Request) {
	address := WebParams(r).Get("address")
	totalCreated := store.DBCount(model.C_Token, bson.M{"creator": address})
	// totalTrading := store.DBCount(model.C_Trade, bson.M{"wallet": address})
	totalOwned := store.DBCount(model.C_Token_holder, bson.M{"wallet": address, "amount": bson.M{"$gt": logic.C_Min_HolderAmount}})

	totalTrading := decimal.Zero
	tradeList, err := logic.GetModelList[model.Trade](model.C_Trade, bson.M{"wallet": address})
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	for _, trade := range tradeList {
		totalTrading = totalTrading.Add(trade.QuoteAmount)
	}
	quotePrice := monitor.GetQuotePrice()
	totalTrading = totalTrading.Mul(quotePrice)
	global.Debug("get total trading for user successed", "quotePrice", quotePrice, "totalTrading", totalTrading)

	ret := &model.UserStatistics{
		TotalOwned:   totalOwned,
		TotalCreated: totalCreated,
		TotalTrading: fmt.Sprintf("$%s", totalTrading.StringFixed(3)),
	}
	WebResponseJson(w, r, ApiResponse(ret, true), http.StatusOK)
}

func webUserReferrals(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	order := WebParams(r).Get("order")
	limit := WebParams(r).Get("limit")
	offset := WebParams(r).Get("offset")
	orderField := WebParams(r).Get("orderField")

	iLimit, _ := strconv.Atoi(limit)
	iOffset, _ := strconv.Atoi(offset)

	userAddress := fmt.Sprintf("%v", userID)
	referList, total, err := logic.GetModelListPageEx[model.Refferal](model.C_Referral, bson.M{"wallet": userAddress}, order, orderField, iOffset, iLimit)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponseList(referList, total), http.StatusOK)
}
