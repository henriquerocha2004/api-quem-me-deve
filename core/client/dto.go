package client

type ClientRequestDto struct {
	Name       string              `json:"name" validate:"required"`
	LastName   string              `json:"last_name" validate:"required"`
	BirthDay   string              `json:"birthday" validate:"required,dateFormat:YYYY-MM-DD"`
	EntityType string              `json:"entity_type" validate:"required"`
	Document   string              `json:"document" validate:"required"`
	Phones     []PhoneRequestDto   `json:"phones,omitempty"`
	Addresses  []AddressRequestDto `json:"addresses,omitempty"`
}

type AddressRequestDto struct {
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
}

type PhoneRequestDto struct {
	Description string `json:"description"`
	Number      string `json:"number"`
}

type PaginationResult struct {
	TotalRecords int       `json:"total_records"`
	Data         []*Client `json:"data"`
}
