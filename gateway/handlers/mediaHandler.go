package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"youtube-clone/database/pbs/helper"
	"youtube-clone/database/pbs/media"
	"youtube-clone/database/pbs/playlist"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/gateway/client"
	"youtube-clone/gateway/helpers"
)

//0:video 1:music 2:photo
// func convertTypeToStr(t media.MediaType) string {
// 	switch t {
// 	case 0:
// 		return "video"
// 	case 1:
// 		return "music"
// 	case 2:
// 		return "photo"
// 	}
// 	panic("incorrect media type")
// }

func toPage(perPage int, pageNumber int) *helper.Paging {
	return &helper.Paging{
		PerPage:    int32(perPage),
		PageNumber: int32(pageNumber),
	}
}

///////// MEDIA

func GetMediaByUrl(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var currentUserId int64 = 0
	if currentUser != nil {
		currentUserId = currentUser.Id
	}

	url := chi.URLParam(r, "url")
	res, err := client.MediaService.GetMediaByUrl(context.Background(), &media.MediaUrl{
		MediaUrl:      url,
		CurrentUserId: currentUserId,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	helpers.WriteProtoJson(w, res.GetMedia(), true, 200)
}

func SearchMedias(w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	var body map[string]any = make(map[string]any)
	helpers.ParseReq(r, body)

	helpers.ValidateAllowedParams(body, "type", "username", "perpage", "page", "sort")
	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	username := helpers.ValidateStr(body["username"], "username", "")
	mediaType := helpers.ValidateAllMediaType(helpers.ValidateRequiredStr(body["type"], "type"))
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateMediasSortTypes(sortTypeStr)

	res, err := client.MediaService.SearchMedias(context.Background(), &media.MediaReq{
		Type:       mediaType,
		Page:       toPage(perpage, page),
		SearchTerm: term,
		UserName:   username,
		Sort:       sortType,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetMedias())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No media found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetMedias(), true, 200)
}

func GetMedias(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]any)
	helpers.ParseReq(r, body)

	helpers.ValidateAllowedParams(body, "type", "username", "perpage", "page", "sort")
	mediaType := helpers.ValidateAllMediaType(helpers.ValidateRequiredStr(body["type"], "type"))
	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	username := helpers.ValidateStr(body["username"], "username", "")
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateMediasSortTypes(sortTypeStr)

	res, err := client.MediaService.GetMedias(context.Background(), &media.MediaReq{
		Type:     mediaType,
		Page:     toPage(perpage, page),
		UserName: username,
		Sort:     sortType,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetMedias())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No media found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetMedias(), true, 200)
}

func CreateMedia(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "type", "title", "description", "url")
	Type := helpers.ValidateMediaType(helpers.ValidateRequiredStr(body["type"], "type"))
	title := helpers.ValidateRequiredStr(body["title"], "title")
	url := helpers.ValidateRequiredStr(body["url"], "url")
	text := helpers.ValidateStr(body["description"], "description", "") ///optional

	res, err := client.MediaService.CreateMedia(context.Background(), &media.EidtMediaData{
		TypeId:        Type,
		Title:         title,
		Text:          text,
		Url:           url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media with url:`"+url+"` created.", 201)
		return
	}
	panic("CreateMedia should return empty or httpError!!!")
}

func EditMedia(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "title", "description")
	title := helpers.ValidateRequiredStr(body["title"], "title")
	text := helpers.ValidateRequiredStr(body["description"], "description")

	res, err := client.MediaService.EditMedia(context.Background(), &media.EidtMediaData{
		Title:         title,
		Text:          text,
		Url:           url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media with url:`"+url+"` edited.", 200)
		return
	}
	panic("EditMedia should return empty or httpError!!!")
}

func DeleteMedia(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")

	res, err := client.MediaService.DeleteMedia(context.Background(), &media.EidtMediaData{
		Url:           url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media with url:`"+url+"` deleted.", 200)
		return
	}
	panic("EditMedia should return empty or httpError!!!")
}

///////// TAG

func AddTagToVideo(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")
	tagName := chi.URLParam(r, "name")

	res, err := client.MediaService.AddTagToMedia(context.Background(), &media.TagMedia{
		TagName:       tagName,
		MediaUrl:      url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Tag added to media.", 201)
		return
	}
	panic("AddTagToMedia should return empty or httpError!!!")
}

func RemoveTagFromVideo(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")
	tagName := chi.URLParam(r, "name")

	res, err := client.MediaService.RemoveTagFromMedia(context.Background(), &media.TagMedia{
		TagName:       tagName,
		MediaUrl:      url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Tag removed from media.", 200)
		return
	}
	panic("RemoveTagFromMedia should return empty or httpError!!!")
}

///// videos of playlist

func GetMediasOfPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistUrl := chi.URLParam(r, "url")

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "perpage", "page")
	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)

	res, err := client.PlaylistService.GetPlaylistMedias(context.Background(), &playlist.PlaylistReq{
		Page:        toPage(perpage, page),
		PlaylistUrl: playlistUrl,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylists())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No media found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetMedias(), true, 200)
}

func AddMediaToPlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")
	playlistUrl := chi.URLParam(r, "playlistUrl")

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "note", "order")
	note := helpers.ValidateRequiredStr(body["note"], "note")
	order := helpers.ValidateRequiredPositiveInt(body["order"], "order")

	res, err := client.PlaylistService.AddMediaToPlaylist(context.Background(), &playlist.PlaylistMediaReq{
		Note:        note,
		Order:       int64(order),
		MediaUrl:    url,
		PlaylistUrl: playlistUrl,
		UserId:      currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media added to playlist.", 201)
		return
	}
	panic("AddMediaToPlaylist should return empty or httpError!!!")
}

func EditMediaFromPlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")
	playlistUrl := chi.URLParam(r, "playlistUrl")

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "new_note", "new_order")
	note := helpers.ValidateRequiredStr(body["new_note"], "new_note")
	order := helpers.ValidateRequiredPositiveInt(body["new_order"], "new_order")

	res, err := client.PlaylistService.EditMediaFromPlaylist(context.Background(), &playlist.PlaylistMediaReq{
		Note:        note,
		Order:       int64(order),
		MediaUrl:    url,
		PlaylistUrl: playlistUrl,
		UserId:      currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media from playlist edited.", 200)
		return
	}
	panic("EditMediaFromPlaylist should return empty or httpError!!!")
}

func DeleteMediaFromPlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	url := chi.URLParam(r, "url")
	playlistUrl := chi.URLParam(r, "playlistUrl")

	res, err := client.PlaylistService.RemoveMediaFromPlaylist(context.Background(), &playlist.PlaylistMediaReq{
		MediaUrl:    url,
		PlaylistUrl: playlistUrl,
		UserId:      currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media from playlist deleted.", 200)
		return
	}
	panic("RemoveMediaFromPlaylist should return empty or httpError!!!")
}
