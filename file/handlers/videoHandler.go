package handlers

import (
	"fmt"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	pbHelper "github.com/rzaf/youtube-clone/database/pbs/helper"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/file/helpers"
	"github.com/rzaf/youtube-clone/file/models"
	"github.com/rzaf/youtube-clone/file/queue"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

const (
	MaxVideoUploadSize = 50 << 20
)

// upload video
//
//	@Summary		upload video
//	@Description	upload video
//	@Tags			videos
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		application/json
//	@Param			file			formData	file	true	"video file"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/videos/upload	[post]
func UploadVideo(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	err := r.ParseMultipartForm(MaxVideoUploadSize)
	if err != nil {
		panic(err)
	}
	mf := r.MultipartForm
	fmt.Printf("%+v\n", mf.Value)
	uploadedFile, headers, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	fmt.Printf("filename:%v\n", headers.Filename)
	fmt.Printf("Header:%v\n", headers.Header)
	fmt.Printf("Size:%v\n", headers.Size)

	var content []byte
	content, err = io.ReadAll(uploadedFile)
	if err != nil {
		panic(err)
	}

	contentType := http.DetectContentType(content)
	fmt.Printf("Type:%v\n", contentType)
	helpers.ValidateVideoType(contentType)
	helpers.CheckUserUploadBandwidth(headers.Size, currentUser.Id)

	url := getUniqueFileUrl()
	CreateAndWriteUrl(content, url, pbHelper.MediaType_VIDEO, headers.Size, currentUser.Id)
	queue.Push(queue.NewFileFormat(url, pbHelper.MediaType_VIDEO))
	helpers.WriteJson(w, map[string]any{
		"Message": "Video uploaded. Will be ready after a while.",
		"Url":     url,
		"Link":    os.Getenv("URL") + "/api/videos/" + url,
	}, 200)
}

// get video
//
//	@Summary		get video
//	@Description	get video
//	@Tags			videos
//	@Produce		application/x-mpegURL
//	@Param			url				path		string	true	"url"	1
//	@Success		200				{string}	string	"ok"
//	@Success		204				{string}	string	"no content"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/videos/{url}	[get]
func GetVideo(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	helpers.ValidateVideoUrl(url)
	urlM := models.GetUrl(url[:16])
	if urlM == nil || urlM.State == models.Removed {
		panic(helpers.NewServerError(fmt.Sprintf("url:'%s' not found", url), 404))
	}
	if urlM.Type != pbHelper.MediaType_VIDEO {
		panic(helpers.NewServerError(fmt.Sprintf("url:'%s' is not video", url), 400))
	}
	if urlM.State == models.Pinned {
		panic(helpers.NewServerError("video is not ready yet", 500))
	}
	file, err := os.Open("storage/videos/" + url)
	if err != nil {
		if os.IsNotExist(err) {
			panic("url exist in Db but not in storage")
		}
		panic(err)
	}
	fmt.Println(url)
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/x-mpegURL")
	w.Write(content)
}
