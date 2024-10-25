package services

import (
	"context"
	"github.com/rzaf/youtube-clone/database/client"
	"github.com/rzaf/youtube-clone/database/models"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
	"github.com/rzaf/youtube-clone/database/pbs/media"
	fileModels "github.com/rzaf/youtube-clone/file/models"
	filePb "github.com/rzaf/youtube-clone/file/pbs/file"
	"log"
)

type mediaServiceServer struct {
	media.MediaServiceServer
}

func newResponseFromMediaData(m *media.MediaData) *media.Response {
	return &media.Response{
		Res: &media.Response_Media{
			Media: m,
		},
	}
}

func mediaTypeToStr(t helper.MediaType) string {
	switch t {
	case helper.MediaType_PHOTO:
		return "photo"
	case helper.MediaType_VIDEO:
		return "video"
	case helper.MediaType_MUSIC:
		return "music"
	case helper.MediaType_ALL:
		return "any"
	}
	log.Println("mediaTypeToStr should not return empty !!!")
	return ""
}

func newResponseFromMedias(medias []*models.Media, pageInfo *helper.PagesInfo) *media.Response {
	if medias == nil {
		return newMediaResponseFromEmpty()
	}
	var ms []*media.MediaShortData
	for _, m := range medias {
		med := &media.MediaShortData{
			Title:           m.Title,
			Url:             m.Url,
			CreatorUsername: m.UserName,
			ChannelName:     m.ChannelName,
			CreatorProfile:  m.UserProfile,
			ViewsCount:      m.ViewsCount,
			Thumbnail:       m.Thumbnail,
			CreatedAt:       m.CreatedAt.Unix(),
			UserLike:        m.CurrentUserLike.String(),
			MediaType:       mediaTypeToStr(m.Type),
		}
		ms = append(ms, med)
	}
	return &media.Response{
		Res: &media.Response_Medias{
			Medias: &media.MediasData{
				Medias:    ms,
				PagesInfo: pageInfo,
			},
		},
	}
}

func newMediaResponseFromError(e *helper.HttpError) *media.Response {
	return &media.Response{
		Res: &media.Response_Err{
			Err: e,
		},
	}
}

func newMediaResponseFromEmpty() *media.Response {
	return &media.Response{
		Res: &media.Response_Empty{},
	}
}

func (*mediaServiceServer) GetMediaByUrl(c context.Context, mu *media.MediaUrl) (*media.Response, error) {
	var m *models.Media
	var err error
	if mu.CurrentUserId == 0 {
		m, err = models.GetMediaByUrl(mu.MediaUrl)
	} else {
		m, err = models.AuthGetMediaByUrl(mu.MediaUrl, mu.CurrentUserId)
	}

	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	m2 := &media.MediaData{
		Title:               m.Title,
		Text:                m.Text,
		Url:                 m.Url,
		ChannelName:         m.ChannelName,
		UserSubscribesCount: m.UserSubscribesCount,
		CommentsCount:       m.CommentsCount,
		LikesCount:          m.LikesCount,
		DislikesCount:       m.DislikesCount,
		ViewsCount:          m.ViewsCount,
		Tags:                m.Tags,
		Thumbnail:           m.Thumbnail,
		Type:                mediaTypeToStr(m.Type),
		CreatorUsername:     m.UserName,
		CreatorProfile:      m.UserProfile,
		CreatedAt:           m.CreatedAt.Unix(),
		UserLike:            m.CurrentUserLike.String(),
		IsUserSubbed:        m.IsCurrentUserSubbedToCreator,
	}
	if m.UpdatedAt != nil {
		m2.UpdatedAt = m.UpdatedAt.Unix()
	}
	return newResponseFromMediaData(m2), nil
}

