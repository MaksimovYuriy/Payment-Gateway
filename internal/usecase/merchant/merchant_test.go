package merchant_test

import (
	"context"
	"errors"
	"testing"

	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apikey"
	"payment_gateway/internal/lib/apperr"
	merchantusecase "payment_gateway/internal/usecase/merchant"

	"github.com/go-playground/validator/v10"
)

type mockMerchantRepo struct {
	created      *entity.Merchant
	createCalled bool
	err          error
}

func (r *mockMerchantRepo) Create(_ context.Context, m *entity.Merchant) error {
	r.createCalled = true
	if r.err != nil {
		return r.err
	}
	r.created = m
	m.Id = 1
	return nil
}

func (r *mockMerchantRepo) List(_ context.Context) ([]*entity.Merchant, error) {
	return nil, nil
}

func (r *mockMerchantRepo) GetById(_ context.Context, _ int64) (*entity.Merchant, error) {
	return nil, nil
}

func TestRegistrationReturnsAPIKeyAndStoresHash(t *testing.T) {
	repo := &mockMerchantRepo{}
	uc := merchantusecase.NewUseCase(repo, validator.New())

	merchant := &entity.Merchant{
		Name:               "Test Merchant",
		Domain:             "example.com",
		WebhookUrl:         "https://example.com/webhook",
		SuccessRedirectUrl: "https://example.com/success",
		FailureRedirectUrl: "https://example.com/failure",
		IsActive:           true,
	}

	created, rawAPIKey, err := uc.Registration(context.Background(), merchant)
	if err != nil {
		t.Fatalf("Registration() error = %v", err)
	}

	if created == nil {
		t.Fatal("Registration() returned nil merchant")
	}
	if rawAPIKey == "" {
		t.Fatal("Registration() returned empty api key")
	}
	if repo.created == nil {
		t.Fatal("Registration() did not create merchant in repo")
	}
	if repo.created.ApiKeyHash == "" {
		t.Fatal("Registration() did not store api key hash")
	}
	if repo.created.ApiKeyHash == rawAPIKey {
		t.Fatal("Registration() stored raw api key instead of hash")
	}
	if !apikey.CheckApiKey(rawAPIKey, repo.created.ApiKeyHash) {
		t.Fatal("Registration() stored hash that does not match returned api key")
	}
}

func TestRegistrationRejectsInvalidMerchant(t *testing.T) {
	repo := &mockMerchantRepo{}
	uc := merchantusecase.NewUseCase(repo, validator.New())

	_, _, err := uc.Registration(context.Background(), &entity.Merchant{})
	if err == nil {
		t.Fatal("Registration() error = nil, want error")
	}

	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("Registration() error type = %T, want *apperr.Error", err)
	}
	if appErr.Code != apperr.CodeInvalidInput {
		t.Fatalf("Registration() error code = %s, want %s", appErr.Code, apperr.CodeInvalidInput)
	}
	if repo.createCalled {
		t.Fatal("Registration() called repo.Create for invalid merchant")
	}
}

func TestRegistrationReturnsRepoError(t *testing.T) {
	wantErr := apperr.AlreadyExists("merchant already exists")
	repo := &mockMerchantRepo{err: wantErr}
	uc := merchantusecase.NewUseCase(repo, validator.New())

	merchant := &entity.Merchant{
		Name:               "Test Merchant",
		Domain:             "example.com",
		WebhookUrl:         "https://example.com/webhook",
		SuccessRedirectUrl: "https://example.com/success",
		FailureRedirectUrl: "https://example.com/failure",
		IsActive:           true,
	}

	_, _, err := uc.Registration(context.Background(), merchant)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Registration() error = %v, want %v", err, wantErr)
	}
}
