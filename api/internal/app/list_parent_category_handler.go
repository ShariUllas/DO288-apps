package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

func (h *Handler) ListCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := request.NewListCategory(r)
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusBadRequest, http.StatusBadRequest)
		return
	}
	res, err := h.Category.ListCategory(r.Context(), categoryID)
	if err != nil {
		log.Println(err)
		api.Fail(w, http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	api.Send(w, http.StatusOK, res)
}
