package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// AddVariant handler for adding variants
func (h *Handler) AddVariant(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewAddVariant(r)
	if err != nil {
		if err == request.ErrInvalidMRP {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidMRP.Error())
			return
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	res, err := h.Variant.AddVariant(r.Context(), req)
	if err != nil {
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	api.Send(w, http.StatusOK, res)
}
