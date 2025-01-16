// Package rest rest holds server and handler(controller) layer
package rest

import (
	"net/http"
	"strings"

	"message-app/models"
	"message-app/services"
)

type MiddleWare interface {
	ValidateAuthToken(next http.Handler) http.Handler
}

type middleWare struct {
	jwtService services.JWTService
}

// NewAuthController return new instance of authcontroller
func NewMiddleWare(
	jwtService services.JWTService,

) MiddleWare {
	return &middleWare{
		jwtService: jwtService,
	}
}

func (m *middleWare) ValidateAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, jsonErrorResponse(models.INVALID_TOKEN, "Token not found"), http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authorization, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || tokenParts[1] == "" {
			http.Error(w, jsonErrorResponse(models.INVALID_TOKEN, "Invalid token format"), http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		UID, err := m.jwtService.VerifyAuthToken(token)
		if err != nil || UID == "" {
			http.Error(w, jsonErrorResponse(models.INVALID_TOKEN, models.INVALID_TOKEN_MESSAGE), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
