package dto

type CreateUserSessionRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,oneof=active completed"`
}

type UpdateUserSessionRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,oneof=active completed"`
}

type GetUserSessionRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type ListUserSessionsRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Limit  int32  `json:"limit" validate:"gte=0,lte=100"`
	Offset int32  `json:"offset" validate:"gte=0"`
}

type ListActiveUserSessionsRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type DeleteUserSessionRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}
