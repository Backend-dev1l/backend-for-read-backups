package dto

type CreateUserProgressRequest struct {
	UserID         string `json:"user_id" validate:"required,uuid"`
	WordID         string `json:"word_id" validate:"required,uuid"`
	CorrectCount   int32  `json:"correct_count" validate:"gte=0"`
	IncorrectCount int32  `json:"incorrect_count" validate:"gte=0"`
}

type UpdateUserProgressRequest struct {
	ID             string `json:"id" validate:"required,uuid"`
	CorrectCount   int32  `json:"correct_count" validate:"gte=0"`
	IncorrectCount int32  `json:"incorrect_count" validate:"gte=0"`
}

type GetUserProgressRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type GetUserProgressByUserAndWordRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	WordID string `json:"word_id" validate:"required,uuid"`
}

type ListUserProgressRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Limit  int32  `json:"limit" validate:"gte=0,lte=100"`
	Offset int32  `json:"offset" validate:"gte=0"`
}

type DeleteUserProgressRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}
