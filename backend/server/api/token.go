package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"meme3/global"
	"meme3/server/web"
	"meme3/service/logic"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/types"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	imageTypes = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
)

var (
	tags = map[string]string{
		"0": "",
		"1": "Meme",
		"2": "AI",
		"3": "Game",
	}
	orderFields = map[string]string{
		"0": "createdAt",
		"1": "volume",
		"2": "progress",
		"3": "marketcap",
	}
	advanceOrderFields = map[string]string{
		"0": "marketcap",
		"1": "volume",
		"2": "holders",
	}
)

func initToken(api *swag.API) {
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/datas/kline",
			endpoint.Handler(webKLines),
			endpoint.Tags("Token"),
			endpoint.Summary("token list"),
			endpoint.Description("get token list"),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.TokenOHLCV{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/token/create",
			endpoint.Tags("Token"),
			endpoint.Handler(webTokenCreate),
			endpoint.Summary("Token Create"),
			endpoint.Description("Token Create"),
			endpoint.Body(model.Token{}, "Help object that needs to be added to the store", true),
			endpoint.Response(http.StatusOK, "Successfully added help", endpoint.SchemaResponseOption(model.Token{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/token/create/ok",
			endpoint.Tags("Token"),
			endpoint.Handler(webTokenCreateOK),
			endpoint.Summary("Token Create OK"),
			endpoint.Description("Token Create OK"),
			endpoint.Query("tokenID", "string", "token ID", true),
			endpoint.Query("txhash", "string", "transaction hash", true),
			endpoint.Response(http.StatusOK, "Successfully added help", endpoint.SchemaResponseOption(model.Token{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/token/image",
			endpoint.Tags("Token"),
			endpoint.Handler(webTokenImage),
			endpoint.Summary("Token image upload"),
			endpoint.Description("Token image upload"),
			endpoint.Query("tokenID", "string", "token ID", false),
			endpoint.Response(http.StatusOK, "Successfully added help", endpoint.SchemaResponseOption("imageUrl")),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/token/image/upload",
			endpoint.Tags("Token"),
			endpoint.Handler(webTokenImageUpload),
			endpoint.Summary("Token image upload"),
			endpoint.Description("Token image upload"),
			endpoint.FormData("file", types.File, "upload token image", true),
			endpoint.Response(http.StatusOK, "Successfully added help", endpoint.SchemaResponseOption("imageUrl")),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/token/aggregate/list",
			endpoint.Handler(webAggregateTokens),
			endpoint.Tags("Token"),
			endpoint.Summary("token list"),
			endpoint.Description("get token list"),
			endpoint.Query("tag", "string", "token tag", false),
			endpoint.Query("creator", "string", "token creator", false),
			endpoint.Query("address", "string", "token creator", false),
			endpoint.Query("progress", "string", "token progress", false),
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
			http.MethodGet, "/token/list",
			endpoint.Handler(webTokens),
			endpoint.Tags("Token"),
			endpoint.Summary("token list"),
			endpoint.Description("get token list"),
			endpoint.Query("tag", "string", "token tag", false),
			endpoint.Query("creator", "string", "token creator", false),
			endpoint.Query("address", "string", "token creator", false),
			endpoint.Query("progress", "string", "token progress", false),
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
			http.MethodGet, "/token/search",
			endpoint.Handler(webTokensAdvance),
			endpoint.Tags("Token"),
			endpoint.Summary("token search"),
			endpoint.Description("get token list"),
			endpoint.Query("rangeStart", "string", "range start", false),
			endpoint.Query("rangeEnd", "string", "range start", false),
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
			http.MethodGet, "/token/detail",
			endpoint.Tags("Token"),
			endpoint.Handler(webToken),
			endpoint.Summary("token information"),
			endpoint.Description("get token information"),
			endpoint.Query("address", "string", "token address", true),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption(model.Token{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	// api.AddEndpoint(
	// 	endpoint.New(
	// 		http.MethodGet, "/token/subscribe",
	// 		endpoint.Tags("Token"),
	// 		endpoint.Handler(webTokenSubscribe),
	// 		endpoint.Summary("token subscribe"),
	// 		endpoint.Description("token subscribe"),
	// 		endpoint.Query("address", "string", "token address", true),
	// 		endpoint.Query("topic", "string", "subscribe topic", true),
	// 		endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption(model.Token{})),
	// 		// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
	// 	),
	// )
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/token/holders",
			endpoint.Tags("Token"),
			endpoint.Handler(webTokenHolders),
			endpoint.Summary("token holders"),
			endpoint.Description("get token holders"),
			endpoint.Query("address", "string", "token address", true),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.TokenHolder{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/token/price",
			endpoint.Handler(webTokenPrice),
			endpoint.Summary("token price"),
			endpoint.Tags("Token"),
			endpoint.Description("get token price"),
			endpoint.Query("address", "string", "token address", true),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption(model.TokenPrice{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/token/comment",
			endpoint.Tags("Token"),
			endpoint.Handler(webTokenComment),
			endpoint.Summary("Token comment"),
			endpoint.Description("Token comment"),
			endpoint.Body(model.TokenComment{}, "Help object that needs to be added to the store", true),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.Token{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/token/hot",
			endpoint.Tags("Token"),
			endpoint.Handler(webHotTokens),
			endpoint.Summary("hot tokens"),
			endpoint.Description("get hot tokens"),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]model.Token{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
}

func webHotTokens(w http.ResponseWriter, r *http.Request) {
	hotList := make([]*model.Token, 0)
	if system, err := store.DBGet[model.System](model.C_System, bson.M{}); system != nil && err == nil {
		for _, tokenAddress := range system.Hot {
			if token, _ := logic.GetToken(tokenAddress); token != nil {
				price, percent := logic.GetTokenPricePercent(token)
				token.MarketCap = price.Mul(decimal.NewFromInt(model.C_total_supply))
				token.Percent = percent
				hotList = append(hotList, token)
			}
		}
	}
	remain := 6 - len(hotList)
	if remain <= 0 {
		WebResponseJson(w, r, ApiResponse(hotList, true), http.StatusOK)
		return
	}

	tokenList, err := logic.GetTokenListFromCache(web.PopFromSession(r, web.C_Session_User), logic.C_tokens_sort_marketcap, 1, 6, 0, nil)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	if len(tokenList) >= remain {
		ret := append(hotList, tokenList[:remain]...)
		for i, _ := range ret {
			tokenPriceUsd, percent := logic.GetTokenPricePercent(ret[i])
			ret[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
			ret[i].Volume24 = logic.GetVolume24hForToken(ret[i], false)
			ret[i].Percent = percent
		}
		WebResponseJson(w, r, ApiResponse(ret, true), http.StatusOK)
	} else if len(tokenList) == 0 {
		WebResponseJson(w, r, ApiResponse(hotList, true), http.StatusOK)
		return
	} else {
		ret := append(hotList, tokenList...)
		WebResponseJson(w, r, ApiResponse(ret, true), http.StatusOK)
	}
}

func webKLines(w http.ResponseWriter, r *http.Request) {
	klinesBuf, _ := json.Marshal(nil)
	WebResponseAny(w, r, klinesBuf, "application/json", http.StatusOK)
}

func webTokenImage(w http.ResponseWriter, r *http.Request) {
	fileName := WebParams(r).Get("fileName")
	if _, err := os.Stat(filepath.Join(logic.ImageDir, fileName)); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filepath.Join(logic.ImageDir, fileName))
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func webTokenImageUpload(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
		return
	}
	global.Debug("webTokenImageUpload", "userID", userID, "paramsAddress", WebParams(r).Get("address"))

	r.ParseMultipartForm(32 << 20)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileBuf, err := io.ReadAll(file)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType == "" {
		mimeType = http.DetectContentType(fileBuf)
	}
	if !slices.Contains(imageTypes, mimeType) {
		WebResponseJson(w, r, ApiError("invalid image"), http.StatusInternalServerError)
		return
	}
	fileName := fmt.Sprintf("%v%s", fmt.Sprintf("%v", time.Now().UnixNano()), filepath.Ext(fileHeader.Filename))

	// tokenID := fmt.Sprintf("%v", time.Now().UnixNano())
	// fileName := fmt.Sprintf("%v%s", tokenID, filepath.Ext(fileHeader.Filename))
	// slog.Info("get image successed", "mimeType", mimeType, "fileName", fileName)

	// token := &model.Token{
	// 	ID:            tokenID,
	// 	ImageFileBuf:  fileBuf,
	// 	ImageFileName: fileName,
	// }
	// store.CacheSetByTime(fmt.Sprintf("wait_create_%v", userID), token, true, time.Duration(5)*time.Minute, nil)
	// WebResponseJson(w, r, ApiResponse(fmt.Sprintf("%s/images/%s/%s", global.Config.HostURL, global.Config.RPC.Chain, fileName), true), http.StatusOK)
	// // WebResponseJson(w, r, ApiResponse(fmt.Sprintf("%s/api/v1/token/image?fileName=%s", global.Config.HostURL, fileName), true), http.StatusOK)

	fullPath := filepath.Join(filepath.Join(global.Config.WebPath, "images", global.Config.RPC.Chain), fileName)
	f, err := os.Create(fullPath)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = f.Write(fileBuf); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponse(fmt.Sprintf("%s/images/%s/%s", global.Config.HostURL, global.Config.RPC.Chain, fileName), true), http.StatusOK)
}

func webTokenCreate(w http.ResponseWriter, r *http.Request) {
	_, requestToken, err := WebBody[model.Token](r)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusServiceUnavailable)
		return
	}
	global.Debug("webTokenCreate", "vid", WebParams(r).Get("vid"))

	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
		return
	}
	global.Debug("webTokenCreate", "userID", userID, "paramsAddress", WebParams(r).Get("address"), "vid", WebParams(r).Get("vid"))

	if token, err := logic.CreateToken(userID, requestToken, WebParams(r).Get("vid")); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusServiceUnavailable)
		return
	} else {
		WebResponseJson(w, r, ApiResponse(token, true), http.StatusOK)
	}
}

func webTokenCreateOK(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
		return
	}
	txhash := WebParams(r).Get("txhash")
	tokenID := WebParams(r).Get("tokenID")

	var blockNumber int64
	if txhash != "" {
		timeout := 300
		timer := 1
		for {
			if timer >= timeout {
				WebResponseJson(w, r, ApiResponse("successed", true), http.StatusOK)
				return
			}
			if bOK, number := monitor.IsTxOK(txhash); bOK {
				blockNumber = number
				break
			}
			time.Sleep(time.Duration(1) * time.Second)
			timer += 1
		}
	} else {
		WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
		return
	}

	tid, _ := decimal.NewFromString(tokenID)
	tokenCreated := &monitor.TokenCreated{
		TokenId:     tid,
		Time:        time.Now(),
		BlockNumber: blockNumber,
		Creator:     common.HexToAddress(fmt.Sprintf("%v", userID)),
		TotalSupply: decimal.NewFromInt(1000000000),
	}

	if bHandle := logic.UpdateTokenCreated(tokenCreated); !bHandle {
		WebResponseJson(w, r, ApiError("create token failed with invalid params"), http.StatusServiceUnavailable)
		return
	}
	WebResponseJson(w, r, ApiResponse("successed", true), http.StatusOK)
}

func webToken(w http.ResponseWriter, r *http.Request) {
	tokenAddress := WebParams(r).Get("address")
	global.Debug("get token", tokenAddress)

	token, err := logic.GetToken(tokenAddress)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}

	if token != nil {
		token.RaisedAmount, _ = logic.GetProgress(token.TotalAmount)
		price, percent := logic.GetTokenPricePercent(token)
		token.MarketCap = price.Mul(decimal.NewFromInt(model.C_total_supply))
		token.Volume24 = logic.GetVolume24hForToken(token, false)
		token.LaunchMC = global.Config.LaunchMC
		token.Percent = percent
	}
	model.SortTokenComments(token.Comments)
	WebResponseJson(w, r, ApiResponse(token, true), http.StatusOK)
}

func webTokenHolders(w http.ResponseWriter, r *http.Request) {
	tokenAddress := WebParams(r).Get("address")

	token, err := logic.GetToken(tokenAddress)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}

	// // 池子流动性
	// tokenTrading := monitor.GetReverse(tokenAddress, global.Config.ContractAddress, token.Launch)

	// percent := decimal.NewFromInt(100)
	// if tokenTrading != nil {
	// 	percent = decimal.NewFromInt(token.Supply).Sub(tokenTrading.TokenAmount).Div(decimal.NewFromInt(token.Supply)).Mul(decimal.NewFromInt(100))
	// 	global.Debug("set pool liquidity", "tokenTrading", tokenTrading, "percent", percent)
	// 	if token.Launch {
	// 		percent = decimal.NewFromInt(20)
	// 	}
	// }
	devHolders := make([]*model.TokenHolder, 0)
	devHolders = append(devHolders, &model.TokenHolder{
		Address: tokenAddress,
		Wallet:  global.Config.ContractAddress,
		Amount:  token.TokenAmount,
		Percent: decimal.NewFromInt(20),
	})

	whereModel := bson.M{"address": tokenAddress, "amount": bson.M{"$gte": logic.C_Min_HolderAmount}}
	tokenHolders, _, err := logic.GetModelListPageEx[model.TokenHolder](model.C_Token_holder, whereModel, "", "percent", 0, 10)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}

	ret := make([]*model.TokenHolder, 0)
	if !token.Launch {
		ret = append(ret, devHolders...)
		ret = append(ret, tokenHolders...)
	} else {
		devHolder := &model.TokenHolder{
			Address: token.Address,
			Wallet:  token.Creator,
			Amount:  monitor.GetTokenBalance(token.Address, token.Creator),
		}
		if devHolder.Amount.GreaterThan(decimal.NewFromFloat(logic.C_Min_HolderAmount)) {
			devHolder.Percent = devHolder.Amount.Div(decimal.NewFromInt(monitor.C_supply)).Mul(decimal.NewFromInt(100))
			devHolders = append(devHolders, devHolder)
		}

		for _, holder := range tokenHolders {
			if strings.ToUpper(holder.Wallet) != strings.ToUpper(token.Creator) {
				holder.Amount = monitor.GetTokenBalance(holder.Address, holder.Wallet)
				holder.Percent = holder.Amount.Div(decimal.NewFromInt(monitor.C_supply)).Mul(decimal.NewFromInt(100))
				ret = append(ret, holder)
			}
		}
		ret = append(devHolders, model.SortTokenHolders(ret)...)
	}
	WebResponseJson(w, r, ApiResponse(ret, true), http.StatusOK)
}

func webTokenPrice(w http.ResponseWriter, r *http.Request) {
	tokenAddress := WebParams(r).Get("address")

	token, err := logic.GetToken(tokenAddress)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}
	tokenPrice := logic.GetTokenPrice(tokenAddress, token.Launch, false)
	tokenTrading := monitor.GetReverse(tokenAddress, global.Config.ContractAddress, token.Launch)

	ret := struct {
		model.TokenPrice
		Percent     float64         `json:"percent"`
		Progress    float64         `json:"prograss"` // 写错了  progress
		TotalAmount decimal.Decimal `json:"totalAmount"`
	}{}
	if token != nil {
		progress := 100.0
		if !token.Launch {
			ret.TotalAmount = tokenTrading.TotalAmount
			_, progress = logic.GetProgress(tokenTrading.TotalAmount)
		} else {
			ret.TotalAmount = decimal.NewFromFloat(global.Config.LaunchMC)
		}
		_, percent := logic.GetTokenPricePercent(token)
		ret.Progress = progress
		ret.Percent = percent
	}
	ret.TokenPrice = *tokenPrice
	WebResponseJson(w, r, ApiResponse(ret, true), http.StatusOK)
}

func webTokenComment(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
		return
	}

	_, comment, err := WebBody[model.TokenComment](r)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	comment.Creator = fmt.Sprintf("%v", userID)

	if token, err := logic.AddComment(comment); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	} else {
		WebResponseJson(w, r, ApiResponse(token, true), http.StatusOK)
	}
}

