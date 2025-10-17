package web

type UserRequest struct{
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required"`
	Password_Hash string `json:"password" validate:"required"`
	Role string `json:"role" validate:"required"`
}

