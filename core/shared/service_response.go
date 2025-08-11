package shared

type ServiceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
