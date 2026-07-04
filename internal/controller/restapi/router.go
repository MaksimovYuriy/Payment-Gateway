package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(bankController *BankController) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", healthHandler)

	r.Post("/admin/banks", bankController.Registration)
	r.Get("/admin/banks", bankController.List)

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "ok",
		"service": "payment-gateway",
	}

	_ = json.NewEncoder(w).Encode(response)
}
