package logic

import (
	"errors"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateToken(userID any, requestToken *model.Token, vid string) (*model.Token, error) {
	if requestToken == nil {
		return nil, errors.New("invalid request token")
	}

	// check tokenID
	imageNames := strings.Split(requestToken.ImageUrl, "/")
	if len(imageNames) > 0 {
		imageFilePath := filepath.Join(global.Config.WebPath, "images", global.Config.RPC.Chain, imageNames[len(imageNames)-1])
		if f, err := os.Stat(imageFilePath); f != nil && err == nil {
			requestToken.ID = strings.TrimSuffix(imageNames[len(imageNames)-1], filepath.Ext(imageNames[len(imageNames)-1]))
			if requestToken.Salt, requestToken.Address = monitor.GetSalt(requestToken.ID, requestToken.Name, requestToken.Symbol); requestToken.Salt == -1 {
				slog.Error("generate salt failed", "token", requestToken)
				return nil, errors.New("create token failed with invalid salt")
			}
		} else {
			slog.Error("invalid image url", "imageUrl", requestToken.ImageUrl, "imageFilePath", imageFilePath)
			return nil, errors.New("invalid token without image")
		}
	} else {
		return nil, errors.New("invalid token without image")
	}

	requestToken.Chain = global.Config.RPC.Chain
	requestToken.Comments = make([]*model.TokenComment, 0)

	requestToken.DBID = primitive.NewObjectID()
	requestToken.CreatedAt = time.Now().UnixMilli()

	store.CacheSetByTime(requestToken.ID, requestToken, true, time.Duration(30)*time.Minute, nil)
	slog.Info("create token successed", "token", requestToken)

	return requestToken, nil
}

func UpdateTokenCreated(data any) (bHandle bool) {
	var token *model.Token
	tokenCreated := data.(*monitor.TokenCreated)

	if global.Config.Sync {
		token = &model.Token{}
		token.Chain = global.Config.RPC.Chain
		token.Tag = model.C_Tag_Meme
		token.Name = tokenCreated.Name
		token.Symbol = tokenCreated.Symbol
		token.DBID = primitive.NewObjectID()
		token.ID = tokenCreated.TokenId.String()
		token.CreatedAt = tokenCreated.Time.UnixMilli()
		token.Creator = tokenCreated.Creator.String()
		token.ImageUrl = fmt.Sprintf("%s/images/%s/%s", global.Config.HostURL, global.Config.RPC.Chain, FindImage(tokenCreated.TokenId.String()))
		token.Comments = make([]*model.TokenComment, 0)
		slog.Info("before create token for syncer", "token", token, "tokenCreated", tokenCreated)
	} else {
		memoryT := store.CacheGet(tokenCreated.TokenId.String(), false, nil)
		if memoryT == nil {
			slog.Error("update token failed with invalid id", slog.Any("tokenCreated", tokenCreated))
			return
		}
		if memoryT.(*model.Token).ID != tokenCreated.TokenId.String() {
			slog.Error("update token failed with invalid id", slog.Any("tokenCreated", tokenCreated))
			return
		}
		tempToken := *(memoryT.(*model.Token))
		token = &tempToken
		defer store.CacheGet(tokenCreated.TokenId, true, nil)
		global.Info("get memory token ok", "token", token, "tokenCreated", tokenCreated)
		// tokenPrice := GetTokenPrice(string(tokenCreated.Token.String()), false, true)
	}
	quotePrice := monitor.GetQuotePrice()

	token.StartTime = tokenCreated.Time.UnixMilli()
	token.Decimals = model.C_decimal_18
	token.BlockNumber = tokenCreated.BlockNumber
	token.StartPrice = decimal.NewFromFloat(global.Config.StartPrice)
	// token.StartPrice = decimal.NewFromFloat(0.2).Div(decimal.NewFromInt(100000000)).Mul(quotePrice)
	token.Supply = tokenCreated.TotalSupply.BigInt().Int64()
	if tokenCreated.Token.String() != monitor.C_zero_address && tokenCreated.Token.String() != "" && tokenCreated.Token.String() != token.Address {
		token.Address = tokenCreated.Token.String()
	}
	if token.Supply <= 0 {
		token.Supply = model.C_total_supply
	}
	token.MarketCap = token.StartPrice.Mul(decimal.NewFromInt(token.Supply).Mul(quotePrice))

	// update or create
	if dbToken, dbErr := GetToken(token.Address); dbToken != nil && dbErr == nil {
		token.DBID = primitive.ObjectID{}
	}

	global.Debug("before update token to db", "token", token)
	if err := store.DBSet(model.C_Token, token, bson.M{"id": token.ID}); err != nil {
		slog.Error("update token failed", "token", token, "error", err.Error())
		return
	}
	slog.Info("update token successed", "token", token)
	return true
}

func GetToken(tokenAddress string) (*model.Token, error) {
	return store.DBGet[model.Token](model.C_Token, bson.M{"address": tokenAddress})
}

func MustGetToken(tokenAddress, pairAddress string, timeout time.Duration) *model.Token {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if tokenAddress != "" && tokenAddress != "0x0000000000000000000000000000000000000000" {
				if token, err := store.DBGet[model.Token](model.C_Token, bson.M{"address": tokenAddress}); err == nil && token != nil {
					return token
				}
			} else if pairAddress != "" && pairAddress != "0x0000000000000000000000000000000000000000" {
				if token, err := store.DBGet[model.Token](model.C_Token, bson.M{"pair": pairAddress}); err == nil && token != nil {
					return token
				}
			}
		case <-time.After(timeout):
			ticker.Stop()
			return nil
		}
	}
}

func GetTokenEx(tokenAddress, pairAddress string) *model.Token {
	if tokenAddress != "" && tokenAddress != "0x0000000000000000000000000000000000000000" {
		if token, err := store.DBGet[model.Token](model.C_Token, bson.M{"address": tokenAddress}); err == nil && token != nil {
			return token
		}
	} else if pairAddress != "" && pairAddress != "0x0000000000000000000000000000000000000000" {
		if token, err := store.DBGet[model.Token](model.C_Token, bson.M{"pair": pairAddress}); err == nil && token != nil {
			return token
		}
	}
	return nil
}

func AddComment(tokenComment *model.TokenComment) (*model.Token, error) {
	token, err := store.DBGet[model.Token](model.C_Token, bson.M{"address": tokenComment.Address})
	if err != nil {
		return nil, err
	}
	dbID := token.DBID

	token.DBID = primitive.ObjectID{}
	tokenComment.CreatedAt = time.Now().UnixMilli()
	token.Comments = append(token.Comments, tokenComment)
	if err := store.DBSet(model.C_Token, token, bson.M{"_id": dbID}); err != nil {
		slog.Error("update token failed", "token", token, "error", err.Error())
		return nil, err
	}
	slog.Info("update token successed", "token", token)
	return token, nil
}

func TokenLaunchNotify(data any) bool {
	tokenLaunchLiquid := data.(*monitor.LaunchLiquid)

	var token *model.Token
	if token = MustGetToken(tokenLaunchLiquid.Token.String(), tokenLaunchLiquid.Pair.String(), time.Duration(5)*time.Minute); token == nil {
		slog.Error("get token from database failed with timeout 5 minutes", "tokenTrading", tokenLaunchLiquid)
		return true
	}
	global.Debug("get token successed", "token", token)

	return true
}
