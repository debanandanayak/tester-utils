package internal

type FriendlyError struct {
	UserError string
}

func (e *FriendlyError) Error() string {
	return e.UserError
}
