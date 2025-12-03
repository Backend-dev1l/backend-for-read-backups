package dto

import "github.com/jackc/pgtype"

type CreateUserProgressRequest struct {
	UserID         pgtype.UUID `json:"user_id" validate:"required,uuid"`
	WordID         pgtype.UUID `json:"word_id" validate:"required,uuid"`
	CorrectCount   int32       `json:"correct_count" validate:"gte=0"`
	IncorrectCount int32       `json:"incorrect_count" validate:"gte=0"`
}

type UpdateUserProgressRequest struct {
	ID             pgtype.UUID `json:"id" validate:"required,uuid"`
	CorrectCount   int32       `json:"correct_count" validate:"gte=0"`
	IncorrectCount int32       `json:"incorrect_count" validate:"gte=0"`
}

type GetUserProgressRequest struct {
	ID pgtype.UUID `json:"id" validate:"required,uuid"`
}

type GetUserProgressByUserAndWordRequest struct {
	UserID pgtype.UUID `json:"user_id" validate:"required,uuid"`
	WordID pgtype.UUID `json:"word_id" validate:"required,uuid"`
}

type ListUserProgressRequest struct {
	UserID pgtype.UUID `json:"user_id" validate:"required,uuid"`
	Limit  int32       `json:"limit" validate:"gte=0,lte=100"`
	Offset int32       `json:"offset" validate:"gte=0"`
}

type DeleteUserProgressRequest struct {
	ID pgtype.UUID `json:"id" validate:"required,uuid"`
}
