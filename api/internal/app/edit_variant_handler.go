package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/core/service"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// EditVariant handler for editing variant
func (h *Handler) EditVariant(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewEditVariant(r)
	if err != nil {
		if err == request.ErrInvalidVariantID {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidVariantID.Error())
			return
		}
		if err == request.ErrInvalidJSON {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidJSON.Error())
			return
		}
		if err == request.ErrInvalidProductID {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidProductID.Error())
			return
		}
		if err == request.ErrInvalidMRP {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidMRP.Error())
			return
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, "Request body is invalid")
		return
	}
	err = h.Variant.EditVariant(r.Context(), req)
	if err != nil {
		if err == service.ErrVariantNotFound {
			log.Printf(err.Error())
			api.Fail(w, http.StatusNotFound, http.StatusNotFound, service.ErrVariantNotFound.Error())
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}

	api.Send(w, http.StatusOK, nil)
}
