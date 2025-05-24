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
	r.Get("/", debtController.GetDebts())
	r.Get("/{clientId}", debtController.GetClientUserDebts())
	r.Get("/{clientId}/{debtId}/installments", debtController.GetDebtInstallments())
	return r
}
