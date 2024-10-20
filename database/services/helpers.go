package services

import (
	"context"
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
