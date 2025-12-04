// Package fault предоставляет механизм типизированных ошибок для HTTP и других слоев приложения.
package fault

import (
	"errors"

	"encoding/json"

	"net/http"
)

// Arg — пара ключ-значение для параметризации ошибок.
type Arg struct {
	K string // ключ аргумента
	V string // значение аргумента
}

// Fault — структура бизнес-ошибки.
type Fault struct {
	Code string            `json:"code"`           // уникальный код ошибки
	Args map[string]string `json:"args,omitempty"` // дополнительные аргументы
}

// Error возвращает строковое представление ошибки (код).
func (f *Fault) Error() string {
	return f.Code
}

// Data возвращает аргументы ошибки.
func (f *Fault) Data() map[string]string {
	return f.Args
}

// Code — тип для объявления кодов ошибок.
type Code string

// Err создаёт новый экземпляр Fault с указанными аргументами.
func (c Code) Err(args ...*Arg) *Fault {
	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.K] = arg.V
	}
	return &Fault{
		Code: string(c),
		Args: argsMap,
	}
}

// New создает Fault с кастомным сообщением и аргументами.
func (c Code) New(args ...*Arg) *Fault {
	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.K] = arg.V
	}
	return &Fault{
		Code: string(c),
		Args: argsMap,
	}
}

// HandleErr преобразует ошибку в Fault (если это возможно).
// Если ошибка не Fault, возвращает UnhandledError.
func HandleErr(err error) *Fault {
	var f *Fault
	if errors.As(err, &f) {
		return f
	}
	return UnhandledError.Err()
}

// --- Базовые системные ошибки ---
const (
	UnhandledError Code = "UNHANDLED_ERROR"
)

func HTTPError(w http.ResponseWriter, err error) error {
	var f *Fault

	if errors.As(err, &f) {
		writeFaultResponse(w, f)
		return nil
	}

	// Use a generic unhandled fault as a fallback.
	fallbackFault := UnhandledError.Err()
	writeFaultResponse(w, fallbackFault)
	return nil
}

func writeFaultResponse(w http.ResponseWriter, f *Fault) {

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code": f.Code,
			"args": f.Args,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}
