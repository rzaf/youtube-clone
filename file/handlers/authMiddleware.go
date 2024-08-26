package handlers

import (
	"context"
	"fmt"
	"net/http"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/file/client"
	"youtube-clone/file/helpers"
)

type authUser string

// func getUserFromHeader(r *http.Request) *user_pb.UserData {
// 	apiKey := r.Header.Get("X-API-KEY")
// 	if apiKey == "" {
// 		return nil
// 	}
// 	res, err := client.UserService.GetUserByApikey(context.Background(), &user_pb.UserApikey{Apikey: apiKey})
// 	if err != nil {
// 		panic(err)
// 	}
// 	if e := res.GetErr(); e != nil {
// 		panic(e)
// 		// panic(helpers.NewServerError("invalid api key", 401))
// 		// helpers.WriteJsonError(w, "invalid api key", 401)
// 		// return
// 	}
// 	return res.GetUser()
// }

// func GetAuthUser(r *http.Request) *user_pb.UserData {
// 	u := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
// 	if u != nil {
// 		panic("GetAuthUser should be called in authenticated routes")
// 	}
// 	return u
// }

func AuthApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			panic(helpers.NewServerError("header X-API-KEY required", 401))
			// helpers.WriteJsonError(w, "header X-API-KEY required", 401)
			// return
		}
		res, err := client.UserService.GetUserByApikey(context.Background(), &user_pb.UserApikey{Apikey: apiKey})
		if err != nil {
			panic(err)
		}
		if e := res.GetErr(); e != nil {
			panic(helpers.NewServerError("invalid api key", 401))
			// helpers.WriteJsonError(w, "invalid api key", 401)
			// return
		}
		user := res.GetAuthUser()
		fmt.Printf("%v\n", user)
		fmt.Printf("%v\n", user)
		if !user.IsVerified {
			panic(helpers.NewServerError("email is not verfied!", 401))
		}
		ctx := context.WithValue(r.Context(), authUser("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
