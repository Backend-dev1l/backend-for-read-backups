package errors_pkg

import "test-http/pkg/fault"

const (
	ValidationError                         fault.Code = "VALIDATION_ERROR"
	DecodeFailed                            fault.Code = "DECODE_FAILED"
	ContextCreatingUserMissing              fault.Code = "CONTEXT_CREATING_USER_MISSING"
	ContextCreatingUserStatisticsMissing    fault.Code = "CONTEXT_CREATING_USER_STATISTICS_MISSING"
	ContextGettingUserStatisticsMissing     fault.Code = "CONTEXT_GETTING_USER_STATISTICS_MISSING"
	ContextUpdatingUserStatisticsMissing    fault.Code = "CONTEXT_UPDATING_USER_STATISTICS_MISSING"
	ContextDeletingUserStatisticsMissing    fault.Code = "CONTEXT_DELETING_USER_STATISTICS_MISSING"
	UUIDParsingFailed                       fault.Code = "UUID_PARSING_FAILED"
	ContextGettingUserMissing               fault.Code = "CONTEXT_GETTING_USER_MISSING"
)
