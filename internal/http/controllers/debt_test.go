package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	clientMemory "github.com/henriquerocha2004/quem-me-deve-api/client/memory"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	debtMemory "github.com/henriquerocha2004/quem-me-deve-api/debt/memory"
	"github.com/henriquerocha2004/quem-me-deve-api/debt/mocks"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/controllers"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDebtController(t *testing.T) {
	t.Run("Deve criar uma divida com sucesso", func(t *testing.T) {
		service := debt.NewDebtService(debtMemory.NewMemoryRepository(), clientMemory.NewClientDebtReader())
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt", controller.CreateDebt())
		dueDate := time.Now().AddDate(0, 0, 1)
		requestBody := debt.DebtDto{
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              dueDate.Format(time.DateOnly),
			InstallmentsQuantity: 12,
			UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
			ServiceIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("deve retornar erro 422 (Unprocessable Entity). Validação dos dados de entrada", func(t *testing.T) {
		service := debt.NewDebtService(debtMemory.NewMemoryRepository(), clientMemory.NewClientDebtReader())
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt", controller.CreateDebt())

		testCases := []struct {
			name       string
			request    debt.DebtDto
			wantFields []string
			wantStatus int
		}{
			{
				name: "quando valor total não for fornecido",
				request: debt.DebtDto{
					Description:          "Test Debt",
					DueDate:              time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
					InstallmentsQuantity: 12,
					UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
					ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
				},
				wantFields: []string{"TotalValue"},
				wantStatus: http.StatusUnprocessableEntity,
			},
			{
				name: "quando data de vencimento não for fornecida",
				request: debt.DebtDto{
					Description:          "Test Debt",
					TotalValue:           1000,
					InstallmentsQuantity: 12,
					UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
					ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
				},
				wantFields: []string{"DueDate"},
				wantStatus: http.StatusUnprocessableEntity,
			},
			{
				name: "quando ID do cliente não for fornecido",
				request: debt.DebtDto{
					Description:          "Test Debt",
					TotalValue:           1000,
					DueDate:              time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
					InstallmentsQuantity: 12,
					ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
				},
				wantFields: []string{"UserClientId"},
				wantStatus: http.StatusUnprocessableEntity,
			},
			{
				name: "quando não for fornecido o campo quantidade de parcelas",
				request: debt.DebtDto{
					Description:  "Test Debt",
					TotalValue:   1000,
					DueDate:      time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
					UserClientId: "01F8Z5G4J6K7N3J4X2G4J6K7N3",
					ProductIds:   []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
					ServiceIds:   []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
				},
				wantFields: []string{"InstallmentsQuantity"},
				wantStatus: http.StatusUnprocessableEntity,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonBody, err := json.Marshal(tc.request)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/v1/debt", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				// Verifica o status code
				assert.Equal(t, tc.wantStatus, w.Code)

				// Verifica a estrutura da resposta
				var response struct {
					Errors []struct {
						Field   string `json:"field"`
						Message string `json:"message"`
					} `json:"errors"`
				}

				err = json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)

				// Verifica se a resposta contém os campos esperados
				assert.NotEmpty(t, response.Errors)

				// Verifica se os campos com erro esperados estão presentes
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

	t.Run("deve retornar os debitos de um cliente pelo ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		duedate, _ := time.Parse(time.DateOnly, "2023-10-01")
		now := time.Now()
		debts := []*debt.Debt{
			{
				Description:          "Test Debt",
				Id:                   ulid.Make(),
				TotalValue:           1000,
				DueDate:              &duedate,
				Status:               debt.Pending,
				UserClientId:         clientId,
				InstallmentsQuantity: 2,
				ServiceIds:           []ulid.ULID{ulid.Make()},
				ProductIds:           []ulid.ULID{ulid.Make()},
				DebtDate:             &now,
			},
		}

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().ClientUserDebts(gomock.Any(), clientId).Return(debts, nil)
		cliRepository := mocks.NewMockClientReader(ctrl)

		service := debt.NewDebtService(debtRepository, cliRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Get("/v1/debt/{clientId}", controller.GetClientUserDebts())

		req := httptest.NewRequest(http.MethodGet, "/v1/debt/01F8Z5G4J6K7N3J4X2G4J6K7N3", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response debt.DebtResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Len(t, response.Data, 1)
	})

	t.Run("deve retornar um erro quando não é passado o clientId", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().ClientUserDebts(gomock.Any(), gomock.Any()).Times(0)
		clientRepository := mocks.NewMockClientReader(ctrl)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Get("/v1/debt/{clientId}", controller.GetClientUserDebts())

		testCases := []struct {
			Url                  string
			ErrorMessageExpected string
			ExpectedCodeRequest  int16
		}{
			{
				Url:                  "/v1/debt/invalidClientId",
				ErrorMessageExpected: "Invalid clientId",
				ExpectedCodeRequest:  http.StatusBadRequest,
			},
		}

		for _, test := range testCases {
			req := httptest.NewRequest(http.MethodGet, test.Url, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("deve retornar a lista de parcelas de uma divida", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		debtId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		duedateFirstInstallment, _ := time.Parse(time.DateOnly, "2025-05-25")
		duedateSecondInstallment, _ := time.Parse(time.DateOnly, "2025-06-25")

		installments := []*debt.Installment{
			{
				Id:            ulid.Make(),
				Description:   "Referente a compra de CD",
				Value:         600,
				DueDate:       &duedateFirstInstallment,
				DebDate:       &duedateFirstInstallment,
				Status:        debt.Pending,
				Number:        1,
				PaymentDate:   nil,
				PaymentMethod: "Parcelado",
			},
			{
				Id:            ulid.Make(),
				Description:   "Referente a compra de CD",
				Value:         600,
				DueDate:       &duedateSecondInstallment,
				DebDate:       &duedateSecondInstallment,
				Status:        debt.Pending,
				Number:        2,
				PaymentDate:   nil,
				PaymentMethod: "Parcelado",
			},
		}

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().DebtInstallments(gomock.Any(), debtId).Return(installments, nil)
		clientRepository := mocks.NewMockClientReader(ctrl)
		clientRepository.EXPECT().ClientExists(gomock.Any(), clientId).Return(true, nil)
		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Get("/v1/debt/{clientId}/{debtId}/installments", controller.GetDebtInstallments())

		req := httptest.NewRequest(http.MethodGet, "/v1/debt/01F8Z5G4J6K7N3J4X2G4J6K7N3/01F8Z5G4J6K7N3J4X2G4J6K7N3/installments", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response debt.DebtResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Len(t, response.Data, 2)
	})

	t.Run("Deve retornar um erro caso seja informado um client id inválido ao tentar obter as parcelas de uma divida", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepository := mocks.NewMockRepository(ctrl)
		clientRepository := mocks.NewMockClientReader(ctrl)

		debtRepository.EXPECT().DebtInstallments(gomock.Any(), gomock.Any()).Times(0)
		clientRepository.EXPECT().ClientExists(gomock.Any(), gomock.Any()).Times(0)
		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)
		r := chi.NewRouter()
		r.Get("/v1/debt/{clientId}/{debtId}/installments", controller.GetDebtInstallments())

		req := httptest.NewRequest(http.MethodGet, "/v1/debt/invalidClientId/01F8Z5G4J6K7N3J4X2G4J6K7N3/installments", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response string
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "clientId invalid", response)
	})

	t.Run("Deve retornar um erro caso seja informado um debt id inválido ao tentar obter as parcelas de uma divida", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepository := mocks.NewMockRepository(ctrl)
		clientRepository := mocks.NewMockClientReader(ctrl)

		debtRepository.EXPECT().DebtInstallments(gomock.Any(), gomock.Any()).Times(0)
		clientRepository.EXPECT().ClientExists(gomock.Any(), gomock.Any()).Times(0)
		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)
		r := chi.NewRouter()
		r.Get("/v1/debt/{clientId}/{debtId}/installments", controller.GetDebtInstallments())

		req := httptest.NewRequest(http.MethodGet, "/v1/debt/01F8Z5G4J6K7N3J4X2G4J6K7N3/debtIdinvalid/installments", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response string
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "debtId invalid", response)
	})
}
