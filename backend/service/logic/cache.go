package logic

import (
	"context"
	"fmt"
	"io/ioutil"
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	C_cache_tokens_total  = 105
	C_cache_trades_total  = 10
	C_cache_ranking_total = 20
	// C_cache_tokens_volume_total = 20
	C_cache_tokens_volume_total = 100

	C_tokens_sort_hot       = "hot"
	C_tokens_sort_volume    = "volume"
	C_tokens_sort_volume24  = "volume24"
	C_tokens_sort_createdAt = "createdAt"
	C_tokens_sort_marketcap = "marketcap"
)

var (
	V_hot_total   = 3
	v_cache_list  = &CacheList{}
	v_lock_hot    = new(sync.RWMutex)
	v_lock_token  = new(sync.RWMutex)
	v_lock_trade  = new(sync.RWMutex)
	v_lock_tokenv = new(sync.RWMutex)
	v_lock_tokenm = new(sync.RWMutex)
	V_hot_tokens  = make([]*model.Token, 0)

	V_cache_tokens21 = new(sync.Map)
)

type CacheList struct {
	MarketCapList []*model.Token
	VolumeList    []*model.Token
	Tokens        []*model.Token
	Trades        []*model.WebTrade
	Hot           []*model.Token

	InitSuccessed bool
}

func GetCacheList() *CacheList {
	return v_cache_list
}

func InitCacheList(hotList string) error {
	var resetErr error
	if v_cache_list.Tokens, resetErr = ResetCacheList(v_lock_token, C_tokens_sort_createdAt, C_cache_tokens_total); resetErr != nil {
		return resetErr
	}
	if v_cache_list.VolumeList, resetErr = ResetCacheList(v_lock_tokenv, C_tokens_sort_volume24, C_cache_tokens_volume_total); resetErr != nil {
		return resetErr
	}
	if v_cache_list.MarketCapList, resetErr = ResetCacheList(v_lock_tokenm, C_tokens_sort_marketcap, C_cache_ranking_total); resetErr != nil {
		return resetErr
	}

	// trades
	tradeList, _, err := GetModelListPageEx[model.Trade](model.C_Trade, bson.M{}, "", "", 0, C_cache_trades_total)
	if err != nil {
		return err
	}
	for _, trade := range tradeList {
		if err := PushTradeToCache(nil, trade, nil); err != nil {
			return err
		}
	}
	// global.Debug("init trade list successed", slog.Any("tradeList", tradeList))

	go GetLaunchPageList()
	ResetTokens21(v_lock_token, nil, false)

	v_cache_list.InitSuccessed = true
	return nil
}

func ResetTokens21(lock *sync.RWMutex, token *model.Token, bUpdate bool) {
	if bUpdate {
		V_cache_tokens21.Store(token.ID, *token)
	} else {
		if lock != nil {
			lock.Lock()
			defer lock.Unlock()
		}
		V_cache_tokens21.Clear()
		for i := 0; i < len(v_cache_list.Tokens) && i < 21; i++ {
			V_cache_tokens21.Store(v_cache_list.Tokens[i].ID, *v_cache_list.Tokens[i])
		}
	}
}

