package routes

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

var (
	publicRoutes = map[string]string{
		"/api/users":                          "get",
		"/api/users/search/{term}":            "get",
		"/api/users/{username}":               "get",
		"/api/users/resend-email":             "post",
		"/api/users/{username}/verify/{code}": "get",

		"/api/medias":               "get",
		"/api/medias/search/{term}": "get",
		"/api/medias/{url}":         "get",

		"/api/playlists":               "get",
		"/api/playlists/search/{term}": "get",
		"/api/playlists/{url}":         "get",
		"/api/playlists/{url}/medias":  "get",

		"/api/comments/{commentUrl}":         "get",
		"/api/comments/medias/{url}":         "get",
		"/api/comments/{commentUrl}/replies": "get",
	}
	protectedRoutes = map[string]string{
		"/api/users/{username}/followings":    "get",
		"/api/users/{username}/profile-photo": "put",
		"/api/users/{username}/channel-photo": "put",
		"/api/users/{username}":               "put,delete",

		"/api/follows/{username}": "post,delete",

		"/api/medias":                               "post",
		"/api/medias/{url}":                         "put,delete",
		"/api/medias/{url}/tag/{name}":              "post,delete",
		"/api/medias/{url}/playlists/{playlistUrl}": "post,delete,put",

		"/api/comments/medias/{url}": "post",
		"/api/comments/{commentUrl}": "put,delete",

		"/api/medias/{url}/likes":   "post,delete",
		"/api/comments/{url}/likes": "post,delete",
	}
)

func Test_routes_exist(t *testing.T) {
	allRoutes := GetRoutes()
	for route, methods := range publicRoutes {
		for _, method := range strings.Split(methods, ",") {
			routeExist(t, allRoutes, route, method)
		}
	}
	for route, methods := range protectedRoutes {
		for _, method := range strings.Split(methods, ",") {
			routeExist(t, allRoutes, route, method)
		}
	}
}

func routeExist(t *testing.T, routes chi.Router, route string, method string) {
	method = strings.ToUpper(method)
	f := false
	_ = chi.Walk(routes, func(foundMethod, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if foundRoute == route && foundMethod == method {
			f = true
		}
		return nil
	})
	if !f {
		t.Errorf("route: %s:'%s' was not found in registered routes", method, route)
	}
}

func Test_protected_routes(t *testing.T) {
	allRoutes := GetRoutes()
	for route, methods := range protectedRoutes {
		for _, method := range strings.Split(methods, ",") {
			method = strings.ToUpper(method)
			f := false
			_ = chi.Walk(allRoutes, func(foundMethod, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
				if foundRoute == route && foundMethod == method {
					fmt.Printf("route:%s , %v\n", route, middlewares)
					f = true
				}
				return nil
			})
			if !f {
				t.Errorf("route: %s:'%s' was not found in registered routes", method, route)
			}
		}
	}
}
