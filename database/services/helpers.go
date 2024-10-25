package services

import (
	"context"
	"fmt"
	"time"

	"github.com/rzaf/youtube-clone/database/client"
	"github.com/rzaf/youtube-clone/database/models"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
	"github.com/rzaf/youtube-clone/notification/pbs/emailPb"
)

func newPagesInfo(PageCount, CurrentPage int32) *helper.PagesInfo {
	return &helper.PagesInfo{
		TotalPages:  PageCount,
		CurrentPage: CurrentPage,
	}
}

func getPage(p *helper.Paging) (int, int) {
	if p == nil {
		panic("*helper.Paging shuld not be nill")
	}
	limit := p.PerPage
	offset := (p.PageNumber - 1) * p.PerPage
	return int(limit), int(offset)
}

// func toPage(perPage int, pageNumber int) *helper.Paging {
// 	return &helper.Paging{
// 		PerPage:    int32(perPage),
// 		PageNumber: int32(pageNumber),
// 	}
// }

func sendEmailVerificationNotification(username string, email string, code string) error {
	r, err := client.EmailService.SendVerifcation(context.Background(), &emailPb.UserVerifyReq{
		Username:  username,
		UserEmail: email,
		Link:      "http://127.0.0.1:6070/api/users/" + username + "/verify/" + code,
	})
	if err != nil {
		return err
	}
	if r.GetErr() != nil {
		return &models.ModelError{Message: "failed to send email", Status: 500}
	}
	return nil
}

func sendNotificationEmail(username string, email string, title string, message string) error {
	retries := 3
	for i := 0; i < retries; i++ {
		r, err := client.EmailService.SendNotification(context.Background(), &emailPb.NotificationData{
			Username:  username,
			UserEmail: email,
			Title:     title,
			Message:   message,
		})
		if err != nil {
			fmt.Printf("failed to send notification to notification service. err:%v\n", err)
			fmt.Printf("retrying after 1500 ms\n")
			time.Sleep(1500 * time.Millisecond)
			continue
		}
		if r.GetErr() != nil {
			fmt.Printf("error recieved from notification service:%v\n", r.GetErr())
			fmt.Printf("retrying after 1500 ms\n")
			time.Sleep(1500 * time.Millisecond)
			continue
		}
		return nil
	}
	fmt.Printf("failed to send email after %d tries!\n", retries)
	return &models.ModelError{Message: "failed to send email", Status: 500}
}

func followingNotification(followerId int64, followingUsername string) error {
	follower, _ := models.Helper_UserById(followerId)
	if follower.Username == followingUsername {
		return nil
	}
	following, _ := models.Helper_UserByUsername(followingUsername)
	if !follower.IsVerified {
		return nil
	}
	title := "Notification: New follower"
	message := fmt.Sprintf("User with username:`%s` and channel_name:`%s` is now following you.", follower.Username, follower.ChannelName)
	err := sendNotificationEmail(followingUsername, following.Email, title, message)
	if err != nil {
		return err
	}
	return nil
}

func newMediaNotification(creatorId int64, mediaUrl string, mediaTitle string, mediaTypeStr string, mediaText string) error {
	mediaCreator, _ := models.Helper_UserById(creatorId)
	followers, _ := models.Helper_FollowersById(creatorId)
	for i := 0; i < len(followers); i++ {
		if followers[i].Username == mediaCreator.Username {
			continue
		}
		title := fmt.Sprintf("Notification: you'r subscribing:`%s` published a new %s titled:`%s`", mediaCreator.ChannelName, mediaTypeStr, mediaTitle)
		message := fmt.Sprintf("You'r subscribing channel:`%s` with username:`%s` published a new %s with title:`%s` , url:`%s` and description:`%s`",
			mediaCreator.ChannelName,
			mediaCreator.Username,
			mediaTypeStr,
			mediaTitle,
			mediaUrl,
			mediaText,
		)
		sendNotificationEmail(followers[i].Username, followers[i].Email, title, message)
	}
	return nil
}

