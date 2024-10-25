package services

import (
	"context"
	"fmt"
	"github.com/rzaf/youtube-clone/database/models"
	"github.com/rzaf/youtube-clone/database/pbs/comment"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
)

type commentServiceServer struct {
	comment.CommentServiceServer
}

func newResponseFromFullCommentData(c *comment.FullCommentData) *comment.Response {
	return &comment.Response{
		Res: &comment.Response_FullComment{
			FullComment: c,
		},
	}
}

func newResponseFromCommentData(c *comment.CommentData) *comment.Response {
	return &comment.Response{
		Res: &comment.Response_Comment{
			Comment: c,
		},
	}
}

func newResponseFromComments(comments []*models.Comment, pageInfo *helper.PagesInfo) *comment.Response {
	if comments == nil {
		return newCommentResponseFromEmpty()
	}
	var commentsData []*comment.CommentData
	for _, c := range comments {
		c2 := &comment.CommentData{
			Url:                c.Url,
			Text:               c.Text,
			LikesCount:         c.LikesCount,
			DislikesCount:      c.DislikesCount,
			CreatorUsername:    c.UserName,
			CreatorUserProfile: c.UserProfile,
			CreatedAt:          c.CreatedAt.Unix(),
			UserLike:           c.CurrentUserLike.String(),
			RepliesCount:       c.RepliesCount,
			ReplyUrl:           c.ReplyUrl,
			ReplyUserName:      c.ReplyUsername,
		}
		if c.UpdatedAt != nil {
			c2.UpdatedAt = c.UpdatedAt.Unix()
		}
		commentsData = append(commentsData, c2)
	}
	return &comment.Response{
		Res: &comment.Response_Comments{
			Comments: &comment.CommentsData{
				Comments:  commentsData,
				PagesInfo: pageInfo,
			},
		},
	}
}

func newCommentResponseFromError(e *helper.HttpError) *comment.Response {
	return &comment.Response{
		Res: &comment.Response_Err{
			Err: e,
		},
	}
}

func newCommentResponseFromEmpty() *comment.Response {
	return &comment.Response{
		Res: &comment.Response_Empty{},
	}
}

func (*commentServiceServer) GetCommentByUrl(con context.Context, c *comment.CommentUrl) (*comment.Response, error) {
	var commentM *models.Comment
	var err error
	if c.UserId == 0 {
		commentM, err = models.GetComment(c.Url)
	} else {
		commentM, err = models.AuthGetComment(c.Url, c.UserId)
	}
	fmt.Printf("%v\n", commentM)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	c2 := &comment.FullCommentData{
		Url:           commentM.Url,
		Text:          commentM.Text,
		RepliesCount:  commentM.RepliesCount,
		LikesCount:    commentM.LikesCount,
		DislikesCount: commentM.DislikesCount,
		CreatedAt:     commentM.CreatedAt.Unix(),
		Username:      commentM.UserName,
		UserProfile:   commentM.UserProfile,

		ReplyText:        commentM.ReplyText,
		ReplyUrl:         commentM.ReplyUrl,
		ReplyUserName:    commentM.ReplyUsername,
		ReplyUserProfile: commentM.ReplyUserProfile,

		UserLike:            commentM.CurrentUserLike.String(),
		MediaUrl:            commentM.MediaUrl,
		MediaTitle:          commentM.MediaTitle,
		MediaCreator:        commentM.MediaCreator,
		MediaCraetorProfile: commentM.MediaCreatorProfile,
		MediaType:           mediaTypeToStr(commentM.MediaType),
		MediaCreatedAt:      commentM.MediaCreatedAt.Unix(),
		MediaThumbnail:      commentM.MediaThumbnail,
	}
	if commentM.UpdatedAt != nil {
		c2.UpdatedAt = commentM.UpdatedAt.Unix()
	}
	return newResponseFromFullCommentData(c2), nil
}

