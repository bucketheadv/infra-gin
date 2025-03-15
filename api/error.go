package api

type BizError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BizError) Error() string {
	return e.Message
}

func NewError(code int, message string) error {
	return &BizError{Code: code, Message: message}
}
