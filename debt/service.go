package debt

import (
	"context"
	"log"
	"time"

	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	CreateDebt(ctx context.Context, debt *DebtDto) DebtResponse
	GetUserDebts(ctx context.Context, userId ulid.ULID) DebtResponse
	GetDebtInstallments(ctx context.Context, clientId, debtId ulid.ULID) DebtResponse
	Debts(ctx context.Context, params paginate.PaginateRequest) DebtResponse
	PayInstallment(ctx context.Context, pgInfo *PaymentInfoDto) DebtResponse
}

type debtService struct {
	debtRepo   Repository
	clientRepo ClientReader
}

func NewDebtService(debtRepo Repository, cliRepo ClientReader) *debtService {
	return &debtService{
		debtRepo:   debtRepo,
		clientRepo: cliRepo,
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

	err = s.debtRepo.Save(ctx, debt)
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
	debts, err := s.debtRepo.ClientUserDebts(ctx, userId)
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

	debtsDto := s.convertToDebtDto(debts)

	return DebtResponse{
		Status:  "success",
		Message: "debts retrieved successfully",
		Data:    debtsDto,
	}
}

func (s *debtService) GetDebtInstallments(ctx context.Context, clientId, debtId ulid.ULID) DebtResponse {
	cliExists, err := s.clientRepo.ClientExists(ctx, clientId)
	if err != nil {
		log.Println(err)
		return DebtResponse{
			Status:  "error",
			Message: "error in validate clientid provided",
		}
	}

	if !cliExists {
		return DebtResponse{
			Status:  "error",
			Message: "client not found",
		}
	}

	installments, err := s.debtRepo.DebtInstallments(ctx, debtId)
	if err != nil {
		log.Println(err)
		return DebtResponse{
			Status:  "error",
			Message: "failed to get debt installments",
		}
	}

	var installmentsDto []InstallmentDto

	for _, installment := range installments {
		paymentDate := ""

		if installment.PaymentDate != nil {
			paymentDate = installment.PaymentDate.Format(time.DateOnly)
		}

		installmentsDto = append(installmentsDto, InstallmentDto{
			Id:            installment.Id.String(),
			Description:   installment.Description,
			Value:         installment.Value,
			DueDate:       installment.DueDate.Format(time.DateOnly),
			DebDate:       installment.DebDate.Format(time.DateOnly),
			Status:        installment.Status.String(),
			PaymentDate:   paymentDate,
			PaymentMethod: installment.PaymentMethod,
			Number:        installment.Number,
		})
	}

	return DebtResponse{
		Status:  "success",
		Message: "installments retrieved",
		Data:    installmentsDto,
	}
}

func (s *debtService) Debts(ctx context.Context, params paginate.PaginateRequest) DebtResponse {
	pagDto := paginate.SearchDto{
		Limit:         params.Limit,
		TermSearch:    params.SearchTerm,
		SortField:     params.SortField,
		SortDirection: params.SortDirection,
	}

	pagDto.SetPage(params.Page)
	pagDto.AddColumnSearch(params.ColumnSearch)

	result, err := s.debtRepo.GetDebts(ctx, pagDto)

	if err != nil {
		log.Println("Error retrieving debts:", err)
		return DebtResponse{
			Status:  "error",
			Message: "error retrieving debts",
		}
	}

	debtsDto := s.convertToDebtDto(result.Data)

	return DebtResponse{
		Status:  "success",
		Message: "debts retrieved successfully",
		Data: paginate.Result{
			TotalRecords: result.TotalRecords,
			Data:         debtsDto,
		},
	}
}

func (s *debtService) PayInstallment(ctx context.Context, pgInfo *PaymentInfoDto) DebtResponse {
	debtId, err := ulid.Parse(pgInfo.DebtId)
	if err != nil {
		log.Println("Error parsing debt ID:", err)
		return DebtResponse{
			Status:  "error",
			Message: "invalid debt ID",
		}
	}

	debt, err := s.debtRepo.GetDebt(ctx, debtId)
	if err != nil {
		log.Println("Error retrieving debt:", err)
		return DebtResponse{
			Status:  "error",
			Message: "error retrieving debt",
		}
	}

	if debt == nil {
		return DebtResponse{
			Status:  "error",
			Message: "debt not found",
		}
	}

	err = debt.PayInstallment(pgInfo)
	if err != nil {
		log.Println("Error paying installment:", err)
		return DebtResponse{
			Status:  "error",
			Message: err.Error(),
		}
	}

	err = s.debtRepo.Update(ctx, debt)
	if err != nil {
		log.Println("Error updating debt:", err)
		return DebtResponse{
			Status:  "error",
			Message: "error updating debt",
		}
	}

	return DebtResponse{
		Status:  "success",
		Message: "installment paid successfully",
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

func (s *debtService) convertToDebtDto(debts []*Debt) []DebtDto {
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

	return debtsDto
}
