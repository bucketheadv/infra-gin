package api

import "fmt"

type BizError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BizError) Error() string {
	return e.Message
}

func NewBizError(code int, message string) error {
	return &BizError{Code: code, Message: message}
}

type ParamError struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

func (e *ParamError) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf(e.Message, e.Fields)
	}
	return e.Message
}

func (e *ParamError) Format(fields ...string) *ParamError {
	return &ParamError{
		Fields:  fields,
		Message: e.Message,
	}
}

func NewParamError(message string) *ParamError {
	return &ParamError{Message: message}
}
