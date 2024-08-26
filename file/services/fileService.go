package services

import (
	"context"
	"fmt"
	"youtube-clone/file/models"
	"youtube-clone/file/pbs/file"
)

type fileServiceServer struct {
	file.FileServiceServer
}

func newFileResponseFromEmpty() *file.Response {
	return &file.Response{
		Res: &file.Response_Empty{},
	}
}

func newFileResponseFromError(m string, s int32) *file.Response {
	return &file.Response{
		Res: &file.Response_Err{
			Err: &file.HttpError{
				Message:    m,
				StatusCode: s,
			},
		},
	}
}

func newResponseFromFileData(f *file.FileData) *file.Response {
	return &file.Response{
		Res: &file.Response_File{
			File: f,
		},
	}
}

func (*fileServiceServer) GetFileByUrl(c context.Context, f *file.FileUrl) (*file.Response, error) {
	url := models.GetUrl(f.Url)
	if url == nil {
		return newFileResponseFromError("file with url:`"+f.Url+"` not found", 404), nil
	}
	fmt.Println(url)
	return newResponseFromFileData(&file.FileData{
		Url:    url.Url,
		UserId: url.UserId,
		Size:   url.Size,
		Type:   int32(url.Type),
		Owner:  int32(url.Owner),
		State:  int32(url.State),
	}), nil
}

func (*fileServiceServer) SetFileOwner(c context.Context, f *file.FileOwner) (*file.Response, error) {
	err := models.SetUrlOwner(f.Url, models.MediaOwner(f.Owner))
	if err != nil {
		return newFileResponseFromError("url not found", 404), nil
	}
	return newFileResponseFromEmpty(), nil
}

func (*fileServiceServer) DeleteFile(context.Context, *file.FileUrl) (*file.Response, error) {
	return newFileResponseFromEmpty(), nil
}

func (*fileServiceServer) DeleteUserFiles(context.Context, *file.UserId) (*file.Response, error) {
	return newFileResponseFromEmpty(), nil
}