func (*commentServiceServer) GetCommentsOfMedia(con context.Context, c *comment.CommentReq) (*comment.Response, error) {
	PerPage, PageNumber := getPage(c.Page)
	var comments []*models.Comment
	var err error
	var totalPages int64
	if c.UserId == 0 {
		totalPages, comments, err = models.GetCommentsOfMedia(c.MediaUrl, PerPage, PageNumber, c.Sort)
	} else {
		totalPages, comments, err = models.AuthGetCommentsOfMedia(c.MediaUrl, PerPage, PageNumber, c.UserId, c.Sort)
	}
	fmt.Printf("%v\n", comments)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromComments(comments, newPagesInfo(int32(totalPages), c.Page.PageNumber)), nil
}

func (*commentServiceServer) GetRepliesOfComment(con context.Context, c *comment.CommentReq) (*comment.Response, error) {
	PerPage, PageNumber := getPage(c.Page)
	var comments []*models.Comment
	var err error
	var totalPages int64
	if c.UserId == 0 {
		totalPages, comments, err = models.GetRepliesOfComment(c.CommentUrl, PerPage, PageNumber, c.Sort)
	} else {
		totalPages, comments, err = models.AuthGetRepliesOfComment(c.CommentUrl, PerPage, PageNumber, c.UserId, c.Sort)
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromComments(comments, newPagesInfo(int32(totalPages), c.Page.PageNumber)), nil
}

// authenticated request
func (*commentServiceServer) GetAllCommentsOfUser(con context.Context, c *comment.CommentReq) (*comment.Response, error) {
	PerPage, PageNumber := getPage(c.Page)
	totalPages, comments, err := models.GetAllCommentsOfUser(c.UserId, PerPage, PageNumber, c.Sort)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromComments(comments, newPagesInfo(int32(totalPages), c.Page.PageNumber)), nil
}

// authenticated request
func (*commentServiceServer) GetCommentsOfUserInMedia(con context.Context, c *comment.CommentReq) (*comment.Response, error) {
	PerPage, PageNumber := getPage(c.Page)
	totalPages, comments, err := models.GetCommentsOfUserInMedia(c.UserId, c.MediaUrl, PerPage, PageNumber, c.Sort)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromComments(comments, newPagesInfo(int32(totalPages), c.Page.PageNumber)), nil
}

func (*commentServiceServer) CreateComment(con context.Context, c *comment.EditCommentData) (*comment.Response, error) {
	var createdComment *models.Comment
	var err error
	if c.ReplyUrl == "" {
		createdComment, err = models.CreateComment(c.Text, c.CurrentUserId, c.MediaUrl)
	} else {
		createdComment, err = models.CreateReply(c.Text, c.CurrentUserId, c.MediaUrl, c.ReplyUrl)
	}
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	if c.ReplyUrl == "" {
		go newCommentNotification(c.CurrentUserId, c.Text, c.MediaUrl)
	} else {
		go newReplyNotification(c.CurrentUserId, c.Text, c.MediaUrl, c.ReplyUrl)
	}
	return newResponseFromCommentData(&comment.CommentData{Url: createdComment.Url}), nil
}

func (*commentServiceServer) EditComment(con context.Context, c *comment.EditCommentData) (*comment.Response, error) {
	err := models.EditComment(c.Url, c.Text, c.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newCommentResponseFromEmpty(), nil
}

func (*commentServiceServer) DeleteComment(con context.Context, c *comment.EditCommentData) (*comment.Response, error) {
	err := models.DeleteComment(c.Url, c.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newCommentResponseFromEmpty(), nil
}

func (*commentServiceServer) CreateLikeComment(con context.Context, l *helper.LikeReq) (*comment.Response, error) {
	err := models.CreateCommentLike(l.UserId, l.Url, l.IsLike)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	go newCommentLikeNotification(l.UserId, l.Url, l.IsLike)
	return newCommentResponseFromEmpty(), nil
}

func (*commentServiceServer) DeleteLikeComment(con context.Context, l *helper.LikeReq) (*comment.Response, error) {
	err := models.DeleteCommentLike(l.UserId, l.Url)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newCommentResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newCommentResponseFromEmpty(), nil
}
