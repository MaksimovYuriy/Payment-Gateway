package restapi

import (
	"encoding/json"
	"net/http"
	"payment_gateway/internal/controller/restapi/request"
	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/usecase/merchant"
)

type MerchantController struct {
	usecase *merchant.UseCase
}

func NewMerchantController(uc *merchant.UseCase) *MerchantController {
	return &MerchantController{usecase: uc}
}

func (mc *MerchantController) Registration(w http.ResponseWriter, r *http.Request) {
	var req request.RegistrationMerchant
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
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

	merchant, err := mc.usecase.Registration(r.Context(), &merchant_entity)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.NewMerchant(*merchant))
}

func (mc *MerchantController) List(w http.ResponseWriter, r *http.Request) {
	merchants, err := mc.usecase.List(r.Context())
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewListMerchant(merchants))
}
