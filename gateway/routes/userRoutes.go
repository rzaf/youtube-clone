package routes

import (
	"github.com/rzaf/youtube-clone/gateway/handlers"

	"github.com/go-chi/chi"
)

func setUserRoutes(router *chi.Mux) {
	router.Get("/users/{username}", handlers.GetUserByUsername)
	router.Get("/users", handlers.GetUsers)
	router.Get("/users/search/{term}", handlers.SearchUsers)

	router.Post("/users", handlers.CreateUser)
	router.Post("/users/sign-in", handlers.SignIn)

	router.Get("/users/{username}/verify/{code}", handlers.VerifyEmail)
	router.Post("/users/resend-email", handlers.ResendEmailVerfication)
}

func setUserAuthRoutes(router chi.Router) {
	router.Get("/users/{username}/followings", handlers.GetFollowings)

	router.Put("/users/{username}/profile-photo", handlers.SetProfilePhoto)
	router.Put("/users/{username}/channel-photo", handlers.SetChannelPhoto)
	router.Put("/users/{username}/newApiKey", handlers.NewUserApiKey)
	router.Put("/users/{username}", handlers.EditUser)
	router.Delete("/users/{username}", handlers.DeleteUser)

	// router.Get("/users/{username}/comments", handlers.GetAllCommentsOfUser)
}
