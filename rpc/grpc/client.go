package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/wanmendaxue/microservice/errors"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"time"
)

func Dial(addr string, dialOption grpcpkg.DialOption) (*grpcpkg.ClientConn, error) {
	return grpcpkg.Dial(
		addr,
		dialOption,
		grpcpkg.WithUnaryInterceptor(newErrorHandlingUnaryClientInterceptor()),
		grpcpkg.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Minute,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)
}

func newErrorHandlingUnaryClientInterceptor() grpcpkg.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpcpkg.ClientConn, invoker grpcpkg.UnaryInvoker, opts ...grpcpkg.CallOption) error {
		t := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		logrus.Debugf("invoke %s, taken %d millis", method, time.Since(t).Milliseconds())

		if ex := status.Convert(err); ex != nil {
			code := uint32(ex.Code())
			if code > uint32(errors.MinimalBusinessCode) {
				err = errors.NewBusinessError(code, ex.Message())
			}
		}

		return err
	}
}
