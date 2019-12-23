package models

type Payment struct {
	Direction   Direction `json:"direction"`
	Account     AccountID `json:"account"`
	FromAccount AccountID `json:"from_account,omitempty"`
	ToAccount   AccountID `json:"to_account,omitempty"`
	Amount      Money     `json:"amount"`
}
