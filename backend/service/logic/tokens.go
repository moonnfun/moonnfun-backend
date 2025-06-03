package logic

import (
	"fmt"
	"log/slog"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"sort"
	"sync"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	c_Max_CacheListSize = 10

	c_Kind_Tokens_Hot           = "Hot"
	c_Kind_Tokens_MarketCap     = "MarketCap"
	c_Kind_Tokens_TradingVolume = "TradingVolume"
)

var (
	v_tradings         = new(sync.Map)
	v_cache_tokens     = &CacheTokens{}
	v_Lock_CacheTokens = new(sync.Mutex)
)

type TokenVolumeList []*model.Token
type TokenMarketCapList []*model.Token

func (t TokenVolumeList) Len() int           { return len(t) }
func (t TokenVolumeList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenVolumeList) Less(i, j int) bool { return t[i].Volume.LessThan(t[j].Volume) }

func (t TokenMarketCapList) Len() int           { return len(t) }
func (t TokenMarketCapList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenMarketCapList) Less(i, j int) bool { return t[i].MarketCap.LessThan(t[j].MarketCap) }

type CacheTokens struct {
	Hot           TokenVolumeList
	MarketCap     TokenMarketCapList
	TradingVolume TokenVolumeList
}

func InitCacheTokens(bUpdate bool) error {
	tokens, err := GetTokenListDB("")
	if err != nil {
		return err
	}
	for i, token := range tokens {
		updatedToken, err := initTokenTrading(token, token.Address, nil)
		if err != nil {
			slog.Error("GetTokenTxList failed", "address", token.Address, "error", err.Error())
			continue
		}
		tokens[i] = updatedToken
	}
	// 分别按MarketCap和Volume排序
	tokenVolumes := TokenVolumeList(tokens)
	tokenMarketCaps := TokenMarketCapList(tokens)
	sort.Sort(tokenMarketCaps)
	sort.Sort(tokenVolumes)

	// Hot
	if v_cache_tokens.Hot == nil {
		v_cache_tokens.Hot = make(TokenVolumeList, 0)
		v_cache_tokens.Hot = tokenVolumes[len(tokenVolumes)-c_Max_CacheListSize-1:]
	}

	// TradingVolume
	if v_cache_tokens.TradingVolume == nil {
		v_cache_tokens.TradingVolume = make(TokenVolumeList, 0)
		v_cache_tokens.TradingVolume = tokenVolumes[len(tokenVolumes)-c_Max_CacheListSize-1:]
	}

	// MarketCap
	if v_cache_tokens.MarketCap == nil {
		v_cache_tokens.MarketCap = make(TokenMarketCapList, 0)
		v_cache_tokens.MarketCap = tokenMarketCaps[len(tokenMarketCaps)-c_Max_CacheListSize-1:]
	}
	return nil
}

func initTokenTrading(token *model.Token, tokenAddress string, tokenTrading *monitor.TokenTrading) (*model.Token, error) {
	vTotal, vVolume, vMarketCap := decimal.Zero, decimal.Zero, decimal.Zero
	if token == nil {
		dbToken, err := store.DBGet[model.Token](model.C_Token, bson.M{"address": tokenAddress})
		if err != nil {
			return nil, fmt.Errorf("store.DBGet token failed, address: %s, error: %s", tokenAddress, err.Error())
		}
		token = dbToken
		if ttrading := GetTokenTrading(tokenAddress); ttrading != nil {
			vVolume = ttrading.Volume.Add(tokenTrading.Amount)
			vMarketCap = tokenTrading.Price.Mul(decimal.NewFromInt(token.Supply))
			vTotal = token.TotalAmount.Add(tokenTrading.QuoteAmount)
		}
	} else {
		// update trading
		vTotal, vVolume, vMarketCap = token.TokenTrading.TotalAmount, token.TokenTrading.Volume, token.TokenTrading.MarketCap
		tradings, err := monitor.GetTokenTxList(token.Address)
		if err != nil {
			return nil, fmt.Errorf("GetTokenTxList failed, address: %s, error: %s", token.Address, err.Error())
		}
		for _, trading := range tradings {
			if trading.Buy {
				vVolume = vVolume.Add(trading.Amount)
				vTotal = vTotal.Add(trading.QuoteAmount)
			} else {
				vVolume = vVolume.Sub(trading.Amount)
				vTotal = vTotal.Sub(trading.QuoteAmount)
			}
			vMarketCap = trading.Price.Mul(decimal.NewFromInt(token.Supply))
		}
	}

	// cache tradings
	tTrading := &model.TokenTrading{
		Volume:      vVolume,
		MarketCap:   vMarketCap,
		TotalAmount: vTotal,
	}
	v_tradings.Store(tokenAddress, tTrading)

	token.TokenTrading = *tTrading
	return token, nil
}

