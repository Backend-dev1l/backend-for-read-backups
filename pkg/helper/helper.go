package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"test-http/pkg/fault"

	"github.com/jackc/pgx/v5/pgtype"
)

func HTTPError(w http.ResponseWriter, err error) error {
	var f *fault.Fault

	if errors.As(err, &f) {
		writeFaultResponse(w, f)
		return nil
	}

	// Use a generic unhandled fault as a fallback.
	fallbackFault := fault.UnhandledError.Err()
	writeFaultResponse(w, fallbackFault)
	return nil
}

func writeFaultResponse(w http.ResponseWriter, f *fault.Fault) {

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code": f.Code,
			"args": f.Args,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}

func ToUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func ToNumeric(accuracy float64) (pgtype.Numeric, error) {
	var num pgtype.Numeric
	if err := num.Scan(fmt.Sprintf("%f", accuracy)); err != nil {
		return pgtype.Numeric{}, err
	}
	return num, nil
}
