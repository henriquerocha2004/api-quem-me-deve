package debt

import "github.com/oklog/ulid/v2"

type DebtDto struct {
	Id                   string           `json:"id,omitempty"`
	Description          string           `json:"description" validate:"required"`
	TotalValue           float64          `json:"total_value" validate:"required,gt=0"`
	DueDate              string           `json:"due_date" validate:"required,dateFormat:YYYY-MM-DD"`
	InstallmentsQuantity int              `json:"installments_quantity" validate:"required,gt=0"`
	UserClientId         string           `json:"user_client_id" validate:"required"`
	ProductIds           []string         `json:"product_ids"`
	ServiceIds           []string         `json:"service_ids"`
	Status               string           `json:"status,omitempty"`
	Intallments          []InstallmentDto `json:"intallments,omitempty"`
	DebtDate             string           `json:"debt_date,omitempty"`
}

type InstallmentDto struct {
	Id            string  `json:"id,omitempty"`
	Description   string  `json:"description"`
	Value         float64 `json:"value"`
	DueDate       string  `json:"due_date"`
	DebDate       string  `json:"debt_date"`
	Status        string  `json:"status"`
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method"`
	Number        int     `json:"number"`
}

type PaginationResult struct {
	TotalRecords int     `json:"total_records"`
	Data         []*Debt `json:"data"`
}

type PaymentInfoDto struct {
	DebtId        string  `json:"debt_id" validate:"required,ulid"`
	InstallmentId string  `json:"installment_id" validate:"required,ulid"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
}

type CancelInfoDto struct {
	DebtId      string `json:"debt_id" validate:"required,ulid"`
	Reason      string `json:"reason" validate:"required"`
	CancelledBy ulid.ULID
}

type ReversalInfoDto struct {
	DebtId     string `json:"debt_id" validate:"required,ulid"`
	Reason     string `json:"reason" validate:"required"`
	ReversedBy ulid.ULID
}
