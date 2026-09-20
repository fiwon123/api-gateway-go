package auth

import (
	"net/http"
	"strings"

	"github.com/fiwon123/api-gateway-go/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthentication(
	secret []byte,
	expectedIssuer string,
	expectedAudience string,
) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
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

			if claims.Issuer != expectedIssuer {
				http.Error(
					w,
					`{"error":"invalid token issuer"}`,
					http.StatusUnauthorized,
				)
				return
			}

			if !hasAudience(claims.Audience, expectedAudience) {
				http.Error(
					w,
					`{"error":"invalid token audience"}`,
					http.StatusUnauthorized,
				)
				return
			}

			if claims.Subject != "" {
				r.Header.Set("X-User-ID", claims.Subject)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasAudience(
	audiences jwt.ClaimStrings,
	expected string,
) bool {
	for _, audience := range audiences {
		if audience == expected {
			return true
		}
	}

	return false
}
