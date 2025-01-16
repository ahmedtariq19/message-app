package rest

import "encoding/json"

func jsonErrorResponse(code uint, message string) string {
	response := NewStandardResponse(false, code, message, nil)
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

type StandardResponse struct {
	Result  bool        `json:"result"`
	Code    uint        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Assuming NewStandardResponse is defined elsewhere to create a standardized response
func NewStandardResponse(result bool, code uint, msg string, data interface{}) *StandardResponse {

	if data == nil {
		data = ""
	}
	return &StandardResponse{
		Result:  result,
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

type getTokenVo struct {
	Token string `json:"token"`
}

func newtokenVo(token string) *getTokenVo {
	return &getTokenVo{
		Token: token,
	}
}
