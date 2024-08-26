package client

import (
	"log"
	"youtube-clone/database/helpers"
	"youtube-clone/file/pbs/file"
	"youtube-clone/notification/pbs/emailPb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	FileServiceGrpcClient         *grpc.ClientConn
	NotificationServiceGrpcClient *grpc.ClientConn

	FileService  file.FileServiceClient
	EmailService emailPb.EmailServiceClient
)

func ConnectToFileServer() {
	host := helpers.FatalIfEmptyVar("FILE_SERVICE_HOST")
	port := helpers.FatalIfEmptyVar("FILE_SERVICE_PORT")
	fileServerAddress := host + ":" + port
	var err error

	log.Println("connecting to FileService:" + fileServerAddress)
	FileServiceGrpcClient, err = grpc.NewClient(fileServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("connection failed:", err)
	}

	FileService = file.NewFileServiceClient(FileServiceGrpcClient)

}

func ConnectToNotificationServer() {
	host := helpers.FatalIfEmptyVar("NOTIFICATION_SERVICE_HOST")
	port := helpers.FatalIfEmptyVar("NOTIFICATION_SERVICE_PORT")
	notifcationServerAddress := host + ":" + port
	var err error

	log.Println("connecting to NotificationService:" + notifcationServerAddress)
	NotificationServiceGrpcClient, err = grpc.NewClient(notifcationServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("connection failed:", err)
	}
	EmailService = emailPb.NewEmailServiceClient(NotificationServiceGrpcClient)
}

func DisconnectFromFileServer() {
	FileServiceGrpcClient.Close()
}

func DisconnectFromNotificationServer() {
	NotificationServiceGrpcClient.Close()
}
