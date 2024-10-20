package server

import (
	"fmt"
	"github.com/rzaf/youtube-clone/database/helpers"
	"github.com/rzaf/youtube-clone/notification/services"
	"log"
	"net"

	// "github.com/rzaf/youtube-clone/database/helpers"
	// "github.com/rzaf/youtube-clone/file/services"

	"google.golang.org/grpc"
)

var (
	listener   net.Listener
	grpcServer *grpc.Server
)

func StartGrpcServer() {
	addr := helpers.FatalIfEmptyVar("GRPC_HOST") + ":" + helpers.FatalIfEmptyVar("GRPC_PORT")
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