func ResetCacheList(lock *sync.RWMutex, orderField string, totalSize int) ([]*model.Token, error) {
	lock.Lock()
	defer lock.Unlock()

	tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", orderField, 0, totalSize)
	if err != nil {
		return nil, err
	}

	for i := len(tokenList) - 1; i >= 0; i-- {
		if orderField == C_tokens_sort_volume24 {
			tokenList[i].Volume24 = GetVolume24hForToken(tokenList[i], true)
		}
		tokenPriceUsd, percent := GetTokenPricePercent(tokenList[i])
		tokenList[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
		tokenList[i].Chain = global.Config.RPC.Chain
		tokenList[i].Percent = percent
	}

	if orderField == C_tokens_sort_marketcap {
		return model.SortMarketCaps(tokenList)[:], nil
	} else if orderField == C_tokens_sort_volume24 {
		return model.SortVolumes(tokenList)[:], nil
	} else {
		return tokenList[:], nil
	}
}

func GetTokenListForAggregate(userID any, field string, offset, limit int) (ret []*model.Token, total int, err error) {
	switch field {
	case C_tokens_sort_hot:
		v_lock_hot.RLock()
		ret = model.CopyTokens(v_cache_list.Hot[:])
		v_lock_hot.RUnlock()
	case C_tokens_sort_createdAt:
		if offset == 1 {
			V_cache_tokens21.Range(func(key, value any) bool {
				vtoken := value.(model.Token)
				ret = append(ret, &vtoken)
				return true
			})
		} else {
			v_lock_token.RLock()
			ret = model.CopyTokens(v_cache_list.Tokens[:])
			v_lock_token.RUnlock()
		}
	case C_tokens_sort_volume:
		v_lock_tokenv.RLock()
		ret = model.CopyTokens(v_cache_list.VolumeList[:])
		v_lock_tokenv.RUnlock()
	case C_tokens_sort_marketcap:
		v_lock_tokenm.RLock()
		ret = model.CopyTokens(v_cache_list.MarketCapList[:])
		v_lock_tokenm.RUnlock()
	default:
		return nil, 0, fmt.Errorf("invalid field: %s", field)
	}
	if ret == nil {
		ret = make([]*model.Token, 0)
	}
	total = len(ret)
	pageList, pageListErr := GetPageList(ret, offset, 0, limit)

	for i, _ := range pageList {
		if userID != nil {
			pageList[i].Follow = IsFavoriteToken(userID, pageList[i].ID)
		}
		_, percent := GetTokenPricePercent(pageList[i])
		// ret[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
		// ret[i].Volume24 = GetVolume24hForToken(ret[i], false)
		pageList[i].Percent = percent
	}
	return pageList, total, pageListErr
}

func GetTokenListFromCache(userID any, field string, offset, limit, total int, tokens []*model.Token) (ret []*model.Token, err error) {
	if tokens == nil {
		tokens = make([]*model.Token, 0)
	}
	switch field {
	case C_tokens_sort_hot:
		v_lock_hot.RLock()
		ret = model.CopyTokens(v_cache_list.Hot[:])
		v_lock_hot.RUnlock()
	case C_tokens_sort_createdAt:
		if offset == 1 {
			V_cache_tokens21.Range(func(key, value any) bool {
				vtoken := value.(model.Token)
				ret = append(ret, &vtoken)
				return true
			})
		} else {
			v_lock_token.RLock()
			ret = model.CopyTokens(v_cache_list.Tokens[:])
			v_lock_token.RUnlock()
		}
	case C_tokens_sort_volume:
		v_lock_tokenv.RLock()
		ret = model.CopyTokens(v_cache_list.VolumeList[:])
		v_lock_tokenv.RUnlock()
	case C_tokens_sort_marketcap:
		v_lock_tokenm.RLock()
		ret = model.CopyTokens(v_cache_list.MarketCapList[:])
		v_lock_tokenm.RUnlock()
	default:
		return nil, fmt.Errorf("invalid field: %s", field)
	}
	if ret == nil {
		ret = make([]*model.Token, 0)
	}

	if field == C_tokens_sort_createdAt {
		ret = append(ret, tokens...)
		ret = model.SortCreatedAtList(ret)
		slog.Debug("after append fetch tokens", "ret", ret)
	} else if field == C_tokens_sort_marketcap {
		ret = append(ret, tokens...)
		ret = model.SortMarketCaps(ret)
		slog.Debug("after append fetch tokens", "ret", ret)
	}
	retTokens, err := GetPageList(ret, offset, limit, total)

	for i, _ := range retTokens {
		if retTokens[i].Chain == global.Config.RPC.Chain || global.Config.RPC.Chain == "" {
			if userID != nil {
				retTokens[i].Follow = IsFavoriteToken(userID, retTokens[i].ID)
			}
			_, percent := GetTokenPricePercent(retTokens[i])
			// ret[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
			// ret[i].Volume24 = GetVolume24hForToken(ret[i], false)
			retTokens[i].Percent = percent
		}
	}
	return retTokens, err
}

func GetListTradesFromCache() ([]*model.WebTrade, error) {
	v_lock_trade.RLock()
	defer v_lock_trade.RUnlock()
	if v_cache_list.Trades == nil {
		v_cache_list.Trades = make([]*model.WebTrade, 0)
	}
	return v_cache_list.Trades[:], nil
}

func PushTradeToCache(token *model.Token, trade *model.Trade, tokenPrice *model.TokenPrice) error {
	v_lock_trade.Lock()
	defer v_lock_trade.Unlock()

	var terr error
	if token == nil {
		if token, terr = GetToken(trade.Address); terr != nil {
			return terr
		}
	}

	webTrade := &model.WebTrade{}
	if tokenPrice == nil {
		tokenPrice = GetTokenPrice(token.Address, token.Launch, true)
	}
	webTrade.Trade = *trade
	webTrade.Token = *token
	webTrade.TokenAddress = token.Address

	totalSize := C_cache_trades_total
	for i, _ := range v_cache_list.Trades {
		if strings.ToUpper(trade.TxHash) == strings.ToUpper(v_cache_list.Trades[i].TxHash) {
			v_cache_list.Trades[i] = webTrade
			return nil
		}
	}
	if len(v_cache_list.Trades) == 0 {
		totalSize = 1
		v_cache_list.Trades = make([]*model.WebTrade, 0)
	} else if len(v_cache_list.Trades) < totalSize {
		totalSize = len(v_cache_list.Trades) + 1
	}
	v_cache_list.Trades = append([]*model.WebTrade{webTrade}, v_cache_list.Trades[0:totalSize-1]...)
	return nil
}

func InitHotTokens(hotList string) error {
	v_cache_list.Hot = make([]*model.Token, 0)
	initDefaultHotList := func() error {
		tokenList, _, err := GetModelListPageEx[model.Token](model.C_Token, bson.M{}, "desc", "marketcap", 0, 3)
		if err != nil {
			return err
		}
		v_cache_list.Hot = tokenList[:]
		return nil
	}
	if hotList == "" {
		return initDefaultHotList()
	} else {
		file, err := os.Open(hotList)
		if err != nil {
			return initDefaultHotList()
		}
		defer file.Close()

		fileBuf, err := ioutil.ReadAll(file)
		if err != nil {
			return initDefaultHotList()
		}

		V_hot_total = -1
		hotTokens := strings.Split(string(fileBuf), ",")
		for _, address := range hotTokens {
			if token, err := GetToken(strings.TrimSpace(address)); err == nil && token != nil {
				v_cache_list.Hot = append(v_cache_list.Hot, token)
			}
		}
	}
	return nil
}

func WatchHotFile(hotList string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("watch hot file failed", "error", err.Error())
		return
	}
	defer watcher.Close()

	err = watcher.Add(hotList)
	if err != nil {
		slog.Error("watch hot file failed", "error", err.Error())
		return
	}
	slog.Info("start watch hot file", "hotList", hotList)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				InitHotTokens(hotList)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("watch hot file failed", "error", err.Error())
		}
	}
}