func webAggregateTokens(w http.ResponseWriter, r *http.Request) {
	tag := WebParams(r).Get("tag")
	page := WebParams(r).Get("page")
	// order := WebParams(r).Get("order")
	limit := WebParams(r).Get("limit")
	offset := WebParams(r).Get("offset")
	creator := WebParams(r).Get("creator")
	address := WebParams(r).Get("address")
	progress := WebParams(r).Get("progress")
	// insurance := WebParams(r).Get("insurance")
	orderField := WebParams(r).Get("orderField")

	iLimit, _ := strconv.Atoi(limit)
	iOffset, _ := strconv.Atoi(offset)
	totalTokens := int(store.DBCount(model.C_Token, bson.M{}))

	// // return cache
	global.Debug("check return cache or not", "page", page)
	bCanCache := orderField == "0" || orderField == "1" || orderField == "3"
	bDefaultCache := iOffset <= 5 && progress == "" && address == "" && creator == "" //&& insurance == ""
	if orderField == "1" || orderField == "3" {
		bCanCache = page == "ranking"
	}
	if page != "ranking" && tag != "0" {
		bDefaultCache = false
	}
	if bCanCache && bDefaultCache {
		tokenList, _, err := logic.GetTokenListForAggregate(web.PopFromSession(r, web.C_Session_User), orderFields[orderField], iOffset, iLimit)
		global.Debug("after GetTokenListFromCache", "params", WebParams(r), "tokenList", tokenList, "cacheErr", err)
		if tokenList != nil && len(tokenList) > 0 {
			if page == "ranking" && len(tokenList) > logic.C_cache_ranking_total {
				tokenList = tokenList[:logic.C_cache_ranking_total]
			}
			WebResponseJson(w, r, ApiResponseList(tokenList, totalTokens), http.StatusOK)
			return
		}
	}
	WebResponseJson(w, r, ApiResponseList([]*model.Token{}, 0), http.StatusOK)
}

