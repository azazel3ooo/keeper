package models

import "github.com/google/uuid"

// GenerateUserID возвращает уникальный id пользователя
func GenerateUserID() string {
	return uuid.New().String()
}