func reverse[T any](s []*T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func RefreshSortList(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(global.Config.RefreshSortListTimer) * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				// v_lock_tokenm.Lock()
				// for i, token := range v_cache_list.VolumeList {
				// 	tokenPriceUsd, percent := GetTokenPricePercent(token)
				// 	v_cache_list.MarketCapList[i].MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
				// 	v_cache_list.MarketCapList[i].Percent = percent
				// }
				// v_lock_tokenm.Unlock()

				v_lock_tokenv.RLock()
				tokenList := v_cache_list.VolumeList[:]
				v_lock_tokenv.RUnlock()
				for i, token := range tokenList {
					tokenList[i].Volume24 = GetVolume24hForTokenForTimer(token, true)
				}
			case <-ctx.Done():
				slog.Info("stop RefreshSortList")
				return
			}
		}
	}()
}

func GetVolume24hForToken(token *model.Token, bReset bool) decimal.Decimal {
	return doGetVolume24hForToken(token, bReset, monitor.C_Web_Script_Url)
}

func GetVolume24hForTokenForTimer(token *model.Token, bReset bool) decimal.Decimal {
	return doGetVolume24hForToken(token, bReset, monitor.C_Web_Script_Url) //+"1")
}

func doGetVolume24hForToken(token *model.Token, bReset bool, scriptUrl string) decimal.Decimal {
	timestamp1 := time.Now().Add(-24 * time.Hour).UnixMilli()
	timestamp2 := time.Now().Add(-48 * time.Hour).UnixMilli()
	// timestamp1 := time.Now().Add(-10 * time.Minute).UnixMilli()
	// timestamp2 := time.Now().Add(-20 * time.Minute).UnixMilli()
	if !bReset {
		if token.VolumeTime > timestamp1 {
			return token.Volume24
		} else if token.VolumeTime > timestamp2 {
			return token.Volume24
		} else {
			return decimal.Zero
		}
	}

	volume := decimal.Zero
	quotePrice := monitor.GetQuotePrice()
	whereModel := bson.D{{"address", token.Address}, {"timestamp", bson.D{{"$gte", timestamp1}}}}
	tradeList, _, _ := store.DBListPage[model.Trade](model.C_Trade, whereModel, "asc", "timestamp", 0, 0, false)
	if tradeList != nil && len(tradeList) > 0 {
		for _, trade := range tradeList {
			volume = volume.Add(trade.QuoteAmount.Mul(quotePrice))
		}
	}
	global.Debug("get volume24h successet", "token", token.Address, "volume", volume)
	v24, _ := decimal.NewFromString(volume.StringFixed(20))
	return v24
}

