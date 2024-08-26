package handlers

import (
	"context"
	"fmt"
	"net/http"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/gateway/client"
	"youtube-clone/gateway/helpers"
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
		panic(err)
	}
	if e := res.GetErr(); e != nil {
		panic(e)
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
			panic(err)
		}
		if e := res.GetErr(); e != nil {
			panic(helpers.NewServerError("invalid api key", 401))
		}
		user := res.GetAuthUser()
		fmt.Printf("%v\n", user)
		if !user.IsVerified {
			panic(helpers.NewServerError("email is not verfied!", 401))
		}
		ctx := context.WithValue(r.Context(), authUser("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
