package models

import (
	"github.com/rzaf/youtube-clone/notification/pbs/notificationHelperPb"
)

type ModelError struct {
	Status  int
	Message string
}

func (e *ModelError) Error() string {
	return e.Message
}

func (e *ModelError) ToHttpError() *notificationHelperPb.HttpError {
	return &notificationHelperPb.HttpError{
		Message:    e.Message,
		StatusCode: int32(e.Status),
	}
}

func NewModelError(message string, status int) *ModelError {
	return &ModelError{Message: message, Status: status}
}

func ConvertError(e error) (*ModelError, bool) {
	me, ok := e.(*ModelError)
	return me, ok
}
