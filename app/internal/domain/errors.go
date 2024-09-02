package domain

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFieldRequired      = errors.New("field is required")
	ErrRoomNotFound       = errors.New("room not found")
	ErrNotUserInRoom      = errors.New("user is not in the room")
)

type ErrorResponse struct {
	Error string `json:"error"`
}
