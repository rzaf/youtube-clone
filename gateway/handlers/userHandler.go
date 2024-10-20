package handlers

import (
	"context"
	"fmt"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/helpers"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		helpers.LogPanic(err)
	}
	return string(bytes)
}

func ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// get specified user
//
//	@Summary		get specified user
//	@Description	get user with specified username
//	@Tags			users
//	@Accept			json
//	@Produce		application/json
//	@Param			username			path		string	true	"username of user"
//	@Param			X-API-KEY			header		string	false	"optional authentication"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/users/{username}	[get]
func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var currentUserId int64 = 0
	if currentUser != nil {
		currentUserId = currentUser.Id
	}

	userName := chi.URLParam(r, "username")
	res, err := client.UserService.GetUserByUsername(context.Background(), &user_pb.UsernameAndId{
		Username:      userName,
		CurrentUserId: currentUserId,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	helpers.WriteProtoJson(w, res.GetUser(), true, 200)
}

// get users
//
//	@Summary		get users
//	@Description	get users
//	@Tags			users
//	@Produce		application/json
//	@Param			page		query		int		false	"page number"	default(1)
//	@Param			perpage		query		int		false	"items perpage"	default(10)
//	@Param			sort		query		string	false	"sort type"		default(newest)	Enums(newest, oldest,most-viewed,least-viewed,most-subbed,least-subbed)
//	@Param			X-API-KEY	header		string	false	"optional authentication"
//	@Success		200			{string}	string	"ok"
//	@Success		204			{string}	string	"no content"
//	@Failure		400			{string}	string	"request failed"
//	@Failure		500			{string}	string	"server error"
//	@Router			/users/																																																																																																																																																	[get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var currentUserId int64 = 0
	if currentUser != nil {
		currentUserId = currentUser.Id
	}
	body := make(map[string]any)
	helpers.ParseReq(r, body)

	helpers.ValidateAllowedParams(body, "perpage", "page", "sort")
	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 5)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateUsersSortTypes(sortTypeStr)

	res, err := client.UserService.GetUsers(context.Background(), &user_pb.UserReq{
		Page:          toPage(perpage, page),
		CurrentUserId: currentUserId,
		Sort:          sortType,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetUsers())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		return
	}
	helpers.WriteProtoJson(w, res.GetUsers(), true, 200)
}

// search users
//
//	@Summary		search users
//	@Description	search users
//	@Tags			users
//	@Produce		application/json
//	@Param			term					path		string	true	"search term"
//	@Param			page					query		int		false	"page number"	default(1)
//	@Param			perpage					query		int		false	"items perpage"	default(10)
//	@Param			sort					query		string	false	"sort type"		default(newest)	Enums(newest, oldest,most-viewed,least-viewed,most-subbed,least-subbed)
//	@Param			X-API-KEY				header		string	false	"optional authentication"
//	@Success		200						{string}	string	"ok"
//	@Success		204						{string}	string	"no content"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/users/search/{term}	[get]
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var currentUserId int64 = 0
	if currentUser != nil {
		currentUserId = currentUser.Id
	}

	term := chi.URLParam(r, "term")
	body := make(map[string]any)
	helpers.ParseReq(r, body)

	helpers.ValidateAllowedParams(body, "perpage", "page", "sort")
	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateUsersSortTypes(sortTypeStr)

	res, err := client.UserService.SearchUsers(context.Background(), &user_pb.UserReq{
		Page:          toPage(perpage, page),
		SearchTerm:    term,
		CurrentUserId: currentUserId,
		Sort:          sortType,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	fmt.Println(res.GetUsers())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		return
	}
	helpers.WriteProtoJson(w, res.GetUsers(), true, 200)
}

// sign up
//
//	@Summary		sign up
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
//	@Router			/users																																																																																																																																																	[post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("error type:%T,error:%v \n", err, err)
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	createdUser := res.GetUserApi()
	helpers.WriteJson(w, map[string]any{
		"Message":  "Validate your email to be able to use your api key",
		"Username": username,
		"ApiKey":   createdUser.Apikey,
	}, 201)
}

// sign in
//
//	@Summary		sign in
//	@Description	getting user api token
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Param			usernameOrEmail	formData	string	true	"usernmae or email"
//	@Param			password		formData	string	true	"password"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		500				{string}	string	"server error"
//	@Router			/users/sign-in	[post]
func SignIn(w http.ResponseWriter, r *http.Request) {
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
	helpers.WriteProtoJson(w, res.GetSignedInUser(), true, 200)
}

