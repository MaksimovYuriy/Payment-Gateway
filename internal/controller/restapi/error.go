package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/lib/apperr"
)

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	errorCode := string(apperr.CodeInternal)
	message := "internal server error"

	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		status = statusByAppCode(appErr.Code)
		errorCode = string(appErr.Code)
		message = appErr.Message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResponse(errorCode, message))
}

func statusByAppCode(code apperr.Code) int {
	switch code {
	case apperr.CodeInvalidInput:
		return http.StatusBadRequest
	case apperr.CodeAlreadyExists:
		return http.StatusConflict
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
