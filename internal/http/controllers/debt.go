package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/customvalidate"
	"github.com/oklog/ulid/v2"
)

type DebtController struct {
	DebtService debt.Service
}

func NewDebtController(debtService debt.Service) *DebtController {
	return &DebtController{
		DebtService: debtService,
	}
}

func (c *DebtController) CreateDebt() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request debt.DebtDto
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response(w, http.StatusBadRequest, "Invalid request")
			return
		}

		v := customvalidate.Validate(request)
		if len(v.Errors) > 0 {
			response(w, http.StatusUnprocessableEntity, v)
			return
		}

		output := c.DebtService.CreateDebt(r.Context(), &request)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output.Data)
			return
		}

		response(w, http.StatusCreated, output)
	})
}

func (c *DebtController) GetClientUserDebts() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := chi.URLParam(r, "clientId")
		if clientId == "" {
			response(w, http.StatusBadRequest, "clientId is required")
			return
		}

		parsedClientId, err := ulid.Parse(clientId)
		if err != nil {
			response(w, http.StatusBadRequest, "Invalid clientId")
			return
		}

		output := c.DebtService.GetUserDebts(r.Context(), parsedClientId)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output.Data)
			return
		}

		response(w, http.StatusOK, output)
	})
}
