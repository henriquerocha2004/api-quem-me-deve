package gorm_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	setupdbtests "github.com/henriquerocha2004/quem-me-deve-api/config/setupDbTests"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/debt/gorm"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/helpers"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
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

	gormDB, _ = gorm.NewGorm(dsn)
	sql, _ := gormDB.DB()
	defer sql.Close()
	m.Run()
}

type DebtRepositorySuiteTest struct {
	suite.Suite
}

func (s *DebtRepositorySuiteTest) TearDownTest() {
	err := setupdbtests.TruncateTables(gormDB)
	if err != nil {
		s.Fail("Failed to truncate tables: %v", err)
	}
}

func TestDebtSuit(t *testing.T) {
	suite.Run(t, new(DebtRepositorySuiteTest))
}

func (s *DebtRepositorySuiteTest) TestShouldCreateDebt() {
	repo := gorm.NewGormDebtRepository(gormDB)
	dueDate := time.Now().AddDate(0, 0, 30)

	debt := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Test Debt",
		TotalValue:   100.0,
		DueDate:      &dueDate,
		UserClientId: ulid.Make(),
		Intallments: []debt.Installment{
			{
				Id:          ulid.Make(),
				Description: "First Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      1,
			},
			{
				Id:          ulid.Make(),
				Description: "Second Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      2,
			},
		},
	}
	err := repo.Save(context.Background(), debt)
	s.Assert().NoError(err)

	savedDebt, err := repo.GetDebt(context.Background(), debt.Id)
	s.Assert().NoError(err)
	s.Assert().NotNil(savedDebt)
	s.Assert().Equal(debt.Description, savedDebt.Description)
	s.Assert().Equal(debt.TotalValue, savedDebt.TotalValue)
	s.Assert().Equal(debt.UserClientId, savedDebt.UserClientId)
	s.Assert().Equal(len(debt.Intallments), len(savedDebt.Intallments))
	s.Assert().Empty(savedDebt.CancelInfo)
	s.Assert().Empty(savedDebt.ReversalInfo)
}

func (s *DebtRepositorySuiteTest) TestShouldUpdateDebtInformation() {
	repo := gorm.NewGormDebtRepository(gormDB)
	dueDate := time.Now().AddDate(0, 0, 30)

	debt := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Test Debt",
		TotalValue:   100.0,
		DueDate:      &dueDate,
		UserClientId: ulid.Make(),
		Intallments: []debt.Installment{
			{
				Id:          ulid.Make(),
				Description: "First Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      1,
			},
			{
				Id:          ulid.Make(),
				Description: "Second Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      2,
			},
		},
	}
	err := repo.Save(context.Background(), debt)
	s.Assert().NoError(err)

	debt.Description = "Updated Debt"
	err = repo.Update(context.Background(), debt)
	s.Assert().NoError(err)

	savedDebt, err := repo.GetDebt(context.Background(), debt.Id)
	s.Assert().NoError(err)
	s.Assert().Equal("Updated Debt", savedDebt.Description)
}

func (s *DebtRepositorySuiteTest) TestShouldCancelDebt() {
	repo := gorm.NewGormDebtRepository(gormDB)
	dueDate := time.Now().AddDate(0, 0, 30)

	d := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Test Debt",
		TotalValue:   100.0,
		DueDate:      &dueDate,
		UserClientId: ulid.Make(),
		Intallments: []debt.Installment{
			{
				Id:          ulid.Make(),
				Description: "First Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      1,
			},
			{
				Id:          ulid.Make(),
				Description: "Second Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      2,
			},
		},
	}
	err := repo.Save(context.Background(), d)
	s.Assert().NoError(err)

	cancelInfo := debt.CancelInfo{
		Reason:      "Test Cancel",
		CancelDate:  &dueDate,
		CancelledBy: ulid.Make(),
	}

	d.CancelInfo = &cancelInfo
	err = repo.Update(context.Background(), d)
	s.Assert().NoError(err)

	savedDebt, err := repo.GetDebt(context.Background(), d.Id)
	s.Assert().NoError(err)
	s.Assert().NotNil(savedDebt.CancelInfo)
	s.Assert().Equal(cancelInfo.Reason, savedDebt.CancelInfo.Reason)
}

