package model

import "github.com/shopspring/decimal"

const C_Trade = "trade"

type Trade struct {
	BaseModel   `bson:",inline"`
	Address     string          `json:"address"`
	Wallet      string          `json:"wallet"`
	Buy         bool            `json:"buy"`
	Liquid      bool            `json:"liquid"`
	Amount      decimal.Decimal `json:"amount"`
	Price       decimal.Decimal `json:"price"`
	TxHash      string          `json:"txhash"`
	Timestamp   int64           `json:"timestamp"`
	QuoteAmount decimal.Decimal `json:"quoteAmount"`
}

type Statistics struct {
	TotalCreated int `json:"totalCreated"`
	TotalOwned   int `json:"totalOwned"`
	TotalTrading int `json:"totalTrading"`
}
