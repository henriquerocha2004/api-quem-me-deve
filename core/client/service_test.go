package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/core/client"
	"github.com/henriquerocha2004/quem-me-deve-api/core/client/mocks"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClientService(t *testing.T) {
	t.Run("should create client with success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cliRepo := mocks.NewMockRepository(ctrl)
		cliRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
		cliRepo.EXPECT().FindByDocument(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

		clientRequest := client.ClientRequestDto{
			Name:       "Nome",
			LastName:   "Sobrenome",
			BirthDay:   "2000-01-01",
			EntityType: "PF",
			Document:   "510.091.940-03",
			Phones: []client.PhoneRequestDto{
				{
					Description: "Pessoal",
					Number:      "99929992929",
				},
			},
			Addresses: []client.AddressRequestDto{
				{
					Street:       "Rua dos Bobos 0",
					Neighborhood: "Bairro dos Bobos",
					City:         "Cidade dos Bobos",
					State:        "BB",
					ZipCode:      "21433-011",
				},
			},
		}

		service := client.NewClientService(cliRepo)
		result := service.Create(context.Background(), &clientRequest)

		assert.Equal(t, result.Status, "success")
	})

	t.Run("should not create client if already exists client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cli := client.Client{
			Name:     "Atreus",
			LastName: "Da Guerra",
		}

		cliRepo := mocks.NewMockRepository(ctrl)
		cliRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0).Return(nil)
		cliRepo.EXPECT().FindByDocument(gomock.Any(), gomock.Any()).Times(1).Return(&cli, nil)

		clientRequest := client.ClientRequestDto{
			Name:       "Nome",
			LastName:   "Sobrenome",
			BirthDay:   "2000-01-01",
			EntityType: "PF",
			Document:   "510.091.940-03",
			Phones: []client.PhoneRequestDto{
				{
					Description: "Pessoal",
					Number:      "99929992929",
				},
			},
			Addresses: []client.AddressRequestDto{
				{
					Street:       "Rua dos Bobos 0",
					Neighborhood: "Bairro dos Bobos",
					City:         "Cidade dos Bobos",
					State:        "BB",
					ZipCode:      "21433-011",
				},
			},
		}

		service := client.NewClientService(cliRepo)
		result := service.Create(context.Background(), &clientRequest)

		assert.Equal(t, result.Status, "error")
		assert.Equal(t, result.Message, "client with this document already exists")
	})

	t.Run("Should udpate client data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cliRepo := mocks.NewMockRepository(ctrl)
		cliRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil)

		clientRequest := client.ClientRequestDto{
			Name:       "Nome",
			LastName:   "Sobrenome",
			BirthDay:   "2000-01-01",
			EntityType: "PF",
			Document:   "510.091.940-03",
			Phones: []client.PhoneRequestDto{
				{
					Description: "Pessoal",
					Number:      "99929992929",
				},
			},
			Addresses: []client.AddressRequestDto{
				{
					Street:       "Rua dos Bobos 0",
					Neighborhood: "Bairro dos Bobos",
					City:         "Cidade dos Bobos",
					State:        "BB",
					ZipCode:      "21433-011",
				},
			},
		}

		service := client.NewClientService(cliRepo)
		result := service.Update(context.Background(), ulid.Make(), &clientRequest)

		assert.Equal(t, result.Status, "success")
	})

	t.Run("should retrieve one client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cli := client.Client{
			Name:     "Atreus",
			LastName: "Da Guerra",
		}

		cliRepo := mocks.NewMockRepository(ctrl)
		cliRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Times(1).Return(&cli, nil)

		service := client.NewClientService(cliRepo)
		result := service.FindById(context.Background(), ulid.Make())

		data := result.Data.(*client.Client)

		assert.Equal(t, result.Status, "success")
		assert.Equal(t, data.Name, "Atreus")
		assert.Equal(t, data.LastName, "Da Guerra")
	})

	t.Run("should retrive clients by filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		now := time.Now()

		c := client.Client{
			Name:     "Atreus",
			LastName: "Da Guerra",
			BirthDay: &now,
			Document: "510.091.940-03",
		}

		clients := []*client.Client{
			&c,
		}

		resultSearch := client.PaginationResult{
			TotalRecords: 1,
			Data:         clients,
		}

		cliRepo := mocks.NewMockRepository(ctrl)
		cliRepo.EXPECT().FindAll(gomock.Any(), gomock.Any()).Times(1).Return(&resultSearch, nil)

		service := client.NewClientService(cliRepo)
		result := service.FindByCriteria(context.Background(), paginate.PaginateRequest{
			Page:  1,
			Limit: 10,
		})

		data := result.Data.(paginate.Result).Data.([]client.ClientRequestDto)
		assert.Equal(t, result.Status, "success")
		assert.Len(t, data, 1)
	})
}
