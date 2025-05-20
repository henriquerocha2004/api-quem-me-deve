package debt

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Reader interface {
	ClientUserDebts(ctx context.Context, clientUserId ulid.ULID) ([]*Debt, error)
	DebtInstallments(ctx context.Context, debtId ulid.ULID) ([]*Installment, error)
}

type Writer interface {
	Save(ctx context.Context, debt *Debt) error
}

type Repository interface {
	Writer
	Reader
}

type ClientReader interface {
	ClientExists(ctx context.Context, id ulid.ULID) (bool, error)
}
