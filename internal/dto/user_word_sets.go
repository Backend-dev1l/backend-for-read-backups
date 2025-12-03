package dto

import "github.com/jackc/pgtype"

type CreateUserWordSetRequest struct {
	UserID    pgtype.UUID `json:"user_id" validate:"required,uuid"`
	WordSetID pgtype.UUID `json:"word_set_id" validate:"required,uuid"`
}

type GetUserWordSetRequest struct {
	ID pgtype.UUID `json:"id" validate:"required,uuid"`
}

type ListUserWordSetsRequest struct {
	UserID pgtype.UUID `json:"user_id" validate:"required,uuid"`
	Limit  int32       `json:"limit" validate:"gte=0,lte=100"`
	Offset int32       `json:"offset" validate:"gte=0"`
}

type UpdateUserWordSetRequest struct {
	ID        pgtype.UUID `json:"id" validate:"required,uuid"`
	WordSetID pgtype.UUID `json:"word_set_id" validate:"required,uuid"`
}

type DeleteUserWordSetRequest struct {
	ID pgtype.UUID `json:"id" validate:"required,uuid"`
}
