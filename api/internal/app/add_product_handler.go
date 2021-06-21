package app

import (
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/helper/api"
	"log"
	"net/http"
)

// AddProduct handler for adding products
func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewAddProduct(r)
	if err != nil {
		if err == request.ErrInvalidProductName {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidProductName.Error())
			return
		} else if err == request.ErrInvalidCategoryID {
			log.Printf(err.Error())
			api.Fail(w, http.StatusBadRequest, http.StatusBadRequest, request.ErrInvalidCategoryID.Error())
			return
		}
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	res, err := h.Product.AddProduct(r.Context(), req)
	if err != nil {
		log.Printf(err.Error())
		api.Fail(w, http.StatusInternalServerError, api.ErrCodeInternalServiceError)
		return
	}
	api.Send(w, http.StatusOK, res)
}
