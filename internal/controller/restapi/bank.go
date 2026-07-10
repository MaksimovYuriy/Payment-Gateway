package restapi

import (
	"encoding/json"
	"net/http"
	"payment_gateway/internal/controller/restapi/request"
	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/usecase/bank"

	"github.com/go-playground/validator/v10"
)

type BankController struct {
	usecase  *bank.UseCase
	validate *validator.Validate
}

func NewBankController(uc *bank.UseCase, vd *validator.Validate) *BankController {
	return &BankController{usecase: uc, validate: vd}
}

func (bc *BankController) Registration(w http.ResponseWriter, r *http.Request) {
	var req request.RegistrationBank
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	if err := bc.validate.Struct(req); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	bankEntity := entity.Bank{
		Code:     req.Code,
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	bank, err := bc.usecase.Registration(r.Context(), &bankEntity)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.NewBank(*bank))
}

func (bc *BankController) List(w http.ResponseWriter, r *http.Request) {
	banks, err := bc.usecase.List(r.Context())
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewListBank(banks))
}