func webTokens(w http.ResponseWriter, r *http.Request) {
	tag := WebParams(r).Get("tag")
	page := WebParams(r).Get("page")
	order := WebParams(r).Get("order")
	limit := WebParams(r).Get("limit")
	offset := WebParams(r).Get("offset")
	creator := WebParams(r).Get("creator")
	address := WebParams(r).Get("address")
	progress := WebParams(r).Get("progress")
	// insurance := WebParams(r).Get("insurance")
	orderField := WebParams(r).Get("orderField")

	iLimit, _ := strconv.Atoi(limit)
	iOffset, _ := strconv.Atoi(offset)

	// // return cache
	global.Debug("check return cache or not", "page", page)
	totalTokens := int(store.DBCount(model.C_Token, bson.M{}))
	bCanCache := orderField == "0" || orderField == "1" || orderField == "3"
	bDefaultCache := iOffset <= 5 && progress == "" && address == "" && creator == "" //&& insurance == ""
	if orderField == "1" || orderField == "3" {
		bCanCache = page == "ranking"
	}
	if page != "ranking" && tag != "0" {
		bDefaultCache = false
	}
	if bCanCache && bDefaultCache {
		fetchTokens, fetchTotal := getTokensForAllChains(w, r)
		tokenList, err := logic.GetTokenListFromCache(web.PopFromSession(r, web.C_Session_User), orderFields[orderField], iOffset, iLimit, totalTokens+fetchTotal, fetchTokens)
		global.Debug("after GetTokenListFromCache", "params", WebParams(r), "tokenList", tokenList, "cacheErr", err)
		if tokenList != nil && len(tokenList) > 0 {
			if page == "ranking" && len(tokenList) > logic.C_cache_ranking_total {
				tokenList = tokenList[:logic.C_cache_ranking_total]
			}
			if totalTokens+fetchTotal > logic.C_cache_tokens_total && WebParams(r).Get("Chainid") != "" {
				WebResponseJson(w, r, ApiResponseList(tokenList, logic.C_cache_tokens_total), http.StatusOK)
			} else {
				WebResponseJson(w, r, ApiResponseList(tokenList, totalTokens+fetchTotal), http.StatusOK)
			}
			return
		}
	} else if iOffset > 5 {
		// 聚合数据限定为最多5页, 不再需要重置页码
		pageLast := totalTokens / iLimit
		totalLastPage := totalTokens % iLimit
		if totalLastPage > 0 {
			pageLast += 1
		}
		iOffset = pageLast
		global.Debug("reset offset successed", "params", WebParams(r), "iOffset", iOffset, "totalTokens", totalTokens)
	}

	whereModel := bson.D{}
	if progress == "launch" {
		whereModel = append(whereModel, bson.E{"launch", true})
	}
	// if insurance == "true" {
	// 	whereModel = append(whereModel, bson.E{"insurance", true})
	// }
	if tag != "" {
		tag = tags[tag]
		if tag != "" {
			// whereModel["tag"] = tag
			whereModel = append(whereModel, bson.E{"tag", tag})
		}
	}
	if orderField != "" {
		orderField = orderFields[orderField]
	}

	if creator != "" {
		whereModel = append(whereModel, bson.E{"creator", creator})
	}
	if address != "" {
		if !strings.HasPrefix(address, "0x") {
			whereModel = append(whereModel, bson.E{"symbol", address})
		} else {
			whereModel = append(whereModel, bson.E{"address", address})
		}
	}

	if order == "" {
		order = "desc"
	}
	if orderField == "" {
		orderField = "createdAt"
	}

	var err error
	var total int64
	var tokenList []*model.Token
	tokenList, total, err = store.DBListPage[model.Token](model.C_Token, whereModel, order, orderField, iOffset, iLimit, false)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	}

	for i, token := range tokenList {
		tokenList[i].Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), tokenList[i].ID)
		tokenPriceUsd, percent := logic.GetTokenPricePercent(token)
		tokenList[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
		tokenList[i].Volume24 = logic.GetVolume24hForToken(tokenList[i], false)
		tokenList[i].Percent = percent
	}
	if orderField == "marketcap" {
		ret := model.TokenMarketCapList(tokenList)
		sort.Sort(ret)
		WebResponseJson(w, r, ApiResponseList(ret, int(total)), http.StatusOK)
	} else {
		WebResponseJson(w, r, ApiResponseList(tokenList, int(total)), http.StatusOK)
	}
}

