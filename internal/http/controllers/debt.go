package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/customvalidate"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
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

func (c *DebtController) GetDebtInstallments() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := chi.URLParam(r, "clientId")
		if clientId == "" {
			response(w, http.StatusBadRequest, "clientId is required")
			return
		}

		clientIdParsed, err := ulid.Parse(clientId)
		if err != nil {
			log.Println(err)
			response(w, http.StatusBadRequest, "clientId invalid")
			return
		}

		debtId := chi.URLParam(r, "debtId")
		if debtId == "" {
			response(w, http.StatusBadRequest, "debtId is required")
		}

		debtIdParsed, err := ulid.Parse(debtId)
		if err != nil {
			log.Println(err)
			response(w, http.StatusBadRequest, "debtId invalid")
			return
		}

		result := c.DebtService.GetDebtInstallments(r.Context(), clientIdParsed, debtIdParsed)

		if result.Status == "error" {
			response(w, http.StatusInternalServerError, result.Message)
			return
		}

		response(w, http.StatusOK, result)
	})
}

func (c *DebtController) GetDebts() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pgRequest, err := paginate.GetPaginateParams(r)
		if err != nil {
			log.Println("Error getting pagination params:", err)
			response(w, http.StatusBadRequest, "Invalid pagination params")
			return
		}

		result := c.DebtService.Debts(r.Context(), *pgRequest)

		if result.Status == "error" {
			response(w, http.StatusInternalServerError, result.Message)
			return
		}

		response(w, http.StatusOK, result)
	})
}

func (c *DebtController) PayInstallment() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var paymentInfo debt.PaymentInfoDto
		if err := json.NewDecoder(r.Body).Decode(&paymentInfo); err != nil {
			response(w, http.StatusBadRequest, "Invalid request")
			return
		}

		v := customvalidate.Validate(paymentInfo)
		if len(v.Errors) > 0 {
			response(w, http.StatusUnprocessableEntity, v)
			return
		}

		output := c.DebtService.PayInstallment(r.Context(), &paymentInfo)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output.Message)
			return
		}

		response(w, http.StatusOK, output)
	})
}

func (c *DebtController) CancelDebt() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cancelInfo debt.CancelInfoDto
		if err := json.NewDecoder(r.Body).Decode(&cancelInfo); err != nil {
			response(w, http.StatusBadRequest, "Invalid request")
			return
		}

		v := customvalidate.Validate(cancelInfo)
		if len(v.Errors) > 0 {
			response(w, http.StatusUnprocessableEntity, v)
			return
		}

		// TODO: Ajustar para colocar o ID do usuário autenticado
		cancelInfo.CancelledBy = ulid.Make()

		output := c.DebtService.CancelDebt(r.Context(), &cancelInfo)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output.Message)
			return
		}

		response(w, http.StatusOK, output)
	})
}

func (c *DebtController) ReversalDebt() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reversalInfo debt.ReversalInfoDto
		if err := json.NewDecoder(r.Body).Decode(&reversalInfo); err != nil {
			response(w, http.StatusBadRequest, "Invalid request")
			return
		}

		v := customvalidate.Validate(reversalInfo)
		if len(v.Errors) > 0 {
			response(w, http.StatusUnprocessableEntity, v)
			return
		}

		// TODO: Ajustar para colocar o ID do usuário autenticado
		reversalInfo.ReversedBy = ulid.Make()

		output := c.DebtService.ReverseDebt(r.Context(), &reversalInfo)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output.Message)
			return
		}

		response(w, http.StatusOK, output)
	})
}
