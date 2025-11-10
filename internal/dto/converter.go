package dto

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

func ParseUUID(uuidStr string) (pgtype.UUID, error) {
	parsed, err := uuid.FromString(uuidStr)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{
		Bytes: parsed,
		Valid: true,
	}, nil
}

func ParseNumeric(value float64) pgtype.Numeric {

	bigInt := big.NewInt(int64(value * 100))

	return pgtype.Numeric{
		Int:   bigInt,
		Exp:   -2,
		Valid: true,
	}
}
func UUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}

	u := uuid.UUID(pgUUID.Bytes)
	return u.String()
}

func NumericToFloat64(numeric pgtype.Numeric) float64 {
	if !numeric.Valid {
		return 0
	}

	if numeric.Int == nil {
		return 0
	}

	floatVal := float64(numeric.Int.Int64())
	exp := float64(numeric.Exp)

	if exp < 0 {
		divisor := 1.0
		for i := int32(0); i < -numeric.Exp; i++ {
			divisor *= 10
		}
		return floatVal / divisor
	}

	multiplier := 1.0
	for i := int32(0); i < numeric.Exp; i++ {
		multiplier *= 10
	}
	return floatVal * multiplier
}
