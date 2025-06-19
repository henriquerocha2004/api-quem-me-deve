package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
)

type SqlcDebtRepository struct {
	db    *Queries
	rawDb *sql.DB
}

func NewSqlcDebtRepository(db *Queries, rdb *sql.DB) *SqlcDebtRepository {
	return &SqlcDebtRepository{
		db:    db,
		rawDb: rdb,
	}
}

func (s *SqlcDebtRepository) Save(ctx context.Context, debt *debt.Debt) error {

	now := time.Now()

	tx, err := s.rawDb.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction:", err)
		return errors.New("failed to start transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := s.db.WithTx(tx)

	params := SaveDebtParams{
		ID:          debt.Id.String(),
		Description: debt.Description,
		TotalValue:  strconv.FormatFloat(debt.TotalValue, 'f', -1, 64),
		DueDate: sql.NullTime{
			Time:  *debt.DueDate,
			Valid: debt.DueDate != nil,
		},
		InstallmentsQuantity: int32(debt.InstallmentsQuantity),
		DebtDate: sql.NullTime{
			Time:  *debt.DebtDate,
			Valid: debt.DebtDate != nil,
		},
		Status:       debt.Status.String(),
		UserClientID: debt.UserClientId.String(),
		ProductIds:   s.pushProducts(debt),
		ServiceIds:   s.pushServices(debt),
		CreatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	err = qtx.SaveDebt(ctx, params)
	if err != nil {
		log.Println("Error saving debt:", err)
		return errors.New("failed to save debt")
	}

	if len(debt.Intallments) < 1 {
		return nil
	}

	for _, installment := range debt.Intallments {
		installmentParams := CreateInstallmentParams{
			ID:          installment.Id.String(),
			Description: installment.Description,
			Value:       strconv.FormatFloat(installment.Value, 'f', -1, 64),
			DueDate: sql.NullTime{
				Time:  *installment.DueDate,
				Valid: installment.DueDate != nil,
			},
			DebDate: sql.NullTime{
				Time:  *installment.DebDate,
				Valid: installment.DebDate != nil,
			},
			Status: installment.Status.String(),
			Number: int32(installment.Number),
			DebtID: debt.Id.String(),
			CreatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
		}

		err = qtx.CreateInstallment(ctx, installmentParams)
		if err != nil {
			log.Println("Error saving installment:", err)
			return errors.New("failed to save installment")
		}

	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Println("Error committing transaction:", commitErr)
		err = errors.New("failed to commit transaction")
	}

	return nil
}

func (s *SqlcDebtRepository) Update(ctx context.Context, debt *debt.Debt) error {
	return nil
}

func (s *SqlcDebtRepository) ClientUserDebts(ctx context.Context, clientUserId string) ([]*debt.Debt, error) {
	return nil, nil
}

func (s *SqlcDebtRepository) DebtInstallments(ctx context.Context, debtId string) ([]*debt.Installment, error) {
	return nil, nil
}

func (s *SqlcDebtRepository) GetDebts(ctx context.Context, pagData paginate.SearchDto) (*debt.PaginationResult, error) {
	return nil, nil
}

func (s *SqlcDebtRepository) GetDebt(ctx context.Context, debtId string) (*debt.Debt, error) {
	return nil, nil
}

func (s *SqlcDebtRepository) pushProducts(debt *debt.Debt) []string {

	var productIds []string

	if len(debt.ProductIds) > 0 {
		for _, productId := range debt.ProductIds {
			productIds = append(productIds, productId.String())
		}
	}

	return productIds
}

func (s *SqlcDebtRepository) pushServices(debt *debt.Debt) []string {
	var serviceIds []string

	if len(debt.ServiceIds) > 0 {
		for _, serviceId := range debt.ServiceIds {
			serviceIds = append(serviceIds, serviceId.String())
		}
	}

	return serviceIds
}
