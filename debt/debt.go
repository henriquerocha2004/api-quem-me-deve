package debt

import (
	"errors"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/pkg/validateErrors"
	"github.com/oklog/ulid/v2"
)

type Installment struct {
	Id            ulid.ULID
	Description   string
	Value         float64
	DueDate       *time.Time
	DebDate       *time.Time
	Status        status
	PaymentDate   *time.Time
	PaymentMethod string
	Number        int
}

type Debt struct {
	Id                   ulid.ULID
	Description          string
	TotalValue           float64
	DueDate              *time.Time
	InstallmentsQuantity int
	DebtDate             *time.Time
	Status               status
	UserClientId         ulid.ULID
	ProductIds           []ulid.ULID
	ServiceIds           []ulid.ULID
	Intallments          []Installment
}

func (d *Debt) Validate() validateErrors.ValidationErrors {
	var validationErrors validateErrors.ValidationErrors

	if err := d.ValidateTotalValue(); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, validateErrors.ValidationError{
			Field:   "totalValue",
			Message: err.Error(),
		})
	}

	if err := d.ValidateDueDate(); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, validateErrors.ValidationError{
			Field:   "dueDate",
			Message: err.Error(),
		})
	}

	if err := d.ValidateServiceOrProduct(); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, validateErrors.ValidationError{
			Field:   "serviceOrProduct",
			Message: err.Error(),
		})
	}
	if err := d.ValidateClientId(); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, validateErrors.ValidationError{
			Field:   "userClientId",
			Message: err.Error(),
		})
	}

	return validationErrors
}

func (d *Debt) ValidateTotalValue() error {
	if d.TotalValue > 0 {
		return nil
	}

	return errors.New("totalValue must be greater than 0")
}

func (d *Debt) ValidateDueDate() error {
	if d.DueDate == nil {
		return errors.New("dueDate is required")
	}

	if d.DueDate.Before(time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0, 0, 0, 0,
		time.Now().Location(),
	)) {
		return errors.New("dueDate must be in the future")
	}

	return nil
}

func (d *Debt) ValidateServiceOrProduct() error {
	if len(d.ProductIds) == 0 && len(d.ServiceIds) == 0 {
		return errors.New("at least one productId or serviceId is required")
	}

	return nil
}

func (d *Debt) ValidateClientId() error {
	if d.UserClientId == (ulid.ULID{}) {
		return errors.New("userClientId is required")
	}

	return nil
}

func (d *Debt) GenerateInstallments() error {
	if d.InstallmentsQuantity <= 0 {
		return nil
	}

	now := time.Now()
	currentDueDate := d.DueDate

	baseValue := d.TotalValue / float64(d.InstallmentsQuantity)
	totalAllocated := 0.0

	for i := range d.InstallmentsQuantity {
		if i > 0 {
			futureDate := currentDueDate.AddDate(0, 0, 30)
			currentDueDate = &futureDate
		}

		value := baseValue
		if i == d.InstallmentsQuantity-1 {
			value = d.TotalValue - totalAllocated
		}

		installment := Installment{
			Id:            ulid.Make(),
			Description:   d.Description,
			Value:         value,
			DueDate:       currentDueDate,
			DebDate:       &now,
			Status:        Pending,
			PaymentDate:   nil,
			PaymentMethod: "",
			Number:        i + 1,
		}

		d.Intallments = append(d.Intallments, installment)
		totalAllocated += value
	}

	return nil
}
