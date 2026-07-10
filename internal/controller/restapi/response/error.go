package response

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewErrorResponse(errorCode, message string) ErrorResponse {
	return ErrorResponse{Error: errorCode, Message: message}
}
