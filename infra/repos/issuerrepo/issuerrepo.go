package issuerrepo

import (
	"database/sql"

	"github.com/lastrust/issuing-service/domain/issuer"
)

const (
	queryFirstByName = `SELECT uuid, name, updated_at, created_at
FROM issuers
WHERE name=$2
LIMIT 1;`
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) issuer.Repository {
	return &repo{db}
}

func (r *repo) FirstByName(name string) (*issuer.Issuer, error) {
	var i issuer.Issuer

	err := r.db.QueryRow(queryFirstByName, name).
		Scan(&i.Uuid, &i.Name, &i.UpdatedAt, &i.CreatedAt)
	return &i, err
}
