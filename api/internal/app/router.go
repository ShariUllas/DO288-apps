package app

import (
	"catalogue/api/internal/config"
	"catalogue/api/internal/core/service"

	"github.com/go-chi/chi"
	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi"
)

// Handler holds tableview sevice interface
type Handler struct {
	*config.Config
	service.Category
	service.Product
	service.Variant
}

// InitRouter sets up tableview's router.
func (h *Handler) InitRouter() chi.Router {
	r := chi.NewRouter()
	if h.DatadogEnableAPM {
		r.Use(chitrace.Middleware(
			chitrace.WithServiceName(h.DatadogServiceName),
		))
	}
	r.Route("/variant", func(r chi.Router) {
		r.Get("/", h.GetVariant)
		r.Post("/add", h.AddVariant)
		r.Route("/{variant_id}", func(r chi.Router) {
			r.Put("/", h.EditVariant)
			r.Delete("/", h.DeleteVariant)
		})
	})
	r.Route("/category", func(r chi.Router) {
		r.Get("/", h.GetCategory)

		r.Post("/add", h.AddCategory)
		r.Route("/{category_id}", func(r chi.Router) {
			r.Get("/list", h.ListCategory)
			r.Put("/", h.EditCategory)
			r.Delete("/", h.DeleteCategory)
		})
	})
	r.Route("/product", func(r chi.Router) {
		r.Get("/", h.GetProduct)
		r.Post("/add", h.AddProduct)
	})
	return r
}
