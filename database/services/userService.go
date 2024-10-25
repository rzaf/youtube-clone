package services

import (
	"context"
	"log"

	"fmt"
	"github.com/rzaf/youtube-clone/database/client"
	"github.com/rzaf/youtube-clone/database/models"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	fileModels "github.com/rzaf/youtube-clone/file/models"
	filePb "github.com/rzaf/youtube-clone/file/pbs/file"
)

type userServiceServer struct {
	user_pb.UserServiceServer
}

func newResponseFromUserData(u *user_pb.UserData) *user_pb.Response {
	return &user_pb.Response{
		Res: &user_pb.Response_User{
			User: u,
		},
	}
}

func newResponseFromCurrentUserData(u *user_pb.CurrentUserData) *user_pb.Response {
	return &user_pb.Response{
		Res: &user_pb.Response_AuthUser{
			AuthUser: u,
		},
	}
}

func newResponseFromSignedInUserData(u *user_pb.SignInUserData) *user_pb.Response {
	return &user_pb.Response{
		Res: &user_pb.Response_SignedInUser{
			SignedInUser: u,
		},
	}
}

func newResponseFromUsers(users []models.User, pageInfo *helper.PagesInfo) *user_pb.Response {
	if users == nil {
		return newUserResponseFromEmpty()
	}
	fmt.Println("newResponseFromUsers")
	log.Println(users)
	var users2 []*user_pb.UserShortData
	for _, u := range users {
		u2 := user_pb.UserShortData{
			Username:            u.Username,
			ChannelName:         u.ChannelName,
			CreatedAt:           u.Created_at.Unix(),
			ProfilePhoto:        u.ProfilePhoto,
			TotalViewCount:      u.TotalViews,
			SubscribersCount:    u.Subscribers,
			IsCurrentUserSubbed: u.IsCurrentUserSubbed,
		}
		users2 = append(users2, &u2)
	}
	return &user_pb.Response{
		Res: &user_pb.Response_Users{
			Users: &user_pb.UsersData{
				Users:     users2,
				PagesInfo: pageInfo,
			},
		},
	}
}

func newUserResponseFromEmpty() *user_pb.Response {
	return &user_pb.Response{
		Res: &user_pb.Response_Empty{},
	}
}

func newUserResponseFromError(e *helper.HttpError) *user_pb.Response {
	return &user_pb.Response{
		Res: &user_pb.Response_Err{
			Err: e,
		},
		// Res: &user_pb.Response_Empty{},
	}
}

func newResponseFromApikey(a *user_pb.UserApikey) *user_pb.Response {
	return &user_pb.Response{
		Res: &user_pb.Response_UserApi{
			UserApi: a,
		},
	}
}

func (s *userServiceServer) GetUserByUsername(c context.Context, u *user_pb.UsernameAndId) (*user_pb.Response, error) {
	var user *models.User
	var err error
	if u.CurrentUserId != 0 {
		user, err = models.AuthGetUserByUsername(u.Username, u.CurrentUserId)
	} else {
		user, err = models.GetUserByUsername(u.Username)
	}

	fmt.Println(user, err)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	u2 := &user_pb.UserData{
		Username:            user.Username,
		ChannelName:         user.ChannelName,
		ProfilePhoto:        user.ProfilePhoto,
		ChannelPhoto:        user.ChannelPhoto,
		Email:               user.Email,
		CreatedAt:           user.Created_at.Unix(),
		UploadCount:         user.UploadCount,
		SubscribersCount:    user.Subscribers,
		SubscringsCount:     user.Subscribings,
		AboutMe:             user.AboutMe,
		IsCurrentUserSubbed: user.IsCurrentUserSubbed,
		TotalViewCount:      user.TotalViews,
	}
	if user.Updated_at != nil {
		u2.UpdatedAt = user.Updated_at.Unix()
	}
	return newResponseFromUserData(u2), nil
}

func (s *userServiceServer) GetUserByNameAndPassword(c context.Context, u *user_pb.UsernameAndPassword) (*user_pb.Response, error) {
	user, err := models.GetUserByUsernameAndPassword(u.UserName, u.Password)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromSignedInUserData(&user_pb.SignInUserData{
		Username:     user.Username,
		Email:        user.Email,
		ChannelName:  user.ChannelName,
		AboutMe:      user.AboutMe,
		ApiKey:       user.ApiKey,
		IsVerified:   user.IsVerified,
		ProfilePhoto: user.ProfilePhoto,
		ChannelPhoto: user.ChannelPhoto,
	}), nil
}

// should be called only in authentication middleware
func (s *userServiceServer) GetUserByApikey(c context.Context, u *user_pb.UserApikey) (*user_pb.Response, error) {
	user, err := models.GetUserByApikey(u.Apikey)
	fmt.Println(user, err)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	// fmt.Printf("%v\n", user)
	return newResponseFromCurrentUserData(&user_pb.CurrentUserData{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		ChannelName:    user.ChannelName,
		AboutMe:        user.AboutMe,
		HashedPassword: user.HashedPassword,
		ProfilePhoto:   user.ProfilePhoto,
		ChannelPhoto:   user.ChannelPhoto,
		IsVerified:     user.IsVerified,
	}), nil
}

