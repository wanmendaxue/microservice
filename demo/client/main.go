package main

import (
	"context"
	"github.com/wanmendaxue/microservice/demo/pkg/test"
	"github.com/wanmendaxue/microservice/rpc/grpc"
	grpcpkg "google.golang.org/grpc"
	"log"
)

func main() {
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpcpkg.WithInsecure())

	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()


	client := test.NewTestServiceClient(conn)

	_, errx := client.Hello(context.Background(), &test.HelloRequest{
		Msg: "lurongkai",
	})
	if errx != nil {
		log.Printf("%+v", errx)
	}
}