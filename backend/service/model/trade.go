package model

import "github.com/shopspring/decimal"

const C_Trade = "trade"
const C_TxRecord = "txrecord"

type Trade struct {
	BaseModel   `bson:",inline"`
	Address     string          `json:"address"`
	Wallet      string          `json:"wallet"`
	Buy         bool            `json:"buy"`
	Liquid      bool            `json:"liquid"`
	Amount      decimal.Decimal `json:"amount"`
	Price       decimal.Decimal `json:"price"`
	MPrice      decimal.Decimal `json:"mprice"`
	TxHash      string          `json:"txhash"`
	Number      int64           `json:"number"`
	Source      decimal.Decimal `json:"-" bson:"source"`
	Timestamp   int64           `json:"timestamp"`
	QuoteAmount decimal.Decimal `json:"quoteAmount"`
}

type Statistics struct {
	TotalCreated int `json:"totalCreated"`
	TotalOwned   int `json:"totalOwned"`
	TotalTrading int `json:"totalTrading"`
}

type WebTrade struct {
	Trade
	Token
	TokenAddress string `json:"tokenAddress"`
}

type TxRecord struct {
	BaseModel `bson:",inline"`
	Hash      string `json:"hash"`
	Number    int64  `json:"number"`
	IsPanic   bool   `json:"isPanic"`
	IsCreate  bool   `json:"isCreate"`
}
