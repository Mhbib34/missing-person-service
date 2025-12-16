package exception

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{Message: message}
}