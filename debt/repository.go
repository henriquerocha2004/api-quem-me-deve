package debt

import (
	"context"

	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
)

type Reader interface {
	ClientUserDebts(ctx context.Context, clientUserId ulid.ULID) ([]*Debt, error)
	DebtInstallments(ctx context.Context, debtId ulid.ULID) ([]*Installment, error)
	GetDebts(ctx context.Context, pagData paginate.SearchDto) (*PaginationResult, error)
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
