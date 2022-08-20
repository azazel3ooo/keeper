package server_repo

import (
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/gofiber/fiber/v2"
)

const (
	SetOperation    = "set"
	UpdateOperation = "upd"
	DeleteOperation = "del"
)

var ProcessingOperations = map[string]int{
	SetOperation:    1,
	UpdateOperation: 2,
	DeleteOperation: 3,
}

type Server struct {
	storage        models.Storable4Server
	cfg            models.Config
	app            *fiber.App
	lg             any
	lgChan         LogChan
	processingChan ProcessingChan
}

// ProcessingTuple содержит необходимые данные для выполнения операции над данными пользователя
type ProcessingTuple struct {
	Operation int
	Data      models.Validatable
	User      string
}

type LogChan chan any // unused
type ProcessingChan chan ProcessingTuple
