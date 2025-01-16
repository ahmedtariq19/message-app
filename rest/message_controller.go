package rest

import (
	"encoding/json"
	"io"
	"message-app/models"
	"message-app/services"
	"net/http"
)

type MessageController interface {
	CreateMessage() func(w http.ResponseWriter, r *http.Request)
	HandlerWithMiddleware(middleware MiddleWare) func(w http.ResponseWriter, r *http.Request)
}

type messageController struct {
	messageService services.MessageService
}

func NewMessageController(messageService services.MessageService) MessageController {
	return &messageController{messageService: messageService}
}

func (a *messageController) CreateMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateMessageReq
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

		err = a.messageService.CreateMessage(&req)
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

		response := NewStandardResponse(true, models.SUCCESS, models.SUCCESS_MSG, nil)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (p *messageController) HandlerWithMiddleware(middleware MiddleWare) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.ValidateAuthToken(http.HandlerFunc(p.CreateMessage())).
			ServeHTTP(w, r)
	}
}
