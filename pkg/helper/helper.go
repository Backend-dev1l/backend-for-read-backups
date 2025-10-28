package helper

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"test-http/errors_pkg"
)

// HTTPError преобразует ошибку в соответствующий HTTP ответ
func HTTPError(ctx context.Context, w http.ResponseWriter, err error) {
	lang, _ := ctx.Value(localeKey{}).(string)
	var f *fault.Fault

	if errors.As(err, &f) {
		writeFaultResponse(w, f.ToProto(lang))
		return
	}

	// fallback ошибка
	fallbackFault := errorsPkg.InfrastructureUnexpected.Err()
	writeFaultResponse(w, fallbackFault.ToProto(lang))
}

// writeFaultResponse записывает ответ на основе fault
func writeFaultResponse(w http.ResponseWriter, faultProto *fault.Proto) {
	statusCode := faultToHTTPStatus(faultProto)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    faultProto.Code,
			"message": faultProto.Message,
			"details": faultProto.Details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// faultToHTTPStatus преобразует код ошибки в HTTP статус
func faultToHTTPStatus(faultProto *fault.Proto) int {
	switch faultProto.Code {
	case string(errors_pkg.ContextUserIDMissing):
		return http.StatusUnauthorized
	case string(errors_pkg.ValidationFailed):
		return http.StatusBadRequest
	case string(errors_pkg.JSONDecodeFailed):
		return http.StatusBadRequest
	case string(errors_pkg.InfrastructureUnexpected):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
