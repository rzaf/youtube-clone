package models

import (
	"github.com/rzaf/youtube-clone/database/pbs/helper"
)

type ModelError struct {
	Status  int
	Message string
}

func (e *ModelError) Error() string {
	return e.Message
}

func (e *ModelError) ToHttpError() *helper.HttpError {
	// err := &user_pb.Response_Err{
	// 	Err: &helper.HttpError{
	// 		Message:    e.Message,
	// 		StatusCode: int32(e.Status),
	// 	},
	// }
	// return &user_pb.Response{Res: err}
	return &helper.HttpError{
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
