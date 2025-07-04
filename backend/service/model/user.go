package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const C_User = "user"
const C_Referral = "referral"

type BaseModel struct {
	DBID      primitive.ObjectID `bson:"_id,omitempty" json:"gid"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
}

type User struct {
	BaseModel `bson:",inline"`

	Address     string   `json:"address"`
	Avatar      string   `json:"avatar"`
	FromID      string   `json:"fromID"`
	RefferalID  string   `json:"refferalID"`
	RefferalUrl string   `json:"refferalUrl"`
	FollowList  []string `json:"followList"`

	Project `bson:",inline"`
}

type Refferal struct {
	BaseModel `bson:",inline"`

	Wallet  string `json:"wallet"`
	Address string `json:"address"`
}

type UserStatistics struct {
	TotalCreated int64  `json:"totalCreated"`
	TotalOwned   int64  `json:"totalOwned"`
	TotalTrading string `json:"totalTrading"`
}

type UserTokenHolder struct {
	TokenAddress string `json:"address"`
	TokenHolder
	Token
}

type Project struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Twitter     string `json:"twitter"`
	Description string `json:"description"`
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
