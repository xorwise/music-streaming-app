package domain

import "context"

type UserUpdateAvatarUsecase interface {
	Update(ctx context.Context, user *User) error
	SaveFile(ctx context.Context, fileData string, filename string) (string, error)
}

type UserUpdateAvatarRequest struct {
	Data     string `json:"data"`
	Filename string `json:"filename"`
}
