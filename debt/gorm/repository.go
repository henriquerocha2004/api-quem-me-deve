package gorm

import (
	"context"

	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/lib/pq"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GormDebtRepository struct {
	db *gorm.DB
}

func NewGormDebtRepository(db *gorm.DB) *GormDebtRepository {
	return &GormDebtRepository{db: db}
}

func (g *GormDebtRepository) ClientUserDebts(ctx context.Context, clientUserId ulid.ULID) ([]*debt.Debt, error) {
	var models []Debt
	result := g.db.Where("user_client_id = ?", clientUserId.String()).
		Preload("Installments").
		Preload("CancelInfo").
		Preload("ReversalInfo").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	var debts []*debt.Debt
	for _, model := range models {
		debts = append(debts, &debt.Debt{
			Id:                   ulid.MustParse(model.ID),
			Description:          model.Description,
			TotalValue:           model.TotalValue,
			DueDate:              model.DueDate,
			InstallmentsQuantity: model.InstallmentsQuantity,
			UserClientId:         ulid.MustParse(model.UserClientId),
			ProductIds:           g.parseToUlidSlice(model.ProductIds),
			ServiceIds:           g.parseToUlidSlice(model.ServiceIds),
			Status:               debt.StatusValue[model.Status],
			DebtDate:             model.DebtDate,
			Intallments:          g.parseInstallments(model.Installments),
			CancelInfo:           g.parseCancelInfo(model.CancelInfo),
			ReversalInfo:         g.parseReversalInfo(model.ReversalInfo),
		})
	}

	return debts, nil
}
func (g *GormDebtRepository) DebtInstallments(ctx context.Context, debtId ulid.ULID) ([]*debt.Installment, error) {
	var installments []Installment
	result := g.db.Where("debt_id = ?", debtId.String()).
		Find(&installments)

	if result.Error != nil {
		return nil, result.Error
	}

	var parsedInstallments []*debt.Installment
	for _, installment := range installments {
		i := debt.Installment{
			Id:            ulid.MustParse(installment.Id),
			Description:   installment.Description,
			Value:         installment.Value,
			DueDate:       installment.DueDate,
			DebDate:       installment.DebDate,
			Status:        debt.StatusValue[installment.Status],
			PaymentDate:   installment.PaymentDate,
			PaymentMethod: installment.PaymentMethod,
			Number:        installment.Number,
		}
		parsedInstallments = append(parsedInstallments, &i)
	}

	return parsedInstallments, nil
}
func (g *GormDebtRepository) GetDebts(ctx context.Context, pagData paginate.SearchDto) (*debt.PaginationResult, error) {
	var models []Debt
	var total int64

	query := g.db.Model(&Debt{}).
		Count(&total).
		Offset(pagData.Offset()).
		Limit(pagData.Limit).
		Order("created_at DESC").
		Preload("Installments").
		Preload("CancelInfo").
		Preload("ReversalInfo")

	if pagData.TermSearch != "" {
		query = query.Where("description LIKE ?", "%"+pagData.TermSearch+"%")
	}

	result := query.Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	var debts []*debt.Debt
	for _, model := range models {
		debts = append(debts, &debt.Debt{
			Id:                   ulid.MustParse(model.ID),
			Description:          model.Description,
			TotalValue:           model.TotalValue,
			DueDate:              model.DueDate,
			InstallmentsQuantity: model.InstallmentsQuantity,
			UserClientId:         ulid.MustParse(model.UserClientId),
			ProductIds:           g.parseToUlidSlice(model.ProductIds),
			ServiceIds:           g.parseToUlidSlice(model.ServiceIds),
			Status:               debt.StatusValue[model.Status],
			DebtDate:             model.DebtDate,
			Intallments:          g.parseInstallments(model.Installments),
			CancelInfo:           g.parseCancelInfo(model.CancelInfo),
			ReversalInfo:         g.parseReversalInfo(model.ReversalInfo),
		})
	}

	return &debt.PaginationResult{
		TotalRecords: int(total),
		Data:         debts,
	}, nil
}
func (g *GormDebtRepository) GetDebt(ctx context.Context, debtId ulid.ULID) (*debt.Debt, error) {

	var model Debt
	result := g.db.Where("id = ?", debtId.String()).
		Preload("Installments").
		Preload("CancelInfo").
		Preload("ReversalInfo").
		First(&model)

	if result.Error != nil {
		return nil, result.Error
	}

	debt := &debt.Debt{
		Id:                   ulid.MustParse(model.ID),
		Description:          model.Description,
		TotalValue:           model.TotalValue,
		DueDate:              model.DueDate,
		InstallmentsQuantity: model.InstallmentsQuantity,
		UserClientId:         ulid.MustParse(model.UserClientId),
		ProductIds:           g.parseToUlidSlice(model.ProductIds),
		ServiceIds:           g.parseToUlidSlice(model.ServiceIds),
		Status:               debt.StatusValue[model.Status],
		DebtDate:             model.DebtDate,
		Intallments:          g.parseInstallments(model.Installments),
		CancelInfo:           g.parseCancelInfo(model.CancelInfo),
		ReversalInfo:         g.parseReversalInfo(model.ReversalInfo),
	}

	return debt, nil
}
func (g *GormDebtRepository) Save(ctx context.Context, debt *debt.Debt) error {

	products := g.pushProducts(debt)
	services := g.pushServices(debt)
	installments := []Installment{}

	if len(debt.Intallments) > 0 {
		for _, installment := range debt.Intallments {
			installments = append(installments, Installment{
				Id:            installment.Id.String(),
				Description:   installment.Description,
				Value:         installment.Value,
				DueDate:       installment.DueDate,
				DebDate:       installment.DebDate,
				Status:        installment.Status.String(),
				PaymentDate:   installment.PaymentDate,
				PaymentMethod: installment.PaymentMethod,
				Number:        installment.Number,
				DebtId:        debt.Id.String(),
			})
		}
	}

	model := &Debt{
		ID:                   debt.Id.String(),
		Description:          debt.Description,
		TotalValue:           debt.TotalValue,
		DueDate:              debt.DueDate,
		InstallmentsQuantity: debt.InstallmentsQuantity,
		UserClientId:         debt.UserClientId.String(),
		ProductIds:           pq.StringArray(products),
		ServiceIds:           pq.StringArray(services),
		Status:               debt.Status.String(),
		DebtDate:             debt.DebtDate,
		Installments:         installments,
	}

	tx := g.db.Begin()

	result := tx.WithContext(ctx).Create(model)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}
