package client

import "github.com/oklog/ulid/v2"

type Phone struct {
	Id          ulid.ULID
	Description string
	Number      string
}
