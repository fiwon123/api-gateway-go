package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fiwon123/api-gateway-go/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuthentication(t *testing.T) {
	secret := []byte("test-secret")
	issuer := "test-issuer"
	audience := "test-audience"

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}

	handler := auth.JwtAuthentication(
		secret,
		issuer,
		audience,
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/users/42",
		nil,
	)

	request.Header.Set(
		"Authorization",
		"Bearer "+signedToken,
	)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestJWTAuthenticationRejectsMissingToken(t *testing.T) {
	handler := auth.JwtAuthentication(
		[]byte("test-secret"),
		"test-issuer",
		"test-audience",
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/users/42",
		nil,
	)

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			response.Code,
		)
	}
}
