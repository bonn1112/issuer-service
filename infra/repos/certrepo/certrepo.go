package certrepo

import (
	"context"
	"database/sql"

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

func (r *repo) StartBulkCreation(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *repo) AppendToBulkCreation(tx *sql.Tx, c *cert.Cert) error {
	_, err := tx.Exec(queryCreate,
		c.Uuid, c.Password, c.AuthorizeRequired, c.IssuerId, c.IssuingProcessId)
	return err
}
