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
	"github.com/oklog/ulid/v2"
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

	now := time.Now()

	tx, err := s.rawDb.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return errors.New("failed to update debt")
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

	params := UpdateDebtParams{
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
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	err = qtx.UpdateDebt(ctx, params)
	if err != nil {
		log.Println("Error updating debt:", err)
		return errors.New("failed to update debt")
	}

	if debt.CancelInfo != nil {
		cancelParams := CreateCancelInfoParams{
			ID:     ulid.Make().String(),
			Reason: debt.CancelInfo.Reason,
			CancelDate: sql.NullTime{
				Time:  now,
				Valid: true,
			},
			CancelledBy: debt.CancelInfo.CancelledBy.String(),
			DebtID:      debt.Id.String(),
			CreatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
		}

		err = qtx.CreateCancelInfo(ctx, cancelParams)
		if err != nil {
			log.Println("Error saving cancel info:", err)
			return errors.New("failed to save cancel info")
		}
	}

	if debt.ReversalInfo != nil {
		reversalParams := CreateReversalInfoParams{
			ID:     ulid.Make().String(),
			Reason: debt.ReversalInfo.Reason,
			ReversalDate: sql.NullTime{
				Time:  now,
				Valid: true,
			},
			ReversedBy:              debt.ReversalInfo.ReversedBy.String(),
			ReversedInstallmentQtd:  int32(debt.ReversalInfo.ReversedInstallmentQtd),
			CancelledInstallmentQtd: int32(debt.ReversalInfo.CancelledInstallmentQtd),
			DebtID:                  debt.Id.String(),
			CreatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
		}

		err = qtx.CreateReversalInfo(ctx, reversalParams)
		if err != nil {
			log.Println("Error saving reversal info:", err)
			return errors.New("failed to save reversal info")
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Println("Error committing transaction:", commitErr)
		err = errors.New("failed to commit transaction")
	}

	return nil
}

func (s *SqlcDebtRepository) ClientUserDebts(ctx context.Context, clientUserId string) ([]*debt.Debt, error) {

	debts, err := s.db.ClientUserDebts(ctx, clientUserId)
	if err != nil {
		log.Println("Error fetching client user debts:", err)
		return nil, errors.New("failed to fetch client user debts")
	}

	if len(debts) < 1 {
		return nil, errors.New("no debts found for the specified user client")
	}

	var result []*debt.Debt

	for _, d := range debts {
		debtId, err := ulid.Parse(d.ID)
		if err != nil {
			log.Println("Error parsing debt ID:", err)
			return nil, errors.New("invalid debt ID format")
		}

		totalValue, err := strconv.ParseFloat(d.TotalValue, 64)
		if err != nil {
			log.Println("Error parsing total value:", err)
			return nil, errors.New("failed to retrieve total value")
		}

		debt := &debt.Debt{
			Id:                   debtId,
			Description:          d.Description,
			TotalValue:           totalValue,
			DueDate:              &d.DueDate.Time,
			InstallmentsQuantity: int(d.InstallmentsQuantity),
			UserClientId:         ulid.MustParse(d.UserClientID),
			ProductIds:           s.convertToUlidSlice(d.ProductIds),
			ServiceIds:           s.convertToUlidSlice(d.ServiceIds),
			Status:               debt.StatusValue[d.Status],
			DebtDate:             &d.DebtDate.Time,
		}

		cancelInfo, err := s.getCancelInfo(ctx, debt.Id)
		if err != nil {
			log.Println("Error fetching cancel info:", err)
			return nil, errors.New("failed to fetch cancel info")
		}

		if cancelInfo != nil {
			debt.CancelInfo = cancelInfo
		}

		reversalInfo, err := s.getReversalInfo(ctx, debt.Id)
		if err != nil {
			log.Println("Error fetching reversal info:", err)
			return nil, errors.New("failed to fetch reversal info")
		}

		if reversalInfo != nil {
			debt.ReversalInfo = reversalInfo
		}

		result = append(result, debt)
	}

	return result, nil
}

func (s *SqlcDebtRepository) DebtInstallments(ctx context.Context, debtId string) ([]*debt.Installment, error) {
	installments, err := s.db.DebtInstallments(ctx, debtId)
	if err != nil {
		log.Println("Error fetching debt installments:", err)
		return nil, errors.New("failed to fetch debt installments")
	}

	if len(installments) < 1 {
		return nil, errors.New("no installments found for the specified debt")
	}

	var result []*debt.Installment

	for _, i := range installments {
		installmentId, err := ulid.Parse(i.ID)
		if err != nil {
			log.Println("Error parsing installment ID:", err)
			return nil, errors.New("invalid installment ID format")
		}

		value, err := strconv.ParseFloat(i.Value, 64)
		if err != nil {
			log.Println("Error parsing installment value:", err)
			return nil, errors.New("failed to retrieve installment value")
		}

		dueDate := i.DueDate.Time
		debDate := i.DebDate.Time

		installment := &debt.Installment{
			Id:          installmentId,
			Description: i.Description,
			Value:       value,
			DueDate:     &dueDate,
			DebDate:     &debDate,
			Status:      debt.StatusValue[i.Status],
			Number:      int(i.Number),
		}

		result = append(result, installment)
	}

	return result, nil
}

func (s *SqlcDebtRepository) GetDebts(ctx context.Context, pagData paginate.SearchDto) (*debt.PaginationResult, error) {
	query := `
		SELECT 
			id,
			description,
			total_value,
			due_date,
			debt_date,
			status,
			user_client_id,
			product_ids,
			service_ids,
			finished_at,
			installments_quantity
		FROM public.debts
	`
	if pagData.TermSearch != "" {
		query += " WHERE description ILIKE '%' || $1 || '%' "
	}

	if len(pagData.ColumnSearch) > 0 {
		for _, column := range pagData.ColumnSearch {
			query += " AND "
			query += column.ColumnName + " = " + column.ColumnValue + " "
		}
	}

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

func (s *SqlcDebtRepository) convertToUlidSlice(ids []string) []ulid.ULID {
	if len(ids) < 1 {
		return nil
	}

	var ulids []ulid.ULID
	for _, id := range ids {
		ulidId, err := ulid.Parse(id)
		if err != nil {
			log.Println("Error parsing ID to ULID:", err)
			return nil
		}
		ulids = append(ulids, ulidId)
	}
	return ulids
}

func (s *SqlcDebtRepository) getCancelInfo(ctx context.Context, debtId ulid.ULID) (*debt.CancelInfo, error) {
	cancel, err := s.db.DebtCancelInfo(ctx, debtId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Println("Error fetching cancel info:", err)
		return nil, errors.New("failed to fetch cancel info")
	}

	cancelInfo := &debt.CancelInfo{
		Reason:      cancel.Reason,
		CancelDate:  &cancel.CancelDate.Time,
		CancelledBy: ulid.MustParse(cancel.CancelledBy),
	}

	return cancelInfo, nil
}

func (s *SqlcDebtRepository) getReversalInfo(ctx context.Context, debtId ulid.ULID) (*debt.ReversalInfo, error) {
	reversal, err := s.db.DebtReversalInfo(ctx, debtId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Println("Error fetching reversal info:", err)
		return nil, errors.New("failed to fetch reversal info")
	}

	reversalInfo := &debt.ReversalInfo{
		Reason:                  reversal.Reason,
		ReversalDate:            &reversal.ReversalDate.Time,
		ReversedBy:              ulid.MustParse(reversal.ReversedBy),
		ReversedInstallmentQtd:  int(reversal.ReversedInstallmentQtd),
		CancelledInstallmentQtd: int(reversal.CancelledInstallmentQtd),
	}

	return reversalInfo, nil
}
