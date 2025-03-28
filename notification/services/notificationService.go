package services

import (
	"context"

	"github.com/rzaf/youtube-clone/notification/handlers"
	"github.com/rzaf/youtube-clone/notification/models"
	"github.com/rzaf/youtube-clone/notification/pbs/notificationHelperPb"
	"github.com/rzaf/youtube-clone/notification/pbs/notificationPb"
)

type notificationServiceServer struct {
	notificationPb.NotificationServiceServer
}

func newNotifcationResponseFromError(e *notificationHelperPb.HttpError) *notificationPb.Response {
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
	n, err := models.CreateNotification(e.UserId, e.Title, e.Message)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newNotifcationResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	handlers.SendWsMessage(n.Id, e.UserId, e.Title, e.Message)
	return newNotifcationResponseFromEmpty(), nil
}
