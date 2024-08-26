package handlers

import (
	"context"
	"net/http"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/gateway/client"
	"youtube-clone/gateway/helpers"

	"github.com/go-chi/chi"
)

////// follows

func AddFollowing(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	username := chi.URLParam(r, "username")
	res, err := client.UserService.CreateFollow(context.Background(), &user_pb.FollowData{
		FollowerId:        currentUser.Id,
		FollowingUsername: username,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Following created.", 201)
		return
	}
	panic("CreateFollow should return empty or httpError!!!")
}

func DeleteFollowing(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	username := chi.URLParam(r, "username")
	res, err := client.UserService.DeleteFollow(context.Background(), &user_pb.FollowData{
		FollowerId:        currentUser.Id,
		FollowingUsername: username,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Following deleted.", 200)
		return
	}
	panic("DeleteFollow should return empty or httpError!!!")
}
