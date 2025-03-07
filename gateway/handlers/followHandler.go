package handlers

import (
	"context"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/helpers"
	"net/http"

	"github.com/go-chi/chi"
)

////// follows

// follow user
//
//	@Summary		follow user
//	@Description	follow user
//	@Tags			follows
//	@Produce		application/json
//	@Security		ApiKeyAuth
//	@Param			username			path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		500					{string}	string	"server error"
//	@Router			/follows/{username}	[post]
func AddFollowing(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	username := chi.URLParam(r, "username")
	res, err := client.UserService.CreateFollow(context.Background(), &user_pb.FollowData{
		FollowerId:        currentUser.Id,
		FollowingUsername: username,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Following created.", 201)
		return
	}
	helpers.LogPanic("CreateFollow should return empty or httpError!!!")
}

// unfollow user
//
//	@Summary		unfollow user
//	@Description	unfollow user
//	@Tags			follows
//	@Produce		application/json
//	@Security		ApiKeyAuth
//	@Param			username			path		string	true	"url"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		500					{string}	string	"server error"
//	@Router			/follows/{username}	[delete]
func DeleteFollowing(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	username := chi.URLParam(r, "username")
	res, err := client.UserService.DeleteFollow(context.Background(), &user_pb.FollowData{
		FollowerId:        currentUser.Id,
		FollowingUsername: username,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Following deleted.", 200)
		return
	}
	helpers.LogPanic("DeleteFollow should return empty or httpError!!!")
}