func GetTokenTrading(tokenAddress string) *model.TokenTrading {
	trading, _ := v_tradings.Load(tokenAddress)
	if trading != nil {
		return trading.(*model.TokenTrading)
	}
	return nil
}

func GetTokensFromCache(kind string) (any, error) {
	switch kind {
	case c_Kind_Tokens_Hot:
		return v_cache_tokens.Hot, nil
	case c_Kind_Tokens_MarketCap:
		return v_cache_tokens.MarketCap, nil
	case c_Kind_Tokens_TradingVolume:
		return v_cache_tokens.TradingVolume, nil
	default:
		return []any{}, fmt.Errorf("invalid params, kind: %s", kind)
	}
}

func UpdateCacheToken(tokenTrading *monitor.TokenTrading) error {
	v_Lock_CacheTokens.Lock()
	defer v_Lock_CacheTokens.Unlock()

	// Hot
	bExistHot := false
	for i, token := range v_cache_tokens.Hot {
		if token.Address == tokenTrading.Token.String() {
			v_cache_tokens.Hot[i].Volume = v_cache_tokens.Hot[i].Volume.Add(tokenTrading.Amount)
			bExistHot = true
			break
		}
	}
	if !bExistHot {
		updatedToken, err := initTokenTrading(nil, tokenTrading.Token.String(), tokenTrading)
		if err != nil {
			return fmt.Errorf("initTokenTrading failed for Hot, tokenTrading: %+v, error: %s", tokenTrading, err.Error())
		}
		v_cache_tokens.Hot = append(v_cache_tokens.Hot, updatedToken)
		sort.Sort(v_cache_tokens.Hot)
	}

	// TradingVolume
	bExistVolume := false
	for i, token := range v_cache_tokens.TradingVolume {
		if token.Address == tokenTrading.Token.String() {
			v_cache_tokens.TradingVolume[i].Volume = v_cache_tokens.TradingVolume[i].Volume.Add(tokenTrading.Amount)
			bExistVolume = true
			break
		}
	}
	if !bExistVolume {
		updatedToken, err := initTokenTrading(nil, tokenTrading.Token.String(), tokenTrading)
		if err != nil {
			return fmt.Errorf("initTokenTrading failed for Volume, tokenTrading: %+v, error: %s", tokenTrading, err.Error())
		}
		v_cache_tokens.TradingVolume = append(v_cache_tokens.TradingVolume, updatedToken)
		sort.Sort(v_cache_tokens.TradingVolume)
	}

	// MarketCap
	bExistMarketCap := false
	for i, token := range v_cache_tokens.MarketCap {
		if token.Address == tokenTrading.Token.String() {
			v_cache_tokens.MarketCap[i].Volume = v_cache_tokens.MarketCap[i].Volume.Add(tokenTrading.Amount)
			bExistMarketCap = true
			break
		}
	}
	if !bExistMarketCap {
		updatedToken, err := initTokenTrading(nil, tokenTrading.Token.String(), tokenTrading)
		if err != nil {
			return fmt.Errorf("initTokenTrading failed for MarketCap, tokenTrading: %+v, error: %s", tokenTrading, err.Error())
		}
		v_cache_tokens.MarketCap = append(v_cache_tokens.MarketCap, updatedToken)
		sort.Sort(v_cache_tokens.MarketCap)
	}

	return nil
}

// func PushCacheToken(token *model.Token, marketCap, volume decimal.Decimal) error {
// 	if token == nil {
// 		return errors.New("invalid token")
// 	}

// 	// Hot
// 	if volume.GreaterThan(decimal.Zero) {
// 		token.Volume = token.Volume.Add(volume)
// 		if err := pushTokenToCache(v_cache_tokens.Hot, token); err != nil {
// 			return err
// 		}
// 	}

// 	// TradingVolume
// 	if volume.GreaterThan(decimal.Zero) {
// 		token.Volume = token.Volume.Add(volume)
// 		if err := pushTokenToCache(v_cache_tokens.TradingVolume, token); err != nil {
// 			return err
// 		}
// 	}

// 	// MarketCap
// 	if marketCap.GreaterThan(decimal.Zero) {
// 		token.MarketCap = token.MarketCap.Add(marketCap)
// 		if err := pushTokenToCache(v_cache_tokens.MarketCap, token); err != nil {
// 			return err
// 		}
// 	}

// 	// PushToDB, BigGames, TopWins
// 	createGamesLog(player, glog)
// 	// InitGameLogs(true)
// 	updateGamesLog(c_Kind_GameLogs_BigGames, glog)
// 	updateGamesLog(c_Kind_GameLogs_TopWins, glog)
// 	return nil
// }
