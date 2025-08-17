package client

import (
	"context"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/core/shared"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/document"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	Create(ctx context.Context, dto *ClientRequestDto) shared.ServiceResponse
	Update(ctx context.Context, id ulid.ULID, dto *ClientRequestDto) shared.ServiceResponse
	Delete(ctx context.Context, id ulid.ULID) shared.ServiceResponse
	FindById(ctx context.Context, id ulid.ULID) shared.ServiceResponse
	FindByCriteria(ctx context.Context, criteria *paginate.PaginateRequest) shared.ServiceResponse
}

type ClientService struct {
	repository Repository
}

func NewClientService(repository Repository) *ClientService {
	return &ClientService{
		repository: repository,
	}
}

func (s *ClientService) Create(ctx context.Context, dto *ClientRequestDto) shared.ServiceResponse {

	c, err := s.repository.FindByDocument(ctx, dto.Document)

	if err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "error in create client",
		}
	}

	if c != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "client with this document already exists",
		}
	}

	birth, _ := time.Parse(time.DateOnly, dto.BirthDay)

	client := &Client{
		Id:         ulid.Make(),
		Name:       dto.Name,
		LastName:   dto.LastName,
		EntityType: EntityType(dto.EntityType),
		Document:   document.Document(dto.Document),
		BirthDay:   &birth,
	}

	err = client.validate()
	if err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: err.Error(),
		}
	}

	if len(dto.Addresses) > 0 {
		for _, address := range dto.Addresses {
			client.addAddress(address.Street, address.Neighborhood, address.City, address.State, address.ZipCode)
		}
	}

	if len(dto.Phones) > 0 {
		for _, phone := range dto.Phones {
			client.addPhone(phone.Description, phone.Number)
		}
	}

	if err = s.repository.Create(ctx, client); err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "error in create client",
		}
	}

	return shared.ServiceResponse{
		Status:  "success",
		Message: "client created successfully",
	}
}

func (s *ClientService) Update(ctx context.Context, id ulid.ULID, dto *ClientRequestDto) shared.ServiceResponse {

	birth, _ := time.Parse(time.DateOnly, dto.BirthDay)

	client := &Client{
		Id:         id,
		Name:       dto.Name,
		LastName:   dto.LastName,
		EntityType: EntityType(dto.EntityType),
		Document:   document.Document(dto.Document),
		BirthDay:   &birth,
	}

	err := client.validate()
	if err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: err.Error(),
		}
	}

	if len(dto.Addresses) > 0 {
		for _, address := range dto.Addresses {
			client.addAddress(address.Street, address.Neighborhood, address.City, address.State, address.ZipCode)
		}
	}

	if len(dto.Phones) > 0 {
		for _, phone := range dto.Phones {
			client.addPhone(phone.Description, phone.Number)
		}
	}

	if err = s.repository.Update(ctx, client); err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "error in update client",
		}
	}

	return shared.ServiceResponse{
		Status:  "success",
		Message: "client updated successfully",
	}
}

func (s *ClientService) Delete(ctx context.Context, id ulid.ULID) shared.ServiceResponse {
	if err := s.repository.Delete(ctx, id); err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "error in delete client",
		}
	}

	return shared.ServiceResponse{
		Status:  "success",
		Message: "client deleted successfully",
	}
}

func (s *ClientService) FindById(ctx context.Context, id ulid.ULID) shared.ServiceResponse {
	client, err := s.repository.FindById(ctx, id)
	if err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "error in find client",
		}
	}

	if client == nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "client not found",
		}
	}

	return shared.ServiceResponse{
		Status:  "success",
		Message: "client found successfully",
		Data:    client,
	}
}

func (s *ClientService) FindByCriteria(ctx context.Context, criteria *paginate.PaginateRequest) shared.ServiceResponse {
	pagDto := paginate.SearchDto{
		Limit:         criteria.Limit,
		SortField:     criteria.SortField,
		TermSearch:    criteria.SearchTerm,
		SortDirection: criteria.SortDirection,
	}

	pagDto.SetPage(criteria.Page)
	pagDto.AddColumnSearch(criteria.ColumnSearch)

	result, err := s.repository.FindAll(ctx, pagDto)

	if err != nil {
		return shared.ServiceResponse{
			Status:  "error",
			Message: "error in find clients",
		}
	}

	clientsDto := s.convertToClientDto(result.Data)

	return shared.ServiceResponse{
		Status:  "success",
		Message: "clients found successfully",
		Data: paginate.Result{
			TotalRecords: result.TotalRecords,
			Data:         clientsDto,
		},
	}
}

func (s *ClientService) convertToClientDto(clients []*Client) []ClientRequestDto {
	var clientsDto []ClientRequestDto

	for _, c := range clients {

		cliDto := ClientRequestDto{
			Name:       c.Name,
			LastName:   c.LastName,
			BirthDay:   c.BirthDay.Format(time.DateOnly),
			EntityType: string(c.EntityType),
			Document:   string(c.Document),
		}

		if len(c.Addresses) >= 1 {
			var addressesDto []AddressRequestDto

			for _, address := range c.Addresses {
				addressesDto = append(addressesDto, AddressRequestDto{
					Street:  address.Street,
					City:    address.City,
					State:   address.State,
					ZipCode: address.ZipCode,
				})
			}

			cliDto.Addresses = addressesDto
		}

		if len(c.Phones) >= 1 {
			var phonesDto []PhoneRequestDto

			for _, phone := range c.Phones {
				phonesDto = append(phonesDto, PhoneRequestDto{
					Description: phone.Description,
					Number:      phone.Number,
				})
			}

			cliDto.Phones = phonesDto
		}

		clientsDto = append(clientsDto, cliDto)
	}

	return clientsDto
}
