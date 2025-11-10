package dto

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type GetUserByIDRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type ListUsersRequest struct {
	Limit  int32 `json:"limit" validate:"gte=0,lte=100"`
	Offset int32 `json:"offset" validate:"gte=0"`
}
