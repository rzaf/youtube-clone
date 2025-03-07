package handlers

import (
	"context"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/helpers"
	"net/http"

	"github.com/go-chi/chi"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
)

////// media likes

// like/dislike media
//
//	@Summary		like/dislike media
//	@Description	like/dislike media
//	@Tags			likes
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			is_like				formData	bool	true	"is_like"
//	@Param			url					path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		500					{string}	string	"server error"
//	@Router			/medias/{url}/likes	[post]
func SetMediaLike(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "is_like")

	isLike := helpers.ValidateBool(body["is_like"], "is_like")

	res, err := client.MediaService.CreateLikeMedia(context.Background(), &helper.LikeReq{
		IsLike: isLike,
		Url:    url,
		UserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Like or dislike created.", 201)
		return
	}
	helpers.LogPanic("LikeMedia should return empty or httpError!!!")
}

// remove like/dislike media
//
//	@Summary		remove like/dislike media
//	@Description	remove like/dislike media
//	@Tags			likes
//	@Produce		application/json
//	@Security		ApiKeyAuth
//	@Param			url					path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		500					{string}	string	"server error"
//	@Router			/medias/{url}/likes	[delete]
func DeleteMediaLike(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	res, err := client.MediaService.DeleteLikeMedia(context.Background(), &helper.LikeReq{
		Url:    url,
		UserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Like or dislike deleted.", 200)
		return
	}
	helpers.LogPanic("DislikeMedia should return empty or httpError!!!")
}

////// comments likes

// like/dislike comment
//
//	@Summary		like/dislike comment
//	@Description	like/dislike comment
//	@Tags			likes
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			is_like					formData	bool	true	"is_like"
//	@Param			url						path		string	true	"url"
//	@Success		200						{string}	string	"ok"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/{url}/likes	[post]
func SetCommentLike(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "is_like")
	isLike := helpers.ValidateBool(body["is_like"], "is_like")

	res, err := client.CommentService.CreateLikeComment(context.Background(), &helper.LikeReq{
		IsLike: isLike,
		Url:    url,
		UserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Like or dislike created.", 201)
		return
	}
	helpers.LogPanic("CreateLikeComment should return empty or httpError!!!")
}

// like/dislike comment
//
//	@Summary		remove like/dislike comment
//	@Description	remove like/dislike comment
//	@Tags			likes
//	@Produce		application/json
//	@Security		ApiKeyAuth
//	@Param			url						path		string	true	"url"
//	@Success		200						{string}	string	"ok"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/{url}/likes	[delete]
func DeleteCommentLike(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	res, err := client.CommentService.DeleteLikeComment(context.Background(), &helper.LikeReq{
		Url:    url,
		UserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Like or dislike deleted.", 200)
		return
	}
	helpers.LogPanic("DeleteLikeComment should return empty or httpError!!!")
}
