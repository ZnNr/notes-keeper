package handler

import (
	"context"
	"github.com/ZnNr/notes-keeper.git/intenal/service"
	"net/http"
)

func JWTMiddleware(authService service.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			tokenString := cookie.Value

			userId, err := authService.ParseToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			r = r.WithContext(context.WithValue(r.Context(), "userID", userId))
			next.ServeHTTP(w, r)
		})
	}
}
