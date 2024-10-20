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

// get media
//
//	@Summary		get media
//	@Description	get media
//	@Tags			medias
//	@Accept			json
//	@Produce		application/json
//	@Param			url				path		string	true	"url"
//	@Param			X-API-KEY		header		string	false	"optional authentication"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/medias/{url}	[get]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	helpers.WriteProtoJson(w, res.GetMedia(), true, 200)
}

// search media
//
//	@Summary		search medias
//	@Description	search medias
//	@Tags			medias
//	@Produce		application/json
//	@Param			term					path		string	true	"search term"
//	@Param			page					query		int		false	"page number"	default(1)
//	@Param			perpage					query		int		false	"items perpage"	default(10)
//	@Param			username				query		string	false	"media creator"
//	@Param			type					query		string	false	"media type"	default(video)	Enums(photo, video, music, any)
//	@Param			sort					query		string	false	"sort type"		default(newest)	Enums(newest, oldest,most-viewed,least-viewed)
//	@Success		200						{string}	string	"ok"
//	@Success		204						{string}	string	"no content"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/medias/search/{term}	[get]
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
		helpers.LogPanic(err)
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

// get medias
//
//	@Summary		get medias
//	@Description	get medias
//	@Tags			medias
//	@Produce		application/json
//	@Param			page		query		int		false	"page number"	default(1)
//	@Param			perpage		query		int		false	"items perpage"	default(10)
//	@Param			username	query		string	false	"media creator"
//	@Param			type		query		string	false	"media type"	default(video)	Enums(photo, video, music, any)
//	@Param			sort		query		string	false	"sort type"		default(newest)	Enums(newest, oldest,most-viewed,least-viewed)
//	@Success		200			{string}	string	"ok"
//	@Success		204			{string}	string	"no content"
//	@Failure		400			{string}	string	"request failed"
//	@Failure		500			{string}	string	"server error"
//	@Router			/medias																																																																																																																																					[get]
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
		helpers.LogPanic(err)
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

// create media
//
//	@Summary		create media
//	@Description	create media
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			type		formData	string	true	"type"	Enums(photo, video, music)
//	@Param			title		formData	string	true	"title"
//	@Param			url			formData	string	true	"url"
//	@Param			description	formData	string	false	"description"
//	@Success		200			{string}	string	"ok"
//	@Failure		400			{string}	string	"request failed"
//	@Failure		500			{string}	string	"server error"
//	@Router			/medias																																																																																																																																					[post]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media with url:`"+url+"` created.", 201)
		return
	}
	helpers.LogPanic("CreateMedia should return empty or httpError!!!")
}

// edit media
//
//	@Summary		edit media
//	@Description	edit media
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			title			formData	string	true	"title"
//	@Param			description		formData	string	true	"description"
//	@Param			url				path		string	true	"url"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		401				{string}	string	"not authenticated"
//	@Failure		403				{string}	string	"not authorized"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/medias/{url}	[put]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media with url:`"+url+"` edited.", 200)
		return
	}
	helpers.LogPanic("EditMedia should return empty or httpError!!!")
}

// deleting media
//
//	@Summary		deleting media
//	@Description	deleting media
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			url				path		string	true	"url"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		401				{string}	string	"not authenticated"
//	@Failure		403				{string}	string	"not authorized"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/medias/{url}	[delete]
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

// add tag to media
//
//	@Summary		add tag to media
//	@Description	add tag to media
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			url							path		string	true	"url"
//	@Param			name						path		string	true	"name"
//	@Success		200							{string}	string	"ok"
//	@Failure		400							{string}	string	"request failed"
//	@Failure		500							{string}	string	"server error"
//	@Router			/medias/{url}/tag/{name}	[post]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Tag added to media.", 201)
		return
	}
	helpers.LogPanic("AddTagToMedia should return empty or httpError!!!")
}

// remove tag from media
//
//	@Summary		remove tag from media
//	@Description	remove tag from media
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			url							path		string	true	"url"
//	@Param			name						path		string	true	"name"
//	@Success		200							{string}	string	"ok"
//	@Failure		400							{string}	string	"request failed"
//	@Failure		500							{string}	string	"server error"
//	@Router			/medias/{url}/tag/{name}	[delete]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Tag removed from media.", 200)
		return
	}
	helpers.LogPanic("RemoveTagFromMedia should return empty or httpError!!!")
}

///// playlist medias

// get medias of playlist
//
//	@Summary		get medias of playlist
//	@Description	get medias of playlist
//	@Tags			medias
//	@Produce		application/json
//	@Param			url						path		string	true	"url"
//	@Param			page					query		int		false	"page number"	default(1)
//	@Param			perpage					query		int		false	"items perpage"	default(10)
//	@Success		200						{string}	string	"ok"
//	@Success		204						{string}	string	"no content"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/playlists/{url}/medias	[get]
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
		helpers.LogPanic(err)
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

// add media to playlist
//
//	@Summary		add media to playlist
//	@Description	add media to playlist
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			note									formData	string	true	"note"
//	@Param			order									formData	int		true	"order"	default(1)
//	@Param			url										path		string	true	"url"
//	@Param			playlistUrl								path		string	true	"playlistUrl"
//	@Success		200										{string}	string	"ok"
//	@Failure		400										{string}	string	"request failed"
//	@Failure		500										{string}	string	"server error"
//	@Router			/medias/{url}/playlists/{playlistUrl}	[post]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media added to playlist.", 201)
		return
	}
	helpers.LogPanic("AddMediaToPlaylist should return empty or httpError!!!")
}

// edit media from playlist
//
//	@Summary		edit media from playlist
//	@Description	edit media from playlist
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			new_note								formData	string	true	"new_note"
//	@Param			new_note								formData	int		true	"new_note"	default(1)
//	@Param			url										path		string	true	"url"
//	@Param			playlistUrl								path		string	true	"playlistUrl"
//	@Success		200										{string}	string	"ok"
//	@Failure		400										{string}	string	"request failed"
//	@Failure		500										{string}	string	"server error"
//	@Router			/medias/{url}/playlists/{playlistUrl}	[put]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media from playlist edited.", 200)
		return
	}
	helpers.LogPanic("EditMediaFromPlaylist should return empty or httpError!!!")
}

// delete media from playlist
//
//	@Summary		delete media from playlist
//	@Description	delete media from playlist
//	@Tags			medias
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			url										path		string	true	"url"
//	@Param			playlistUrl								path		string	true	"playlistUrl"
//	@Success		200										{string}	string	"ok"
//	@Failure		400										{string}	string	"request failed"
//	@Failure		500										{string}	string	"server error"
//	@Router			/medias/{url}/playlists/{playlistUrl}	[delete]
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
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Media from playlist deleted.", 200)
		return
	}
	helpers.LogPanic("RemoveMediaFromPlaylist should return empty or httpError!!!")
}
