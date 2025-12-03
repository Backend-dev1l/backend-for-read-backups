package uuidconv

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func SetPgUUID(u uuid.UUID) (pgtype.UUID, error) {
	var pg pgtype.UUID
	copy(pg.Bytes[:], u[:])
	pg.Valid = true
	return pg, nil
}
