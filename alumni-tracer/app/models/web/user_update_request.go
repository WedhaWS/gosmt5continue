package web

type UserUpdateRequest struct{
	Id int64 `json:"-" validate:"required"`
	Name string `json:"username" validate:"required,max=200,min=1"`
	Email string `json:"email" validate:"required"`
	Password_Hash string `json:"password" validate:"required"`
	Role string `json:"role" validate:"required"`
}