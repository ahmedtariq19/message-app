package rest

import (
	"encoding/json"
	"io"
	"message-app/models"
	"message-app/services"
	"net/http"
)

type AuthenticationController interface {
	CreateToken() func(w http.ResponseWriter, r *http.Request)
}

type authenticationController struct {
	authenticationService services.AuthenticationService
}

func NewAuthenticationController(authenticationService services.AuthenticationService) AuthenticationController {
	return &authenticationController{authenticationService: authenticationService}
}

func (a *authenticationController) CreateToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthenticationReq
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(NewStandardResponse(false, models.INVALID_INPUT, err.Error(), nil))
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(NewStandardResponse(false, models.INVALID_INPUT, err.Error(), nil))
			return
		}

		token, err := a.authenticationService.CreateAuthToken(&req)
		if err != nil {
			if standardError, ok := err.(*models.StandardError); ok {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(NewStandardResponse(false, standardError.Code, standardError.Message, nil))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(NewStandardResponse(false, models.INTERNAL_SERVER_ERROR, err.Error(), nil))
				return
			}
		}

		response := NewStandardResponse(true, models.SUCCESS, models.SUCCESS_MSG, newtokenVo(token))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