// resend verification email
//
//	@Summary		resend verification email
//	@Description	resend verification email
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		500					{string}	string	"server error"
//	@Router			/users/resend-email	[post]
func ResendEmailVerfication(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	if currentUser == nil {
		helpers.LogPanic(helpers.NewServerError("header X-API-KEY required", 401))
	}
	fmt.Printf("%v\n", currentUser)
	if currentUser.IsVerified {
		helpers.WriteJsonMessage(w, "email already verified!", 400)
		return
	}
	res, err := client.UserService.ResendEmailVerification(context.Background(), &user_pb.UsernameAndEmail{
		Username: currentUser.Username,
		Email:    currentUser.Email,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "verification email resent succesfully.", 200)
		return
	}
	helpers.LogPanic("ResendEmailVerification should return empty or httpError!!!")
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	verficationCode := chi.URLParam(r, "code")
	username := chi.URLParam(r, "username")
	helpers.ValidateRequiredStr(verficationCode, "verficationCode")

	res, err := client.UserService.VerifyUserEmail(context.Background(), &user_pb.EmailCode{
		Username: username,
		Code:     verficationCode,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "user email verfied successfully.", 200)
		return
	}
	helpers.LogPanic("VerifyUserEmail should return empty or httpError!!!")
}

// edit user profile photo
//
//	@Summary		edit user profile photo
//	@Description	edit user profile photo
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			photo_url						formData	string	true	"photo_url"
//	@Param			username						path		string	true	"username"
//	@Success		200								{string}	string	"ok"
//	@Failure		400								{string}	string	"request failed"
//	@Failure		401								{string}	string	"not authenticated"
//	@Failure		403								{string}	string	"not authorized"
//	@Failure		404								{string}	string	"not found"
//	@Failure		500								{string}	string	"server error"
//	@Router			/users/{username}/profile-photo	[put]
func SetProfilePhoto(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	userName := chi.URLParam(r, "username")
	fmt.Printf("%v\n", currentUser)
	if userName != currentUser.Username {
		helpers.WriteJsonError(w, "Not allowed to edit user with username:`"+userName+"`", 403)
		return
	}

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "photo_url")
	photoUrl := helpers.ValidateRequiredStr(body["photo_url"], "photo_url")

	if photoUrl == currentUser.ProfilePhoto {
		helpers.WriteJsonError(w, "photo_url is already as user profile photo", 400)
		return
	}
	res, err := client.UserService.SetUserPhoto(context.Background(), &user_pb.UserPhoto{
		Id:           currentUser.Id,
		ProfilePhoto: photoUrl,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "user profile photo edited.", 200)
		return
	}
	helpers.LogPanic("SetUserPhoto should return empty or httpError!!!")
}

// edit user channel photo
//
//	@Summary		edit user channel photo
//	@Description	edit user channel photo
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			photo_url						formData	string	true	"photo_url"
//	@Param			username						path		string	true	"username"
//	@Success		200								{string}	string	"ok"
//	@Failure		400								{string}	string	"request failed"
//	@Failure		401								{string}	string	"not authenticated"
//	@Failure		403								{string}	string	"not authorized"
//	@Failure		404								{string}	string	"not found"
//	@Failure		500								{string}	string	"server error"
//	@Router			/users/{username}/channel-photo	[put]
func SetChannelPhoto(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	userName := chi.URLParam(r, "username")
	fmt.Printf("%v\n", currentUser)
	if userName != currentUser.Username {
		helpers.WriteJsonError(w, "Not allowed to edit user with username:`"+userName+"`", 403)
		return
	}

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "photo_url")
	photoUrl := helpers.ValidateRequiredStr(body["photo_url"], "photo_url")

	if photoUrl == currentUser.ChannelPhoto {
		helpers.WriteJsonError(w, "photo_url is already as user cahnnel photo", 400)
		return
	}

	res, err := client.UserService.SetUserPhoto(context.Background(), &user_pb.UserPhoto{
		Id:           currentUser.Id,
		ChannelPhoto: photoUrl,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "user channel photo edited.", 200)
		return
	}
	helpers.LogPanic("SetUserPhoto should return empty or httpError!!!")
}

