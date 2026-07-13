package payment_test

import (
	"context"
	"errors"
	"testing"

	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
	"payment_gateway/internal/lib/bankprocessor"
	paymentusecase "payment_gateway/internal/usecase/payment"

	"github.com/go-playground/validator/v10"
)

type mockPaymentRepo struct {
	created       *entity.Payment
	updates       []*entity.Payment
	createCalled  bool
	createErr     error
	updateErr     error
	nextPaymentID int64
}

func (r *mockPaymentRepo) Create(_ context.Context, p *entity.Payment) error {
	r.createCalled = true
	if r.createErr != nil {
		return r.createErr
	}
	if r.nextPaymentID == 0 {
		r.nextPaymentID = 1
	}
	p.Id = r.nextPaymentID
	p.Status = entity.PaymentStatusCreated
	r.created = clonePayment(p)
	return nil
}

func (r *mockPaymentRepo) Update(_ context.Context, p *entity.Payment) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	r.updates = append(r.updates, clonePayment(p))
	return nil
}

func (r *mockPaymentRepo) GetById(_ context.Context, _ int64) (*entity.Payment, error) {
	return nil, nil
}

func (r *mockPaymentRepo) List(_ context.Context) ([]*entity.Payment, error) {
	return nil, nil
}

type mockPaymentAttemptRepo struct {
	created       *entity.PaymentAttempt
	updates       []*entity.PaymentAttempt
	nextAttemptID int64
	createErr     error
	updateErr     error
}

func (r *mockPaymentAttemptRepo) Create(_ context.Context, pa *entity.PaymentAttempt) error {
	if r.createErr != nil {
		return r.createErr
	}
	if r.nextAttemptID == 0 {
		r.nextAttemptID = 1
	}
	pa.Id = r.nextAttemptID
	pa.Status = entity.PAttemptStatusCreated
	r.created = clonePaymentAttempt(pa)
	return nil
}

func (r *mockPaymentAttemptRepo) Update(_ context.Context, pa *entity.PaymentAttempt) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	r.updates = append(r.updates, clonePaymentAttempt(pa))
	return nil
}

type mockBankRepo struct {
	bank *entity.Bank
	err  error
}

func (r *mockBankRepo) Create(_ context.Context, _ *entity.Bank) error {
	return nil
}

func (r *mockBankRepo) List(_ context.Context) ([]*entity.Bank, error) {
	return nil, nil
}

func (r *mockBankRepo) GetById(_ context.Context, _ int64) (*entity.Bank, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.bank, nil
}

type mockMerchantRepo struct {
	merchant *entity.Merchant
	err      error
}

func (r *mockMerchantRepo) Create(_ context.Context, _ *entity.Merchant) error {
	return nil
}

func (r *mockMerchantRepo) List(_ context.Context) ([]*entity.Merchant, error) {
	return nil, nil
}

