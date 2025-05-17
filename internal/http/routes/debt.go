package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/container"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/controllers"
)

func DebtRoutes(d *container.Dependencies) http.Handler {
	r := chi.NewRouter()
	debtController := controllers.NewDebtController(d.DebtService)
	r.Post("/", debtController.CreateDebt())
	r.Get("/{clientId}", debtController.GetClientUserDebts())
	return r
}
