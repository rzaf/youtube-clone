package client

import (
	"log"
	"youtube-clone/database/pbs/comment"
	"youtube-clone/database/pbs/media"
	"youtube-clone/database/pbs/playlist"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/gateway/helpers"
	// "youtube-clone/notification/pbs/emailPb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GrpcClient *grpc.ClientConn
	// NotificationServiceGrpcClient *grpc.ClientConn

	serverAddress string

	UserService     user_pb.UserServiceClient
	MediaService    media.MediaServiceClient
	CommentService  comment.CommentServiceClient
	PlaylistService playlist.PlaylistServiceClient
	// EmailService    emailPb.EmailServiceClient
)

func ConnectToDatabaseServer() {
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
	MediaService = media.NewMediaServiceClient(GrpcClient)
	CommentService = comment.NewCommentServiceClient(GrpcClient)
	PlaylistService = playlist.NewPlaylistServiceClient(GrpcClient)

	log.Println(GrpcClient.GetState())
	GrpcClient.Connect()
	log.Printf("grpc state:%s\n", GrpcClient.GetState().String())
}

// func ConnectToNotificationServer() {
// 	host := helpers.FatalIfEmptyVar("NOTIFICATION_SERVICE_HOST")
// 	port := helpers.FatalIfEmptyVar("NOTIFICATION_SERVICE_PORT")
// 	notifcationServerAddress := host + ":" + port
// 	var err error

// 	log.Println("connecting to NotificationService:" + notifcationServerAddress)
// 	NotificationServiceGrpcClient, err = grpc.NewClient(notifcationServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalln("connection failed:", err)
// 	}
// 	EmailService = emailPb.NewEmailServiceClient(NotificationServiceGrpcClient)
// }

// func DisconnectFromNotificationServer() {
// 	NotificationServiceGrpcClient.Close()
// }

func DisconnectFromDatabaseServer() {
	GrpcClient.Close()
}
