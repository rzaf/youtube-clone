package handlers

import (
	"context"
	"fmt"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/database/pbs/playlist"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/helpers"
	"net/http"

	"github.com/go-chi/chi"
)

// search playlists
//
//	@Summary		search playlists
//	@Description	search playlists
//	@Tags			playlists
//	@Accept			json
//	@Produce		application/json
//	@Param			page						query		int		false	"page number"	default(1)
//	@Param			perpage						query		int		false	"items perpage"	default(10)
//	@Param			username					query		string	false	"playlist creator"
//	@Param			sort						query		string	false	"sort type"	default(newest)	Enums(newest, oldest,most-viewed,least-viewed)
//	@Param			term						path		string	true	"search term"
//	@Success		200							{string}	string	"ok"
//	@Success		204							{string}	string	"no content"
//	@Failure		400							{string}	string	"request failed"
//	@Failure		404							{string}	string	"not found"
//	@Failure		500							{string}	string	"server error"
//	@Router			/playlists/search/{term}	[get]
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
		helpers.LogPanic(err)
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

// get playlists
//
//	@Summary		get playlists
//	@Description	get playlists
//	@Tags			playlists
//	@Accept			json
//	@Produce		application/json
//	@Param			page		query		int		false	"page number"	default(1)
//	@Param			perpage		query		int		false	"items perpage"	default(10)
//	@Param			username	query		string	false	"playlist creator"
//	@Param			sort		query		string	false	"sort type"	default(newest)	Enums(newest, oldest,most-viewed,least-viewed)
//	@Success		200			{string}	string	"ok"
//	@Failure		400			{string}	string	"request failed"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"server error"
//	@Router			/playlists/	[get]
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
		helpers.LogPanic(err)
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

// get playlist
//
//	@Summary		get playlist
//	@Description	get playlist
//	@Tags			playlists
//	@Accept			json
//	@Produce		application/json
//	@Param			url					path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/playlists/{url}	[get]
func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	res, err := client.PlaylistService.GetPlaylist(context.Background(), &playlist.PlaylistReq{
		PlaylistUrl: url,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylists())
	helpers.WriteProtoJson(w, res.GetPlaylist(), true, 200)
}

// create a playlist
//
//	@Summary		create a playlist
//	@Description	create a playlist
//	@Tags			playlists
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			name		formData	string	true	"name"
//	@Param			text		formData	string	true	"text"
//	@Param			type		formData	string	true	"type"	Enums(photo, video, music, any)
//	@Success		200			{string}	string	"ok"
//	@Failure		400			{string}	string	"request failed"
//	@Failure		500			{string}	string	"server error"
//	@Router			/playlists	[post]
func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetPlaylist())
	p := res.GetPlaylist()
	if p == nil {
		helpers.LogPanic("CreatePlaylist should return error or playlistData")
	}
	helpers.WriteJson(w, map[string]string{
		"message": "Playlist created",
		"url":     p.Url,
	}, 201)
}

// edit playlist
//
//	@Summary		edit playlist
//	@Description	edit playlist
//	@Tags			playlists
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			name				formData	string	true	"name"
//	@Param			text				formData	string	true	"text"
//	@Param			url					path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		401					{string}	string	"not authenticated"
//	@Failure		403					{string}	string	"not authorized"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/playlists/{url}	[put]
func EditPlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "playlist Edited.", 200)
		return
	}
	helpers.LogPanic("EditPlaylist should return empty or httpError!!!")

}

// delete playlist
//
//	@Summary		delete playlist
//	@Description	delete playlist
//	@Tags			playlists
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			url					path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		401					{string}	string	"not authenticated"
//	@Failure		403					{string}	string	"not authorized"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/playlists/{url}	[delete]
func DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")

	res, err := client.PlaylistService.DeletePlaylist(context.Background(), &playlist.EditPlaylistData{
		Url:           url,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "playlist Deleted.", 200)
		return
	}
	helpers.LogPanic("DeletePlaylist should return empty or httpError!!!")
}
