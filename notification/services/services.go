package services

import (
	"github.com/rzaf/youtube-clone/notification/pbs/notificationPb"

	"google.golang.org/grpc"
)

var (
	grpcServer *grpc.Server
)

func RegisterAllServices(server *grpc.Server) {
	grpcServer = server

	notificationPb.RegisterNotificationServiceServer(grpcServer, &notificationServiceServer{})
}

// func newHttpError(m string, s int) *helper.HttpError {
// 	return &helper.HttpError{
// 		Message:    m,
// 		StatusCode: int32(s),
// 	}
// }
