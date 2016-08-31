package models

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message, omitempty"`
}

var OK_RESPONSE StatusResponse = StatusResponse{
	Status: "ok",
}

func ErrorResponse(message string) StatusResponse {
	return StatusResponse{
		Status:  "error",
		Message: message,
	}
}
