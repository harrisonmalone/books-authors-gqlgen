package middleware

import (
	"log"
	"net/http"
)

func Authentication() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			log.Println("auth", auth)

			next.ServeHTTP(w, r)
		})
	}
}
