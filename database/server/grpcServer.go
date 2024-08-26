package server

import (
	"fmt"
	"log"
	"net"
	"youtube-clone/database/helpers"
	"youtube-clone/database/services"

	"google.golang.org/grpc"
)

var (
	listener   net.Listener
	grpcServer *grpc.Server
)

func StartGrpcServer() {
	host := helpers.FatalIfEmptyVar("GRPC_HOST")
	port := helpers.FatalIfEmptyVar("GRPC_PORT")
	addr := host + ":" + port
	var err error
	listener, err = net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("Failed to listen at:%v\nError:%v", addr, err)
	}
	grpcServer = grpc.NewServer()
	services.RegisterAllServices(grpcServer)
	fmt.Printf("grpc server is running at %s\n", addr)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to listen at:%v\n", addr)
	}
}

func StopGrpcServer() {
	grpcServer.Stop()
}
