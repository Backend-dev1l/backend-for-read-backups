package dto

import "github.com/jackc/pgx/v5/pgtype"

type CreateUserSessionRequest struct {
	UserID pgtype.UUID `json:"user_id" validate:"required,uuid"`
	Status string      `json:"status" validate:"required,oneof=active completed"`
}

type UpdateUserSessionRequest struct {
	ID      pgtype.UUID        `json:"id" validate:"required,uuid"`
	Status  string             `json:"status" validate:"required,oneof=active completed"`
	EndedAt pgtype.Timestamptz `json:"ended_at"`
}

type GetUserSessionRequest struct {
	ID pgtype.UUID `json:"id" validate:"required,uuid"`
}

type ListUserSessionsRequest struct {
	UserID pgtype.UUID `json:"user_id" validate:"required,uuid"`
	Limit  int32       `json:"limit" validate:"gte=0,lte=100"`
	Offset int32       `json:"offset" validate:"gte=0"`
}

type ListActiveUserSessionsRequest struct {
	UserID pgtype.UUID `json:"user_id" validate:"required,uuid"`
}

type DeleteUserSessionRequest struct {
	ID pgtype.UUID `json:"id" validate:"required,uuid"`
}
