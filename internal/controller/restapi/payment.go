package restapi

import (
	"encoding/json"
	"net/http"
	"payment_gateway/internal/controller/restapi/request"
	"payment_gateway/internal/controller/restapi/response"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/usecase/payment"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PaymentController struct {
	usecase *payment.UseCase
}

func NewPaymentController(uc *payment.UseCase) *PaymentController {
	return &PaymentController{usecase: uc}
}

func (pc *PaymentController) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreatePayment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	paymentEntity := entity.Payment{
		MerchantId: req.MerchantId,
		OrderId:    req.OrderId,
		Amount:     req.Amount,
		Currency:   req.Currency,
	}

	payment, err := pc.usecase.Create(r.Context(), &paymentEntity, req.BankId)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.NewPayment(*payment))
}

func (pc *PaymentController) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	payment, err := pc.usecase.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewPayment(*payment))
}

func (pc *PaymentController) List(w http.ResponseWriter, r *http.Request) {
	payments, err := pc.usecase.List(r.Context())
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewListPayment(payments))
}
