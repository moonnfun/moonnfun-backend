package model

const (
	C_System  = "system"
	C_Listing = "listing"

	C_Status_paid      = "paid"
	C_Status_pending   = "pending"
	C_Status_verifying = "verifying"
)

type System struct {
	BaseModel `bson:",inline"`
	Hot       []string `json:"hot"`
	Banner    []string `json:"banner"`
}

type Listing struct {
	BaseModel `bson:",inline"`

	Status         string `json:"status"`
	Listed         bool   `json:"listed" bson:"-"`
	Entry          bool   `json:"entry"`
	System         bool   `json:"system"`
	Start          int64  `json:"start"`
	Token          string `json:"token"`
	Wallet         string `json:"wallet"`
	BannerImageUrl string `json:"bannerImageUrl"`
}

type ListingWait struct {
	Wait  int64 `json:"wait"`
	Total int64 `json:"total"`
}
