package server_repo

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	logic "github.com/azazel3ooo/keeper/internal/logic/server"
	"github.com/azazel3ooo/keeper/internal/models/testing_repos_server"
	"github.com/stretchr/testify/assert"
)

func TestServer_registration(t *testing.T) {
	var store testing_repos_server.TestingServerStorage
	store.Init()
	s := NewServer(WithStorage(store))
	s.SetupApp()

	tests := []struct {
		description  string
		req          string
		expectedCode int
	}{
		{
			description:  "success",
			expectedCode: http.StatusOK,
			req:          "{\"login\":\"q\",\"password\":\"q\"}",
		},
		{
			description:  "bad request",
			expectedCode: http.StatusBadRequest,
			req:          "{",
		},

		// копия 1 запроса. Должна быть ошибка, поскольку создали ранее
		{
			description:  "conflict",
			expectedCode: http.StatusConflict,
			req:          "{\"login\":\"q\",\"password\":\"q\"}",
		},
	}
	for _, tt := range tests {
		b := bytes.NewBuffer([]byte(tt.req))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/registration", b)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req, -1)
		if err != nil {
			log.Println(err)
			continue
		}
		assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func TestServer_authorization(t *testing.T) {
	var store testing_repos_server.TestingServerStorage
	store.Init()
	s := NewServer(WithStorage(store))
	s.SetupApp()

	store.CreateUser("q", "q")
	tests := []struct {
		description  string
		req          string
		expectedCode int
	}{
		{
			description:  "success",
			expectedCode: http.StatusOK,
			req:          "{\"login\":\"q\",\"password\":\"q\"}",
		},
		{
			description:  "bad request",
			expectedCode: http.StatusBadRequest,
			req:          "{",
		},
		{
			description:  "forbidden",
			expectedCode: http.StatusForbidden,
			req:          "{\"login\":\"q\",\"password\":\"qqqqq\"}",
		},
	}
	for _, tt := range tests {
		b := bytes.NewBuffer([]byte(tt.req))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth", b)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req, -1)
		if err != nil {
			log.Println(err)
			continue
		}
		assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func TestServer_delete(t *testing.T) {
	procChan := make(ProcessingChan)
	s := NewServer(WithProcessingChan(procChan))
	s.SetupApp()

	userID := "user"
	testToken, _ := logic.GenerateToken(userID, 5.0)
	expiredToken, _ := logic.GenerateToken("user_2", 0.0)

	tests := []struct {
		description  string
		req          string
		token        string
		expectedCode int
	}{
		{
			description:  "success",
			expectedCode: http.StatusOK,
			token:        testToken,
			req:          "{\"id\":\"test_id\"}",
		},
		{
			description:  "bad request",
			expectedCode: http.StatusBadRequest,
			token:        testToken,
			req:          "{",
		},
		{
			description:  "forbidden",
			expectedCode: http.StatusForbidden,
			token:        testToken[:len(testToken)-2],
		},
		{
			description:  "expired token",
			expectedCode: http.StatusUnauthorized,
			token:        expiredToken,
		},
	}

	go func() {
		for el := range procChan {
			assert.Equalf(t, ProcessingOperations[DeleteOperation], el.Operation, "success")
		}
	}()

	for _, tt := range tests {
		b := bytes.NewBuffer([]byte(tt.req))
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items", b)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tt.token)

		resp, err := s.app.Test(req, -1)
		if err != nil {
			log.Println(err)
			continue
		}
		assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
	close(procChan)
}

func TestServer_getAll(t *testing.T) {
	var store testing_repos_server.TestingServerStorage
	store.Init()
	s := NewServer(WithStorage(store))
	s.SetupApp()

	userID := "user"
	testToken, _ := logic.GenerateToken(userID, 5.0)
	expiredToken, _ := logic.GenerateToken("user_2", 0.0)

	tests := []struct {
		description  string
		token        string
		expectedCode int
	}{
		{
			description:  "success",
			expectedCode: http.StatusOK,
			token:        testToken,
		},
		{
			description:  "forbidden",
			expectedCode: http.StatusForbidden,
			token:        testToken[:len(testToken)-2],
		},
		{
			description:  "expired token",
			expectedCode: http.StatusUnauthorized,
			token:        expiredToken,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
		req.Header.Set("Authorization", tt.token)

		resp, err := s.app.Test(req, -1)
		if err != nil {
			log.Println(err)
			continue
		}
		assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func TestServer_set(t *testing.T) {
	procChan := make(ProcessingChan)
	s := NewServer(WithProcessingChan(procChan))
	s.SetupApp()

	userID := "user"
	testToken, _ := logic.GenerateToken(userID, 5.0)
	expiredToken, _ := logic.GenerateToken("user_2", 0.0)

	tests := []struct {
		description  string
		req          string
		token        string
		expectedCode int
	}{
		{
			description:  "success",
			expectedCode: http.StatusOK,
			token:        testToken,
			req:          "{\"id\":\"test_id\",\"data\":\"test_data\"}",
		},
		{
			description:  "bad request",
			expectedCode: http.StatusBadRequest,
			token:        testToken,
			req:          "{",
		},
		{
			description:  "forbidden",
			expectedCode: http.StatusForbidden,
			token:        testToken[:len(testToken)-2],
		},
		{
			description:  "expired token",
			expectedCode: http.StatusUnauthorized,
			token:        expiredToken,
		},
	}

	go func() {
		for el := range procChan {
			assert.Equalf(t, ProcessingOperations[SetOperation], el.Operation, "success")
		}
	}()

	for _, tt := range tests {
		b := bytes.NewBuffer([]byte(tt.req))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", b)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tt.token)

		resp, err := s.app.Test(req, -1)
		if err != nil {
			log.Println(err)
			continue
		}
		assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
	close(procChan)
}

func TestServer_update(t *testing.T) {
	procChan := make(ProcessingChan)
	s := NewServer(WithProcessingChan(procChan))
	s.SetupApp()

	userID := "user"
	testToken, _ := logic.GenerateToken(userID, 5.0)
	expiredToken, _ := logic.GenerateToken("user_2", 0.0)

	tests := []struct {
		description  string
		req          string
		token        string
		expectedCode int
	}{
		{
			description:  "success",
			expectedCode: http.StatusOK,
			token:        testToken,
			req:          "{\"id\":\"test_id\",\"data\":\"test_data\"}",
		},
		{
			description:  "bad request",
			expectedCode: http.StatusBadRequest,
			token:        testToken,
			req:          "{",
		},
		{
			description:  "forbidden",
			expectedCode: http.StatusForbidden,
			token:        testToken[:len(testToken)-2],
		},
		{
			description:  "expired token",
			expectedCode: http.StatusUnauthorized,
			token:        expiredToken,
		},
	}

	go func() {
		for el := range procChan {
			assert.Equalf(t, ProcessingOperations[UpdateOperation], el.Operation, "success")
		}
	}()

	for _, tt := range tests {
		b := bytes.NewBuffer([]byte(tt.req))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/items", b)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tt.token)

		resp, err := s.app.Test(req, -1)
		if err != nil {
			log.Println(err)
			continue
		}
		assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
	close(procChan)
}
