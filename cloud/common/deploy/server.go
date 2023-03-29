package deploy

import (
	"fmt"
	"log"
	"net"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"google.golang.org/grpc"
)

func StartServer(deploySrv v1.DeployServiceServer) {
	port := utils.GetEnv("PORT", "50051")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("error listening on port %s %v", port, err)
	}

	srv := grpc.NewServer()

	v1.RegisterDeployServiceServer(srv, deploySrv)

	fmt.Printf("Deployment server started on %s\n", lis.Addr().String())
	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("error serving requests %v", err)
	}
}
