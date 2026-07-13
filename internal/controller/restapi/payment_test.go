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
	"payment_gateway/internal/lib/bankprocessor"
	paymentusecase "payment_gateway/internal/usecase/payment"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type mockPaymentRepo struct {
	created *entity.Payment
	list    []*entity.Payment
	get     *entity.Payment
	err     error
	listErr error
	getErr  error
}

func (r *mockPaymentRepo) Create(_ context.Context, p *entity.Payment) error {
	if r.err != nil {
		return r.err
	}
	p.Id = 1
	p.Status = entity.PaymentStatusCreated
	r.created = p
	return nil
}

func (r *mockPaymentRepo) Update(_ context.Context, _ *entity.Payment) error {
	return r.err
}

func (r *mockPaymentRepo) GetById(_ context.Context, _ int64) (*entity.Payment, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.get, nil
}

func (r *mockPaymentRepo) List(_ context.Context) ([]*entity.Payment, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.list, nil
}

type mockPaymentAttemptRepo struct {
	err error
}

func (r *mockPaymentAttemptRepo) Create(_ context.Context, pa *entity.PaymentAttempt) error {
	if r.err != nil {
		return r.err
	}
	pa.Id = 1
	pa.Status = entity.PAttemptStatusCreated
	return nil
}

func (r *mockPaymentAttemptRepo) Update(_ context.Context, _ *entity.PaymentAttempt) error {
	return r.err
}

type mockUnitOfWork struct {
	repos paymentusecase.Repositories
}

func (u *mockUnitOfWork) Do(ctx context.Context, fn func(context.Context, paymentusecase.Repositories) error) error {
	return fn(ctx, u.repos)
}

type mockBankProcessor struct {
	result bankprocessor.Result
}

func (p mockBankProcessor) Process(_ context.Context, _ *entity.Payment, _ int64) bankprocessor.Result {
	return p.result
}

func TestPaymentCreateReturnsCreated(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{get: &entity.Bank{Id: 1, IsActive: true}},
		&mockMerchantRepo{get: &entity.Merchant{Id: 2, IsActive: true}},
	)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(validPaymentJSON()))
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var got response.Payment
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Id != 1 {
		t.Fatalf("payment id = %d, want 1", got.Id)
	}
	if got.Status != string(entity.PaymentStatusCompleted) {
		t.Fatalf("payment status = %q, want %q", got.Status, entity.PaymentStatusCompleted)
	}
}

func TestPaymentCreateRejectsInvalidJSON(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{get: &entity.Bank{Id: 1, IsActive: true}},
		&mockMerchantRepo{get: &entity.Merchant{Id: 2, IsActive: true}},
	)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestPaymentCreateRejectsInvalidRequest(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{get: &entity.Bank{Id: 1, IsActive: true}},
		&mockMerchantRepo{get: &entity.Merchant{Id: 2, IsActive: true}},
	)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(`{"merchant_id":0,"bank_id":1,"order_id":"order-1","amount":"10.50","currency":"USD"}`))
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestPaymentCreateMapsUseCaseError(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{getByErr: apperr.NotFound("bank not found")},
		&mockMerchantRepo{get: &entity.Merchant{Id: 2, IsActive: true}},
	)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(validPaymentJSON()))
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assertErrorResponse(t, rec, http.StatusNotFound, apperr.CodeNotFound)
}

func TestPaymentGetByIdReturnsOK(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{get: &entity.Payment{Id: 1, MerchantId: 2, OrderId: "order-1", Amount: "10.50", Currency: "USD", Status: entity.PaymentStatusCompleted}},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{},
		&mockMerchantRepo{},
	)
	req := requestWithURLParam("id", "1")
	rec := httptest.NewRecorder()

	controller.GetById(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got response.Payment
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Id != 1 {
		t.Fatalf("payment id = %d, want 1", got.Id)
	}
}

func TestPaymentGetByIdRejectsInvalidID(t *testing.T) {
	controller := newPaymentController(&mockPaymentRepo{}, &mockPaymentAttemptRepo{}, &mockBankRepo{}, &mockMerchantRepo{})
	req := requestWithURLParam("id", "bad-id")
	rec := httptest.NewRecorder()

	controller.GetById(rec, req)

	assertErrorResponse(t, rec, http.StatusBadRequest, apperr.CodeInvalidInput)
}

func TestPaymentGetByIdMapsUseCaseError(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{getErr: apperr.NotFound("payment not found")},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{},
		&mockMerchantRepo{},
	)
	req := requestWithURLParam("id", "1")
	rec := httptest.NewRecorder()

	controller.GetById(rec, req)

	assertErrorResponse(t, rec, http.StatusNotFound, apperr.CodeNotFound)
}

func TestPaymentListReturnsOK(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{list: []*entity.Payment{
			{Id: 1, MerchantId: 2, OrderId: "order-1", Amount: "10.50", Currency: "USD", Status: entity.PaymentStatusCompleted},
			{Id: 2, MerchantId: 2, OrderId: "order-2", Amount: "20.00", Currency: "USD", Status: entity.PaymentStatusFailed},
		}},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{},
		&mockMerchantRepo{},
	)
	req := httptest.NewRequest(http.MethodGet, "/payments", nil)
	rec := httptest.NewRecorder()

	controller.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got response.ListPayment
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got.Payments) != 2 {
		t.Fatalf("payments len = %d, want 2", len(got.Payments))
	}
}

func TestPaymentListMapsUseCaseError(t *testing.T) {
	controller := newPaymentController(
		&mockPaymentRepo{listErr: apperr.Internal("failed to list payments", nil)},
		&mockPaymentAttemptRepo{},
		&mockBankRepo{},
		&mockMerchantRepo{},
	)
	req := httptest.NewRequest(http.MethodGet, "/payments", nil)
	rec := httptest.NewRecorder()

	controller.List(rec, req)

	assertErrorResponse(t, rec, http.StatusInternalServerError, apperr.CodeInternal)
}

func newPaymentController(
	paymentRepo *mockPaymentRepo,
	attemptRepo *mockPaymentAttemptRepo,
	bankRepo *mockBankRepo,
	merchantRepo *mockMerchantRepo,
) *restapi.PaymentController {
	validate := validator.New()
	uow := &mockUnitOfWork{
		repos: paymentusecase.Repositories{
			Payment:        paymentRepo,
			PaymentAttempt: attemptRepo,
		},
	}
	uc := paymentusecase.NewUseCase(
		paymentRepo,
		attemptRepo,
		bankRepo,
		merchantRepo,
		mockBankProcessor{result: bankprocessor.Result{Success: true, ExternalPaymentId: "ext-1"}},
		uow,
		validate,
	)
	return restapi.NewPaymentController(uc, validate)
}

func validPaymentJSON() string {
	return `{"merchant_id":2,"bank_id":1,"order_id":"order-1","amount":"10.50","currency":"USD"}`
}

func requestWithURLParam(key, value string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/payments/"+value, nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}
