package tools

type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}
