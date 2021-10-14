package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

func NewGrpcServer() *grpcpkg.Server {
	return grpcpkg.NewServer(
		grpcpkg.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}),
		grpcpkg.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Minute, // because linux ipvs default timeout is 900s
			Timeout: 120 * time.Second,
		}),
		grpcpkg.MaxRecvMsgSize(1024*1024*100),
		grpcpkg.MaxSendMsgSize(1024*1024*100),
		grpcpkg.MaxConcurrentStreams(2000),
		grpcpkg.UnaryInterceptor(newErrorHandlingUnaryServerInterceptor()))
}

func newErrorHandlingUnaryServerInterceptor() grpcpkg.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpcpkg.UnaryServerInfo, handler grpcpkg.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("service unavailable currently: %+v", r)
				logrus.WithField("panic", r).Error(msg)
				err = status.Error(codes.Internal, msg)
			}
		}()

		t := time.Now()
		resp, err = handler(ctx, req)
		logrus.Debugf("%s processed taken %d millis", info.FullMethod, time.Since(t).Milliseconds())

		if err != nil {
			if _, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
				// if err is GRPC status error object, then do nothing
				logrus.WithFields(logrus.Fields{"error": err, "type": "grpc-status"}).Debugf("service response status error: %v", err)
			} else if ex, ok := err.(interface{ Business() (uint32, string) }); ok {
				c, m := ex.Business()
				logrus.WithFields(logrus.Fields{"error": err, "type": "business"}).Debugf("service response business error: [%d] %s", c, m)
				err = status.Error(codes.Code(c), m)
			} else if ex, ok := err.(interface{ Demand() string }); ok {
				msg := ex.Demand()
				logrus.WithFields(logrus.Fields{"error": err, "type": "demand"}).Debugf("service response demand error: %s", msg)
				err = status.Error(codes.InvalidArgument, msg)
			} else {
				logrus.WithFields(logrus.Fields{"error": err, "type": "general"}).Debugf("service response general error: %+v", err)
				err = status.Error(codes.Unknown, err.Error())
			}
		}
		return
	}
}
