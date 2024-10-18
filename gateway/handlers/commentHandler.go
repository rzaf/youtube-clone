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
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	helpers.WriteProtoJson(w, res.GetFullComment(), true, 200)
}

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
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteEmpty(w)
		// helpers.WriteJsonMessage(w, "No comment found!", 404)
		return
	}
	helpers.WriteProtoJson(w, res.GetComments(), true, 200)
}

func GetRepliesOfComment(w http.ResponseWriter, r *http.Request) {
	currentUser := getUserFromHeader(r)
	var userId int64 = 0
	if currentUser != nil {
		userId = currentUser.Id
	}
	url := chi.URLParam(r, "url")
	commentUrl := chi.URLParam(r, "commentUrl")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "perpage", "page", "sort")

	perpage := helpers.ValidatePositiveInt(body["perpage"], "perpage", 10)
	page := helpers.ValidatePositiveInt(body["page"], "page", 1)
	sortTypeStr := helpers.ValidateStr(body["sort"], "sort", "newest")
	sortType := helpers.ValidateCommentsSortTypes(sortTypeStr)

	res, err := client.CommentService.GetRepliesOfComment(context.Background(), &comment.CommentReq{
		CommentUrl: commentUrl, MediaUrl: url,
		UserId: userId,
		Page:   toPage(perpage, page),
		Sort:   sortType,
	})
	if err != nil {
		panic(err)
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
// 		panic(err)
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
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if c := res.GetComment(); c != nil {
		helpers.WriteJson(w, map[string]any{
			"message": "Comment created.",
			"url":     c.Url,
		}, 201)
		return
	}
	panic("CreateComment should return empty or CommentData!!!")
}

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
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Comment Edited.", 200)
		return
	}
	panic("EditComment should return empty or httpError!!!")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authUser("user")).(*user_pb.CurrentUserData)
	// currentUser := GetAuthUser(r)
	url := chi.URLParam(r, "url")
	commentUrl := chi.URLParam(r, "commentUrl")
	body := make(map[string]any)
	helpers.ParseReq(r, body)
	helpers.ValidateAllowedParams(body, "description")

	text := helpers.ValidateRequiredStr(body["description"], "description")
	helpers.ValidateVar(text, "description", "required")
	// helpers.ValidateVar(commentIdStr, "reply_id", "int")
	res, err := client.CommentService.DeleteComment(context.Background(), &comment.EditCommentData{
		MediaUrl:      url,
		Url:           commentUrl,
		CurrentUserId: currentUser.Id,
	})
	if err != nil {
		panic(err)
	}
	PanicIfIsError(res.GetErr())
	if res.GetEmpty() != nil {
		helpers.WriteJsonMessage(w, "Comment Deleted.", 200)
		return
	}
	panic("DeleteComment should return empty or httpError!!!")
}
