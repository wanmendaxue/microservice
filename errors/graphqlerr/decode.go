package graphqlerr

import (
	"github.com/vektah/gqlparser/gqlerror"
	mserror "github.com/wanmendaxue/microservice/errors"
)

func Decode(err error) error {
	if err == nil {
		return nil
	}

	if ex, ok := err.(*gqlerror.Error); ok {
		var rpcCode uint32
		var rpcMsg string

		if code, ok := ex.Extensions[GrpcDefaultErrCodeKey]; ok {
			rpcCode, _ = code.(uint32)
		}

		if msg, ok := ex.Extensions[GrpcDefaultErrMsgKey]; ok {
			rpcMsg, _ = msg.(string)
		} else {
			rpcMsg = ex.Message
		}

		return &mserror.RpcError{
			Code: rpcCode,
			Msg:  rpcMsg,
		}
	}

	return &mserror.RpcError{
		Msg: err.Error(),
	}
}