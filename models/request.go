package models

type AuthenticationReq struct {
	Password string `json:"password"`
}

type CreateMessageReq struct {
	Content string `json:"content"`
}
