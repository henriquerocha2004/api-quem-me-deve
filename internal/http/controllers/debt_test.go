package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	clientMemory "github.com/henriquerocha2004/quem-me-deve-api/client/memory"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/debt/mocks"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/controllers"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/customvalidate"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDebtController(t *testing.T) {
	t.Run("Deve criar uma divida com sucesso", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		debtRepo := mocks.NewMockRepository(ctrl)
		debtRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		service := debt.NewDebtService(debtRepo, clientMemory.NewClientDebtReader())
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

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		debtRepo := mocks.NewMockRepository(ctrl)
		debtRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(0)

		service := debt.NewDebtService(debtRepo, clientMemory.NewClientDebtReader())
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

	t.Run("Deve retornar os débitos paginados", func(t *testing.T) {
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

		pgResult := &debt.PaginationResult{
			TotalRecords: 1,
			Data:         debts,
		}

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().GetDebts(gomock.Any(), gomock.Any()).Return(pgResult, nil)
		clientRepository := mocks.NewMockClientReader(ctrl)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Get("/v1/debt", controller.GetDebts())

		query := url.Values{}
		query.Add("page", "1")
		query.Add("limit", "10")
		query.Add("sort_field", "created_at")
		query.Add("sort_direction", "desc")
		query.Add("column_search[0][name]", "description")
		query.Add("column_search[0][value]", "Test Debt")

		uri := "/v1/debt?" + query.Encode()

		req := httptest.NewRequest(http.MethodGet, uri, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response debt.DebtResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Len(t, response.Data, 2)
	})

	t.Run("Deve realizar o pagamento de uma parcela", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		debtId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		installmentId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        debtId.String(),
			InstallmentId: installmentId.String(),
			Amount:        500,
			PaymentMethod: "Credit Card",
		}

		d := &debt.Debt{
			Id:                   debtId,
			UserClientId:         clientId,
			InstallmentsQuantity: 2,
			Intallments: []debt.Installment{
				{
					Id:            installmentId,
					Description:   "Test Installment",
					Value:         500,
					DueDate:       nil,
					DebDate:       nil,
					Status:        debt.Pending,
					PaymentDate:   nil,
					PaymentMethod: "",
					Number:        1,
				},
				{
					Id:            ulid.Make(),
					Description:   "Test Installment 2",
					Value:         500,
					DueDate:       nil,
					DebDate:       nil,
					Status:        debt.Pending,
					PaymentDate:   nil,
					PaymentMethod: "",
					Number:        2,
				},
			},
			Status: debt.Pending,
		}

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().GetDebt(gomock.Any(), debtId).Return(d, nil)
		debtRepository.EXPECT().Update(gomock.Any(), d).Return(nil)
		clientRepository := mocks.NewMockClientReader(ctrl)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt/pay-installment", controller.PayInstallment())

		jsonBody, err := json.Marshal(paymentInfo)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt/pay-installment", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response debt.DebtResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "installment paid successfully", response.Message)
	})

	t.Run("Deve retornar um erro ao tentar pagar uma parcela com dados inválidos", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepository := mocks.NewMockRepository(ctrl)
		clientRepository := mocks.NewMockClientReader(ctrl)

		debtRepository.EXPECT().GetDebt(gomock.Any(), gomock.Any()).Times(0)
		debtRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
		clientRepository.EXPECT().ClientExists(gomock.Any(), gomock.Any()).Times(0)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt/pay-installment", controller.PayInstallment())

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        "invalidDebtId",
			InstallmentId: "invalidInstallmentId",
			Amount:        0,
			PaymentMethod: "",
		}

		jsonBody, err := json.Marshal(paymentInfo)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt/pay-installment", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response customvalidate.ValidationResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, 4, len(response.Errors))
		assert.Equal(t, "DebtId", response.Errors[0].Field)
		assert.Equal(t, "Invalid ULID format", response.Errors[0].Message)
		assert.Equal(t, "InstallmentId", response.Errors[1].Field)
		assert.Equal(t, "Invalid ULID format", response.Errors[1].Message)
		assert.Equal(t, "Amount", response.Errors[2].Field)
		assert.Equal(t, "This field is required", response.Errors[2].Message)
		assert.Equal(t, "PaymentMethod", response.Errors[3].Field)
		assert.Equal(t, "This field is required", response.Errors[3].Message)
	})

	t.Run("Deve cancelar uma dívida com sucesso", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		debtId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")

		cancelInfo := &debt.CancelInfoDto{
			DebtId: debtId.String(),
			Reason: "Client requested cancellation",
		}

		d := &debt.Debt{
			Id:                   debtId,
			UserClientId:         clientId,
			InstallmentsQuantity: 2,
			Status:               debt.Pending,
		}

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().GetDebt(gomock.Any(), debtId).Return(d, nil)
		debtRepository.EXPECT().Update(gomock.Any(), d).Return(nil)
		clientRepository := mocks.NewMockClientReader(ctrl)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt/cancel", controller.CancelDebt())

		jsonBody, err := json.Marshal(cancelInfo)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt/cancel", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response debt.DebtResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debt cancelled successfully", response.Message)
	})

	t.Run("Deve retornar um erro ao tentar cancelar uma dívida com dados inválidos", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepository := mocks.NewMockRepository(ctrl)
		clientRepository := mocks.NewMockClientReader(ctrl)

		debtRepository.EXPECT().GetDebt(gomock.Any(), gomock.Any()).Times(0)
		debtRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
		clientRepository.EXPECT().ClientExists(gomock.Any(), gomock.Any()).Times(0)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt/cancel", controller.CancelDebt())

		cancelInfo := &debt.CancelInfoDto{
			DebtId: "invalidDebtId",
			Reason: "",
		}

		jsonBody, err := json.Marshal(cancelInfo)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt/cancel", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response customvalidate.ValidationResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, 2, len(response.Errors))
		assert.Equal(t, "DebtId", response.Errors[0].Field)
		assert.Equal(t, "Invalid ULID format", response.Errors[0].Message)
		assert.Equal(t, "Reason", response.Errors[1].Field)
		assert.Equal(t, "This field is required", response.Errors[1].Message)
	})

	t.Run("Deve reverter uma dívida com sucesso", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		debtId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")

		reversalInfo := &debt.ReversalInfoDto{
			DebtId: debtId.String(),
			Reason: "Client requested reversal",
		}

		d := &debt.Debt{
			Id:                   debtId,
			UserClientId:         clientId,
			InstallmentsQuantity: 2,
			Status:               debt.Pending,
		}

		debtRepository := mocks.NewMockRepository(ctrl)
		debtRepository.EXPECT().GetDebt(gomock.Any(), debtId).Return(d, nil)
		debtRepository.EXPECT().Update(gomock.Any(), d).Return(nil)
		clientRepository := mocks.NewMockClientReader(ctrl)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt/reversal", controller.ReversalDebt())

		jsonBody, err := json.Marshal(reversalInfo)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt/reversal", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response debt.DebtResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debt reversed successfully", response.Message)
	})

	t.Run("Deve retornar um erro ao tentar reverter uma dívida com dados inválidos", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepository := mocks.NewMockRepository(ctrl)
		clientRepository := mocks.NewMockClientReader(ctrl)

		debtRepository.EXPECT().GetDebt(gomock.Any(), gomock.Any()).Times(0)
		debtRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
		clientRepository.EXPECT().ClientExists(gomock.Any(), gomock.Any()).Times(0)

		service := debt.NewDebtService(debtRepository, clientRepository)
		controller := controllers.NewDebtController(service)

		r := chi.NewRouter()
		r.Post("/v1/debt/reversal", controller.ReversalDebt())

		reversalInfo := &debt.ReversalInfoDto{
			DebtId: "invalidDebtId",
			Reason: "",
		}

		jsonBody, err := json.Marshal(reversalInfo)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/debt/reversal", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response customvalidate.ValidationResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, 2, len(response.Errors))
		assert.Equal(t, "DebtId", response.Errors[0].Field)
		assert.Equal(t, "Invalid ULID format", response.Errors[0].Message)
		assert.Equal(t, "Reason", response.Errors[1].Field)
		assert.Equal(t, "This field is required", response.Errors[1].Message)
	})
}