func (*mediaServiceServer) SearchMedias(c context.Context, m *media.MediaReq) (*media.Response, error) {
	PerPage, PageNumber := getPage(m.Page)
	var err error
	var totalPages int64
	var medias []*models.Media

	if m.UserName != "" {
		totalPages, medias, err = models.SearchUserMedias(PerPage, PageNumber, m.SearchTerm, m.UserName, m.Type, m.Sort)
	} else {
		totalPages, medias, err = models.SearchMedias(PerPage, PageNumber, m.SearchTerm, m.Type, m.Sort)
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromMedias(medias, newPagesInfo(int32(totalPages), m.Page.PageNumber)), nil
}

func (*mediaServiceServer) GetMedias(c context.Context, m *media.MediaReq) (*media.Response, error) {
	PerPage, PageNumber := getPage(m.Page)
	var err error
	var totalPages int64
	var medias []*models.Media

	if m.UserName != "" {
		totalPages, medias, err = models.GetUserMedias(PerPage, PageNumber, m.UserName, m.Type, m.Sort)
	} else {
		totalPages, medias, err = models.GetMedias(PerPage, PageNumber, m.Type, m.Sort)
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromMedias(medias, newPagesInfo(int32(totalPages), m.Page.PageNumber)), nil
}

func (*mediaServiceServer) CreateMedia(c context.Context, m *media.EidtMediaData) (*media.Response, error) {
	r, err := client.FileService.GetFileByUrl(context.Background(), &filePb.FileUrl{Url: m.Url})
	if err != nil {
		return nil, err
	}
	if err2 := r.GetErr(); err2 != nil {
		return newMediaResponseFromError(&helper.HttpError{
			Message:    err2.Message,
			StatusCode: err2.StatusCode,
		}), nil
	}
	file := r.GetFile()
	if file == nil {
		log.Println("fileService GetFileByUrl should return HttpError or FileData")
		return newMediaResponseFromError(&helper.HttpError{
			Message:    "Something went wrong",
			StatusCode: 500,
		}), nil
	}
	log.Println(m)
	log.Println(file)
	if file.UserId != m.CurrentUserId {
		return newMediaResponseFromError(&helper.HttpError{
			Message:    "url not belong to user",
			StatusCode: 403,
		}), nil
	}
	if file.State == int32(fileModels.Removed) {
		return newMediaResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + m.Url + "` is removed",
			StatusCode: 400,
		}), nil
	}
	if file.State == int32(fileModels.Pinned) {
		return newMediaResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + m.Url + "` is not ready yet",
			StatusCode: 400,
		}), nil
	}
	if file.Owner != int32(fileModels.None) {
		return newMediaResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + m.Url + "` is already used",
			StatusCode: 400,
		}), nil
	}
	if file.Type != int32(m.TypeId) {
		return newMediaResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + m.Url + "` is not of type " + m.TypeId.String(),
			StatusCode: 400,
		}), nil
	}
	err = models.CreateMedia(m.Title, m.Text, m.Url, m.CurrentUserId, m.TypeId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}

	var owner fileModels.MediaOwner
	switch m.TypeId {
	case helper.MediaType_PHOTO:
		owner = fileModels.PhotoMedia
	case helper.MediaType_MUSIC:
		owner = fileModels.MusicMedia
	case helper.MediaType_VIDEO:
		owner = fileModels.VideoMedia
	}
	_, err = client.FileService.SetFileOwner(context.Background(), &filePb.FileOwner{
		Owner: int32(owner),
		Url:   m.Url,
	})
	if err != nil {
		return nil, err
	}
	go newMediaNotification(m.CurrentUserId, m.Url, m.Title, mediaTypeToStr(m.TypeId), m.Text)
	return newMediaResponseFromEmpty(), nil
}

func (*mediaServiceServer) EditMedia(c context.Context, m *media.EidtMediaData) (*media.Response, error) {
	err := models.EditMedia(m.Url, m.Title, m.Text, m.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newMediaResponseFromEmpty(), nil
}

// TODO
func (*mediaServiceServer) SetMediaThumbnail(c context.Context, m *media.EidtMediaData) (*media.Response, error) {
	return nil, nil
}

func (*mediaServiceServer) DeleteMedia(c context.Context, m *media.EidtMediaData) (*media.Response, error) {
	err := models.DeleteMedia(m.Url, m.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newMediaResponseFromEmpty(), nil
}

////// TAGS

func (*mediaServiceServer) AddTagToMedia(c context.Context, m *media.TagMedia) (*media.Response, error) {
	models.CreateTag(m.TagName) // assuring tag exist
	err := models.AddTagToMedia(m.MediaUrl, m.TagName, m.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newMediaResponseFromEmpty(), nil
}

func (*mediaServiceServer) RemoveTagFromMedia(c context.Context, m *media.TagMedia) (*media.Response, error) {
	err := models.RemoveTagFromMedia(m.MediaUrl, m.TagName, m.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newMediaResponseFromEmpty(), nil
}

/////////LIKES

func (*mediaServiceServer) CreateLikeMedia(con context.Context, l *helper.LikeReq) (*media.Response, error) {
	err := models.CreateMediaLike(l.UserId, l.Url, l.IsLike)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	go newMediaLikeNotification(l.UserId, l.Url, l.IsLike)
	return newMediaResponseFromEmpty(), nil
}

func (*mediaServiceServer) DeleteLikeMedia(con context.Context, l *helper.LikeReq) (*media.Response, error) {
	err := models.DeleteMediaLike(l.UserId, l.Url)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newMediaResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newMediaResponseFromEmpty(), nil
}
