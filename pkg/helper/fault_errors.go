package helper

import (
	"encoding/json"
	"errors"

	"net/http"

	"test-http/pkg/fault"
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
