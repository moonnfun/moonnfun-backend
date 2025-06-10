package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const C_User = "user"

type BaseModel struct {
	DBID      primitive.ObjectID `bson:"_id,omitempty" json:"gid"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
}

type User struct {
	BaseModel `bson:",inline"`

	Address string `json:"address"`
	Avatar  string `json:"avatar"`
}

type UserStatistics struct {
	TotalCreated int64  `json:"totalCreated"`
	TotalOwned   int64  `json:"totalOwned"`
	TotalTrading string `json:"totalTrading"`
}

// func (p User) MarshalJSON() ([]byte, error) {
// 	ret := &struct {
// 		Address string `json:"address"`
// 		Avatar  string `json:"avatar"`

// 		OwnedTokens   []*Token `json:"ownedTokens"`
// 		CreatedTokens []*Token `json:"createdTokens"`
// 	}{
// 		Address: p.Address,
// 		Avatar:  p.Avatar,
// 	}

// 	// token

// 	return json.Marshal(ret)
// }