func webTokensAdvance(w http.ResponseWriter, r *http.Request) {
	rangeParams := WebParams(r).Get("rangeParams")
	order := WebParams(r).Get("order")
	limit := WebParams(r).Get("limit")
	launch := WebParams(r).Get("launch")
	offset := WebParams(r).Get("offset")
	orderField := WebParams(r).Get("orderField")
	if launch == "true" {
		orderField = "launch"
	}

	// 单页100条
	limit = "100"
	offset = "1"

	rangeParamsMap := make(map[string]map[string]string, 0)
	if paramsBuf, err := hex.DecodeString(rangeParams); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
		return
	} else if jerr := json.Unmarshal(paramsBuf, &rangeParamsMap); jerr != nil {
		WebResponseJson(w, r, ApiError(jerr.Error()), http.StatusOK)
		return
	}
	slog.Info("get rangeParamsMap successed", "rangeParamsMap", rangeParamsMap)

	var total int
	iLimit, _ := strconv.Atoi(limit)
	iOffset, _ := strconv.Atoi(offset)
	tokenList := make([]*model.Token, 0)
	if launch == "80" {
		seachTokenList, stotal, err := logic.GetAboutLaunchList(launch, order, orderField, rangeParamsMap, iOffset, iLimit)
		if err != nil {
			WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
			return
		}
		total = stotal
		tokenList = seachTokenList[:]
	} else {
		seachTokenList, err := logic.SearchTokenList(orderField, rangeParamsMap)
		if err != nil {
			WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
			return
		}
		tokenList, err = logic.GetPageList(seachTokenList, iOffset, iLimit, 0)
		if err != nil {
			WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
			return
		}
		total = len(seachTokenList)
		global.Debug("after GetTokenListFromCache", "tokenList", tokenList, "orderField", orderField)
	}

	// output
	total = 0
	wlist := make([]*model.WebToken, len(tokenList))
	for i, token := range tokenList {
		if launch == "true" {
			wlist[i] = &model.WebToken{
				Token:      *(tokenList[i]),
				TokenPrice: *(logic.GetTokenPrice(token.Address, true, false)),
			}
			wlist[i].Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), wlist[i].ID)
		} else {
			tokenList[i].Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), tokenList[i].ID)
		}
	}
	var retList any
	if launch == "true" {
		retList = wlist[:]
	} else {
		retList = tokenList[:]
	}
	WebResponseJson(w, r, ApiResponseList(retList, total), http.StatusOK)

	// bMC := false
	// bDefaultParams := true
	// whereModel := bson.M{}
	// var vistart, viend float64
	// for k, vMap := range rangeParamsMap {
	// 	workK := k
	// 	iStart, _ := strconv.ParseFloat(vMap["rangeStart"], 64)
	// 	iEnd, _ := strconv.ParseFloat(vMap["rangeEnd"], 64)
	// 	if workK == "marketCap" {
	// 		workK = "marketcap"
	// 		iStart *= 1000
	// 		iEnd *= 1000
	// 		bMC = true
	// 		if iStart >= 0 || iEnd > 0 {
	// 			vistart, viend = iStart, iEnd
	// 		}
	// 		iStart = 0
	// 		iEnd = 0
	// 	} else if workK == "volume" {
	// 		iStart *= 1000
	// 		iEnd *= 1000
	// 	}

	// 	if iStart > 0 && iEnd > 0 {
	// 		bDefaultParams = false
	// 		whereModel[workK] = bson.M{"$gte": iStart, "$lte": iEnd}
	// 	} else if iStart > 0 {
	// 		bDefaultParams = false
	// 		whereModel[workK] = bson.M{"$gte": iStart}
	// 	} else if iEnd > 0 {
	// 		bDefaultParams = false
	// 		whereModel[workK] = bson.M{"$lte": iEnd}
	// 	}
	// }
	// iLimit, _ := strconv.Atoi(limit)
	// iOffset, _ := strconv.Atoi(offset)
	// if launch == "true" {
	// 	whereModel["launch"] = true
	// 	// orderField = "launchTime"
	// } else if launch == "80" {
	// 	whereModel["launch"] = false
	// 	whereModel["progress"] = bson.M{"$gte": 80}
	// }

	// // return cache
	// bDefaultCache := offset == "1" && launch == ""
	// if bDefaultParams && bDefaultCache {
	// 	tokenList, _ := logic.GetTokenListFromCache(web.PopFromSession(r, web.C_Session_User), orderField, iOffset, iLimit)
	// 	global.Debug("after GetTokenListFromCache", "tokenList", tokenList, "orderField", orderField)
	// 	if tokenList != nil && len(tokenList) > 0 {
	// 		if len(tokenList) >= 21 {
	// 			tokenList = tokenList[:20]
	// 		}
	// 		WebResponseJson(w, r, ApiResponseList(tokenList, len(tokenList)), http.StatusOK)
	// 		// WebResponseJson(w, r, ApiResponseList(tokenList, int(store.DBCount(model.C_Token, bson.M{}))), http.StatusOK)
	// 		return
	// 	}
	// }

	// var err error
	// var total int
	// var tokenList []*model.Token
	// tokenList, total, err = logic.GetModelListPageEx[model.Token](model.C_Token, whereModel, order, orderField, iOffset, iLimit)
	// if err != nil {
	// 	WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
	// 	return
	// }

	// ret := make([]*model.Token, 0)
	// if bMC && vistart >= 0 && viend > 0 {
	// 	global.Debug("check marketcap", "viStart", vistart, "viEnd", viend, "rangeParamsMap", rangeParamsMap, "whereModel", whereModel)
	// 	for _, token := range tokenList {
	// 		tokenPrice, percent := logic.GetTokenPricePercent(token)
	// 		token.MarketCap = tokenPrice.Mul(decimal.NewFromInt(model.C_total_supply))
	// 		token.Percent = percent
	// 		token.Volume24 = logic.GetVolume24hForToken(token, false)
	// 		token.Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), token.ID)
	// 		if token.MarketCap.LessThan(decimal.NewFromFloat(vistart)) || token.MarketCap.GreaterThan(decimal.NewFromFloat(viend)) {
	// 			continue
	// 		}
	// 		ret = append(ret, token)
	// 	}
	// } else {
	// 	for _, token := range tokenList {
	// 		token.Follow = logic.IsFavoriteToken(web.PopFromSession(r, web.C_Session_User), token.ID)
	// 		tokenPrice, percent := logic.GetTokenPricePercent(token)
	// 		token.MarketCap = tokenPrice.Mul(decimal.NewFromInt(model.C_total_supply))
	// 		token.Volume24 = logic.GetVolume24hForToken(token, false)
	// 		token.Percent = percent
	// 		ret = append(ret, token)
	// 	}
	// }
	// WebResponseJson(w, r, ApiResponseList(ret, total), http.StatusOK)
}

