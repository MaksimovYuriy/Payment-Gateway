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
	bankusecase "payment_gateway/internal/usecase/bank"

	"github.com/go-playground/validator/v10"
)

type mockBankRepo struct {
	created  *entity.Bank
	list     []*entity.Bank
	get      *entity.Bank
	err      error
	listErr  error
	getByErr error
}

func (r *mockBankRepo) Create(_ context.Context, b *entity.Bank) error {
	if r.err != nil {
		return r.err
	}
	b.Id = 1
	r.created = b
	return nil
}

func (r *mockBankRepo) List(_ context.Context) ([]*entity.Bank, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.list, nil
}

func (r *mockBankRepo) GetById(_ context.Context, _ int64) (*entity.Bank, error) {
	if r.getByErr != nil {
		return nil, r.getByErr
	}
	return r.get, nil
}

func TestBankRegistrationReturnsCreated(t *testing.T) {
	controller := newBankController(&mockBankRepo{})
	req := httptest.NewRequest(http.MethodPost, "/admin/banks", bytes.NewBufferString(validBankJSON()))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var got response.Bank
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Id != 1 {
		t.Fatalf("bank id = %d, want 1", got.Id)
	}
	if got.Code != "test-bank" {
		t.Fatalf("bank code = %q, want test-bank", got.Code)
	}
}

func TestBankRegistrationRejectsInvalidJSON(t *testing.T) {
	controller := newBankController(&mockBankRepo{})
	req := httptest.NewRequest(http.MethodPost, "/admin/banks", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestBankRegistrationRejectsInvalidRequest(t *testing.T) {
	controller := newBankController(&mockBankRepo{})
	req := httptest.NewRequest(http.MethodPost, "/admin/banks", bytes.NewBufferString(`{"name":"Test Bank","is_active":true}`))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestBankRegistrationMapsUseCaseError(t *testing.T) {
	controller := newBankController(&mockBankRepo{
		err: apperr.AlreadyExists("bank already exists"),
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/banks", bytes.NewBufferString(validBankJSON()))
	rec := httptest.NewRecorder()

	controller.Registration(rec, req)

	assertErrorResponse(t, rec, http.StatusConflict, apperr.CodeAlreadyExists)
}

func TestBankListReturnsOK(t *testing.T) {
	controller := newBankController(&mockBankRepo{
		list: []*entity.Bank{
			{Id: 1, Code: "first-bank", Name: "First Bank", IsActive: true},
			{Id: 2, Code: "second-bank", Name: "Second Bank", IsActive: false},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/admin/banks", nil)
	rec := httptest.NewRecorder()

	controller.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got response.ListBank
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got.Banks) != 2 {
		t.Fatalf("banks len = %d, want 2", len(got.Banks))
	}
}

func TestBankListMapsUseCaseError(t *testing.T) {
	controller := newBankController(&mockBankRepo{
		listErr: apperr.Internal("failed to list banks", nil),
	})
	req := httptest.NewRequest(http.MethodGet, "/admin/banks", nil)
	rec := httptest.NewRecorder()

	controller.List(rec, req)

	assertErrorResponse(t, rec, http.StatusInternalServerError, apperr.CodeInternal)
}

func newBankController(repo *mockBankRepo) *restapi.BankController {
	validate := validator.New()
	uc := bankusecase.NewUseCase(repo, validate)
	return restapi.NewBankController(uc, validate)
}

func validBankJSON() string {
	return `{"code":"test-bank","name":"Test Bank","is_active":true}`
}
