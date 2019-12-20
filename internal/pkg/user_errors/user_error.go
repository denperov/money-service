package user_errors

var (
	ErrorWrongAmount             = new("negative or zero amount")
	ErrorWrongSourceAccount      = new("source account not exists")
	ErrorWrongDestinationAccount = new("destination account not exists")
	ErrorDifferentCurrencies     = new("accounts have different currencies")
	ErrorSameAccount             = new("transfer between the same accounts")
	ErrorNotEnoughMoney          = new("not enough money")
)

func IsUserFriendlyError(err error) bool {
	_, ok := err.(*userFriendlyError)
	return ok
}

type userFriendlyError struct{ string }

func (e *userFriendlyError) Error() string {
	return e.string
}

func new(text string) *userFriendlyError {
	return &userFriendlyError{text}
}
