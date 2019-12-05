package errors

import (
	"fmt"
)

var (
	MinimalBusinessCode = 100000
)

type BusinessError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func (be BusinessError) Error() string {
	return fmt.Sprintf("[%s] %s", be.Code, be.Msg)
}

// create a new business error
func NewBusinessError(code uint32, msg string) *BusinessError {
	if code < uint32(MinimalBusinessCode) {
		panic(fmt.Sprintf("business code should greater than %d", MinimalBusinessCode))
	}

	if msg == "" {
		panic("msg is required")
	}

	return &BusinessError{
		Code: code,
		Msg:  msg,
	}
}
