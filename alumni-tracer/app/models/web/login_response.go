package web

import "WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"

type LoginResponse struct {
	User domain.Users `json:"user"`
	Token string `json:"token"`
}