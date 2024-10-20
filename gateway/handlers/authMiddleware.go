package handlers

import (
	"context"
	"fmt"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/helpers"
	"net/http"
)

type authUser string

// should be called in handlers that dont have auth middleware
func getUserFromHeader(r *http.Request) *user_pb.CurrentUserData {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		return nil
	}
	res, err := client.UserService.GetUserByApikey(context.Background(), &user_pb.UserApikey{Apikey: apiKey})
	if err != nil {
		helpers.LogPanic(err)
	}
	if e := res.GetErr(); e != nil {
		helpers.LogPanic(e)
	}
	return res.GetAuthUser()
}

func AuthApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			panic(helpers.NewServerError("header X-API-KEY required", 401))
		}
		res, err := client.UserService.GetUserByApikey(context.Background(), &user_pb.UserApikey{Apikey: apiKey})
		if err != nil {
			helpers.LogPanic(err)
		}
		if e := res.GetErr(); e != nil {
			helpers.LogPanic(helpers.NewServerError("invalid api key", 401))
		}
		user := res.GetAuthUser()
		fmt.Printf("%v\n", user)
		if !user.IsVerified {
			helpers.LogPanic(helpers.NewServerError("email is not verfied!", 401))
		}
		ctx := context.WithValue(r.Context(), authUser("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
