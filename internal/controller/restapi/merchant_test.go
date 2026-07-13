package restapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"payment_gateway/internal/controller/restapi"
	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
	merchantusecase "payment_gateway/internal/usecase/merchant"

	"github.com/go-playground/validator/v10"
)

type mockMerchantRepo struct {
	created  *entity.Merchant
	list     []*entity.Merchant
	get      *entity.Merchant
	err      error
	listErr  error
	getByErr error
}

func (r *mockMerchantRepo) Create(_ context.Context, m *entity.Merchant) error {
	if r.err != nil {
		return r.err
	}
	m.Id = 1
	r.created = m
	return nil
}

func (r *mockMerchantRepo) List(_ context.Context) ([]*entity.Merchant, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.list, nil
}

func (r *mockMerchantRepo) GetById(_ context.Context, _ int64) (*entity.Merchant, error) {
	if r.getByErr != nil {
		return nil, r.getByErr
	}
	return r.get, nil
}

func TestMerchantRegistrationReturnsCreated(t *testing.T) {
	controller := newMerchantController(&mockMerchantRepo{})
	req := httptest.NewRequest(http.MethodPost, "/admin/merchants", bytes.NewBufferString(validMerchantJSON()))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var got response.MerchantRegistration
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ApiKey == "" {
		t.Fatal("api_key is empty")
	}
	if got.Merchant.Id != 1 {
		t.Fatalf("merchant id = %d, want 1", got.Merchant.Id)
	}
	if got.Merchant.Name != "Test Merchant" {
		t.Fatalf("merchant name = %q, want Test Merchant", got.Merchant.Name)
	}
}

func TestMerchantRegistrationRejectsInvalidJSON(t *testing.T) {
	controller := newMerchantController(&mockMerchantRepo{})
	req := httptest.NewRequest(http.MethodPost, "/admin/merchants", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestMerchantRegistrationRejectsInvalidRequest(t *testing.T) {
	controller := newMerchantController(&mockMerchantRepo{})
	body := `{
		"name":"Test Merchant",
		"domain":"example.com",
		"webhook_url":"not-url",
		"success_redirect_url":"https://example.com/success",
		"failure_redirect_url":"https://example.com/failure",
		"is_active":true
	}`
	req := httptest.NewRequest(http.MethodPost, "/admin/merchants", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestMerchantRegistrationMapsUseCaseError(t *testing.T) {
	controller := newMerchantController(&mockMerchantRepo{
		err: apperr.AlreadyExists("merchant already exists"),
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/merchants", bytes.NewBufferString(validMerchantJSON()))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	assertErrorResponse(t, rec, http.StatusConflict, apperr.CodeAlreadyExists)
}

func TestMerchantListReturnsOK(t *testing.T) {
	controller := newMerchantController(&mockMerchantRepo{
		list: []*entity.Merchant{
			{Id: 1, Name: "First Merchant", Domain: "first.example.com", IsActive: true},
			{Id: 2, Name: "Second Merchant", Domain: "second.example.com", IsActive: false},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/admin/merchants", nil)
	rec := httptest.NewRecorder()

	controller.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got response.ListMerchant
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got.Merchants) != 2 {
		t.Fatalf("merchants len = %d, want 2", len(got.Merchants))
	}
}

func TestMerchantListMapsUseCaseError(t *testing.T) {
	controller := newMerchantController(&mockMerchantRepo{
		listErr: apperr.Internal("failed to list merchants", nil),
	})
	req := httptest.NewRequest(http.MethodGet, "/admin/merchants", nil)
	rec := httptest.NewRecorder()

	controller.List(rec, req)

	assertErrorResponse(t, rec, http.StatusInternalServerError, apperr.CodeInternal)
}

func newMerchantController(repo *mockMerchantRepo) *restapi.MerchantController {
	validate := validator.New()
	uc := merchantusecase.NewUseCase(repo, validate)
	return restapi.NewMerchantController(uc, validate)
}

func validMerchantJSON() string {
	return `{
		"name":"Test Merchant",
		"domain":"example.com",
		"webhook_url":"https://example.com/webhook",
		"success_redirect_url":"https://example.com/success",
		"failure_redirect_url":"https://example.com/failure",
		"is_active":true
	}`
}
