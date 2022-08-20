package client_repo

import (
	server_logic "github.com/azazel3ooo/keeper/internal/logic/server"
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/azazel3ooo/keeper/internal/models/server_repo"
	"github.com/azazel3ooo/keeper/internal/models/testing_repos_client"
	"github.com/azazel3ooo/keeper/internal/models/testing_repos_server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_ActionToServer(t *testing.T) {
	var store testing_repos_server.TestingServerStorage
	store.Init()
	procChan := make(server_repo.ProcessingChan, 100) // чтобы не считывать
	defer close(procChan)
	s := server_repo.NewServer(server_repo.WithStorage(store), server_repo.WithProcessingChan(procChan))
	s.SetupApp()

	c := NewClient(WithClient(testing_repos_client.TestingClient{S: *s}))

	testToken, _ := server_logic.GenerateToken("tmp", 5.0)
	expiredToken, _ := server_logic.GenerateToken("tmp", 0.0)
	tests := []struct {
		description string
		req         models.Validatable
		token       string
		route       string
		method      string
		expectedErr error
	}{
		{
			description: "success",
			req:         models.DeleteRequest{ID: "1"},
			method:      http.MethodDelete,
			token:       testToken,
			route:       "/api/v1/items",
			expectedErr: nil,
		},
		{
			description: "unauthorized",
			req:         models.DeleteRequest{ID: "1"},
			token:       expiredToken,
			method:      http.MethodDelete,
			route:       "/api/v1/items",
			expectedErr: models.ErrExpiredToken,
		},
		{
			description: "forbidden",
			req:         models.DeleteRequest{ID: "1"},
			token:       testToken[:len(testToken)-2],
			method:      http.MethodDelete,
			route:       "/api/v1/items",
			expectedErr: models.ErrForbidden,
		},
		{
			description: "bad request",
			req:         models.DeleteRequest{ID: ""},
			token:       testToken,
			method:      http.MethodDelete,
			route:       "/api/v1/items",
			expectedErr: models.ErrBadRequest,
		},
	}
	for _, tt := range tests {
		c.token = tt.token
		err := c.ActionToServer(tt.req, tt.route, tt.method)
		assert.Equalf(t, tt.expectedErr, err, tt.description)
	}
}

func TestClient_GetActualData(t *testing.T) {
	var store testing_repos_server.TestingServerStorage
	store.Init()
	s := server_repo.NewServer(server_repo.WithStorage(store))
	s.SetupApp()

	c := NewClient(WithClient(testing_repos_client.TestingClient{S: *s}))

	route := "/api/v1/items"
	uid := "tmp"
	testToken, _ := server_logic.GenerateToken(uid, 5.0)
	expiredToken, _ := server_logic.GenerateToken(uid, 0.0)
	store.SetData(models.UserData{
		Data:    "tmp_data",
		Comment: "tmp_metadata",
		ID:      "data_ID",
	}, uid)

	tests := []struct {
		description string
		token       string
		route       string
		expectedErr error
	}{
		{
			description: "success",
			token:       testToken,
			route:       route,
			expectedErr: nil,
		},
		{
			description: "unauthorized",
			token:       expiredToken,
			route:       route,
			expectedErr: models.ErrExpiredToken,
		},
		{
			description: "forbidden",
			token:       testToken[:len(testToken)-2],
			route:       route,
			expectedErr: models.ErrForbidden,
		},
	}
	for _, tt := range tests {
		c.token = tt.token
		_, err := c.GetActualData()
		assert.Equalf(t, tt.expectedErr, err, tt.description)
	}
}

func TestClient_GetToken(t *testing.T) {
	var store testing_repos_server.TestingServerStorage
	store.Init()
	s := server_repo.NewServer(server_repo.WithStorage(store))
	s.SetupApp()

	c := NewClient(WithClient(testing_repos_client.TestingClient{S: *s}))

	authRoute := "/api/v1/auth"
	regRoute := "/api/v1/registration"

	tests := []struct {
		description string
		req         models.UserRequest
		route       string
		expectedErr error
	}{
		{
			description: "success",
			req:         models.UserRequest{Login: "q", Password: "q"},
			route:       regRoute,
			expectedErr: nil,
		},
		{
			description: "conflict", // for registration
			req:         models.UserRequest{Login: "q", Password: "q"},
			route:       regRoute,
			expectedErr: models.ErrUserRegistrationConflict,
		},
		{
			description: "forbidden", // for authorization
			req:         models.UserRequest{Login: "q", Password: "qwerty"},
			route:       authRoute,
			expectedErr: models.ErrForbidden,
		},
		{
			description: "bad request",
			req:         models.UserRequest{Login: "", Password: ""},
			route:       regRoute,
			expectedErr: models.ErrBadRequest,
		},
	}
	for _, tt := range tests {
		err := c.GetToken(tt.req, tt.route)
		assert.Equalf(t, tt.expectedErr, err, tt.description)
		if tt.expectedErr == nil {
			assert.Equalf(t, true, c.token != "", tt.description)
		}
	}
}
