package logic

import (
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/store"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	c_page_max     = 5
	c_page_limit   = 21
	c_page_timeout = time.Duration(5) * time.Minute
)

var (
	// {pageIndex: []*model.Token}
	v_page_tokens = store.NewCache()
)

type PageTokens struct {
	PageIndex  int
	TokenIndex int
	TokenList  []*model.Token
}

func InitPageTokens() (err error) {
	for i := 1; i <= c_page_max; i++ {
		if tokenList, err := doInitPageTokens(i, true); err != nil {
			return err
		} else {
			SavePageTokens(i, tokenList, true)
		}
	}
	return nil
}

func doInitPageTokens(offset int, bInit bool) (ret []*model.Token, err error) {
	tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "createdAt", offset, c_page_limit)
	if err != nil {
		return nil, err
	}
	global.Debug("before init page tokens", "bInit", bInit, "offset", offset, "tokenList", len(tokenList))

	return tokenList, nil
}

func SavePageTokens(index int, tokenList []*model.Token, bInit bool) {
	for i, _ := range tokenList {
		tokenList[i].Volume24 = GetVolume24hForToken(tokenList[i], false)
		tokenPriceUsd, percent := GetTokenPricePercent(tokenList[i])
		tokenList[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
		tokenList[i].Percent = percent
	}
	// v_page_tokens.Store(pageIndex, tokenList)

	if bInit {
		store.CachePush(v_page_tokens, index, tokenList, true)
	} else {
		store.CachePushByTime(v_page_tokens, index, tokenList, true, c_page_timeout, func(key, val any) bool {
			global.Debug("before remove page tokens", "key", key)
			return true
		})
	}
}

func GetPageTokens(userID any, offset, limit int) (ret []*model.Token, err error) {
	pageIndex := offset % c_page_max
	if pageIndex == 0 {
		pageIndex = c_page_max
	}
	if limit != c_page_limit {
		err = fmt.Errorf("GetPageTokens failed with invalid limit, offset: %d, limit: %d", offset, limit)
	} else if pageIndex < 0 || pageIndex > c_page_max {
		err = fmt.Errorf("GetPageTokens failed without page cache, offset: %d, limit: %d", offset, limit)
	} else {
		if pageTokens, _ := v_page_tokens.Load(offset); pageTokens != nil {
			ret = pageTokens.([]*model.Token)
		} else {
			if tokenList, err := doInitPageTokens(offset, false); err != nil {
				err = fmt.Errorf("GetPageTokens failed with init page cache, pageIndex: %d, offset: %d, limit: %d", pageIndex, offset, limit)
			} else {
				return tokenList[:], nil
			}
			// err = fmt.Errorf("GetPageTokens failed with invalid page cache, pageIndex: %d, offset: %d, limit: %d", pageIndex, offset, limit)
		}
	}
	if err != nil {
		return nil, err
	}

	for i, _ := range ret {
		if userID != nil {
			ret[i].Follow = IsFavoriteToken(userID, ret[i].ID)
		}
	}

	go SyncPageTokens(offset)

	return
}

func SyncPageTokens(offset int) {
	total := store.DBCount(model.C_Token, bson.M{})
	pageTotal := total / c_page_limit
	if total%c_page_limit > 0 {
		pageTotal += 1
	}

	// prev1
	offsetPrev1 := offset - 1
	if _, ok := v_page_tokens.Load(offsetPrev1); !ok && offsetPrev1 >= 1 {
		tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "createdAt", offsetPrev1, c_page_limit)
		if err != nil {
			slog.Error("UpdatePageTokens failed", "offsetPrev1", offsetPrev1, "offset", offset)
			return
		}
		global.Debug("before update page tokens", "offset", offset, "offsetPrev1", offsetPrev1, "tokenList", len(tokenList))
		SavePageTokens(offsetPrev1, tokenList, false)
	}

	// prev2
	offsetPrev2 := offset - 2
	if _, ok := v_page_tokens.Load(offsetPrev2); !ok && offsetPrev2 >= 1 {
		tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "createdAt", offsetPrev2, c_page_limit)
		if err != nil {
			slog.Error("UpdatePageTokens failed", "offsetPrev2", offsetPrev2, "offset", offset)
			return
		}
		global.Debug("before update page tokens", "offset", offset, "offsetPrev2", offsetPrev2, "tokenList", len(tokenList))
		SavePageTokens(offsetPrev2, tokenList, false)
	}

	// next1
	offsetNext1 := offset + 1
	if _, ok := v_page_tokens.Load(offsetNext1); !ok && offsetNext1 <= int(pageTotal) {
		tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "createdAt", offsetNext1, c_page_limit)
		if err != nil {
			slog.Error("UpdatePageTokens failed", "offsetNext1", offsetNext1, "offset", offset)
			return
		}
		global.Debug("before update page tokens", "offset", offset, "offsetNext1", offsetNext1, "tokenList", len(tokenList))
		SavePageTokens(offsetNext1, tokenList, false)
	}

	// next2
	offsetNext2 := offset + 2
	if _, ok := v_page_tokens.Load(offsetNext2); !ok && offsetNext2 <= int(pageTotal) {
		tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "createdAt", offsetNext2, c_page_limit)
		if err != nil {
			slog.Error("UpdatePageTokens failed", "offsetNext2", offsetNext2, "offset", offset)
			return
		}
		global.Debug("before update page tokens", "offset", offset, "offsetNext2", offsetNext2, "tokenList", len(tokenList))
		SavePageTokens(offsetNext2, tokenList, false)
	}
}