func getTokensForAllChains(w http.ResponseWriter, r *http.Request) ([]*model.Token, int) {
	total := 0
	tokens := make([]*model.Token, 0)
	if WebParams(r).Get("page") == "ranking" {
		return tokens, total
	}
	// return tokens, total

	chains := strings.Split(WebParams(r).Get("Chainid"), ",")
	for _, chain := range chains {
		chain = strings.TrimSpace(chain)
		if !strings.EqualFold(chain, global.Config.RPC.Chain) && strings.Contains(global.Config.RPC.ValidChains, chain) {
			fetchUrl := fmt.Sprintf("http://%s%s", global.Config.RPC.ChainAddrs[chain], r.RequestURI)
			fetchUrl = strings.ReplaceAll(fetchUrl, global.Config.RPC.Chain, chain)
			fetchUrl = strings.ReplaceAll(fetchUrl, "/token/list", "/token/aggregate/list")
			fetchBuf, err := web.HttpGet(fetchUrl)
			if err != nil {
				slog.Error("webTokensForAllChains: fetch token list failed", "request", r, "currentChain", global.Config.RPC.Chain, "chain", chain, "fetchUrl", fetchUrl)
				continue
			}
			fetchTokenList := make([]*model.Token, 0)
			tokensListBuf := gjson.Get(string(fetchBuf), "data.data").String()
			if err := json.Unmarshal([]byte(tokensListBuf), &fetchTokenList); err != nil {
				slog.Error("webTokensForAllChains: parse token list failed", "currentChain", global.Config.RPC.Chain, "chain", chain, "tokensListBuf", string(tokensListBuf))
				continue
			}
			slog.Info("fetch token list for all chains ok", "currentChain", global.Config.RPC.Chain, "chain", chain, "fetchUrl", fetchUrl, "fetchTokenList", fetchTokenList)
			total += int(gjson.Get(string(fetchBuf), "data.total").Int())
			tokens = append(tokens, fetchTokenList...)
		}
	}
	return tokens[:], total
}
