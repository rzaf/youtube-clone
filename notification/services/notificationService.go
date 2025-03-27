package services

import (
	"context"

	"github.com/rzaf/youtube-clone/notification/handlers"
	"github.com/rzaf/youtube-clone/notification/pbs/notificationPb"
)

type notificationServiceServer struct {
	notificationPb.NotificationServiceServer
}

func newNotifcationResponseFromError(e *notificationPb.HttpError) *notificationPb.Response {
	return &notificationPb.Response{
		Res: &notificationPb.Response_Err{
			Err: e,
		},
	}
}

func newNotifcationResponseFromEmpty() *notificationPb.Response {
	return &notificationPb.Response{
		Res: &notificationPb.Response_Empty{},
	}
}

func (notificationServiceServer) SetNotification(c context.Context, e *notificationPb.NotificationData) (*notificationPb.Response, error) {
	handlers.SendWsMessage(e.UserId, e.Title, e.Message)
	return newNotifcationResponseFromEmpty(), nil
}
