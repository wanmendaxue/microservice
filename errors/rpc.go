package errors

import "fmt"

type RpcError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func (rr RpcError) Error() string {
	return fmt.Sprintf("[%d] %s", rr.Code, rr.Msg)
}