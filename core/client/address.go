package client

import "github.com/oklog/ulid/v2"

type Address struct {
	Id           ulid.ULID
	Street       string
	Neighborhood string
	City         string
	State        string
	ZipCode      string
}
