package models

type Balance struct {
	GormModel
	Amount float64 `json:"amount"`
}
