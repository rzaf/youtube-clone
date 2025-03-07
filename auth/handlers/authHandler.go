package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rzaf/youtube-clone/auth/client"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/gateway/helpers"
)

var TokenExpireTime = time.Minute * 15
var SigningKey []byte

// login
//
//	@Summary		login
//	@Description	getting a new token
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Param			usernameOrEmail	formData	string	true	"usernmae or email"
//	@Param			password		formData	string	true	"password"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		500				{string}	string	"server error"
//	@Router			/login	[post]
func Login(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "usernameOrEmail", "password")

	username := helpers.ValidateRequiredStr(body["usernameOrEmail"], "usernameOrEmail")
	password := helpers.ValidateRequiredStr(body["password"], "password")

	res, err := client.UserService.GetUserByNameAndPassword(context.Background(), &user_pb.UsernameAndPassword{
		UserName: username,
		Password: password,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	user := res.GetAuthUser()

	jwtToken, err := generateJwtToken(user)
	if err != nil {
		helpers.WriteJsonError(w, "failed to create token", 500)
	}

	helpers.WriteJson(w, map[string]any{
		"token":         jwtToken,
		"refresh_token": user.RefreshToken,
	}, 200)
}

// refresh
//
//	@Summary		refresh token
//	@Description	getting a new token using refresh token
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Param			refresh_token	formData	string	true	"refresh_token"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		500				{string}	string	"server error"
//	@Router			/refresh	[post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "refresh_token")

	refreshToken := helpers.ValidateRequiredStr(body["refresh_token"], "refresh_token")

	res, err := client.UserService.GetUserByRefreshToken(context.Background(), &user_pb.UserRefreshToken{
		RefreshToken: refreshToken,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	user := res.GetAuthUser()

	jwtToken, err := generateJwtToken(user)
	if err != nil {
		helpers.WriteJsonError(w, "failed to create token", 500)
	}

	helpers.WriteJson(w, map[string]any{
		"token":         jwtToken,
		"refresh_token": user.RefreshToken,
	}, 200)
}

// Register
//
//	@Summary		Register
//	@Description	creating a user
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Param			email		formData	string	true	"email"
//	@Param			username	formData	string	true	"username"
//	@Param			channelName	formData	string	true	"channel name"
//	@Param			password	formData	string	true	"password"
//	@Param			aboutMe		formData	string	false	"about me"
//	@Success		200			{string}	string	"ok"
//	@Failure		400			{string}	string	"request failed"
//	@Failure		500			{string}	string	"server error"
//	@Router			/register																																																																																																																																																	[post]
func Register(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "email", "username", "channelName", "password", "aboutMe")

	email := helpers.ValidateRequiredStr(body["email"], "email")
	helpers.ValidateVar(email, "email", "email")
	username := helpers.ValidateRequiredStr(body["username"], "username")
	channelName := helpers.ValidateRequiredStr(body["channelName"], "channelName")
	password := helpers.ValidateRequiredStr(body["password"], "password")
	aboutMe := helpers.ValidateStr(body["aboutMe"], "aboutMe", "") //optional

	res, err := client.UserService.CreateUser(context.Background(), &user_pb.EditUserData{
		Email:          email,
		Username:       username,
		HashedPassword: password,
		ChannelName:    channelName,
		AboutMe:        aboutMe,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	createdUser := res.GetAuthUser()
	jwtToken, err := generateJwtToken(createdUser)
	if err != nil {
		helpers.WriteJsonError(w, "failed to create token", 500)
	}

	helpers.WriteJson(w, map[string]any{
		"Message":       "Validate your email to continue using app",
		"token":         jwtToken,
		"refresh_token": createdUser.RefreshToken,
	}, 201)
}

func generateJwtToken(u *user_pb.CurrentUserData) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":            u.Id,
			"email":         u.Email,
			"username":      u.Username,
			"channel_name":  u.ChannelName,
			"is_verified":   u.IsVerified,
			"profile_photo": u.ProfilePhoto,
			"channel_photo": u.ChannelPhoto,

			"expire": time.Now().UTC().Add(TokenExpireTime).Unix(),
		})

	token, err := jwtToken.SignedString(SigningKey)
	if err != nil {
		log.Println("jwtToken.SignedString failed:", err)
		return "", err
	}

	fmt.Printf("user jwt token created: %v \n", token)
	return token, nil
}
