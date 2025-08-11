package client

import (
	"testing"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/pkg/document"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateClient(t *testing.T) {

	birth, _ := time.Parse("2006-01-02", "2000-01-01")

	client := Client{
		Id:         ulid.Make(),
		Name:       "Henrique",
		LastName:   "Souza",
		EntityType: EntityType("PF"),
		Document:   document.Document("61472869001"),
		BirthDay:   &birth,
	}

	err := client.validate()
	assert.NoError(t, err)
}

func TestShouldReturnErrorIfpassInvalidEntity(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")

	client := Client{
		Id:         ulid.Make(),
		Name:       "Henrique",
		LastName:   "Souza",
		EntityType: EntityType("AV"),
		Document:   document.Document("614.728.690-01"),
		BirthDay:   &birth,
	}

	err := client.validate()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "the entity type informed is invalid")
}

func TestShouldReturnErrorIfPassInvalidDocument(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")

	client := Client{
		Id:         ulid.Make(),
		Name:       "Henrique",
		LastName:   "Souza",
		EntityType: EntityType("PF"),
		Document:   document.Document("614.728.690-090"),
		BirthDay:   &birth,
	}

	err := client.validate()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid document")
}

func TestShouldCreateClientWithAddress(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")

	client := Client{
		Id:         ulid.Make(),
		Name:       "Henrique",
		LastName:   "Souza",
		EntityType: EntityType("PF"),
		Document:   document.Document("614.728.690-090"),
		BirthDay:   &birth,
	}

	client.addAddress("Rua dos Bobos", "Bairro", "Cidade", "BA", "32344-032")

	assert.Len(t, client.Addresses, 1)
}

func TestShouldCreateClientWithPhone(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")

	client := Client{
		Id:         ulid.Make(),
		Name:       "Henrique",
		LastName:   "Souza",
		EntityType: EntityType("PF"),
		Document:   document.Document("614.728.690-090"),
		BirthDay:   &birth,
	}

	client.addPhone("Residencial", "712939393939")
	assert.Len(t, client.Phones, 1)
}
