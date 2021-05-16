package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

func JSON() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(resp, req)
		})
	}
}
