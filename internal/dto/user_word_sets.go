package dto

type CreateUserWordSetRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	WordSetID string `json:"word_set_id" validate:"required,uuid"`
}

type GetUserWordSetRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type ListUserWordSetsRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type DeleteUserWordSetRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}
