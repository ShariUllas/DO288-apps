package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// DeleteCategory handler for deleting category
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := request.NewDeleteCategory(r)
	if err != nil {
		if err == request.ErrInvalidCategoryID {
			log.Println(err)
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidCategoryID.Error())
			return
		}
		log.Println(err)
		api.Fail(w, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	err = h.Category.DeleteCategory(r.Context(), id)
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}

	api.Send(w, http.StatusOK, nil)
}
