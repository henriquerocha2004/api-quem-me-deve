package memory

import (
	"context"

	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
)

type DebtMemoryRepository struct {
	Debts []debt.Debt
}

func NewMemoryRepository() *DebtMemoryRepository {
	return &DebtMemoryRepository{
		Debts: []debt.Debt{},
	}
}

func (m *DebtMemoryRepository) Save(ctx context.Context, debt *debt.Debt) error {
	m.Debts = append(m.Debts, *debt)
	return nil
}

func (m *DebtMemoryRepository) ClientUserDebts(ctx context.Context, clientId ulid.ULID) ([]*debt.Debt, error) {
	var clientDebt []*debt.Debt

	for _, debt := range m.Debts {
		if debt.UserClientId != clientId {
			continue
		}

		clientDebt = append(clientDebt, &debt)
	}

	return clientDebt, nil
}

func (m *DebtMemoryRepository) DebtInstallments(ctx context.Context, debtId ulid.ULID) ([]*debt.Installment, error) {
	var installments []*debt.Installment

	for _, debt := range m.Debts {
		if debt.Id != debtId {
			continue
		}

		for _, installment := range debt.Intallments {
			installments = append(installments, &installment)
		}
	}

	return installments, nil
}

func (m *DebtMemoryRepository) GetDebts(ctx context.Context, pagData paginate.SearchDto) (*debt.PaginationResult, error) {
	var debts []*debt.Debt
	for _, debt := range m.Debts {
		debts = append(debts, &debt)
	}

	return &debt.PaginationResult{
		TotalRecords: len(m.Debts),
		Data:         debts,
	}, nil
}
