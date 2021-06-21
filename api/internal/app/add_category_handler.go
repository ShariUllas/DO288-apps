package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// AddCategory handler for adding categories
func (h *Handler) AddCategory(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewAddCategory(r)
	if err != nil {
		if err == request.ErrInvalidCategory {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidCategory.Error())
			return
		} else if err == request.ErrInvalidJSON {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidJSON.Error())
			return
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	res, err := h.Category.AddCategory(r.Context(), req)
	if err != nil {
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	api.Send(w, http.StatusOK, res)
}
