package client_repo

import (
	"github.com/azazel3ooo/keeper/internal/models"
)

// Client структура клиента для отправки запросов к серверу и работы приложения клиента
type Client struct {
	cl    models.ClientHttpInterface // for testing
	store models.ClientStorable
	cfg   models.Config
	token string
}