// // // return cache
// // if handleCacheListForIndex(w, r, orderField, progress, address, creator, offset, limit) {
// // 	return
// // }
// // if handleCacheListForRanking(w, r, orderField, page) {
// // 	return
// // }

// func handleCacheListForIndex(w http.ResponseWriter, r *http.Request, orderField, progress, address, creator, offset, limit string) bool {
// 	bIndex := orderField == "0"
// 	bDefaultCache := progress == "" && address == "" && creator == ""
// 	// // return cache
// 	// bRankingCache := false
// 	// global.Debug("check return cache or not", "page", page)
// 	// bCanCache := orderField == "0" || orderField == "1" || orderField == "3"
// 	// // bDefaultCache := offset == "1" && progress == "" && address == "" && creator == ""
// 	// bDefaultCache := progress == "" && address == "" && creator == ""
// 	// if orderField == "1" || orderField == "3" {
// 	// 	bRankingCache = page == "ranking"
// 	// }
// 	// if page != "ranking" && tag != "0" {
// 	// 	bDefaultCache = false
// 	// }
// 	if bIndex && bDefaultCache {
// 		iLimit, _ := strconv.Atoi(limit)
// 		iOffset, _ := strconv.Atoi(offset)
// 		tokenList, err := logic.GetPageTokens(web.PopFromSession(r, web.C_Session_User), iOffset, iLimit)
// 		global.Debug("after get page tokens", "params", WebParams(r), "tokenList", len(tokenList), "cacheErr", err)
// 		if tokenList != nil && len(tokenList) > 0 {
// 			WebResponseJson(w, r, ApiResponseList(tokenList, int(store.DBCount(model.C_Token, bson.M{}))), http.StatusOK)
// 			return true
// 		}
// 	}
// 	return false
// }

// func handleCacheListForRanking(w http.ResponseWriter, r *http.Request, orderField, page string) bool {
// 	// return cache
// 	bRanking := page == "ranking"

// 	if bRanking {
// 		tokenList, err := logic.GetTokenListFromCache(web.PopFromSession(r, web.C_Session_User), orderFields[orderField])
// 		global.Debug("after GetTokenListFromCache", "params", WebParams(r), "tokenList", len(tokenList), "cacheErr", err)
// 		if tokenList != nil && len(tokenList) > 0 {
// 			WebResponseJson(w, r, ApiResponseList(tokenList, int(store.DBCount(model.C_Token, bson.M{}))), http.StatusOK)
// 			return true
// 		}
// 	}
// 	return false
// }
