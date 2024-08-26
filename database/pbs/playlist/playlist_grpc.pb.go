// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: youtube-clone/database/pbs/playlist.proto

package playlist

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PlaylistServiceClient is the client API for PlaylistService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlaylistServiceClient interface {
	GetPlaylist(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error)
	GetPlaylists(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error)
	SearchPlaylists(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error)
	CreatePlaylist(ctx context.Context, in *EditPlaylistData, opts ...grpc.CallOption) (*Response, error)
	EditPlaylist(ctx context.Context, in *EditPlaylistData, opts ...grpc.CallOption) (*Response, error)
	DeletePlaylist(ctx context.Context, in *EditPlaylistData, opts ...grpc.CallOption) (*Response, error)
	// // playlist medias
	GetPlaylistMedias(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error)
	AddMediaToPlaylist(ctx context.Context, in *PlaylistMediaReq, opts ...grpc.CallOption) (*Response, error)
	EditMediaFromPlaylist(ctx context.Context, in *PlaylistMediaReq, opts ...grpc.CallOption) (*Response, error)
	RemoveMediaFromPlaylist(ctx context.Context, in *PlaylistMediaReq, opts ...grpc.CallOption) (*Response, error)
}

type playlistServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlaylistServiceClient(cc grpc.ClientConnInterface) PlaylistServiceClient {
	return &playlistServiceClient{cc}
}

func (c *playlistServiceClient) GetPlaylist(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/GetPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) GetPlaylists(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/GetPlaylists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) SearchPlaylists(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/SearchPlaylists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) CreatePlaylist(ctx context.Context, in *EditPlaylistData, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/CreatePlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) EditPlaylist(ctx context.Context, in *EditPlaylistData, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/EditPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) DeletePlaylist(ctx context.Context, in *EditPlaylistData, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/DeletePlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) GetPlaylistMedias(ctx context.Context, in *PlaylistReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/GetPlaylistMedias", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) AddMediaToPlaylist(ctx context.Context, in *PlaylistMediaReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/AddMediaToPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) EditMediaFromPlaylist(ctx context.Context, in *PlaylistMediaReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/EditMediaFromPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playlistServiceClient) RemoveMediaFromPlaylist(ctx context.Context, in *PlaylistMediaReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/playlist.PlaylistService/RemoveMediaFromPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlaylistServiceServer is the server API for PlaylistService service.
// All implementations must embed UnimplementedPlaylistServiceServer
// for forward compatibility
type PlaylistServiceServer interface {
	GetPlaylist(context.Context, *PlaylistReq) (*Response, error)
	GetPlaylists(context.Context, *PlaylistReq) (*Response, error)
	SearchPlaylists(context.Context, *PlaylistReq) (*Response, error)
	CreatePlaylist(context.Context, *EditPlaylistData) (*Response, error)
	EditPlaylist(context.Context, *EditPlaylistData) (*Response, error)
	DeletePlaylist(context.Context, *EditPlaylistData) (*Response, error)
	// // playlist medias
	GetPlaylistMedias(context.Context, *PlaylistReq) (*Response, error)
	AddMediaToPlaylist(context.Context, *PlaylistMediaReq) (*Response, error)
	EditMediaFromPlaylist(context.Context, *PlaylistMediaReq) (*Response, error)
	RemoveMediaFromPlaylist(context.Context, *PlaylistMediaReq) (*Response, error)
	mustEmbedUnimplementedPlaylistServiceServer()
}

// UnimplementedPlaylistServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPlaylistServiceServer struct {
}

func (UnimplementedPlaylistServiceServer) GetPlaylist(context.Context, *PlaylistReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) GetPlaylists(context.Context, *PlaylistReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaylists not implemented")
}
func (UnimplementedPlaylistServiceServer) SearchPlaylists(context.Context, *PlaylistReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchPlaylists not implemented")
}
func (UnimplementedPlaylistServiceServer) CreatePlaylist(context.Context, *EditPlaylistData) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) EditPlaylist(context.Context, *EditPlaylistData) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditPlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) DeletePlaylist(context.Context, *EditPlaylistData) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) GetPlaylistMedias(context.Context, *PlaylistReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaylistMedias not implemented")
}
func (UnimplementedPlaylistServiceServer) AddMediaToPlaylist(context.Context, *PlaylistMediaReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMediaToPlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) EditMediaFromPlaylist(context.Context, *PlaylistMediaReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditMediaFromPlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) RemoveMediaFromPlaylist(context.Context, *PlaylistMediaReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMediaFromPlaylist not implemented")
}
func (UnimplementedPlaylistServiceServer) mustEmbedUnimplementedPlaylistServiceServer() {}

// UnsafePlaylistServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlaylistServiceServer will
// result in compilation errors.
type UnsafePlaylistServiceServer interface {
	mustEmbedUnimplementedPlaylistServiceServer()
}

func RegisterPlaylistServiceServer(s grpc.ServiceRegistrar, srv PlaylistServiceServer) {
	s.RegisterService(&PlaylistService_ServiceDesc, srv)
}

func _PlaylistService_GetPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).GetPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/GetPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).GetPlaylist(ctx, req.(*PlaylistReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_GetPlaylists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).GetPlaylists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/GetPlaylists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).GetPlaylists(ctx, req.(*PlaylistReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_SearchPlaylists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).SearchPlaylists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/SearchPlaylists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).SearchPlaylists(ctx, req.(*PlaylistReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_CreatePlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditPlaylistData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).CreatePlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/CreatePlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).CreatePlaylist(ctx, req.(*EditPlaylistData))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_EditPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditPlaylistData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).EditPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/EditPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).EditPlaylist(ctx, req.(*EditPlaylistData))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_DeletePlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditPlaylistData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).DeletePlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/DeletePlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).DeletePlaylist(ctx, req.(*EditPlaylistData))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_GetPlaylistMedias_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).GetPlaylistMedias(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/GetPlaylistMedias",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).GetPlaylistMedias(ctx, req.(*PlaylistReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_AddMediaToPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistMediaReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).AddMediaToPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/AddMediaToPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).AddMediaToPlaylist(ctx, req.(*PlaylistMediaReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_EditMediaFromPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistMediaReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).EditMediaFromPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/EditMediaFromPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).EditMediaFromPlaylist(ctx, req.(*PlaylistMediaReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaylistService_RemoveMediaFromPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaylistMediaReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaylistServiceServer).RemoveMediaFromPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/playlist.PlaylistService/RemoveMediaFromPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaylistServiceServer).RemoveMediaFromPlaylist(ctx, req.(*PlaylistMediaReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PlaylistService_ServiceDesc is the grpc.ServiceDesc for PlaylistService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlaylistService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "playlist.PlaylistService",
	HandlerType: (*PlaylistServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPlaylist",
			Handler:    _PlaylistService_GetPlaylist_Handler,
		},
		{
			MethodName: "GetPlaylists",
			Handler:    _PlaylistService_GetPlaylists_Handler,
		},
		{
			MethodName: "SearchPlaylists",
			Handler:    _PlaylistService_SearchPlaylists_Handler,
		},
		{
			MethodName: "CreatePlaylist",
			Handler:    _PlaylistService_CreatePlaylist_Handler,
		},
		{
			MethodName: "EditPlaylist",
			Handler:    _PlaylistService_EditPlaylist_Handler,
		},
		{
			MethodName: "DeletePlaylist",
			Handler:    _PlaylistService_DeletePlaylist_Handler,
		},
		{
			MethodName: "GetPlaylistMedias",
			Handler:    _PlaylistService_GetPlaylistMedias_Handler,
		},
		{
			MethodName: "AddMediaToPlaylist",
			Handler:    _PlaylistService_AddMediaToPlaylist_Handler,
		},
		{
			MethodName: "EditMediaFromPlaylist",
			Handler:    _PlaylistService_EditMediaFromPlaylist_Handler,
		},
		{
			MethodName: "RemoveMediaFromPlaylist",
			Handler:    _PlaylistService_RemoveMediaFromPlaylist_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "youtube-clone/database/pbs/playlist.proto",
}
