package service

type AccountID string

type Money string

type Currency string

type Direction string

const (
	OutgoingDirection Direction = "outgoing"
	IncomingDirection Direction = "incoming"
)

type Account struct {
	ID       AccountID `json:"id"`
	Currency Currency  `json:"currency"`
	Balance  Money     `json:"balance"`
}

type Payment struct {
	Direction   Direction `json:"direction"`
	Account     AccountID `json:"account"`
	FromAccount AccountID `json:"from_account,omitempty"`
	ToAccount   AccountID `json:"to_account,omitempty"`
	Amount      Money     `json:"amount"`
}

type Transfer struct {
	FromAccount AccountID `json:"from_account,omitempty"`
	ToAccount   AccountID `json:"to_account,omitempty"`
	Amount      Money     `json:"amount"`
}
