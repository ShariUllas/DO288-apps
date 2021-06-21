package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/core/service"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// EditCategory handler for editing variant
func (h *Handler) EditCategory(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewEditCategory(r)
	if err != nil {
		if err == request.ErrInvalidCategoryID {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidCategoryID.Error())
			return
		}
		if err == request.ErrInvalidJSON {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidJSON.Error())
			return
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, "Request body is invalid")
		return
	}
	err = h.Category.EditCategory(r.Context(), req)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			log.Printf(err.Error())
			api.Fail(w, http.StatusNotFound, http.StatusNotFound, service.ErrCategoryNotFound.Error())
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}

	api.Send(w, http.StatusOK, nil)
}
