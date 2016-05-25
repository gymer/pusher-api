package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gymer/pusher-api/models"
)

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="API realm"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.String(), ns1(""), "", 1)

		if strings.Index(path, "/apps/") == 0 {
			var app models.App
			u, p, ok := r.BasicAuth()
			appId := strings.Split(path, "/")[2]

			if !ok {
				unauthorized(w)
				return
			}
			fmt.Printf("AuthMiddleware req: %+v \n", appId)

			err := models.DB.Where("id = ? and client_access_token = ? and server_access_token = ?", appId, u, p).First(&app).Error
			if err != nil {
				unauthorized(w)
				return
			}

			// context.Set(r, "app", controllers.findOrAddApp(appId))
		}

		next.ServeHTTP(w, r)
	})
}

func AddFilter(pattern string, filter func(h http.Handler) http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.Replace(r.URL.String(), ns1(""), "", 1)

			if strings.Index(path, pattern) == 0 {
				// filter.ServeHTTP(w, r)
			}

			next.ServeHTTP(w, r)
		})
	}
}
