package memory

import (
	"context"

	"github.com/henriquerocha2004/quem-me-deve-api/debt"
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
