package web

import "time"

type UserResponse struct{
	Id int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}