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
	// "github.com/h2non/bimg"
)

const (
	MaxPhotoUploadSize = 10 << 20
)

// upload photo
//
//	@Summary		upload photo
//	@Description	upload photo
//	@Tags			photos
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		application/json
//	@Param			file			formData	file	true	"photo file"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/photos/upload	[post]
func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	err := r.ParseMultipartForm(MaxPhotoUploadSize)
	if err != nil {
		panic(err)
	}
	mf := r.MultipartForm
	fmt.Printf("%+v\n", mf.Value)
	// file := mf.File["file"][0]
	uploadedFile, headers, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	helpers.CheckUserUploadBandwidth(headers.Size, currentUser.Id)

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
	helpers.ValidateImageType(contentType)

	// newContent, err := bimg.NewImage(content).Process(bimg.Options{Quality: 5})
	// if err != nil {
	// 	panic(err)
	// }

	url := getUniqueFileUrl()
	CreateAndWriteUrl(content, url, pbHelper.MediaType_PHOTO, headers.Size, currentUser.Id)
	queue.Push(queue.NewFileFormat(url, pbHelper.MediaType_PHOTO))
	helpers.WriteJson(w, map[string]any{
		"Message": "Image uploaded",
		"Url":     url,
		"Link":    os.Getenv("URL") + "/api/photos/" + url,
	}, 200)
}

// get photo
//
//	@Summary		get photo
//	@Description	get photo
//	@Tags			photos
//	@Produce		octet-stream
//	@Param			url				path		string	true	"url"	1
//	@Success		200				{string}	string	"ok"
//	@Success		204				{string}	string	"no content"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/photos/{url}	[get]
func GetPhoto(w http.ResponseWriter, r *http.Request) {
	str1, _ := os.Getwd()
	str2, _ := os.Executable()
	fmt.Println(str1, str2)

	url := chi.URLParam(r, "url")
	helpers.ValidateUrl(url)
	urlM := models.GetUrl(url[:16])
	if urlM == nil || urlM.State == models.Removed {
		panic(helpers.NewServerError(fmt.Sprintf("url:'%s' not found", url), 404))
	}
	if urlM.Type != pbHelper.MediaType_PHOTO {
		panic(helpers.NewServerError(fmt.Sprintf("url:'%s' is not photo", url), 400))
	}
	if urlM.State == models.Pinned {
		panic(helpers.NewServerError("photo is not ready yet", 500))
	}
	file, err := os.Open("storage/photos/" + url)
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
	// w.Header().Add("Content-Type", "image/jpeg")
	w.Write(content)
}
