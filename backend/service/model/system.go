package model

const C_System = "system"

type System struct {
	BaseModel `bson:",inline"`
	Hot       []string `json:"hot"`
	Banner    []string `json:"banner"`
}
