package grpcerr

import (
	mserror "github.com/wanmendaxue/microservice/errors"
	"google.golang.org/grpc/status"
)

func Decode(err error) error {
	if err == nil {
		return nil
	}

	ex := status.Convert(err)

	return &mserror.RpcError{
		Code: uint32(ex.Code()),
		Msg:  ex.Message(),
	}
}