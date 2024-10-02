package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomUtils struct{}

func NewRoomUtils() domain.RoomUtils {
	return &roomUtils{}
}

func (ru *roomUtils) GenerateRoomCode(roomID int64) string {
	h := fmt.Sprintf("%x", roomID)
	hash := sha256.Sum256([]byte(h))

	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
