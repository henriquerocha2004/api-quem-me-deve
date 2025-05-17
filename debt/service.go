package debt

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service interface {
	CreateDebt(ctx context.Context, debt *DebtDto) DebtResponse
	GetUserDebts(ctx context.Context, userId ulid.ULID) DebtResponse
}

type debtService struct {
	repo Repository
}

func NewDebtService(repo Repository) *debtService {
	return &debtService{
		repo: repo,
	}
}

func (s *debtService) CreateDebt(ctx context.Context, d *DebtDto) DebtResponse {

	now := time.Now()
	dueDate, _ := time.Parse(time.DateOnly, d.DueDate)
	clientUserId, _ := ulid.Parse(d.UserClientId)

	serviceIds, err := s.putServiceIds(d.ServiceIds)
	if err != nil {
		log.Println("Error parsing service IDs:", err)
		return DebtResponse{
			Status:  "error",
			Message: "invalid service IDs",
		}
	}

	productIds, err := s.putProductIds(d.ProductIds)
	if err != nil {
		log.Println("Error parsing product IDs:", err)
		return DebtResponse{
			Status:  "error",
			Message: "invalid product IDs",
		}
	}

	debt := &Debt{
		Description:          d.Description,
		Id:                   ulid.Make(),
		TotalValue:           d.TotalValue,
		DueDate:              &dueDate,
		Status:               Pending,
		UserClientId:         clientUserId,
		InstallmentsQuantity: d.InstallmentsQuantity,
		ServiceIds:           serviceIds,
		ProductIds:           productIds,
		DebtDate:             &now,
	}

	validationErrors := debt.Validate()
	if len(validationErrors.Errors) > 0 {
		log.Println("Validation errors:", validationErrors)
		return DebtResponse{
			Status:  "error",
			Message: "validation errors",
			Data:    validationErrors,
		}
	}

	err = debt.GenerateInstallments()
	if err != nil {
		log.Println("Error generating installments:", err)
		return DebtResponse{
			Status:  "error",
			Message: "error generating installments",
		}
	}

	err = s.repo.Save(ctx, debt)
	if err != nil {
		log.Println("Error saving debt:", err)
		return DebtResponse{
			Status:  "error",
			Message: "error saving debt",
		}
	}

	return DebtResponse{
		Status:  "success",
		Message: "debt created successfully",
	}

}

func (s *debtService) GetUserDebts(ctx context.Context, userId ulid.ULID) DebtResponse {
	debts, err := s.repo.ClientUserDebts(ctx, userId)
	if err != nil {
		log.Println("Error retrieving debts:", err)
		return DebtResponse{
			Status:  "error",
			Message: "error retrieving debts",
		}
	}

	if len(debts) == 0 {
		return DebtResponse{
			Status:  "success",
			Message: "no debts found",
		}
	}

	var debtsDto []DebtDto

	for _, d := range debts {
		debtsDto = append(debtsDto, DebtDto{
			Id:                   d.Id.String(),
			Description:          d.Description,
			TotalValue:           d.TotalValue,
			DueDate:              d.DueDate.Format(time.DateOnly),
			InstallmentsQuantity: d.InstallmentsQuantity,
			Status:               d.Status.String(),
			UserClientId:         d.UserClientId.String(),
			ProductIds:           s.getProductIds(d.ProductIds),
			ServiceIds:           s.getServiceIds(d.ServiceIds),
		})
	}
	fmt.Println(debtsDto)
	return DebtResponse{
		Status:  "success",
		Message: "debts retrieved successfully",
		Data:    debtsDto,
	}
}

func (s *debtService) putServiceIds(serviceids []string) ([]ulid.ULID, error) {

	if len(serviceids) == 0 {
		return []ulid.ULID{}, nil
	}

	var ids []ulid.ULID
	for _, id := range serviceids {
		ulidId, err := ulid.Parse(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, ulidId)
	}

	return ids, nil
}

func (s *debtService) putProductIds(productids []string) ([]ulid.ULID, error) {

	if len(productids) == 0 {
		return []ulid.ULID{}, nil
	}

	var ids []ulid.ULID
	for _, id := range productids {
		ulidId, err := ulid.Parse(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, ulidId)
	}

	return ids, nil
}

func (s *debtService) getProductIds(productIds []ulid.ULID) []string {
	var ids []string

	if len(productIds) == 0 {
		return []string{}
	}

	for _, id := range productIds {
		ids = append(ids, id.String())
	}
	return ids
}

func (s *debtService) getServiceIds(serviceIds []ulid.ULID) []string {
	var ids []string

	if len(serviceIds) == 0 {
		return []string{}
	}

	for _, id := range serviceIds {
		ids = append(ids, id.String())
	}
	return ids
}
