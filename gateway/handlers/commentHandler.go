package handlers

import (
	"context"
	"fmt"
	"net/http"
	"youtube-clone/database/pbs/comment"
	user_pb "youtube-clone/database/pbs/user-pb"
	"youtube-clone/gateway/client"
	"youtube-clone/gateway/helpers"

	"github.com/go-chi/chi"
)

// get comment
//
//	@Summary		get comment
//	@Description	get comment
//	@Tags			comments
//	@Accept			json
//	@Produce		application/json
//	@Param			commentUrl				path		string	true	"commentUrl"
//	@Param			X-API-KEY				header		string	false	"optional authentication"
//	@Success		200						{string}	string	"ok"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		404						{string}	string	"not found"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/{commentUrl}	[get]
func GetComment(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var userId int64 = 0
	if currentUser != nil {
		userId = currentUser.Id
	}
	commentUrl := chi.URLParam(r, "commentUrl")

	res, err := client.CommentService.GetCommentByUrl(context.Background(), &comment.CommentUrl{
		Url:    commentUrl,
		UserId: userId,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	helpers.WriteProtoJson(w, res.GetFullComment(), true, 200)
}

// get comments of media
//
//	@Summary		get comments of media
//	@Description	get comments of media
//	@Tags			comments
//	@Produce		application/json
//	@Param			page					query		int		false	"page number"	default(1)
//	@Param			perpage					query		int		false	"items perpage"	default(10)
//	@Param			sort					query		string	false	"sort type"		default(newest)	Enums(newest, oldest,most-liked,least-liked,most-disliked,least-disliked,most-replied,least-replied)
//	@Param			url						path		string	false	"url"
//	@Param			X-API-KEY				header		string	false	"optional authentication"
//	@Success		200						{string}	string	"ok"
//	@Success		204						{string}	string	"no content"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/medias/{url}	[get]
func GetCommentsOfMedia(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var userId int64 = 0
	if currentUser != nil {
		userId = currentUser.Id
	}
	url := chi.URLParam(r, "url")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "perpage", "page", "sort")

	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateCommentsSortTypes(sortTypeStr)

	fmt.Println(perpage, page)
	res, err := client.CommentService.GetCommentsOfMedia(context.Background(), &comment.CommentReq{
		MediaUrl: url,
		Page:     toPage(perpage, page),
		UserId:   userId,
		Sort:     sortType,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No comment found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetComments(), true, 200)
}

// get replies of comments
//
//	@Summary		get replies of comments
//	@Description	get replies of comments
//	@Tags			comments
//	@Produce		application/json
//	@Param			page							query		int		false	"page number"	default(1)
//	@Param			perpage							query		int		false	"items perpage"	default(10)
//	@Param			sort							query		string	false	"sort type"		default(newest)	Enums(newest, oldest,most-liked,least-liked,most-disliked,least-disliked,most-replied,least-replied)
//	@Param			commentUrl						path		string	false	"commentUrl"
//	@Param			X-API-KEY						header		string	false	"optional authentication"
//	@Success		200								{string}	string	"ok"
//	@Success		204								{string}	string	"no content"
//	@Failure		400								{string}	string	"request failed"
//	@Failure		500								{string}	string	"server error"
//	@Router			/comments/{commentUrl}/replies	[get]
func GetRepliesOfComment(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var userId int64 = 0
	if currentUser != nil {
		userId = currentUser.Id
	}
	commentUrl := chi.URLParam(r, "commentUrl")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "perpage", "page", "sort")

	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateCommentsSortTypes(sortTypeStr)

	res, err := client.CommentService.GetRepliesOfComment(context.Background(), &comment.CommentReq{
		CommentUrl: commentUrl,
		UserId:     userId,
		Page:       toPage(perpage, page),
		Sort:       sortType,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No comment found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetComments(), true, 200)
}

// func GetAllCommentsOfUser(w http.ResponseWriter, r *http.Request) {
// 	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
// 	// currentUser := GetAuthUser(r)
// 	url := chi.URLParam(r, "url")
// 	var body map[string]any
// 	helpers.ValidateAllowedParams(body, "perpage", "page")
// 	helpers.ReadJson(r, &body)
// 	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage")
// 	page := helpers.ValidatePositiveInt(body["page"], "page")

// 	res, err := client.CommentService.GetAllCommentsOfUser(context.Background(), &comment.CommentReq{
// 		UserId:   currentUser.Id,
// 		MediaUrl: url,
// 		Page:     toPage(perpage, page),
// 	})
// 	if err != nil {
// 		helpers.LogPanic(err)
// 	}
// 	PanicIfIsError(res.GetErr())
// 	if res.GetEmpty() != nil {
// 		helpers.WriteJsonMessage(w, "No media found!", 404)
// 		return
// 	}
// 	commentsData := res.GetComments()
// 	// comments := commentsData.GetComments()
// 	helpers.WriteJson(w, commentsData, 200)
// }

// create comment/reply
//
//	@Summary		create comment/reply
//	@Description	create comment/reply
//	@Tags			comments
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			description				formData	string	true	"description"
//	@Param			reply_url				formData	string	false	"reply_url"
//	@Param			url						path		string	true	"media url"
//	@Success		200						{string}	string	"ok"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/medias/{url}	[post]
func CreateComment(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url") /// media url
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "reply_url", "description")

	text := helpers.ValidateRequiredStr(body["description"], "description")
	commentUrl := helpers.ValidateStr(body["reply_url"], "reply_url", "") // could be empty

	res, err := client.CommentService.CreateComment(context.Background(), &comment.EditCommentData{
		Text:          text,
		MediaUrl:      url,
		ReplyUrl:      commentUrl,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if c := res.GetComment(); c != nil {
		helpers.WriteJson(w, map[string]any{
			"message": "Comment created.",
			"url":     c.Url,
		}, 201)
		return
	}
	helpers.LogPanic("CreateComment should return empty or CommentData!!!")
}

// edit comment
//
//	@Summary		edit comment
//	@Description	edit comment
//	@Tags			comments
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			description				formData	string	true	"description"
//	@Param			commentUrl				path		string	true	"commentUrl"
//	@Success		200						{string}	string	"ok"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		401						{string}	string	"not authenticated"
//	@Failure		403						{string}	string	"not authorized"
//	@Failure		404						{string}	string	"not found"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/{commentUrl}	[put]
func EditComment(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	commentUrl := chi.URLParam(r, "commentUrl")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "description")

	text := helpers.ValidateRequiredStr(body["description"], "description")
	res, err := client.CommentService.EditComment(context.Background(), &comment.EditCommentData{
		Text:          text,
		Url:           commentUrl,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Comment Edited.", 200)
		return
	}
	helpers.LogPanic("EditComment should return empty or httpError!!!")
}

// delete comment
//
//	@Summary		delete comment
//	@Description	delete comment
//	@Tags			comments
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			commentUrl				path		string	true	"commentUrl"
//	@Success		200						{string}	string	"ok"
//	@Failure		400						{string}	string	"request failed"
//	@Failure		401						{string}	string	"not authenticated"
//	@Failure		403						{string}	string	"not authorized"
//	@Failure		404						{string}	string	"not found"
//	@Failure		500						{string}	string	"server error"
//	@Router			/comments/{commentUrl}	[delete]
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	commentUrl := chi.URLParam(r, "commentUrl")

	res, err := client.CommentService.DeleteComment(context.Background(), &comment.EditCommentData{
		MediaUrl:      url,
		Url:           commentUrl,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		helpers.LogPanic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Comment Deleted.", 200)
		return
	}
	helpers.LogPanic("DeleteComment should return empty or httpError!!!")
}
