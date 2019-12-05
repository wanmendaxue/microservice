package grpcerr

import (
	"github.com/pkg/errors"
	mserror "github.com/wanmendaxue/microservice/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Encode(err error) error {
	return EncodeWrapf(err, "dependent services are currently unavailable")
}

func EncodeWrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	switch ex := err.(type) {
	case *mserror.BusinessError:
		return status.New(codes.Code(ex.Code), ex.Msg).Err()
	case *mserror.RpcError:
		wrap := status.New(codes.Code(ex.Code), ex.Msg).Err()
		if ex.Code < uint32(mserror.MinimalBusinessCode) {
			wrap = errors.Wrapf(wrap, format, args)
		}
		return wrap
	default:
		return status.Error(codes.Unknown, ex.Error())
	}
}

