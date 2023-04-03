package middleware

import (
	"context"
	"log"
	"net/http"

	"tfg/internal/jwt"
	"tfg/internal/users"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user users.User

			for name, values := range r.Header {
				// Loop over all values for the name.
				for _, value := range values {
					log.Printf("Unauthenticated %s %s", name, value)
				}
			}

			log.Printf("Starting middleware")

			header := r.Header.Get("Authorization")

			if header == "" {
				log.Printf("Unauthenticated %v", r.URL.Scheme)
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				log.Printf("Error: Petition forbidden")
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user.Username = username
			err = user.GetUserByUsername()
			if err != nil {
				log.Printf("Error: User not found")
				http.Error(w, "User not found", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, &user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *users.User {
	log.Printf("Get context of petition")
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}