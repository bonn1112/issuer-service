package issuer

import "time"

type Issuer struct {
	Uuid      string
	Name      string
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Repository interface {
	FirstByUuid(uuid string) (*Issuer, error)
}
