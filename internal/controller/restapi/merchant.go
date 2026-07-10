package restapi

import (
	"encoding/json"
	"net/http"
	"payment_gateway/internal/controller/restapi/request"
	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
	"payment_gateway/internal/usecase/merchant"

	"github.com/go-playground/validator/v10"
)

type MerchantController struct {
	usecase  *merchant.UseCase
	validate *validator.Validate
}

func NewMerchantController(uc *merchant.UseCase, vd *validator.Validate) *MerchantController {
	return &MerchantController{usecase: uc, validate: vd}
}

func (mc *MerchantController) Registration(w http.ResponseWriter, r *http.Request) {
	var req request.RegistrationMerchant
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, apperr.InvalidInput(apperr.MessageInvalidInput))
		return
	}
	if err := mc.validate.Struct(req); err != nil {
		WriteError(w, apperr.InvalidInput(apperr.MessageInvalidInput))
		return
	}

	merchant_entity := entity.Merchant{
		Name:               req.Name,
		Domain:             req.Domain,
		WebhookUrl:         req.WebhookUrl,
		SuccessRedirectUrl: req.SuccessRedirectUrl,
		FailureRedirectUrl: req.FailureRedirectUrl,
		IsActive:           req.IsActive,
	}

	merchant, apiKey, err := mc.usecase.Registration(r.Context(), &merchant_entity)
	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response.NewMerchantRegistration(*merchant, apiKey))
}

func (mc *MerchantController) List(w http.ResponseWriter, r *http.Request) {
	merchants, err := mc.usecase.List(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.NewListMerchant(merchants))
}