func (g *GormDebtRepository) Update(ctx context.Context, debt *debt.Debt) error {
	products := g.pushProducts(debt)
	services := g.pushServices(debt)

	model := &Debt{
		ID:                   debt.Id.String(),
		Description:          debt.Description,
		TotalValue:           debt.TotalValue,
		DueDate:              debt.DueDate,
		InstallmentsQuantity: debt.InstallmentsQuantity,
		UserClientId:         debt.UserClientId.String(),
		ProductIds:           products,
		ServiceIds:           services,
		Status:               debt.Status.String(),
		DebtDate:             debt.DebtDate,
		Installments:         g.convertInstallmentsToModel(debt.Intallments),
	}

	tx := g.db.Begin()

	err := tx.WithContext(ctx).Model(&Debt{}).
		Where("id = ?", debt.Id.String()).
		Updates(model).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.WithContext(ctx).Model(model).
		Association("Installments").
		Replace(model.Installments)

	if err != nil {
		tx.Rollback()
		return err
	}

	if debt.CancelInfo != nil {
		cancelInfo := CancelInfo{
			Id:          ulid.Make().String(),
			Reason:      debt.CancelInfo.Reason,
			CancelDate:  debt.CancelInfo.CancelDate,
			CancelledBy: debt.CancelInfo.CancelledBy.String(),
			DebtId:      debt.Id.String(),
		}

		err = tx.WithContext(ctx).Model(&CancelInfo{}).
			Create(&cancelInfo).Error

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if debt.ReversalInfo != nil {
		reversalInfo := ReversalInfo{
			Id:                      ulid.Make().String(),
			Reason:                  debt.ReversalInfo.Reason,
			ReversalDate:            debt.ReversalInfo.ReversalDate,
			ReversedBy:              debt.ReversalInfo.ReversedBy.String(),
			ReversedInstallmentQtd:  debt.ReversalInfo.ReversedInstallmentQtd,
			CancelledInstallmentQtd: debt.ReversalInfo.CancelledInstallmentQtd,
			DebtId:                  debt.Id.String(),
		}
		err = tx.WithContext(ctx).Model(&ReversalInfo{}).
			Create(&reversalInfo).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
func (s *GormDebtRepository) pushProducts(debt *debt.Debt) []string {

	var productIds []string

	if len(debt.ProductIds) > 0 {
		for _, productId := range debt.ProductIds {
			productIds = append(productIds, productId.String())
		}
	}

	return productIds
}
func (s *GormDebtRepository) pushServices(debt *debt.Debt) []string {
	var serviceIds []string

	if len(debt.ServiceIds) > 0 {
		for _, serviceId := range debt.ServiceIds {
			serviceIds = append(serviceIds, serviceId.String())
		}
	}

	return serviceIds
}
func (s *GormDebtRepository) parseToUlidSlice(productIds []string) []ulid.ULID {
	var parsedIds []ulid.ULID

	for _, productId := range productIds {
		parsedId, err := ulid.Parse(productId)
		if err == nil {
			parsedIds = append(parsedIds, parsedId)
		}
	}

	return parsedIds
}
func (s *GormDebtRepository) parseInstallments(installments []Installment) []debt.Installment {
	var parsedInstallments []debt.Installment

	for _, installment := range installments {
		parsedInstallments = append(parsedInstallments, debt.Installment{
			Id:            ulid.MustParse(installment.Id),
			Description:   installment.Description,
			Value:         installment.Value,
			DueDate:       installment.DueDate,
			DebDate:       installment.DebDate,
			Status:        debt.StatusValue[installment.Status],
			PaymentDate:   installment.PaymentDate,
			PaymentMethod: installment.PaymentMethod,
			Number:        installment.Number,
		})
	}

	return parsedInstallments
}

func (g *GormDebtRepository) parseCancelInfo(cancelInfo CancelInfo) *debt.CancelInfo {
	if cancelInfo.Id == "" {
		return nil
	}

	return &debt.CancelInfo{
		Reason:      cancelInfo.Reason,
		CancelDate:  cancelInfo.CancelDate,
		CancelledBy: ulid.MustParse(cancelInfo.CancelledBy),
	}
}

func (g *GormDebtRepository) parseReversalInfo(reversalInfo ReversalInfo) *debt.ReversalInfo {
	if reversalInfo.Id == "" {
		return nil
	}

	return &debt.ReversalInfo{
		Reason:                  reversalInfo.Reason,
		ReversalDate:            reversalInfo.ReversalDate,
		ReversedBy:              ulid.MustParse(reversalInfo.ReversedBy),
		ReversedInstallmentQtd:  reversalInfo.ReversedInstallmentQtd,
		CancelledInstallmentQtd: reversalInfo.CancelledInstallmentQtd,
	}
}

func (g *GormDebtRepository) convertInstallmentsToModel(installments []debt.Installment) []Installment {
	var models []Installment
	for _, installment := range installments {
		models = append(models, Installment{
			Id:            installment.Id.String(),
			Description:   installment.Description,
			Value:         installment.Value,
			DueDate:       installment.DueDate,
			DebDate:       installment.DebDate,
			Status:        installment.Status.String(),
			PaymentDate:   installment.PaymentDate,
			PaymentMethod: installment.PaymentMethod,
			Number:        installment.Number,
		})
	}
	return models
}
