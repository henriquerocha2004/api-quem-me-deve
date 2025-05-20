package memory

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type ClientDebtReader struct{}

func NewClientDebtReader() *ClientDebtReader {
	return &ClientDebtReader{}
}

func (c *ClientDebtReader) ClientExists(ctx context.Context, id ulid.ULID) (bool, error) {
	return true, nil
}
