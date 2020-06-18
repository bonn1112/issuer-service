package certrepo

import (
	"context"
	"database/sql"
	"sync"

	"github.com/lastrust/issuing-service/domain/cert"
)

const (
	queryCreate = `INSERT INTO certificates
(uuid, password, authorize_required, issuer_id, issuing_process_id)
VALUES ($1, $2, $3, $4, $5);`
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) cert.Repository {
	return &repo{db}
}

func (r *repo) StartBulkCreation(ctx context.Context) (*cert.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &cert.Tx{SqlTx: tx, Mu: sync.Mutex{}}, nil
}

func (r *repo) AppendToBulkCreation(tx *cert.Tx, c *cert.Cert) error {
	tx.Mu.Lock()
	_, err := tx.SqlTx.Exec(queryCreate,
		c.Uuid, c.Password, c.AuthorizeRequired, c.IssuerId, c.IssuingProcessId)
	tx.Mu.Unlock()
	return err
}
