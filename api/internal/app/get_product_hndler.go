package app

import (
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// GetProduct handler for listing products
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	result, err := h.Product.GetProduct(r.Context())
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	api.Send(w, http.StatusOK, result)
}
