package services

import (
	// user_pb "youtube-clone/database/pbs/user-pb"

	"youtube-clone/database/pbs/comment"
	"youtube-clone/database/pbs/media"
	"youtube-clone/database/pbs/playlist"
	user_pb "youtube-clone/database/pbs/user-pb"

	"google.golang.org/grpc"
)

var (
	grpcServer *grpc.Server
)

func RegisterAllServices(server *grpc.Server) {
	grpcServer = server
	user_pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})
	media.RegisterMediaServiceServer(grpcServer, &mediaServiceServer{})
	comment.RegisterCommentServiceServer(grpcServer, &commentServiceServer{})
	playlist.RegisterPlaylistServiceServer(grpcServer, &playlistServiceServer{})
}
