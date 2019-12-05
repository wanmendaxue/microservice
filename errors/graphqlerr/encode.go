package graphqlerr

import (
	"github.com/vektah/gqlparser/gqlerror"
	mserror "github.com/wanmendaxue/microservice/errors"
)

var (
	GrpcDefaultErrCodeKey = "code"
	GrpcDefaultErrMsgKey  = "msg"
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
		return &gqlerror.Error{
			Message: ex.Error(),
			Extensions: map[string]interface{}{
				GrpcDefaultErrCodeKey: ex.Code,
				GrpcDefaultErrMsgKey:  ex.Msg,
			},
		}
	case *mserror.RpcError:
		if ex.Code < uint32(mserror.MinimalBusinessCode) {
			return gqlerror.Errorf(format, args)
		}

		return &gqlerror.Error{
			Message: ex.Error(),
			Extensions: map[string]interface{}{
				GrpcDefaultErrCodeKey: ex.Code,
				GrpcDefaultErrMsgKey:  ex.Msg,
			},
		}
	default:
		return &gqlerror.Error{
			Message: ex.Error(),
		}
	}
}
