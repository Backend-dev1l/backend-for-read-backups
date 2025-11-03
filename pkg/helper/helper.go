package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	errorsPkg "test-http/pkg/errors_pkg"
	"test-http/pkg/fault"
)

func HTTPError(w http.ResponseWriter, err error) error {
	var f *fault.Fault

	if errors.As(err, &f) {
		writeFaultResponse(w, f)
		return nil
	}

	fallbackFault := errorsPkg.InfrastructureUnexpected.Err()
	writeFaultResponse(w, fallbackFault)
	return nil
}

// writeFaultResponse формирует JSON-ответ на основе fault.Fault.
func writeFaultResponse(w http.ResponseWriter, f *fault.Fault) {

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    f.Code,
			"message": f.Message,
			"args":    f.Args,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}
