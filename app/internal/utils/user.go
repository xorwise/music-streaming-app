package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUtils struct {
	tokenTTL  int
	jwtSecret string
}

func NewUserUtils(tokenTTL int, jwtSecret string) domain.UserUtils {
	return &userUtils{
		tokenTTL:  tokenTTL,
		jwtSecret: jwtSecret,
	}
}
func (u *userUtils) CreateAccessToken(ctx context.Context, user *domain.User) (string, error) {
	exp := time.Now().Add(time.Duration(u.tokenTTL) * time.Second)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = exp.Unix()

	tokenStr, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
func (u *userUtils) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *userUtils) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (u *userUtils) SaveFile(ctx context.Context, fileData string, filename string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return "", err
	}
	path, err := getUniquePath("media", filename)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filepath.Join("media", path), data, 0644)
	if err != nil {
		return "", err
	}
	return "media/" + path, nil
}

func getUniquePath(directory, filename string) (string, error) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	newFilename := filename
	counter := 1

	for {
		if _, err := os.Stat(filepath.Join(directory, newFilename)); os.IsNotExist(err) {
			return newFilename, nil
		}
		newFilename = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}
}
