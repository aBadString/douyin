package base

func NewBadRequestError(message string) error {
	return NewError(400, message)
}

func NewUnauthorizedError() error {
	return NewError(401, "请先登录")
}

func NewForbiddenError(message string) error {
	return NewError(403, message)
}

func NewNotFoundError(message string) error {
	return NewError(404, message)
}

func NewServerError(message string) error {
	return NewError(500, message)
}

func NewError(code int, message string) error {
	return &serviceError{
		Code:    code,
		Message: message,
	}
}

type serviceError struct {
	Code    int
	Message string
}

func (b *serviceError) Error() string {
	return b.Message
}

func (b *serviceError) GetCode() int {
	return b.Code
}
