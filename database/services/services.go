package services

import (
	// user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"

	"github.com/rzaf/youtube-clone/database/pbs/comment"
	"github.com/rzaf/youtube-clone/database/pbs/media"
	"github.com/rzaf/youtube-clone/database/pbs/playlist"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"

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
