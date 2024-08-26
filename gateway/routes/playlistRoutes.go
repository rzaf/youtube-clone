package routes

import (
	"youtube-clone/gateway/handlers"

	"github.com/go-chi/chi"
)

func setPlaylistRoutes(router *chi.Mux) {
	router.Get("/playlists", handlers.GetPlaylists)
	router.Get("/playlists/search/{term}", handlers.SearchPlaylists)
	router.Get("/playlists/{url}", handlers.GetPlaylist)
	router.Get("/playlists/{url}/medias", handlers.GetMediasOfPlaylist)
}

func setPlaylistAuthRoutes(router chi.Router) {
	router.Post("/playlists", handlers.CreatePlaylist)
	router.Put("/playlists/{url}", handlers.EditPlaylist)
	router.Delete("/playlists/{url}", handlers.DeletePlaylist)
}