func newCommentNotification(commenterId int64, commentText string, mediaUrl string) error {
	media, _ := models.Helper_MediaByUrl(mediaUrl)
	commenter, _ := models.Helper_UserById(commenterId)
	if commenter.Username == media.UserName {
		return nil
	}
	mediaType := mediaTypeToStr(media.Type)
	title := fmt.Sprintf("Notification: New comment on your %s with title:`%s`", mediaType, media.Title)
	message := fmt.Sprintf(
		"User with username:`%s` and channel_name:`%s` has commented:`%s` on your %s media with url:`%s` and title:`%s`.",
		commenter.Username,
		commenter.ChannelName,
		commentText,
		mediaType,
		mediaUrl,
		media.Title,
	)
	err := sendNotificationEmail(media.UserName, media.UserEmail, title, message)
	if err != nil {
		return err
	}
	return nil
}

func newReplyNotification(replierId int64, replyText string, mediaUrl string, commentUrl string) error {
	media, _ := models.Helper_MediaByUrl(mediaUrl)
	repliedTo, _ := models.Helper_CommentByUrl(commentUrl)
	fmt.Println(commentUrl, repliedTo)
	replier, _ := models.Helper_UserById(replierId)
	if replier.Username == repliedTo.UserName {
		return nil
	}
	mediaType := mediaTypeToStr(media.Type)
	var title, message string
	if repliedTo.ReplyId == 0 {
		title = fmt.Sprintf("Notification: New reply to your comment on %s with title:`%s`", mediaType, media.Title)
		message = fmt.Sprintf(
			"User with username:`%s` and channel_name:`%s` has replied:`%s` to your comment:`%s` with url:`%s` on %s with url:`%s` and title:`%s`.",
			replier.Username,
			replier.ChannelName,
			replyText,
			repliedTo.Text,
			commentUrl,
			mediaType,
			mediaUrl,
			media.Title,
		)
	} else {
		title = fmt.Sprintf("Notification: New reply to your reply on %s with title:`%s`", mediaType, media.Title)
		message = fmt.Sprintf(
			"User with username:`%s` and channel_name:`%s` has replied:`%s` to your reply:`%s` with url:`%s` on %s with url:`%s` and title:`%s`.",
			replier.Username,
			replier.ChannelName,
			replyText,
			repliedTo.Text,
			commentUrl,
			mediaType,
			mediaUrl,
			media.Title,
		)
	}

	err := sendNotificationEmail(repliedTo.UserName, repliedTo.UserEmail, title, message)
	if err != nil {
		return err
	}
	return nil
}

func newMediaLikeNotification(likerId int64, mediaUrl string, isLike bool) error {
	media, _ := models.Helper_MediaByUrl(mediaUrl)
	liker, _ := models.Helper_UserById(likerId)
	mediaType := mediaTypeToStr(media.Type)
	if liker.Username == media.UserName {
		return nil
	}
	likeStr := ""
	if isLike {
		likeStr = "like"
	} else {
		likeStr = "dislike"
	}
	title := fmt.Sprintf("Notification: New %s on your %s with title:`%s`", likeStr, mediaType, media.Title)
	message := fmt.Sprintf(
		"User with username:`%s` and channel_name:`%s` has %s'd your %s media with url:`%s` and title:`%s`.",
		liker.Username,
		liker.ChannelName,
		likeStr,
		mediaType,
		mediaUrl,
		media.Title,
	)
	err := sendNotificationEmail(media.UserName, media.UserEmail, title, message)
	if err != nil {
		return err
	}
	return nil
}

func newCommentLikeNotification(likerId int64, commentUrl string, isLike bool) error {
	comment, _ := models.Helper_CommentByUrl(commentUrl)
	commentStr := "comment"
	if comment.ReplyId != 0 {
		commentStr = "reply"
	}
	fmt.Println(commentUrl, comment)
	liker, _ := models.Helper_UserById(likerId)
	mediaType := mediaTypeToStr(comment.MediaType)
	if liker.Username == comment.UserName {
		return nil
	}
	likeStr := ""
	if isLike {
		likeStr = "like"
	} else {
		likeStr = "dislike"
	}
	title := fmt.Sprintf("Notification: New %s on your %s on %s with title:`%s`", likeStr, commentStr, mediaType, comment.MediaTitle)
	message := fmt.Sprintf(
		"User with username:`%s` and channel_name:`%s` has %s'd your %s:`%s` with url:`%s` on %s with url:`%s` and title:`%s`.",
		liker.Username,
		liker.ChannelName,
		likeStr,
		commentStr,
		comment.Text,
		commentUrl,
		mediaType,
		comment.MediaUrl,
		comment.MediaTitle,
	)
	err := sendNotificationEmail(comment.UserName, comment.UserEmail, title, message)
	if err != nil {
		return err
	}
	return nil
}
