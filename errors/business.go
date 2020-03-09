package errors

import (
	"fmt"
)

var (
	MinimalBusinessCode = 100000
)

type businessErr struct {
	code uint32
	msg  string
}

func (be *businessErr) Error() string {
	return fmt.Sprintf("[%d] %s", be.code, be.msg)
}

func (be *businessErr) Business() (uint32, string) {
	return be.code, be.msg
}

type businessErrorDefinition struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func (bed *businessErrorDefinition) BuildErr() error {
	return &businessErr{
		code: bed.Code,
		msg:  bed.Msg,
	}
}

// define a new business error
func DefineBusinessError(code uint32, msg string) *businessErrorDefinition {
	if code < uint32(MinimalBusinessCode) {
		panic(fmt.Sprintf("business code should greater than %d", MinimalBusinessCode))
	}

	if msg == "" {
		panic("msg is required")
	}

	return &businessErrorDefinition{
		Code: code,
		Msg:  msg,
	}
}

func NewBusinessError(code uint32, msg string) error {
	return &businessErr{
		code: code,
		msg:  msg,
	}
}
