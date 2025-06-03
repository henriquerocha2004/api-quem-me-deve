package debt_test

import (
	"testing"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestDebtTests(t *testing.T) {
	t.Run("Should create a new Debt", func(t *testing.T) {
		// Arrange
		d := &debt.Debt{
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              nil,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
		}

		assert.NotNil(t, d)
		assert.Equal(t, "Test Debt", d.Description)
		assert.Equal(t, float64(1000), d.TotalValue)
	})

	t.Run("Should validate Debt with totalValue <= 0", func(t *testing.T) {
		// Arrange
		d := &debt.Debt{
			Description:          "Test Debt",
			TotalValue:           -1000,
			DueDate:              nil,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
		}

		// Act
		err := d.ValidateTotalValue()

		assert.NotNil(t, err)
		assert.Equal(t, "totalValue must be greater than 0", err.Error())

	})

	t.Run("Should validate debt when not provided due date", func(t *testing.T) {
		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              nil,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
		}

		// Act
		err := d.ValidateDueDate()

		assert.NotNil(t, err)
		assert.Equal(t, "dueDate is required", err.Error())
	})

	t.Run("Should validate debt when duedate provided is after than today", func(t *testing.T) {

		now := time.Now()
		dueDate := now.Add(-24 * time.Hour)

		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
		}

		// Act
		err := d.ValidateDueDate()

		assert.NotNil(t, err)
		assert.Equal(t, "dueDate must be in the future", err.Error())
	})

	t.Run("Should validate debt when not provided productId or serviceId", func(t *testing.T) {
		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              nil,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{},
			ServiceIds:           []ulid.ULID{},
		}

		// Act
		err := d.ValidateServiceOrProduct()

		assert.NotNil(t, err)
		assert.Equal(t, "at least one productId or serviceId is required", err.Error())
	})

	t.Run("Should validate debt when not provided userClientId", func(t *testing.T) {
		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              nil,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.ULID{},
			ProductIds:           []ulid.ULID{},
			ServiceIds:           []ulid.ULID{},
		}

		// Act
		err := d.ValidateClientId()

		assert.NotNil(t, err)
		assert.Equal(t, "userClientId is required", err.Error())
	})

	t.Run("Should validate all", func(t *testing.T) {

		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
		}

		// Act
		errs := d.Validate()

		assert.Len(t, errs.Errors, 0)

	})

	t.Run("Should create debt with installments", func(t *testing.T) {

		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 12,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		assert.NotNil(t, d)
		assert.Len(t, d.Intallments, 12)

		var value float64

		for _, installment := range d.Intallments {
			assert.NotNil(t, installment)
			assert.Equal(t, debt.Pending, installment.Status)
			value += installment.Value
		}

		assert.Equal(t, d.TotalValue, value)
	})

	t.Run("Deve realizar o pagamento de uma parcela", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		installmentId := d.Intallments[0].Id.String()

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: installmentId,
			Amount:        500,
			PaymentMethod: "credit_card",
		}

		err := d.PayInstallment(paymentInfo)
		assert.Nil(t, err)
		assert.Equal(t, debt.Paid, d.Intallments[0].Status)
		assert.NotNil(t, d.Intallments[0].PaymentDate)
		assert.Equal(t, paymentInfo.PaymentMethod, d.Intallments[0].PaymentMethod)
	})

	t.Run("Deve retornar erro ao tentar pagar uma parcela de uma divida que nao esta pendente", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Paid,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		installmentId := d.Intallments[0].Id.String()

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: installmentId,
			Amount:        500,
			PaymentMethod: "credit_card",
		}

		err := d.PayInstallment(paymentInfo)
		assert.NotNil(t, err)
		assert.Equal(t, "debt is not in pending status", err.Error())
	})

	t.Run("Deve retornar erro ao tentar pagar uma parcela que nao existe", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		// Arrange
		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: ulid.Make().String(),
			Amount:        500,
			PaymentMethod: "credit_card",
		}

		err := d.PayInstallment(paymentInfo)
		assert.NotNil(t, err)
		assert.Equal(t, "installment not found", err.Error())
	})

	t.Run("Deve retornar erro ao tentar pagar uma parcela que nao esta pendente", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		d.Intallments[0].Status = debt.Paid
		d.Intallments[0].PaymentDate = &now
		d.Intallments[0].PaymentMethod = "credit_card"

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: d.Intallments[0].Id.String(),
			Amount:        500,
			PaymentMethod: "credit_card",
		}

		err := d.PayInstallment(paymentInfo)
		assert.NotNil(t, err)
		assert.Equal(t, "installment is not in pending status", err.Error())
	})

	t.Run("Deve retornar erro ao tentar pagar uma parcela com valor diferente do valor da parcela", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		paymentInfo := &debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: d.Intallments[0].Id.String(),
			Amount:        400,
			PaymentMethod: "credit_card",
		}

		err := d.PayInstallment(paymentInfo)
		assert.NotNil(t, err)
		assert.Equal(t, "amount does not match the installment value", err.Error())
	})

	t.Run("Deve finalizar a divida quando todas as parcelas forem pagas", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		for i := range d.Intallments {
			paymentInfo := &debt.PaymentInfoDto{
				DebtId:        d.Id.String(),
				InstallmentId: d.Intallments[i].Id.String(),
				Amount:        d.Intallments[i].Value,
				PaymentMethod: "credit_card",
			}

			err := d.PayInstallment(paymentInfo)
			assert.Nil(t, err)
		}

		assert.Equal(t, debt.Paid, d.Status)
	})

	t.Run("Deve cancelar a divida", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		cancelInfo := &debt.CancelInfoDto{
			DebtId:      d.Id.String(),
			Reason:      "User requested cancellation",
			CancelledBy: ulid.Make(),
		}

		err := d.Cancel(cancelInfo)
		assert.Nil(t, err)
		assert.Equal(t, debt.Canceled, d.Status)
		assert.NotNil(t, d.CancelInfo)
		assert.Equal(t, d.CancelInfo.Reason, cancelInfo.Reason)
		assert.Equal(t, d.CancelInfo.CancelledBy.String(), cancelInfo.CancelledBy.String())
	})

	t.Run("Deve retornar erro ao tentar cancelar uma divida que nao esta pendente", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Paid,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		cancelInfo := &debt.CancelInfoDto{
			DebtId:      d.Id.String(),
			Reason:      "User requested cancellation",
			CancelledBy: ulid.Make(),
		}

		err := d.Cancel(cancelInfo)
		assert.NotNil(t, err)
		assert.Equal(t, "debt is not in pending status", err.Error())
	})

	t.Run("Deve retornar erro ao tentar cancelar uma divida com parcelas pagas", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		d.Intallments[0].Status = debt.Paid
		d.Intallments[0].PaymentDate = &now
		d.Intallments[0].PaymentMethod = "credit_card"

		cancelInfo := &debt.CancelInfoDto{
			DebtId:      d.Id.String(),
			Reason:      "User requested cancellation",
			CancelledBy: ulid.Make(),
		}

		err := d.Cancel(cancelInfo)
		assert.NotNil(t, err)
		assert.Equal(t, "cannot cancel debt with paid installments", err.Error())
	})

	t.Run("deve realizar um estorno com sucesso", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		reverseInfo := debt.ReversalInfoDto{
			DebtId:     d.Id.String(),
			Reason:     "Cliente desistiu da compra",
			ReversedBy: ulid.Make(),
		}

		err := d.Reverse(&reverseInfo)
		assert.Nil(t, err)
		assert.Equal(t, d.Status.String(), debt.Reversed.String())
		assert.NotNil(t, d.ReversalInfo)
		assert.Equal(t, d.ReversalInfo.Reason, reverseInfo.Reason)
		assert.Equal(t, d.ReversalInfo.ReversedBy.String(), reverseInfo.ReversedBy.String())

		for _, installment := range d.Intallments {
			assert.Equal(t, debt.Canceled.String(), installment.Status.String())
		}
	})

	t.Run("deve estornar uma divida com parcelas pagas", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Pending,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		d.GenerateInstallments()

		paymentInfo := debt.PaymentInfoDto{
			DebtId:        d.Id.String(),
			InstallmentId: d.Intallments[0].Id.String(),
			Amount:        500,
			PaymentMethod: "Cartão de crédito",
		}

		d.PayInstallment(&paymentInfo)

		reverseInfo := debt.ReversalInfoDto{
			DebtId:     d.Id.String(),
			Reason:     "Produto com defeito",
			ReversedBy: ulid.Make(),
		}

		err := d.Reverse(&reverseInfo)
		assert.Nil(t, err)
		assert.Equal(t, debt.Reversed.String(), d.Status.String())
		assert.Equal(t, debt.Reversed.String(), d.Intallments[0].Status.String())
		assert.NotNil(t, d.ReversalInfo)
	})

	t.Run("deve retornar um erro caso a divida esteja cancelada ou já estornada ao tentar realizar o estorno", func(t *testing.T) {
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)

		d := &debt.Debt{
			Id:                   ulid.Make(),
			Description:          "Test Debt",
			TotalValue:           1000,
			DueDate:              &dueDate,
			InstallmentsQuantity: 2,
			DebtDate:             nil,
			Status:               debt.Canceled,
			UserClientId:         ulid.Make(),
			ProductIds:           []ulid.ULID{ulid.Make()},
			ServiceIds:           []ulid.ULID{},
			Intallments:          []debt.Installment{},
		}

		reverseInfo := debt.ReversalInfoDto{
			DebtId:     d.Id.String(),
			Reason:     "Cliente desistiu da compra",
			ReversedBy: ulid.Make(),
		}

		err := d.Reverse(&reverseInfo)
		assert.Error(t, err, "debt is already canceled or reversed")
	})
}
