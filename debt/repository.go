package debt

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Reader interface {
	ClientUserDebts(ctx context.Context, clientUserId ulid.ULID) ([]*Debt, error)
}

type Writer interface {
	Save(ctx context.Context, debt *Debt) error
}

type Repository interface {
	Writer
	Reader
}
