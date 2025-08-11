package gorm

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	setupdbtests "github.com/henriquerocha2004/quem-me-deve-api/config/setupDbTests"
	"github.com/henriquerocha2004/quem-me-deve-api/core/client"
	ormdb "github.com/henriquerocha2004/quem-me-deve-api/core/shared/gorm"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/document"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/helpers"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	orm "gorm.io/gorm"
)

var gormDB *orm.DB = nil

func TestMain(m *testing.M) {
	envPath := helpers.ProjetctRoot() + ".env.testing"
	err := godotenv.Overload(envPath)
	if err != nil {
		log.Println(err)
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	gormDB, err = ormdb.NewGorm(dsn)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	sql, err := gormDB.DB()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	sql.SetMaxIdleConns(10)
	sql.SetMaxOpenConns(100)
	sql.SetConnMaxLifetime(30 * time.Minute)

	defer sql.Close()
	m.Run()
}

type ClientRepositorySuiteTest struct {
	suite.Suite
}

func (s *ClientRepositorySuiteTest) TearDownTest() {
	err := setupdbtests.TruncateTables(gormDB)
	if err != nil {
		s.Fail("Failed to truncate tables: %v", err)
	}
}

func TestClientRepositorySuite(t *testing.T) {
	suite.Run(t, new(ClientRepositorySuiteTest))
}

func (s *ClientRepositorySuiteTest) TestShouldCreateClient() {
	clientRepo := NewGormClientRepository(gormDB)
	now := time.Now()
	client := &client.Client{
		Id:         ulid.Make(),
		Name:       "John",
		LastName:   "Doe",
		EntityType: client.Individual,
		Document:   document.Document("61824136030"),
		BirthDay:   &now,
	}

	err := clientRepo.Create(context.Background(), client)
	s.NoError(err, "Expected no error when creating client")

	cliDb, err := clientRepo.FindById(context.Background(), client.Id)
	s.NoError(err, "Expected no error when finding client by ID")
	s.NotNil(cliDb, "Expected client to be found")
	s.Equal(client.Id, cliDb.Id, "Expected found client ID to match")
}

func (s *ClientRepositorySuiteTest) TestShouldUpdateClient() {
	clientRepo := NewGormClientRepository(gormDB)
	now := time.Now()
	client := &client.Client{
		Id:         ulid.Make(),
		Name:       "John",
		LastName:   "Doe",
		EntityType: client.Individual,
		Document:   document.Document("61824136030"),
		BirthDay:   &now,
		Addresses: []client.Address{
			{
				Id:           ulid.Make(),
				Street:       "123 Main St",
				City:         "Anytown",
				State:        "CA",
				ZipCode:      "12345",
				Neighborhood: "Downtown",
			},
		},
		Phones: []client.Phone{},
	}

	err := clientRepo.Create(context.Background(), client)
	s.NoError(err, "Expected no error when creating client")

	client.Name = "Jane"
	client.LastName = "Smith"

	err = clientRepo.Update(context.Background(), client)
	s.NoError(err, "Expected no error when updating client")

	cliDb, err := clientRepo.FindById(context.Background(), client.Id)
	s.NoError(err, "Expected no error when finding updated client by ID")
	s.NotNil(cliDb, "Expected updated client to be found")
	s.Equal("Jane", cliDb.Name, "Expected updated name to match")
	s.Equal("Smith", cliDb.LastName, "Expected updated last name to match")
}
