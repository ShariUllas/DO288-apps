package app

import (
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// GetVariant handler
func (h *Handler) GetVariant(w http.ResponseWriter, r *http.Request) {
	result, err := h.Variant.GetVariant(r.Context())
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	api.Send(w, http.StatusOK, result)
}
