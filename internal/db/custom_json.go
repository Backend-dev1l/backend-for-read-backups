package db

import (
	"encoding/json"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

func (u *UserStatistic) UnmarshalJSON(data []byte) error {
	type Alias struct {
		UserID            string  `json:"user_id"`
		TotalWordsLearned int32   `json:"total_words_learned"`
		Accuracy          float64 `json:"accuracy"`
		TotalTime         int32   `json:"total_time"`
	}

	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	parsedUUID, err := uuid.FromString(alias.UserID)
	if err != nil {
		return err
	}

	u.UserID = pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}
	u.TotalWordsLearned = alias.TotalWordsLearned
	u.TotalTime = alias.TotalTime

	accuracyInt := big.NewInt(int64(alias.Accuracy * 100))
	u.Accuracy = pgtype.Numeric{
		Int:   accuracyInt,
		Exp:   -2,
		Valid: true,
	}

	return nil
}
