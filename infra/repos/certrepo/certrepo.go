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

func (r *repo) BulkCreate(ctx context.Context, certs []*cert.Cert) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	for _, c := range certs {
		err = r.create(tx, c)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	return tx.Commit()
}

func (r *repo) create(tx *sql.Tx, c *cert.Cert) error {
	_, err := tx.Exec(queryCreate,
		c.Uuid, c.Password, c.AuthorizeRequired, c.IssuerId, c.IssuingProcessId)
	return err
}
