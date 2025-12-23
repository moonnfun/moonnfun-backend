package model

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

const C_Token = "token"
const C_Token_holder = "holder"
const C_Token_Trading = "trading"

const (
	C_decimal_18   = 18
	C_total_supply = 1000000000

	C_Tag_AI   = "AI"
	C_Tag_Defi = "Defi"
	C_Tag_Game = "Game"
	C_Tag_Meme = "Meme"

	C_Tag_Removed = "Removed"
)

var BeforeQueryTokenList func(token *Token) error

type Token struct {
	BaseModel   `bson:",inline"`
	ID          string          `json:"id"`
	Tag         string          `json:"tag"`
	Name        string          `json:"name"`
	Pair        string          `json:"pair"`
	Chain       string          `json:"chain"`
	Symbol      string          `json:"symbol"`
	Supply      int64           `json:"supply"`
	Creator     string          `json:"creator"`
	Address     string          `json:"address"`
	Decimals    int             `json:"decimals"`
	ImageUrl    string          `json:"imageUrl"`
	Website     string          `json:"website"`
	Twitter     string          `json:"twitter"`
	Telegram    string          `json:"telegram"`
	DeBoxUrl    string          `json:"deboxUrl"`
	Progress    float64         `json:"prograss"`
	Launch      bool            `json:"launch"`
	BuyAmount   string          `json:"buyAmount"`
	Description string          `json:"description"`
	Comments    []*TokenComment `json:"comments"`
	Holders     int             `json:"holders"`
	Percent     float64         `json:"percent" bson:"-"`
	Follow      bool            `json:"follow" bson:"-"`
	BlockNumber int64           `json:"blockNumber"`
	StartTime   int64           `json:"startTime"`
	VolumeTime  int64           `json:"volumeTime"`
	Volume24    decimal.Decimal `json:"volume24"`
	StartPrice  decimal.Decimal `json:"startPrice"`

	TokenTrading `bson:",inline"`

	RaisedAmount decimal.Decimal `json:"raisedAmount" bson:"-"`

	Salt       int64   `json:"salt" bson:"-"`
	LaunchMC   float64 `json:"launchMC" bson:"-"`
	LaunchTime int64   `json:"launchTime"`

	ImageFileBuf  []byte `json:"-" bson:"-"`
	ImageFileName string `json:"-" bson:"-"`
}

type WebToken struct {
	Token
	TokenPrice
}

type TokenTrading struct {
	Volume           decimal.Decimal `json:"volume"`
	MarketCap        decimal.Decimal `json:"marketCap"`
	TotalAmount      decimal.Decimal `json:"totalAmount"`
	TokenAmount      decimal.Decimal `json:"tokenAmount"`
	VirtualLiquidity decimal.Decimal `json:"virtualLiquidity"`
}

