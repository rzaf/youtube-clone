package client

import (
	"github.com/rzaf/youtube-clone/database/helpers"
	"github.com/rzaf/youtube-clone/email/pbs/emailPb"
	"github.com/rzaf/youtube-clone/file/pbs/file"
	"github.com/rzaf/youtube-clone/notification/pbs/notificationPb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	FileServiceGrpcClient         *grpc.ClientConn
	EmailServiceGrpcClient        *grpc.ClientConn
	NotificationServiceGrpcClient *grpc.ClientConn

	FileService         file.FileServiceClient
	EmailService        emailPb.EmailServiceClient
	NotificationService notificationPb.NotificationServiceClient
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

func ConnectToEmailServer() {
	host := helpers.FatalIfEmptyVar("EMAIL_SERVICE_HOST")
	port := helpers.FatalIfEmptyVar("EMAIL_SERVICE_PORT")
	emailServerAddress := host + ":" + port
	var err error

	log.Println("connecting to EmailService:" + emailServerAddress)
	EmailServiceGrpcClient, err = grpc.NewClient(emailServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("connection failed:", err)
	}
	EmailService = emailPb.NewEmailServiceClient(EmailServiceGrpcClient)
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
	NotificationService = notificationPb.NewNotificationServiceClient(NotificationServiceGrpcClient)
}

func DisconnectFromFileServer() {
	FileServiceGrpcClient.Close()
}

func DisconnectFromEmailServer() {
	EmailServiceGrpcClient.Close()
}

func DisconnectFromNotificationServer() {
	NotificationServiceGrpcClient.Close()
}
