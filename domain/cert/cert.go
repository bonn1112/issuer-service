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
	Err   error
	Done  bool
}

func (tx *Tx) Commit() error {
	tx.Mu.Lock()
	defer tx.Mu.Unlock()
	tx.Done = true
	return tx.SqlTx.Commit()
}

func (tx *Tx) Rollback(err error) {
	tx.Mu.Lock()
	tx.SqlTx.Rollback()
	tx.Err = err
	tx.Done = true
	tx.Mu.Unlock()
}
