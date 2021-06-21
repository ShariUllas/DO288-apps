package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// DeleteVariant handler for delting variants
func (h *Handler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	id, err := request.NewDeleteVariant(r)
	if err != nil {
		if err == request.ErrInvalidVariantID {
			log.Println(err)
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidVariantID.Error())
			return
		}
		log.Println(err)
		api.Fail(w, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	err = h.Variant.DeleteVariant(r.Context(), id)
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}

	api.Send(w, http.StatusOK, nil)
}
