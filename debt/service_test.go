package debt_test

import (
	context "context"
	"testing"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/debt/mocks"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/validateErrors"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestDebtServiceTests(t *testing.T) {
	t.Run("Should create a new Debt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dueDate := time.Now().AddDate(0, 0, 1)

		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		debtRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		service := debt.NewDebtService(debtRepo, cliRepo)
		d := &debt.DebtDto{
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              dueDate.Format(time.DateOnly),
			InstallmentsQuantity: 12,
			UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
			ServiceIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
		}

		ctx := context.Background()

		response := service.CreateDebt(ctx, d)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debt created successfully", response.Message)
	})

	t.Run("Should return error if provided invalid total value", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dueDate := time.Now().AddDate(0, 0, 1)
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)
		d := &debt.DebtDto{
			Description:          "Test Debt",
			TotalValue:           -1000,
			DueDate:              dueDate.Format(time.DateOnly),
			InstallmentsQuantity: 12,
			UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
			ServiceIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
		}

		ctx := context.Background()

		response := service.CreateDebt(ctx, d)
		errorMessage := response.Data.(validateErrors.ValidationErrors)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "totalValue must be greater than 0", errorMessage.Errors[0].Message)

	})

	t.Run("Should return error if provided invalid due date", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dueDate := time.Now().AddDate(0, 0, -1)
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)
		d := &debt.DebtDto{
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              dueDate.Format(time.DateOnly),
			InstallmentsQuantity: 12,
			UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
			ServiceIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
		}

		ctx := context.Background()

		response := service.CreateDebt(ctx, d)
		errorMessage := response.Data.(validateErrors.ValidationErrors)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "dueDate must be in the future", errorMessage.Errors[0].Message)
	})

	t.Run("Should return error if provided invalid user client id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dueDate := time.Now().AddDate(0, 0, 1)
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)
		d := &debt.DebtDto{
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              dueDate.Format(time.DateOnly),
			InstallmentsQuantity: 12,
			UserClientId:         "",
			ProductIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
			ServiceIds:           []string{"01F8Z5G4J6K7N3J4X2G4J6K7N3"},
		}

		ctx := context.Background()

		response := service.CreateDebt(ctx, d)
		errorMessage := response.Data.(validateErrors.ValidationErrors)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "userClientId is required", errorMessage.Errors[0].Message)
	})

	t.Run("Should return error if provided invalid service ids", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dueDate := time.Now().AddDate(0, 0, 1)
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)
		d := &debt.DebtDto{
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              dueDate.Format(time.DateOnly),
			InstallmentsQuantity: 12,
			UserClientId:         "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			ProductIds:           []string{},
			ServiceIds:           []string{},
		}

		ctx := context.Background()

		response := service.CreateDebt(ctx, d)
		errorMessage := response.Data.(validateErrors.ValidationErrors)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "at least one productId or serviceId is required", errorMessage.Errors[0].Message)
	})

	t.Run("Deve retornar a lista de debitos de um cliente", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		ctx := context.Background()
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

		debtRepo.EXPECT().ClientUserDebts(ctx, clientId).Return(debts, nil)
		response := service.GetUserDebts(ctx, clientId)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debts retrieved successfully", response.Message)
		assert.Len(t, response.Data, 1)
		assert.Equal(t, "Test Debt", response.Data.([]debt.DebtDto)[0].Description)
		assert.Equal(t, 1000.0, response.Data.([]debt.DebtDto)[0].TotalValue)
		assert.Equal(t, "2023-10-01", response.Data.([]debt.DebtDto)[0].DueDate)
		assert.Equal(t, 2, response.Data.([]debt.DebtDto)[0].InstallmentsQuantity)
		assert.Equal(t, debt.Pending.String(), response.Data.([]debt.DebtDto)[0].Status)
		assert.Equal(t, "01F8Z5G4J6K7N3J4X2G4J6K7N3", response.Data.([]debt.DebtDto)[0].UserClientId)
	})

	t.Run("Deve retornar uma lista vazia de debitos de um cliente", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		ctx := context.Background()

		debtRepo.EXPECT().ClientUserDebts(ctx, clientId).Return([]*debt.Debt{}, nil)
		response := service.GetUserDebts(ctx, clientId)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "no debts found", response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("Deve retornar a lista de parcelas de um débito", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)
		duedateFirstInstallment, _ := time.Parse(time.DateOnly, "2025-05-25")
		duedateSecondInstallment, _ := time.Parse(time.DateOnly, "2025-06-25")
		ctx := context.Background()

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

		debtRepo.EXPECT().DebtInstallments(gomock.Any(), debtId).Return(installments, nil)
		cliRepo.EXPECT().ClientExists(gomock.Any(), clientId).Return(true, nil)
		response := service.GetDebtInstallments(ctx, clientId, debtId)
		assert.Len(t, response.Data, 2)
	})

	t.Run("Deve retornar um erro caso o cliente informado não exista", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		debtId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		clientId, _ := ulid.Parse("01F8Z5G4J6K7N3J4X2G4J6K7N3")
		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)
		cliRepo.EXPECT().ClientExists(gomock.Any(), clientId).Return(false, nil)
		debtRepo.EXPECT().DebtInstallments(gomock.Any(), debtId).Times(0)
		response := service.GetDebtInstallments(ctx, clientId, debtId)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "client not found", response.Message)
	})

	t.Run("Deve retornar uma lista de dividas paginada", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		pgRequest := paginate.PaginateRequest{
			Page:  1,
			Limit: 10,
		}

		pagData := paginate.SearchDto{
			Limit: pgRequest.Limit,
		}

		pagData.SetPage(pgRequest.Page)

		dueDate := time.Now().AddDate(0, 0, 1)

		debts := []*debt.Debt{
			{
				Description:          "Test Debt",
				Id:                   ulid.Make(),
				TotalValue:           1000,
				DueDate:              &dueDate,
				Status:               debt.Pending,
				UserClientId:         ulid.Make(),
				InstallmentsQuantity: 2,
				ServiceIds:           []ulid.ULID{ulid.Make()},
				ProductIds:           []ulid.ULID{ulid.Make()},
			},
		}

		paginatedResult := &debt.PaginationResult{
			Data:         debts,
			TotalRecords: 1,
		}

		debtRepo.EXPECT().GetDebts(gomock.Any(), pagData).Return(paginatedResult, nil)
		response := service.Debts(ctx, pgRequest)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debts retrieved successfully", response.Message)
		result := response.Data.(paginate.Result)

		assert.Equal(t, 1, result.TotalRecords)
		data := result.Data.([]debt.DebtDto)
		assert.Len(t, data, 1)
	})

	t.Run("Deve realizar o pagamento de uma parcela", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			UserClientId:         ulid.Make(),
			Status:               debt.Pending,
			DueDate:              nil,
			ProductIds:           []ulid.ULID{},
			ServiceIds:           []ulid.ULID{},
			InstallmentsQuantity: 1,
		}

		d.GenerateInstallments()

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: d.Intallments[0].Id.String(),
			Amount:        500,
			PaymentMethod: "Credit Card",
		}

		debtId, _ := ulid.Parse(paymentInfo.DebtId)
		debtRepo.EXPECT().GetDebt(gomock.Any(), debtId).Return(d, nil)

		debtRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		response := service.PayInstallment(ctx, paymentInfo)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "installment paid successfully", response.Message)
	})

	t.Run("Deve retornar um erro caso seja informado um debtoId inválido", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        "invalid-debt-id",
			InstallmentId: "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			Amount:        500,
			PaymentMethod: "Credit Card",
		}

		response := service.PayInstallment(ctx, paymentInfo)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "invalid debt ID", response.Message)
	})

	t.Run("Deve retornar um erro caso o debito não seja encontrado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			InstallmentId: "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			Amount:        500,
			PaymentMethod: "Credit Card",
		}

		debtId, _ := ulid.Parse(paymentInfo.DebtId)
		debtRepo.EXPECT().GetDebt(gomock.Any(), debtId).Return(nil, nil)

		response := service.PayInstallment(ctx, paymentInfo)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "debt not found", response.Message)
	})

	t.Run("Deve realizar o cancelamento de uma divida", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			UserClientId:         ulid.Make(),
			Status:               debt.Pending,
			DueDate:              nil,
			ProductIds:           []ulid.ULID{},
			ServiceIds:           []ulid.ULID{},
			InstallmentsQuantity: 1,
		}

		debtId, _ := ulid.Parse(d.Id.String())
		debtRepo.EXPECT().GetDebt(gomock.Any(), debtId).Return(d, nil)

		debtRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		cancelInfo := &debt.CancelInfoDto{
			DebtId:      d.Id.String(),
			Reason:      "Client requested cancellation",
			CancelledBy: ulid.Make(),
		}

		response := service.CancelDebt(ctx, cancelInfo)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debt cancelled successfully", response.Message)
		assert.Equal(t, debt.Canceled, d.Status)
	})

	t.Run("Deve retornar um erro caso seja informado um debtoId inválido para cancelamento", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		cancelInfo := &debt.CancelInfoDto{
			DebtId:      "invalid-debt-id",
			Reason:      "Client requested cancellation",
			CancelledBy: ulid.Make(),
		}

		response := service.CancelDebt(ctx, cancelInfo)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "invalid debt ID", response.Message)
	})

	t.Run("Deve retornar um erro caso o debito não seja encontrado para cancelamento", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		cancelInfo := &debt.CancelInfoDto{
			DebtId:      "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			Reason:      "Client requested cancellation",
			CancelledBy: ulid.Make(),
		}

		debtId, _ := ulid.Parse(cancelInfo.DebtId)
		debtRepo.EXPECT().GetDebt(gomock.Any(), debtId).Return(nil, nil)

		response := service.CancelDebt(ctx, cancelInfo)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "debt not found", response.Message)
	})

	t.Run("Deve realizar um estorno de uma divida", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			UserClientId:         ulid.Make(),
			Status:               debt.Pending,
			DueDate:              nil,
			ProductIds:           []ulid.ULID{},
			ServiceIds:           []ulid.ULID{},
			InstallmentsQuantity: 1,
		}

		debtId, _ := ulid.Parse(d.Id.String())
		debtRepo.EXPECT().GetDebt(gomock.Any(), debtId).Return(d, nil)

		debtRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		cancelInfo := &debt.ReversalInfoDto{
			DebtId:     d.Id.String(),
			Reason:     "Client requested Reversal",
			ReversedBy: ulid.Make(),
		}

		response := service.ReverseDebt(ctx, cancelInfo)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "debt reversed successfully", response.Message)
		assert.Equal(t, debt.Reversed, d.Status)
	})

	t.Run("Deve retornar um erro caso seja informado um debtoId inválido para estorno", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		cancelInfo := &debt.ReversalInfoDto{
			DebtId:     "invalid-debt-id",
			Reason:     "Client requested Reversal",
			ReversedBy: ulid.Make(),
		}

		response := service.ReverseDebt(ctx, cancelInfo)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "invalid debt ID", response.Message)
	})

	t.Run("Deve retornar um erro caso o debito não seja encontrado para estorno", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		debtRepo := mocks.NewMockRepository(ctrl)
		cliRepo := mocks.NewMockClientReader(ctrl)
		service := debt.NewDebtService(debtRepo, cliRepo)

		cancelInfo := &debt.ReversalInfoDto{
			DebtId:     "01F8Z5G4J6K7N3J4X2G4J6K7N3",
			Reason:     "Client requested Reversal",
			ReversedBy: ulid.Make(),
		}

		debtId, _ := ulid.Parse(cancelInfo.DebtId)
		debtRepo.EXPECT().GetDebt(gomock.Any(), debtId).Return(nil, nil)

		response := service.ReverseDebt(ctx, cancelInfo)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "debt not found", response.Message)
	})
}
