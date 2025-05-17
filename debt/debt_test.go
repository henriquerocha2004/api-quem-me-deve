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
}
