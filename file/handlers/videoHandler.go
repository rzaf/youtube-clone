package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	pbHelper "youtube-clone/database/pbs/helper"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/file/helpers"
	"youtube-clone/file/models"
	"youtube-clone/file/queue"

	"github.com/go-chi/chi"
)

const (
	MaxVideoUploadSize = 50 << 20
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
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
