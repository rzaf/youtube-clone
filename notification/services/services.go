package services

import (
	"youtube-clone/notification/pbs/emailPb"

	"google.golang.org/grpc"
)

var (
	grpcServer *grpc.Server
)

func RegisterAllServices(server *grpc.Server) {
	grpcServer = server

	emailPb.RegisterEmailServiceServer(grpcServer, &emailServiceServer{})
}

// func newHttpError(m string, s int) *helper.HttpError {
// 	return &helper.HttpError{
// 		Message:    m,
// 		StatusCode: int32(s),
// 	}
// }
