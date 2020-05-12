package issuerrepo

import (
	"database/sql"

	"github.com/lastrust/issuing-service/domain/issuer"
)

const (
	queryFirstByUuid = `SELECT uuid, name, updated_at, created_at
FROM issuers
WHERE uuid=$1
LIMIT 1;`
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) issuer.Repository {
	return &repo{db}
}

func (r *repo) FirstByUuid(uuid string) (*issuer.Issuer, error) {
	var i issuer.Issuer

	err := r.db.QueryRow(queryFirstByUuid, uuid).
		Scan(&i.Uuid, &i.Name, &i.UpdatedAt, &i.CreatedAt)
	return &i, err
}
