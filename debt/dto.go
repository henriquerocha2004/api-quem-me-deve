package debt

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

type DebtResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
