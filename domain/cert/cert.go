package cert

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

//go:generate mockgen -destination=../../mocks/cert_Repository.go -package=mocks github.com/lastrust/issuing-service/domain/cert Repository

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
	StartBulkCreation(ctx context.Context) (*Tx, error)
	AppendToBulkCreation(tx *Tx, c *Cert) error
}

type Tx struct {
	SqlTx *sql.Tx
	Mu    sync.Mutex
}
