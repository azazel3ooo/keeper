package testing_repos_client

import (
	"github.com/azazel3ooo/keeper/internal/models/server_repo"
	"net/http"
)

// TestingClient имитация http.Client для тестирования клиента
type TestingClient struct {
	S server_repo.Server
}

// Do выполняет тестовый запрос на сервере
func (c TestingClient) Do(req *http.Request) (*http.Response, error) {
	return c.S.Test(req, -1)
}
