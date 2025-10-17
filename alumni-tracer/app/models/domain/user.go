package domain

import "time"

type Users struct{
	Id int64
	Username string
	Email string
	Password_Hash string
	Role string
	CreatedAt time.Time
}
