package core

import "fmt"

const (
	ErrCodeBadRequest    = "ERR_BAD_REQUEST"
	ErrCodeInternalError = "ERR_INTERNAL_ERROR"
)

type Error struct {
	ErrCode string `json:"err"`
	Message string `json:"msg"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v - %v", e.ErrCode, e.Message)
}

func NewBadRequestError(msg string) *Error {
	return &Error{
		ErrCode: ErrCodeBadRequest,
		Message: msg,
	}
}

func NewInternalError(err error) *Error {
	return &Error{
		ErrCode: ErrCodeInternalError,
		Message: err.Error(),
	}
}
