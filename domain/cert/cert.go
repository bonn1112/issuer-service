package cert

import (
	"context"
	"time"
)

type Cert struct {
	Uuid              string
	Password          string
	AuthorizeRequired bool
	IssuerId          string
	IssuingProcessId  string
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

type Repository interface {
	BulkCreate(context.Context, []*Cert) error
}