func (s *DebtRepositorySuiteTest) TestShouldReverseDebt() {
	repo := gorm.NewGormDebtRepository(gormDB)
	dueDate := time.Now().AddDate(0, 0, 30)

	d := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Test Debt",
		TotalValue:   100.0,
		DueDate:      &dueDate,
		UserClientId: ulid.Make(),
		Intallments: []debt.Installment{
			{
				Id:          ulid.Make(),
				Description: "First Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      1,
			},
			{
				Id:          ulid.Make(),
				Description: "Second Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      2,
			},
		},
	}
	err := repo.Save(context.Background(), d)
	s.Assert().NoError(err)

	reversalInfo := debt.ReversalInfo{
		Reason:       "Test Reversal",
		ReversedBy:   ulid.Make(),
		ReversalDate: &dueDate,
	}

	d.ReversalInfo = &reversalInfo
	err = repo.Update(context.Background(), d)
	s.Assert().NoError(err)

	savedDebt, err := repo.GetDebt(context.Background(), d.Id)
	s.Assert().NoError(err)
	s.Assert().NotNil(savedDebt.ReversalInfo)
	s.Assert().Equal(reversalInfo.Reason, savedDebt.ReversalInfo.Reason)
}

func (s *DebtRepositorySuiteTest) TestShouldGetClientUserDebts() {
	repo := gorm.NewGormDebtRepository(gormDB)
	clientUserId := ulid.Make()

	debt1 := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Client Debt 1",
		TotalValue:   100.0,
		DueDate:      nil,
		UserClientId: clientUserId,
	}
	err := repo.Save(context.Background(), debt1)
	s.Assert().NoError(err)

	debt2 := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Client Debt 2",
		TotalValue:   200.0,
		DueDate:      nil,
		UserClientId: clientUserId,
	}
	err = repo.Save(context.Background(), debt2)
	s.Assert().NoError(err)

	debts, err := repo.ClientUserDebts(context.Background(), clientUserId)
	s.Assert().NoError(err)
	s.Assert().Len(debts, 2)
}

func (s *DebtRepositorySuiteTest) TestShouldGetDebtInstallments() {
	repo := gorm.NewGormDebtRepository(gormDB)
	dueDate := time.Now().AddDate(0, 0, 30)

	debt := &debt.Debt{
		Id:           ulid.Make(),
		Description:  "Test Debt",
		TotalValue:   100.0,
		DueDate:      &dueDate,
		UserClientId: ulid.Make(),
		Intallments: []debt.Installment{
			{
				Id:          ulid.Make(),
				Description: "First Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      1,
			},
			{
				Id:          ulid.Make(),
				Description: "Second Installment",
				Value:       50.0,
				DueDate:     &dueDate,
				Status:      debt.Pending,
				Number:      2,
			},
		},
	}
	err := repo.Save(context.Background(), debt)
	s.Assert().NoError(err)

	installments, err := repo.DebtInstallments(context.Background(), debt.Id)
	s.Assert().NoError(err)
	s.Assert().Len(installments, 2)
}

func (s *DebtRepositorySuiteTest) TestShouldGetDebtsWithPagination() {
	repo := gorm.NewGormDebtRepository(gormDB)

	for i := 0; i < 10; i++ {
		dueDate := time.Now().AddDate(0, 0, 30)
		debt := &debt.Debt{
			Id:           ulid.Make(),
			Description:  fmt.Sprintf("Debt %d", i),
			TotalValue:   float64(i * 10),
			DueDate:      &dueDate,
			UserClientId: ulid.Make(),
		}
		err := repo.Save(context.Background(), debt)
		s.Assert().NoError(err)
	}

	pagData := paginate.SearchDto{
		Limit: 5,
	}
	pagData.SetPage(1)

	result, err := repo.GetDebts(context.Background(), pagData)
	s.Assert().NoError(err)
	s.Assert().NotNil(result)
	s.Assert().Len(result.Data, 5)
	s.Assert().Equal(10, result.TotalRecords)
}
