package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"
	"github.com/rzaf/youtube-clone/notification/helpers"
	"github.com/rzaf/youtube-clone/notification/models"
)

// get notification
//
//	@Summary		get notification
//	@Description	get notification
//	@Tags			notifications
//	@Accept			json
//	@Produce		application/json
//	@Security		ApiKeyAuth
//	@Param			id					path		string	true	"id"
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		404					{string}	string	"not found"
//	@Failure		500					{string}	string	"server error"
//	@Router			/notifications/{id}	[get]
func GetNotification(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	notificationId := chi.URLParam(r, "id")

	res, err := models.GetNotification(notificationId, currentUser.Id)
	if err != nil {
		helpers.LogPanic(err)
	}
	helpers.WriteJson(w, res, 200)

}

// get all notifications of current user
//
//	@Summary		get all notifications of current user
//	@Description	get all notifications of current user
//	@Tags			notifications
//	@Param			page	query	int		false	"page number"	default(1)
//	@Param			perpage	query	int		false	"items perpage"	default(10)
//	@Param			sort	query	string	false	"sort type"		default(newest)	Enums(newest, oldest)
//	@Param			type	query	string	false	"seen type"		default(any)	Enums(any, seen, not-seen)
//	@Accept			json
//	@Produce		application/json
//	@Security		ApiKeyAuth
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"request failed"
//	@Failure		404				{string}	string	"not found"
//	@Failure		500				{string}	string	"server error"
//	@Router			/notifications	[get]
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)

	var body map[string]any = make(map[string]any)
	helpers.ParseReq(r, body)

	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)

	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateNotificationSortTypes(sortTypeStr)

	seenTypeStr := helpers.ValidateStr(body["type"], "type", "any")
	seenType := helpers.ValidateNotificationSeenTypes(seenTypeStr)

	totalPages, data, err := models.GetAllNotification(currentUser.Id, perpage, page-1, sortType, seenType)
	if err != nil {
		helpers.LogPanic(err)
	}
	helpers.WriteJson(w, map[string]any{
		"total_pages":  totalPages,
		"current_page": data,
	}, 200)

}

// set specified notification as seen
//
//	@Summary		set specified notification as seen
//	@Description	set specified notification as seen
//	@Tags			notifications
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Param			id	path	string	true	"id"
//	@Security		ApiKeyAuth
//	@Success		200							{string}	string	"ok"
//	@Failure		400							{string}	string	"request failed"
//	@Failure		500							{string}	string	"server error"
//	@Router			/notifications/{id}/seen	[post]
func ReadNotification(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)
	notificationId := chi.URLParam(r, "id")

	err := models.ReadNotification(notificationId, currentUser.Id)
	if err != nil {
		helpers.LogPanic(err)
	}
	helpers.WriteJsonMessage(w, "Notification set as seen successfully", 200)
}

// set all notification of current user as seen
//
//	@Summary		set all notification of current user as seen
//	@Description	set all notification of current user as seen
//	@Tags			notifications
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Success		200					{string}	string	"ok"
//	@Failure		400					{string}	string	"request failed"
//	@Failure		500					{string}	string	"server error"
//	@Router			/notifications/seen																						[post]
func ReadAllNotifications(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)

	err := models.ReadAllNotification(currentUser.Id)
	if err != nil {
		helpers.LogPanic(err)
	}
	helpers.WriteJsonMessage(w, "All notifications set as seen successfully", 200)
}
