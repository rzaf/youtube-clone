package client

import (
	"log"
	"youtube-clone/database/helpers"
	user_pb "youtube-clone/database/pbs/user-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GrpcClient    *grpc.ClientConn
	serverAddress string

	UserService user_pb.UserServiceClient
	// MediaService    media.MediaServiceClient
)

func ConnectToDataBaseServer() {
	host := helpers.FatalIfEmptyVar("DATABASE_SERVICE_HOST")
	port := helpers.FatalIfEmptyVar("DATABASE_SERVICE_PORT")
	serverAddress = host + ":" + port
	var err error
	GrpcClient, err = grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("connection to "+serverAddress+" failed:", err)
	}

	// services
	UserService = user_pb.NewUserServiceClient(GrpcClient)

	log.Println(GrpcClient.GetState())
	GrpcClient.Connect()
	log.Printf("grpc state:%s\n", GrpcClient.GetState().String())
}

func DisconnectFromServer() {
	GrpcClient.Close()
}
