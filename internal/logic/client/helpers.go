package client_logic

import (
	"fmt"

	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/google/uuid"
)

// GenerateID создает уникальный  Id
func GenerateID() string {
	return uuid.New().String()
}

// PrintData печатает переданные данные в формате "record_id | record_data with metadata: record_metadata\n"
func PrintData(data []models.UserData) {
	for _, el := range data {
		fmt.Printf("%s | %s with metadata: %s\n", el.ID, el.Data, el.Comment)
	}
	fmt.Println()
}
