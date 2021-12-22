package main

import (
	"fmt"
	"net"

	v1 "github.com/nitrictech/nitric/interfaces/nitric/v1"
	"github.com/nitrictech/nitric/pkg/deploy"
	"google.golang.org/grpc"
)

func main() {
	srv := grpc.NewServer()
	dSrv := deploy.New(deploy.NewApp())

	v1.RegisterResourceServiceServer(srv, dSrv)
	v1.RegisterFaasServiceServer(srv, dSrv)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	fmt.Println("listening @ 50051")

	srv.Serve(lis)
}
