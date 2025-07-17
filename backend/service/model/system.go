package model

const (
	C_System  = "system"
	C_Listing = "listing"

	C_Status_listed        = "listed"
	C_Status_verifying     = "verifying"
	C_Status_verify_failed = "verifyFailed"
)

type System struct {
	BaseModel `bson:",inline"`
	Hot       []string `json:"hot"`
	Banner    []string `json:"banner"`
}

type ListingPrepare struct {
	Token  string `json:"token"`
	Wallet string `json:"wallet"`
}

type Listing struct {
	BaseModel `bson:",inline"`

	ID string `json:"id"`

	Status         string `json:"status"`
	Entry          bool   `json:"entry"`
	System         bool   `json:"system"`
	Start          int64  `json:"start"`
	Token          string `json:"token"`
	Wallet         string `json:"wallet"`
	TxHash         string `json:"txhash"`
	VerifyResult   string `json:"verifyResult"`
	BannerImageUrl string `json:"bannerImageUrl"`

	Listed bool  `json:"listed" bson:"-"`
	Remain int64 `json:"remain" bson:"-"`

	ImageFileBuf  []byte `json:"-" bson:"-"`
	ImageFileName string `json:"-" bson:"-"`
}

type ListingWait struct {
	Wait  int64 `json:"wait"`
	Total int64 `json:"total"`
}

type ListingList []*Listing

func (t ListingList) Len() int           { return len(t) }
func (t ListingList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ListingList) Less(i, j int) bool { return t[i].Start > t[j].Start }