func (*userServiceServer) SearchUsers(c context.Context, u *user_pb.UserReq) (*user_pb.Response, error) {
	PerPage, PageNumber := getPage(u.Page)
	var err error
	var totalPages int64
	var users []models.User

	if u.CurrentUserId != 0 {
		totalPages, users, err = models.AuthSearchUsers(PerPage, PageNumber, u.SearchTerm, u.Sort, u.CurrentUserId)
	} else {
		totalPages, users, err = models.SearchUsers(PerPage, PageNumber, u.SearchTerm, u.Sort)
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromUsers(users, newPagesInfo(int32(totalPages), u.Page.PageNumber)), nil
}

func (*userServiceServer) GetUsers(c context.Context, u *user_pb.UserReq) (*user_pb.Response, error) {
	PerPage, PageNumber := getPage(u.Page)
	var err error
	var totalPages int64
	var users []models.User

	if u.CurrentUserId != 0 {
		totalPages, users, err = models.AuthGetUsers(PerPage, PageNumber, u.Sort, u.CurrentUserId)
	} else {
		totalPages, users, err = models.GetUsers(PerPage, PageNumber, u.Sort)
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromUsers(users, newPagesInfo(int32(totalPages), u.Page.PageNumber)), nil
}

func (s *userServiceServer) EditUser(c context.Context, u *user_pb.EditUserData) (*user_pb.Response, error) {
	err := models.EditUserInfo(u.Id, u.Email, u.Username, u.AboutMe, u.ChannelName, u.HashedPassword)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	return newUserResponseFromEmpty(), nil
}

func (s *userServiceServer) VerifyUserEmail(c context.Context, e *user_pb.EmailCode) (*user_pb.Response, error) {
	err := models.VerifyUserEmail(e.Username, e.Code)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	return newUserResponseFromEmpty(), nil
}

func (s *userServiceServer) ResendEmailVerification(c context.Context, e *user_pb.UsernameAndEmail) (*user_pb.Response, error) {
	code, err := models.SetEmailCode(e.Email)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	err = sendEmailVerificationNotification(e.Username, e.Email, code)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	return newUserResponseFromEmpty(), nil
}

func (s *userServiceServer) SetUserPhoto(c context.Context, u *user_pb.UserPhoto) (*user_pb.Response, error) {
	isProfile := u.ChannelPhoto == ""
	url := ""
	if isProfile {
		url = u.ProfilePhoto
	} else {
		url = u.ChannelPhoto
	}
	r, err := client.FileService.GetFileByUrl(context.Background(), &filePb.FileUrl{Url: url})
	if err != nil {
		return nil, err
	}
	if err2 := r.GetErr(); err2 != nil {
		return newUserResponseFromError(&helper.HttpError{
			Message:    err2.Message,
			StatusCode: err2.StatusCode,
		}), nil
	}
	file := r.GetFile()
	if file == nil {
		log.Println("fileService GetFileByUrl should return HttpError or FileData")
		return newUserResponseFromError(&helper.HttpError{
			Message:    "Something went wrong",
			StatusCode: 500,
		}), nil
	}

	if file.UserId != u.Id {
		return newUserResponseFromError(&helper.HttpError{
			Message: "file with url:`" + file.Url + "` do not belong to user",
			// Message:    "url not belong to user",
			StatusCode: 403,
		}), nil
	}
	if file.State == int32(fileModels.Removed) {
		return newUserResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + file.Url + "` is removed",
			StatusCode: 400,
		}), nil
	}
	if file.State == int32(fileModels.Pinned) {
		return newUserResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + file.Url + "` is not ready yet",
			StatusCode: 400,
		}), nil
	}
	if file.Type != int32(helper.MediaType_PHOTO) {
		return newUserResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + file.Url + "` is not photo",
			StatusCode: 400,
		}), nil
	}
	if file.Owner != int32(fileModels.None) {
		return newUserResponseFromError(&helper.HttpError{
			Message:    "file with url:`" + file.Url + "` is already used",
			StatusCode: 400,
		}), nil
	}
	owner := fileModels.ProfilePhoto
	if isProfile {
		err = models.SetUserProfilePhoto(u.Id, url)
	} else {
		err = models.SetUserChannelPhoto(u.Id, url)
		owner = fileModels.ChannelPhoto
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	_, err = client.FileService.SetFileOwner(context.Background(), &filePb.FileOwner{
		Owner: int32(owner),
		Url:   file.Url,
	})
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	return newUserResponseFromEmpty(), nil
}

func (s *userServiceServer) EditUserApiKey(c context.Context, u *user_pb.UserId) (*user_pb.Response, error) {
	user, err := models.EditUserApikey(u.Id)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	return newResponseFromApikey(&user_pb.UserApikey{Apikey: user.ApiKey}), nil
}

func (s *userServiceServer) DeleteUser(c context.Context, u *user_pb.UserId) (*user_pb.Response, error) {
	err := models.DeleteUser(u.Id)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
	}
	return newUserResponseFromEmpty(), nil
}

func (s *userServiceServer) CreateUser(c context.Context, u *user_pb.EditUserData) (*user_pb.Response, error) {
	fmt.Println(u.Username)
	user, err := models.CreateUser(u.Email, u.ChannelName, u.Username, u.HashedPassword, u.AboutMe)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	go sendEmailVerificationNotification(user.Username, user.Email, user.EmailVerification)
	return newResponseFromApikey(&user_pb.UserApikey{
		Apikey: user.ApiKey,
	}), nil
}

/// follow

func (s *userServiceServer) CreateFollow(c context.Context, f *user_pb.FollowData) (*user_pb.Response, error) {
	err := models.CreateFollowing(f.FollowerId, f.FollowingUsername)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	go followingNotification(f.FollowerId, f.FollowingUsername)
	return newUserResponseFromEmpty(), nil
}

func (s *userServiceServer) DeleteFollow(c context.Context, f *user_pb.FollowData) (*user_pb.Response, error) {
	err := models.DeleteFollowing(f.FollowerId, f.FollowingUsername)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newUserResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newUserResponseFromEmpty(), nil
}
