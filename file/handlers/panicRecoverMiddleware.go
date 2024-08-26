package handlers

import (
	"fmt"

	// "log"
	"net/http"
	"youtube-clone/database/pbs/helper"
	"youtube-clone/file/helpers"
)

func RecoverServerPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("error type:%T,error:%v \n", err, err)
				switch err2 := err.(type) {
				case *helper.HttpError:
					if err2 != nil {
						helpers.WriteJsonError(w, err2.Message, int(err2.StatusCode))
					} else {
						helpers.WriteJsonError(w, "Something went wrong !!!", 500)
					}
				case *helpers.ServerError:
					helpers.WriteJson(w, err2.ErrorMessage(), err2.Status)
				case helpers.ServerError:
					helpers.WriteJson(w, err2.ErrorMessage(), err2.Status)
				case helpers.ValidationFieldError:
					helpers.WriteJson(w, err2.ErrorMessage(), 400)
				case *helpers.ValidationFieldErrors:
					helpers.WriteJson(w, err2.ErrorMessage(), 400)
				default:
					helpers.WriteJsonError(w, "Something went wrong !!!", 500)
				}
				return
			}
		}()

		next.ServeHTTP(w, r)

	})
}
