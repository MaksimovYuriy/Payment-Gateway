package restapi_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/lib/apperr"
)

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantCode apperr.Code) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, wantStatus, rec.Body.String())
	}

	var got response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if got.Error != string(wantCode) {
		t.Fatalf("error code = %q, want %q", got.Error, wantCode)
	}
}
