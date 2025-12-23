package logic

import (
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/store"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	C_cache_launch_max = 100
)

var (
	V_cache_page_launch []*model.Token
	v_lock_page_launch  = new(sync.Mutex)
)

func GetPageList(tokenList []*model.Token, offset, limit, total int) ([]*model.Token, error) {
	if offset == 0 || limit == 0 {
		return tokenList[:], nil
	}
	total = len(tokenList)
	slog.Debug("before page tokens", "tokenList", tokenList)

	pageNO := offset
	pageIndex := (offset - 1) * limit
	totalPages := total / limit
	if pageNO > totalPages && total%limit > 0 {
		limit = total % limit
	}

	if pageIndex < total && (pageIndex+limit) <= total {
		global.DebugForce("GetPageList", "pageNO", pageNO, "pageIndex", pageIndex, "totalPages", totalPages, "limit", limit, "total", total)
		return tokenList[pageIndex : pageIndex+limit], nil
	}
	return tokenList[:], nil
}

func UpdateLaunchPageList(token *model.Token) error {
	if !token.Launch {
		return nil
	}

	v_lock_page_launch.Lock()
	defer v_lock_page_launch.Unlock()
	for i, _ := range V_cache_page_launch {
		global.Debug("UpdateLaunchPageList", "token.Address", strings.ToUpper(token.Address), "cache.Token.Address", strings.ToUpper(V_cache_page_launch[i].Address))
		if strings.ToUpper(token.Address) == strings.ToUpper(V_cache_page_launch[i].Address) {
			V_cache_page_launch[i] = token
			global.Debug("UpdateLaunchPageList", "i", i, "token", V_cache_page_launch[i])
			return nil
		}
	}

	totalSize := C_cache_launch_max
	global.Debug("push token to cache", "token", token.Address)
	if len(V_cache_page_launch) == 0 {
		totalSize = 1
		V_cache_page_launch = make([]*model.Token, 0)
	} else if len(V_cache_page_launch) < C_cache_launch_max {
		totalSize = len(V_cache_page_launch) + 1
	}
	V_cache_page_launch = append([]*model.Token{token}, V_cache_page_launch[0:totalSize-1]...)
	return nil
}

func GetLaunchPageList() ([]*model.Token, error) {
	v_lock_page_launch.Lock()
	defer v_lock_page_launch.Unlock()
	whereModel := bson.M{"launch": true}
	if len(V_cache_page_launch) > 0 { //&& store.DBCount(model.C_Token, whereModel) == int64(len(V_cache_page_launch)) {
		return V_cache_page_launch, nil
	}
	V_cache_page_launch = make([]*model.Token, 0)

	tokenList, _, err := store.DBListPage[model.Token](model.C_Token, whereModel, "desc", "launchtime", 1, 100, false) //launch == "true")
	if err != nil {
		return nil, err
	}
	for i, token := range tokenList {
		tokenPrice, percent := GetTokenPricePercent(token)
		tokenList[i].MarketCap = tokenPrice.Mul(decimal.NewFromInt(model.C_total_supply))
		tokenList[i].Volume24 = GetVolume24hForToken(token, false)
		tokenList[i].Percent = percent
	}
	V_cache_page_launch = tokenList[:]
	return V_cache_page_launch[:], nil
}

func SearchTokenList(orderField string, rangeParamsMap map[string]map[string]string) ([]*model.Token, error) {
	ret := make([]*model.Token, 0)

	var err error
	var tokenList []*model.Token
	if orderField == C_tokens_sort_createdAt {
		v_lock_token.RLock()
		defer v_lock_token.RUnlock()
		tokenList = v_cache_list.Tokens[:]
	} else if orderField == C_tokens_sort_volume {
		v_lock_tokenv.RLock()
		defer v_lock_tokenv.RUnlock()
		tokenList = v_cache_list.VolumeList[:]
	} else if orderField == "launch" {
		if tokenList, err = GetLaunchPageList(); err != nil {
			return nil, err
		}
	}

	for _, token := range tokenList {
		bNeedAdd := true
		for k, vMap := range rangeParamsMap {
			iStart, _ := strconv.ParseFloat(vMap["rangeStart"], 64)
			iEnd, _ := strconv.ParseFloat(vMap["rangeEnd"], 64)
			if k == "marketCap" {
				iEnd *= 1000
				iStart *= 1000
				if bInValid := isInValidToken(token.MarketCap, iStart, iEnd); bInValid {
					bNeedAdd = false
					break
				}
			} else if k == "volume" {
				iEnd *= 1000
				iStart *= 1000
				if bInValid := isInValidToken(token.Volume, iStart, iEnd); bInValid {
					bNeedAdd = false
					break
				}
			} else if k == "holders" {
				if bInValid := isInValidToken(decimal.NewFromInt(int64(token.Holders)), iStart, iEnd); bInValid {
					bNeedAdd = false
					break
				}
			}
		}
		if bNeedAdd {
			ret = append(ret, token)
		}
	}
	return ret, nil
}

func isInValidToken(fieldV decimal.Decimal, iStart, iEnd float64) bool {
	if iStart > 0 && iEnd > 0 {
		if fieldV.GreaterThanOrEqual(decimal.NewFromFloat(iStart)) && fieldV.LessThanOrEqual(decimal.NewFromFloat(iEnd)) {
			return false
		}
	} else if iStart > 0 {
		if fieldV.GreaterThanOrEqual(decimal.NewFromFloat(iStart)) {
			return false
		}
	} else if iEnd > 0 {
		if fieldV.LessThanOrEqual(decimal.NewFromFloat(iEnd)) {
			return false
		}
	} else {
		return false
	}
	return true
}

func GetAboutLaunchList(launch string, order, orderField string, rangeParamsMap map[string]map[string]string, offset, limit int) ([]*model.Token, int, error) {
	whereModel := bson.M{}
	if launch == "true" {
		whereModel["launch"] = true
		orderField = "launchtime"
	} else if launch == "80" {
		whereModel["launch"] = false
		whereModel["progress"] = bson.M{"$gte": 80}
	}
	for k, vMap := range rangeParamsMap {
		workK := k
		iStart, _ := strconv.ParseFloat(vMap["rangeStart"], 64)
		iEnd, _ := strconv.ParseFloat(vMap["rangeEnd"], 64)
		if workK == "marketCap" {
			workK = "marketcap"
			iStart *= 1000
			iEnd *= 1000
		} else if workK == "volume" {
			iStart *= 1000
			iEnd *= 1000
		}

		if iStart > 0 && iEnd > 0 {
			whereModel[workK] = bson.D{{"$gte", iStart}, {"$lte", iEnd}}
		} else if iStart > 0 {
			whereModel[workK] = bson.M{"$gte": iStart}
		} else if iEnd > 0 {
			whereModel[workK] = bson.M{"$lte": iEnd}
		}
	}

	var err error
	var total int64
	var tokenList []*model.Token
	tokenList, total, err = store.DBListPage[model.Token](model.C_Token, whereModel, order, orderField, offset, limit, false) //launch == "true")
	if err != nil {
		return nil, 0, err
	}
	for i, token := range tokenList {
		tokenPrice, percent := GetTokenPricePercent(token)
		tokenList[i].MarketCap = tokenPrice.Mul(decimal.NewFromInt(model.C_total_supply))
		tokenList[i].Volume24 = GetVolume24hForToken(token, false)
		tokenList[i].Percent = percent
	}
	return tokenList[:], int(total), nil
}