func (r *mockMerchantRepo) GetById(_ context.Context, _ int64) (*entity.Merchant, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.merchant, nil
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

func TestCreateCompletesPayment(t *testing.T) {
	uc, paymentRepo, attemptRepo := newPaymentUseCase(
		&entity.Bank{Id: 1, IsActive: true},
		&entity.Merchant{Id: 2, IsActive: true},
		bankprocessor.Result{Success: true, ExternalPaymentId: "ext-1"},
	)

	payment := validPayment()
	created, err := uc.Create(context.Background(), payment, 1)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.Status != entity.PaymentStatusCompleted {
		t.Fatalf("payment status = %s, want %s", created.Status, entity.PaymentStatusCompleted)
	}
	if created.SuccessBankId == nil || *created.SuccessBankId != 1 {
		t.Fatalf("payment success bank id = %v, want 1", created.SuccessBankId)
	}
	if len(paymentRepo.updates) != 2 {
		t.Fatalf("payment updates = %d, want 2", len(paymentRepo.updates))
	}
	if paymentRepo.updates[0].Status != entity.PaymentStatusProcessing {
		t.Fatalf("first payment update status = %s, want %s", paymentRepo.updates[0].Status, entity.PaymentStatusProcessing)
	}
	if paymentRepo.updates[1].Status != entity.PaymentStatusCompleted {
		t.Fatalf("second payment update status = %s, want %s", paymentRepo.updates[1].Status, entity.PaymentStatusCompleted)
	}
	if len(attemptRepo.updates) != 2 {
		t.Fatalf("attempt updates = %d, want 2", len(attemptRepo.updates))
	}
	if attemptRepo.updates[1].Status != entity.PAttemptStatusSucceeded {
		t.Fatalf("final attempt status = %s, want %s", attemptRepo.updates[1].Status, entity.PAttemptStatusSucceeded)
	}
	if attemptRepo.updates[1].ExternalPaymentId == nil || *attemptRepo.updates[1].ExternalPaymentId != "ext-1" {
		t.Fatalf("external payment id = %v, want ext-1", attemptRepo.updates[1].ExternalPaymentId)
	}
}

func TestCreateRejectsInactiveBank(t *testing.T) {
	uc, paymentRepo, _ := newPaymentUseCase(
		&entity.Bank{Id: 1, IsActive: false},
		&entity.Merchant{Id: 2, IsActive: true},
		bankprocessor.Result{Success: true},
	)

	_, err := uc.Create(context.Background(), validPayment(), 1)
	assertAppErrorCode(t, err, apperr.CodeInvalidInput)
	if paymentRepo.createCalled {
		t.Fatal("Create() created payment for inactive bank")
	}
}

func TestCreateRejectsInactiveMerchant(t *testing.T) {
	uc, paymentRepo, _ := newPaymentUseCase(
		&entity.Bank{Id: 1, IsActive: true},
		&entity.Merchant{Id: 2, IsActive: false},
		bankprocessor.Result{Success: true},
	)

	_, err := uc.Create(context.Background(), validPayment(), 1)
	assertAppErrorCode(t, err, apperr.CodeInvalidInput)
	if paymentRepo.createCalled {
		t.Fatal("Create() created payment for inactive merchant")
	}
}

func TestCreateRejectsInvalidPayment(t *testing.T) {
	uc, paymentRepo, _ := newPaymentUseCase(
		&entity.Bank{Id: 1, IsActive: true},
		&entity.Merchant{Id: 2, IsActive: true},
		bankprocessor.Result{Success: true},
	)

	payment := validPayment()
	payment.Amount = "0"

	_, err := uc.Create(context.Background(), payment, 1)
	assertAppErrorCode(t, err, apperr.CodeInvalidInput)
	if paymentRepo.createCalled {
		t.Fatal("Create() created invalid payment")
	}
}

func TestCreateMarksPaymentFailedWhenProcessorFails(t *testing.T) {
	uc, paymentRepo, attemptRepo := newPaymentUseCase(
		&entity.Bank{Id: 1, IsActive: true},
		&entity.Merchant{Id: 2, IsActive: true},
		bankprocessor.Result{
			Success:      false,
			ErrorCode:    "bank_declined",
			ErrorMessage: "bank declined payment",
		},
	)

	created, err := uc.Create(context.Background(), validPayment(), 1)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.Status != entity.PaymentStatusFailed {
		t.Fatalf("payment status = %s, want %s", created.Status, entity.PaymentStatusFailed)
	}
	if paymentRepo.updates[1].Status != entity.PaymentStatusFailed {
		t.Fatalf("final payment update status = %s, want %s", paymentRepo.updates[1].Status, entity.PaymentStatusFailed)
	}
	if attemptRepo.updates[1].Status != entity.PAttemptStatusFailed {
		t.Fatalf("final attempt status = %s, want %s", attemptRepo.updates[1].Status, entity.PAttemptStatusFailed)
	}
	if attemptRepo.updates[1].ErrorCode == nil || *attemptRepo.updates[1].ErrorCode != "bank_declined" {
		t.Fatalf("attempt error code = %v, want bank_declined", attemptRepo.updates[1].ErrorCode)
	}
}

func newPaymentUseCase(
	bank *entity.Bank,
	merchant *entity.Merchant,
	processorResult bankprocessor.Result,
) (*paymentusecase.UseCase, *mockPaymentRepo, *mockPaymentAttemptRepo) {
	paymentRepo := &mockPaymentRepo{}
	attemptRepo := &mockPaymentAttemptRepo{}
	uow := &mockUnitOfWork{
		repos: paymentusecase.Repositories{
			Payment:        paymentRepo,
			PaymentAttempt: attemptRepo,
		},
	}

	uc := paymentusecase.NewUseCase(
		paymentRepo,
		attemptRepo,
		&mockBankRepo{bank: bank},
		&mockMerchantRepo{merchant: merchant},
		mockBankProcessor{result: processorResult},
		uow,
		validator.New(),
	)

	return uc, paymentRepo, attemptRepo
}

func validPayment() *entity.Payment {
	return &entity.Payment{
		MerchantId: 2,
		OrderId:    "order-1",
		Amount:     "10.50",
		Currency:   "USD",
	}
}

func assertAppErrorCode(t *testing.T, err error, code apperr.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want %s", code)
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type = %T, want *apperr.Error", err)
	}
	if appErr.Code != code {
		t.Fatalf("error code = %s, want %s", appErr.Code, code)
	}
}

func clonePayment(p *entity.Payment) *entity.Payment {
	clone := *p
	if p.SuccessBankId != nil {
		successBankID := *p.SuccessBankId
		clone.SuccessBankId = &successBankID
	}
	return &clone
}

func clonePaymentAttempt(pa *entity.PaymentAttempt) *entity.PaymentAttempt {
	clone := *pa
	if pa.ExternalPaymentId != nil {
		externalPaymentID := *pa.ExternalPaymentId
		clone.ExternalPaymentId = &externalPaymentID
	}
	if pa.ErrorMessage != nil {
		errorMessage := *pa.ErrorMessage
		clone.ErrorMessage = &errorMessage
	}
	if pa.ErrorCode != nil {
		errorCode := *pa.ErrorCode
		clone.ErrorCode = &errorCode
	}
	return &clone
}
