package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/core/client"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/customvalidate"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
)

type ClientController struct {
	ClientService client.Service
}

func NewClientController(clientService client.Service) *ClientController {
	return &ClientController{
		ClientService: clientService,
	}
}

func (c *ClientController) Create() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cliRequest client.ClientRequestDto

		if err := json.NewDecoder(r.Body).Decode(&cliRequest); err != nil {
			response(w, http.StatusBadRequest, "Invalid request")
			return
		}

		v := customvalidate.Validate(cliRequest)
		if len(v.Errors) > 0 {
			response(w, http.StatusUnprocessableEntity, v)
			return
		}

		output := c.ClientService.Create(r.Context(), &cliRequest)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output)
			return
		}

		response(w, http.StatusCreated, output)
	})
}

func (c *ClientController) Update() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientId := chi.URLParam(r, "clientId")
		if clientId == "" {
			response(w, http.StatusBadRequest, "Missing client ID")
			return
		}

		clientIdParsed, err := ulid.Parse(clientId)
		if err != nil {
			response(w, http.StatusBadRequest, "Invalid client ID")
			return
		}

		var cliRequest client.ClientRequestDto
		if err := json.NewDecoder(r.Body).Decode(&cliRequest); err != nil {
			response(w, http.StatusBadRequest, "Invalid request")
			return
		}

		v := customvalidate.Validate(cliRequest)
		if len(v.Errors) > 0 {
			response(w, http.StatusUnprocessableEntity, v)
			return
		}

		output := c.ClientService.Update(r.Context(), clientIdParsed, &cliRequest)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output)
			return
		}

		response(w, http.StatusOK, output)
	})
}

func (c *ClientController) Delete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := chi.URLParam(r, "clientId")
		if clientId == "" {
			response(w, http.StatusBadRequest, "Missing client ID")
			return
		}

		clientIdParsed, err := ulid.Parse(clientId)
		if err != nil {
			response(w, http.StatusBadRequest, "Invalid client ID")
			return
		}

		output := c.ClientService.Delete(r.Context(), clientIdParsed)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output)
			return
		}

		response(w, http.StatusNoContent, nil)
	})
}

func (c *ClientController) FindOne() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := chi.URLParam(r, "clientId")
		if clientId == "" {
			response(w, http.StatusBadRequest, "Missing client ID")
			return
		}

		clientIdParsed, err := ulid.Parse(clientId)
		if err != nil {
			response(w, http.StatusBadRequest, "Invalid client ID")
			return
		}

		output := c.ClientService.FindById(r.Context(), clientIdParsed)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output)
			return
		}

		response(w, http.StatusOK, output)
	})
}

func (c *ClientController) FindAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pgRequest, err := paginate.GetPaginateParams(r)
		if err != nil {
			log.Println("Error getting pagination params:", err)
			response(w, http.StatusBadRequest, "Invalid pagination params")
			return
		}

		output := c.ClientService.FindByCriteria(r.Context(), pgRequest)
		if output.Status == "error" {
			response(w, http.StatusInternalServerError, output)
			return
		}

		response(w, http.StatusOK, output)
	})
}
