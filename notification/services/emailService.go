package services

import (
	"context"
	"youtube-clone/notification/email"
	"youtube-clone/notification/pbs/emailPb"
)

type emailServiceServer struct {
	emailPb.EmailServiceServer
}

func newEmailResponseFromError(e *emailPb.HttpError) *emailPb.Response {
	return &emailPb.Response{
		Res: &emailPb.Response_Err{
			Err: e,
		},
	}
}

func newEmailResponseFromEmpty() *emailPb.Response {
	return &emailPb.Response{
		Res: &emailPb.Response_Empty{},
	}
}

func (emailServiceServer) SendVerifcation(c context.Context, e *emailPb.UserVerifyReq) (*emailPb.Response, error) {
	err := email.SendVerifcationEmail(e.UserEmail, e.Username, e.Link)
	if err != nil {
		return newEmailResponseFromError(&emailPb.HttpError{Message: "failed to send email", StatusCode: 500}), nil
	}
	return newEmailResponseFromEmpty(), nil
}

func (emailServiceServer) SendNotification(context.Context, *emailPb.NotificationData) (*emailPb.Response, error) {
	return newEmailResponseFromEmpty(), nil
}
