package debt

type DebtDto struct {
	Id                   string   `json:"id,omitempty"`
	Description          string   `json:"description" validate:"required"`
	TotalValue           float64  `json:"total_value" validate:"required,gt=0"`
	DueDate              string   `json:"due_date" validate:"required,dateFormat:YYYY-MM-DD"`
	InstallmentsQuantity int      `json:"installments_quantity" validate:"required,gt=0"`
	UserClientId         string   `json:"user_client_id" validate:"required"`
	ProductIds           []string `json:"product_ids"`
	ServiceIds           []string `json:"service_ids"`
	Status               string   `json:"status,omitempty"`
	Intallments          []string `json:"intallments,omitempty"`
	DebtDate             string   `json:"debt_date,omitempty"`
}

type DebtResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
