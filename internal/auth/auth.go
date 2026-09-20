package auth

import (
	"net/http"
	"strings"

	"github.com/fiwon123/api-gateway-go/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthentication(secret []byte) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(
					w,
					`{"error":"missing bearer token"}`,
					http.StatusUnauthorized,
				)
				return
			}

			rawToken := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &jwt.RegisteredClaims{}

			token, err := jwt.ParseWithClaims(
				rawToken,
				claims,
				func(token *jwt.Token) (any, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, jwt.ErrSignatureInvalid
					}

					return secret, nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(
					w,
					`{"error":"invalid or expired token"}`,
					http.StatusUnauthorized,
				)
				return
			}

			// Pass the authenticated user's ID to the backend.
			if claims.Subject != "" {
				r.Header.Set("X-User-ID", claims.Subject)
			}

			next.ServeHTTP(w, r)
		})
	}
}
