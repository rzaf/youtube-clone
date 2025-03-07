package handlers

import (
	"fmt"

	"net/http"

	authMiddlewares "github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
	"github.com/rzaf/youtube-clone/gateway/helpers"
)

func RecoverServerPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("RecoverServerPanics middleware caught panic! (type(err):%T,value(err):%v)\n", err, err)
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

				case authMiddlewares.AuthError:
					helpers.WriteJsonError(w, err2.Message, err2.Status)
				case *authMiddlewares.AuthError:
					helpers.WriteJsonError(w, err2.Message, err2.Status)
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
