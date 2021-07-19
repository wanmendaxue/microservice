//go:generate protoc -I ./pkg/test --go_out=plugins=grpc:./pkg/test ./pkg/test/test.proto
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	errors2 "github.com/wanmendaxue/microservice/errors"

	"net"

	"github.com/sirupsen/logrus"
	"github.com/wanmendaxue/microservice/demo/pkg/test"
	"github.com/wanmendaxue/microservice/rpc/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logrus.SetLevel(logrus.DebugLevel)
	s := grpc.NewGrpcServer()
	impl := NewTestServer()

	test.RegisterTestServiceServer(s, impl)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewTestServer() test.TestServiceServer {
	return &testServer{}
}

type testServer struct {
	test.UnimplementedTestServiceServer
}

func (t testServer) Hello(ctx context.Context, req *test.HelloRequest) (*test.HelloReply, error) {
	return nil, e2()
}

func e1() error {
	a := errors.New("a")
	b := fmt.Errorf("b: %w", a)
	c := fmt.Errorf("c: %w", b)

	return c
}

func e2() error {
	return errors2.NewBusinessError(111222, "business err test")
}
