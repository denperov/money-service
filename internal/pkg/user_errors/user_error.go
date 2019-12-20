package user_errors

type UserFriendlyError struct{ string }

func (e *UserFriendlyError) Error() string {
	return e.string
}

func New(text string) *UserFriendlyError {
	return &UserFriendlyError{text}
}
