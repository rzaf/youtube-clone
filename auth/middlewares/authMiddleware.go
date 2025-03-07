package middlewares

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"

	"net/http"
)

var SigningKey []byte

type AuthError struct {
	Message string
	Status  int
}

type AuthUser string

// should be called in handlers that dont have auth middleware
func GetUserFromHeader(r *http.Request) *user_pb.CurrentUserData {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return nil
	}
	token := ""
	if len(bearer) > 7 && bearer[0:6] == "Bearer" {
		token = bearer[7:]
	}

	err, user := parseJwt(token)
	if err != nil {
		panic(err)
	}
	return user
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			panic(AuthError{"Authorization header required", 401})
		}
		token := ""
		if len(bearer) > 7 && bearer[0:6] == "Bearer" {
			token = bearer[7:]
		}
		err, user := parseJwt(token)
		if err != nil {
			panic(err)
		}

		if !user.IsVerified {
			panic(&AuthError{"email is not verfied!", 401})
		}
		ctx := context.WithValue(r.Context(), AuthUser("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseJwt(token string) (*AuthError, *user_pb.CurrentUserData) {

	claims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	if err != nil || !jwtToken.Valid {
		log.Println(err)
		return &AuthError{"invalid token", 401}, nil
	}
	fmt.Printf("%v\n", claims)

	expiresAt := time.Unix(int64(claims["expire"].(float64)), 0)
	fmt.Printf("%v \n", expiresAt)
	if expiresAt.Before(time.Now().UTC()) {
		return &AuthError{"token is expired", 401}, nil
	}

	user := &user_pb.CurrentUserData{
		Id:           int64(claims["id"].(float64)),
		Email:        claims["email"].(string),
		Username:     claims["username"].(string),
		ChannelName:  claims["channel_name"].(string),
		IsVerified:   claims["is_verified"].(bool),
		ProfilePhoto: claims["profile_photo"].(string),
		ChannelPhoto: claims["channel_photo"].(string),
	}

	return nil, user
}
