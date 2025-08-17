package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/container"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/controllers"
)

func ClientRoutes(d *container.Dependencies) http.Handler {
	r := chi.NewRouter()
	clientController := controllers.NewClientController(d.ClientService)

	r.Post("/", clientController.Create())
	r.Put("/{clientId}", clientController.Update())
	r.Delete("/{clientId}", clientController.Delete())
	r.Get("/{clientId}", clientController.FindOne())
	r.Get("/", clientController.FindAll())

	return r
}
