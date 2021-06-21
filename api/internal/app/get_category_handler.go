package app

import (
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// GetCategory handler
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	result, err := h.Category.GetCategory(r.Context())
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	api.Send(w, http.StatusOK, result)
}
