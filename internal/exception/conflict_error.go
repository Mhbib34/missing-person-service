package exception

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

func NewConflictError(message string) ConflictError {
	return ConflictError{Message: message}
}