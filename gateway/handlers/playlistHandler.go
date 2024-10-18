package handlers

import (
	"context"
	"fmt"
	"net/http"
	"youtube-clone/database/pbs/playlist"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/gateway/client"
	"youtube-clone/gateway/helpers"

	"github.com/go-chi/chi"
)

////// playlist

func SearchPlaylists(w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "sort", "username", "perpage", "page")

	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	username := helpers.ValidateStr(body["username"], "username", "")
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidatePlaylistsSortTypes(sortTypeStr)

	res, err := client.PlaylistService.SearchPlaylists(context.Background(), &playlist.PlaylistReq{
		Page:       toPage(perpage, page),
		Username:   username,
		Sort:       sortType,
		SearchTerm: term,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylists())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No playlist found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetPlaylists(), true, 200)
}

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "sort", "username", "perpage", "page")

	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	username := helpers.ValidateStr(body["username"], "username", "")
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidatePlaylistsSortTypes(sortTypeStr)

	res, err := client.PlaylistService.GetPlaylists(context.Background(), &playlist.PlaylistReq{
		Page:     toPage(perpage, page),
		Username: username,
		Sort:     sortType,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylists())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No playlist found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetPlaylists(), true, 200)
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	res, err := client.PlaylistService.GetPlaylist(context.Background(), &playlist.PlaylistReq{
		PlaylistUrl: url,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylists())
	helpers.WriteProtoJson(w, res.GetPlaylist(), true, 200)
}

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "name", "text", "type")

	name := helpers.ValidateRequiredStr(body["name"], "name")
	text := helpers.ValidateRequiredStr(body["text"], "text")
	mediaTypeStr := helpers.ValidateRequiredStr(body["type"], "type")
	mediaType := helpers.ValidateAllMediaType(mediaTypeStr)

	res, err := client.PlaylistService.CreatePlaylist(context.Background(), &playlist.EditPlaylistData{
		Name:          name,
		MediaTypeId:   mediaType,
		Text:          text,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylist())
	p := res.GetPlaylist()
	if p == nil {
		panic("CreatePlaylist should return error or playlistData")
	}
	helpers.WriteJson(w, map[string]string{
		"message": "Playlist created",
		"url":     p.Url,
	}, 201)
}

func EditPlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "name", "text")

	name := helpers.ValidateRequiredStr(body["name"], "name")
	text := helpers.ValidateRequiredStr(body["text"], "text")

	res, err := client.PlaylistService.EditPlaylist(context.Background(), &playlist.EditPlaylistData{
		Text:          text,
		Name:          name,
		Url:           url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "playlist Edited.", 200)
		return
	}
	panic("EditPlaylist should return empty or httpError!!!")

}

func DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")

	res, err := client.PlaylistService.DeletePlaylist(context.Background(), &playlist.EditPlaylistData{
		Url:           url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "playlist Deleted.", 200)
		return
	}
	panic("DeletePlaylist should return empty or httpError!!!")
}
