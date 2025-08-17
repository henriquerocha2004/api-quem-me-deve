package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/henriquerocha2004/quem-me-deve-api/core/client"
	"github.com/henriquerocha2004/quem-me-deve-api/core/client/mocks"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/controllers"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClientController(t *testing.T) {
	t.Run("TestUpdateClient", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClientService := mocks.NewMockRepository(ctrl)
		mockClientService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		service := client.NewClientService(mockClientService)
		r := chi.NewRouter()
		controller := controllers.NewClientController(service)
		r.Put("/v1/client/{clientId}", controller.Update())

		clientId := "01F8Z5G4J6K7N3J4X2G4J6K7N3"
		requestBody := client.ClientRequestDto{
			Name:       "John",
			LastName:   "Doe",
			BirthDay:   "1990-01-01",
			EntityType: "PF",
			Document:   "932.222.900-40",
		}
		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/v1/client/"+clientId, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestDeleteClient", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClientService := mocks.NewMockRepository(ctrl)
		mockClientService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		service := client.NewClientService(mockClientService)
		r := chi.NewRouter()
		controller := controllers.NewClientController(service)
		r.Delete("/v1/client/{clientId}", controller.Delete())

		clientId := "01F8Z5G4J6K7N3J4X2G4J6K7N3"
		req := httptest.NewRequest(http.MethodDelete, "/v1/client/"+clientId, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("TestFindOneClient", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClientService := mocks.NewMockRepository(ctrl)
		mockClientService.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(&client.Client{
			Id:         ulid.Make(),
			Name:       "John",
			LastName:   "Doe",
			EntityType: client.Individual,
			Document:   "932.222.900-40",
		}, nil).Times(1)

		service := client.NewClientService(mockClientService)
		r := chi.NewRouter()
		controller := controllers.NewClientController(service)
		r.Get("/v1/client/{clientId}", controller.FindOne())

		clientId := "01F8Z5G4J6K7N3J4X2G4J6K7N3"
		req := httptest.NewRequest(http.MethodGet, "/v1/client/"+clientId, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestFindAllClients", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		birthDay, _ := time.Parse("2006-01-02", "1990-01-01")
		mockClientService := mocks.NewMockRepository(ctrl)
		mockClientService.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(&client.PaginationResult{
			TotalRecords: 1,
			Data: []*client.Client{{
				Id:         ulid.Make(),
				Name:       "John",
				LastName:   "Doe",
				EntityType: client.Individual,
				Document:   "932.222.900-40",
				BirthDay:   &birthDay,
			}},
		}, nil).Times(1)

		service := client.NewClientService(mockClientService)
		r := chi.NewRouter()
		controller := controllers.NewClientController(service)
		r.Get("/v1/client", controller.FindAll())

		req := httptest.NewRequest(http.MethodGet, "/v1/client?page=1&limit=10", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("TestCreateClient", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClientService := mocks.NewMockRepository(ctrl)
		mockClientService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockClientService.EXPECT().FindByDocument(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

		service := client.NewClientService(mockClientService)
		r := chi.NewRouter()
		controller := controllers.NewClientController(service)
		r.Post("/v1/client", controller.Create())
		requestBody := client.ClientRequestDto{
			Name:       "John",
			LastName:   "Doe",
			BirthDay:   "1990-01-01",
			EntityType: "PF",
			Document:   "932.222.900-40",
			Phones: []client.PhoneRequestDto{
				{
					Number:      "1234-5678",
					Description: "home",
				},
			},
			Addresses: []client.AddressRequestDto{
				{
					Street:  "123 Main St",
					City:    "Anytown",
					State:   "CA",
					ZipCode: "12345",
				},
			},
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/v1/client", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("deve retornar erro 422 (Unprocessable Entity) na validação dos dados de entrada", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClientService := mocks.NewMockRepository(ctrl)
		// Não espera chamada de Create nem FindByDocument

		service := client.NewClientService(mockClientService)
		r := chi.NewRouter()
		controller := controllers.NewClientController(service)
		r.Post("/v1/client", controller.Create())

		testCases := []struct {
			name       string
			request    client.ClientRequestDto
			wantFields []string
			wantStatus int
		}{
			{
				name: "quando nome não for fornecido",
				request: client.ClientRequestDto{
					LastName:   "Doe",
					BirthDay:   "1990-01-01",
					EntityType: "PF",
					Document:   "932.222.900-40",
				},
				wantFields: []string{"Name"},
				wantStatus: http.StatusUnprocessableEntity,
			},
			{
				name: "quando documento não for fornecido",
				request: client.ClientRequestDto{
					Name:       "John",
					LastName:   "Doe",
					BirthDay:   "1990-01-01",
					EntityType: "PF",
				},
				wantFields: []string{"Document"},
				wantStatus: http.StatusUnprocessableEntity,
			},
			{
				name: "quando tipo de entidade não for fornecido",
				request: client.ClientRequestDto{
					Name:     "John",
					LastName: "Doe",
					BirthDay: "1990-01-01",
					Document: "932.222.900-40",
				},
				wantFields: []string{"EntityType"},
				wantStatus: http.StatusUnprocessableEntity,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonBody, err := json.Marshal(tc.request)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/v1/client", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.wantStatus, w.Code)

				var response struct {
					Errors []struct {
						Field   string `json:"field"`
						Message string `json:"message"`
					} `json:"errors"`
				}

				err = json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Errors)

				foundFields := make([]string, 0)
				for _, err := range response.Errors {
					foundFields = append(foundFields, err.Field)
				}

				for _, wantField := range tc.wantFields {
					assert.Contains(t, foundFields, wantField, "Campo esperado não encontrado nos erros: %s", wantField)
				}
			})
		}
	})
}