// edit user info
//
//	@Summary		edit user info
//	@Description	edit user info
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			password			formData	string	false	"password"
//	@Param			new_password		formData	string	false	"new_password"
//	@Param			new_aboutMe			formData	string	false	"new_aboutMe"
//	@Param			new_username		formData	string	false	"new_username"
//	@Param			new_channelName		formData	string	false	"new_channelName"
//	@Param			new_email			formData	string	false	"new_email"
//	@Param			username			path		string	true	"username"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		401					{string}	string	"not authenticated"
//	@Failure		403					{string}	string	"not authorized"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/users/{username}/	[put]
func EditUser(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "username")
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)

	// currentUser := GetAuthUser(r)
	fmt.Printf("%v\n", currentUser)
	if userName != currentUser.Username {
		helpers.WriteJsonError(w, "Not allowed to edit user with username:`"+userName+"`", 403)
		return
	}

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "password", "new_password", "new_aboutMe", "new_username", "new_channelName", "new_email")
	fmt.Println("request body:", body)
	fmt.Printf("user:%v\n", currentUser)

	password := helpers.ValidateRequiredStr(body["password"], "password")
	if err := ComparePassword(password, currentUser.HashedPassword); err != nil {
		fmt.Println("err:", err)
		helpers.WriteJsonError(w, "Incorrect password", 400)
		return
	}
	newPassword := helpers.ValidateStr(body["new_password"], "new_password", "")
	aboutMe := helpers.ValidateStr(body["new_aboutMe"], "new_aboutMe", "")
	username := helpers.ValidateStr(body["new_username"], "new_username", "")
	channelName := helpers.ValidateStr(body["new_channelName"], "new_channelName", "")
	email := helpers.ValidateStr(body["new_email"], "new_email", "")

	if newPassword == "" && aboutMe == "" && username == "" && channelName == "" && email == "" {
		helpers.WriteJsonError(w, "No new field", 400)
		return
	}

	if newPassword != "" {
		newPassword = HashPassword(newPassword)
	} else {
		newPassword = currentUser.HashedPassword
	}
	if aboutMe == "" {
		aboutMe = currentUser.AboutMe
	}
	if username == "" {
		username = currentUser.Username
	}
	if channelName == "" {
		channelName = currentUser.ChannelName
	}
	if email == "" {
		email = currentUser.Email
	}
	helpers.ValidateVar(email, "new_email", "email")

	res, err := client.UserService.EditUser(context.Background(), &user_pb.EditUserData{
		Id:             currentUser.Id,
		Email:          email,
		HashedPassword: newPassword,
		AboutMe:        aboutMe,
		Username:       username,
		ChannelName:    channelName,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res)
	helpers.WriteJsonMessage(w, fmt.Sprintf("user with username:`%s` edited", currentUser.Username), 200)
}

// setting new user api key
//
//	@Summary		setting new user api key
//	@Description	setting new user api key
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			password					formData	string	true	"password"
//	@Param			username					path		string	true	"username"
//	@Success		200							{string}	string	"ok"
//	@Failure		400							{string}	string	"request failed"
//	@Failure		401							{string}	string	"not authenticated"
//	@Failure		403							{string}	string	"not authorized"
//	@Failure		404							{string}	string	"not found"
//	@Failure		500							{string}	string	"server error"
//	@Router			/users/{username}/newApiKey	[put]
func NewUserApiKey(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	userName := chi.URLParam(r, "username")
	// currentUser := GetAuthUser(r)
	fmt.Printf("%v\n", currentUser)
	if userName != currentUser.Username {
		helpers.WriteJsonError(w, "Not allowed to edit user with username:`"+userName+"`", 403)
		return
	}

	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "password")
	fmt.Printf("user:%v\n", currentUser)
	password := helpers.ValidateRequiredStr(body["password"], "password")

	if err := ComparePassword(password, currentUser.HashedPassword); err != nil {
		fmt.Println("err:", err)
		helpers.WriteJsonError(w, "Incorrect password", 400)
		return
	}
	res, err := client.UserService.EditUserApiKey(context.Background(), &user_pb.UserId{Id: currentUser.Id})
	PanicIfIsError(err)
	apikey := res.GetUserApi()
	helpers.WriteJson(w, map[string]any{
		"message":     "New api key created",
		"New Api Key": apikey.Apikey,
	}, 200)
}

// deleting user
//
//	@Summary		deleting user
//	@Description	deleting user
//	@Tags			users
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			username			path		string	true	"username"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		401					{string}	string	"not authenticated"
//	@Failure		403					{string}	string	"not authorized"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/users/{username}	[delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	userName := chi.URLParam(r, "username")
	// currentUser := GetAuthUser(r)
	fmt.Printf("%v\n", currentUser)
	if userName != currentUser.Username {
		helpers.WriteJsonError(w, "Not allowed to delete user with username:`"+userName+"`", 403)
		return
	}
	res, err := client.UserService.DeleteUser(context.Background(), &user_pb.UserId{Id: currentUser.Id})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	helpers.WriteJsonMessage(w, fmt.Sprintf("user with username:`%s` deleted", currentUser.Username), 200)
}
