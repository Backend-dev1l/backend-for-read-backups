package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	errorsPkg "test-http/pkg/errors_pkg"
	"test-http/pkg/fault"
)

// HTTPError преобразует ошибку в корректный HTTP-ответ.
func HTTPError(w http.ResponseWriter, err error) {
	var f *fault.Fault

	// Если это наша бизнес-ошибка
	if errors.As(err, &f) {
		writeFaultResponse(w, f)
		return
	}

	// Иначе возвращаем стандартную системную ошибку
	fallbackFault := errorsPkg.InfrastructureUnexpected.Err()
	writeFaultResponse(w, fallbackFault)
}

// writeFaultResponse формирует JSON-ответ на основе fault.Fault.
func writeFaultResponse(w http.ResponseWriter, f *fault.Fault) {
	statusCode := faultToHTTPStatus(f)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    f.Code,
			"message": f.Message,
			"args":    f.Args,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// faultToHTTPStatus определяет HTTP-статус на основе кода ошибки.
func faultToHTTPStatus(f *fault.Fault) int {
	switch f.Code {
	case string(errorsPkg.ContextGettingUserMissing):
		return http.StatusUnauthorized
	case string(errorsPkg.ValidationError):
		return http.StatusBadRequest
	case string(errorsPkg.DecodeFailed):
		return http.StatusBadRequest
	case string(errorsPkg.InfrastructureUnexpected):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
