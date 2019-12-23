package models

type Transfer struct {
	FromAccount AccountID `json:"from_account,omitempty"`
	ToAccount   AccountID `json:"to_account,omitempty"`
	Amount      Money     `json:"amount"`
}