func handleVolumeList24h(token *model.Token, totalSize int) error {
	v_lock_tokenv.Lock()

	// 1. 针对已存在的token，直接更新
	bExistToken := false
	workingVolumeList := make([]*model.Token, 0)
	if len(v_cache_list.VolumeList) > 0 {
		workingVolumeList = v_cache_list.VolumeList[:]
	}
	for i, _ := range workingVolumeList {
		global.Debug("volume ranking check", "token.Address", strings.ToUpper(token.Address), "cache.Token.Address", strings.ToUpper(v_cache_list.Tokens[i].Address))
		if strings.ToUpper(token.Address) == strings.ToUpper(workingVolumeList[i].Address) {
			workingVolumeList[i] = token
			bExistToken = true
			break
		}
	}

	// 2. 针对不存在的token, 先直接增加
	if !bExistToken {
		workingVolumeList = append(workingVolumeList, token)
	}

	// 3. 再按成交量排序, 因为是24小时滚动排序, 缓存里面只会留存按指定数量排名的token
	workingVolumeList = model.SortVolumes(workingVolumeList)

	// 4. 删除排名最后的token
	if len(workingVolumeList) < totalSize {
		totalSize = len(workingVolumeList)
	}
	v_cache_list.VolumeList = workingVolumeList[:totalSize]

	v_lock_tokenv.Unlock()
	return nil
}

func handleMarketCapList(token *model.Token, totalSize int) error {
	v_lock_tokenm.Lock()

	tokenPriceUsd, percent := GetTokenPricePercent(token)
	token.MarketCap = tokenPriceUsd.Mul(decimal.NewFromInt(model.C_total_supply))
	token.Percent = percent

	go UpdateLaunchPageList(token)

	// 1. 针对已存在的token，直接更新
	bExistToken := false
	workingMCList := make([]*model.Token, 0)
	if len(v_cache_list.MarketCapList) > 0 {
		workingMCList = v_cache_list.MarketCapList[:]
	}
	for i, _ := range workingMCList {
		global.Debug("marketcap ranking check", "token.Address", strings.ToUpper(token.Address), "cache.Token.Address", strings.ToUpper(v_cache_list.Tokens[i].Address))
		if strings.ToUpper(token.Address) == strings.ToUpper(workingMCList[i].Address) {
			workingMCList[i] = token
			bExistToken = true
			break
		}
	}

	// 2. 针对不存在的token, 先直接增加
	if !bExistToken {
		workingMCList = append(workingMCList, token)
	}

	// 3. 再按成交量排序, 因为是24小时滚动排序, 缓存里面只会留存按指定数量排名的token
	workingMCList = model.SortMarketCaps(workingMCList)

	// 4. 删除排名最后的token
	if len(workingMCList) < totalSize {
		totalSize = len(workingMCList)
	}
	v_cache_list.MarketCapList = workingMCList[:totalSize]

	v_lock_tokenm.Unlock()
	return nil
}
