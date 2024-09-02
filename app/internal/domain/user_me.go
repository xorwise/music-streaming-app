package domain

import "context"

type UserMeResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserMeUsecase interface {
	GetByID(ctx context.Context, id int64) (*User, error)
}
