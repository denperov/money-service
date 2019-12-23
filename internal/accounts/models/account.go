package models

type Account struct {
	ID       AccountID `json:"id"`
	Currency Currency  `json:"currency"`
	Balance  Money     `json:"balance"`
}
