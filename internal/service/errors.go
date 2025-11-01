package service

import (
	"test-http/pkg/fault"
)

const (
	BadRequest               fault.Code = "BAD_REQUEST"
	ValidationFailed         fault.Code = "VALIDATION_FAILED"
	NotFound                 fault.Code = "NOT_FOUND"
	PermissionDenied         fault.Code = "PERMISSION_DENIED"
	Conflict                 fault.Code = "CONFLICT"
	InfrastructureUnexpected fault.Code = "INFRASTRUCTURE_UNEXPECTED"

)
