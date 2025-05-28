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

type CancelInfo struct {
	Reason      string
	CancelDate  *time.Time
	CancelledBy ulid.ULID
}

type ReversalInfo struct {
	Reason                  string
	ReversalDate            *time.Time
	ReversedBy              ulid.ULID
	ReversedInstallmentQtd  int
	CancelledInstallmentQtd int
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
	CancelInfo           *CancelInfo
	ReversalInfo         *ReversalInfo
	FinishedAt           *time.Time
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
		d.InstallmentsQuantity = 1
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

func (d *Debt) PayInstallment(payInfo *PaymentInfoDto) error {
	now := time.Now()

	if d.Status != Pending {
		return errors.New("debt is not in pending status")
	}

	installmentExists := false

	for i, installment := range d.Intallments {
		if installment.Id.String() != payInfo.InstallmentId {
			continue
		}

		installmentExists = true

		if installment.Status != Pending {
			return errors.New("installment is not in pending status")
		}

		if payInfo.Amount < installment.Value {
			return errors.New("amount does not match the installment value")
		}

		installment.Status = Paid
		installment.PaymentDate = &now
		installment.PaymentMethod = payInfo.PaymentMethod

		d.Intallments[i] = installment
		break
	}

	if !installmentExists {
		return errors.New("installment not found")
	}

	d.updateDebtStatus()

	return nil
}

func (d *Debt) Cancel(cancelInfo *CancelInfoDto) error {
	if d.Status != Pending {
		return errors.New("debt is not in pending status")
	}

	if d.hasInstallmentPaid() {
		return errors.New("cannot cancel debt with paid installments")
	}

	now := time.Now()

	d.Status = Canceled
	d.FinishedAt = &now
	d.CancelInfo = &CancelInfo{
		Reason:      cancelInfo.Reason,
		CancelDate:  &now,
		CancelledBy: cancelInfo.CancelledBy,
	}

	for i := range d.Intallments {
		d.Intallments[i].Status = Canceled
	}

	return nil
}

func (d *Debt) Reverse(reversalInfo *ReversalInfoDto) error {
	if d.Status == Canceled || d.Status == Reversed {
		return errors.New("debt is already canceled or reversed")
	}

	if d.ReversalInfo != nil {
		return errors.New("debt has already been reversed")
	}

	now := time.Now()

	d.Status = Reversed
	d.FinishedAt = &now
	qtdInstallmentsReversed := 0
	qtdInstallmentsCanceled := 0

	for i := range d.Intallments {
		if d.Intallments[i].Status == Paid {
			d.Intallments[i].Status = Reversed
			qtdInstallmentsReversed++
			continue
		}

		d.Intallments[i].Status = Canceled
		qtdInstallmentsCanceled++
	}

	d.ReversalInfo = &ReversalInfo{
		Reason:                  reversalInfo.Reason,
		ReversalDate:            &now,
		ReversedBy:              reversalInfo.ReversedBy,
		ReversedInstallmentQtd:  qtdInstallmentsReversed,
		CancelledInstallmentQtd: qtdInstallmentsCanceled,
	}

	return nil
}

func (d *Debt) updateDebtStatus() {
	allPaid := true

	for _, installment := range d.Intallments {
		if installment.Status != Paid {
			allPaid = false
			break
		}
	}

	if !allPaid {
		return
	}

	now := time.Now()
	d.Status = Paid
	d.FinishedAt = &now
}

func (d *Debt) hasInstallmentPaid() bool {
	for _, installment := range d.Intallments {
		if installment.Status == Paid {
			return true
		}
	}
	return false
}
