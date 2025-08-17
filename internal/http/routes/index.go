package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/container"
)

func Start(d *container.Dependencies) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/debt", DebtRoutes(d))
		r.Mount("/client", ClientRoutes(d))
	})

	return r
}