type TokenComment struct {
	Address   string `json:"address"`
	Creator   string `json:"creator"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

type TokenPrice struct {
	Price       decimal.Decimal `json:"price"`
	PriceUsd    decimal.Decimal `json:"priceUsd"`
	CurrentTime int64           `json:"currentTime"`

	Method string `json:"-"`
}

type TokenOHLCV struct {
	O       decimal.Decimal `json:"open"`
	C       decimal.Decimal `json:"close"`
	H       decimal.Decimal `json:"high"`
	L       decimal.Decimal `json:"low"`
	V       decimal.Decimal `json:"volume"`
	T       uint64          `json:"timestamp"`
	OT      uint64          `json:"-"`
	StartAt time.Time
}

type TokenHolder struct {
	BaseModel `bson:",inline"`
	Address   string          `json:"address"`
	Wallet    string          `json:"wallet"`
	Amount    decimal.Decimal `json:"amount"`
	Percent   decimal.Decimal `json:"percent"`
	AvgPrice  decimal.Decimal `json:"avgPrice"`
}

// type TokenTrading struct {
// 	BaseModel `bson:",inline"`
// 	Token     string          `bson:"token"`
// 	Wallet    string          `bson:"wallet"`
// 	Buy       bool            `bson:"buy"`
// 	Amount    decimal.Decimal `bson:"amount"`
// 	Price     decimal.Decimal `bson:"price"`
// 	TxHash    string          `json:"txhash"`
// 	Timestamp time.Time       `bson:"timestamp"`
// }

func (p TokenOHLCV) MarshalJSON() ([]byte, error) {
	ret := &struct {
		O decimal.Decimal `json:"open"`
		C decimal.Decimal `json:"close"`
		H decimal.Decimal `json:"high"`
		L decimal.Decimal `json:"low"`
		V float64         `json:"volume"`
		T uint64          `json:"timestamp"`
	}{}
	ret.O = p.O
	ret.H = p.H
	ret.L = p.L
	ret.C = p.C
	ret.V, _ = p.V.Float64()
	ret.T = p.T
	return json.Marshal(ret)
}

// func (t *Token) BeforeQuery(ctx context.Context) error {
// 	slog.Info("before query called", slog.Any("token", t.Address), slog.Any("ctx", ctx))
// 	if BeforeQueryTokenList != nil {
// 		return BeforeQueryTokenList(t)
// 	}
// 	return nil
// }

type TokenList []*Token
type TokenVolumeList []*Token
type TokenMarketCapList []*Token
type TokenProgressList []*Token
type TokenHolderList []*TokenHolder
type TokenCommentList []*TokenComment
type UserTokenHolderList []*UserTokenHolder
type UserTokenHolderListMC []*UserTokenHolder
type UserTokenHolderListVolume []*UserTokenHolder
type UserTokenHolderListProgress []*UserTokenHolder

func (t TokenList) Len() int           { return len(t) }
func (t TokenList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenList) Less(i, j int) bool { return t[i].CreatedAt > t[j].CreatedAt }

func (t TokenVolumeList) Len() int           { return len(t) }
func (t TokenVolumeList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenVolumeList) Less(i, j int) bool { return t[i].Volume24.GreaterThan(t[j].Volume24) }

func (t TokenHolderList) Len() int           { return len(t) }
func (t TokenHolderList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenHolderList) Less(i, j int) bool { return t[i].Percent.GreaterThan(t[j].Percent) }

func (t TokenProgressList) Len() int           { return len(t) }
func (t TokenProgressList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenProgressList) Less(i, j int) bool { return t[i].Progress > t[j].Progress }

func (t TokenMarketCapList) Len() int           { return len(t) }
func (t TokenMarketCapList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenMarketCapList) Less(i, j int) bool { return t[i].MarketCap.GreaterThan(t[j].MarketCap) }

func (t TokenCommentList) Len() int           { return len(t) }
func (t TokenCommentList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenCommentList) Less(i, j int) bool { return t[i].CreatedAt > t[j].CreatedAt }

func (t UserTokenHolderList) Len() int      { return len(t) }
func (t UserTokenHolderList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t UserTokenHolderList) Less(i, j int) bool {
	return t[i].TokenHolder.Amount.GreaterThan(t[j].TokenHolder.Amount)
}

func (t UserTokenHolderListMC) Len() int           { return len(t) }
func (t UserTokenHolderListMC) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t UserTokenHolderListMC) Less(i, j int) bool { return t[i].MarketCap.GreaterThan(t[j].MarketCap) }

func (t UserTokenHolderListVolume) Len() int           { return len(t) }
func (t UserTokenHolderListVolume) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t UserTokenHolderListVolume) Less(i, j int) bool { return t[i].Volume.GreaterThan(t[j].Volume) }

func (t UserTokenHolderListProgress) Len() int           { return len(t) }
func (t UserTokenHolderListProgress) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t UserTokenHolderListProgress) Less(i, j int) bool { return t[i].Progress > t[j].Progress }

func SortCreatedAtList(tokenList []*Token) TokenList {
	ret := TokenList(tokenList)
	sort.Sort(ret)
	return ret
}

func SortMarketCaps(tokenList []*Token) TokenMarketCapList {
	ret := TokenMarketCapList(tokenList)
	sort.Sort(ret)
	return ret
}

func SortVolumes(tokenList []*Token) TokenVolumeList {
	ret := TokenVolumeList(tokenList)
	sort.Sort(ret)
	return ret
}

func SortTokenHolders(tokenHolderList []*TokenHolder) TokenHolderList {
	ret := TokenHolderList(tokenHolderList)
	sort.Sort(ret)
	return ret
}

func SortProgress(tokenList []*Token) TokenProgressList {
	ret := TokenProgressList(tokenList)
	sort.Sort(ret)
	return ret
}

func SortTokenComments(commentList []*TokenComment) TokenCommentList {
	ret := TokenCommentList(commentList)
	sort.Sort(ret)
	return ret
}

func CopyTokens(tokens []*Token) []*Token {
	ret := make([]*Token, 0)
	for _, token := range tokens {
		tmpToken := *token
		ret = append(ret, &tmpToken)
	}
	return ret
}
