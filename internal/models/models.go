package models

import (
	"errors"
	"net/http"
)

var (
	ErrForbidden                = errors.New("status forbidden")
	ErrBadRequest               = errors.New("status bad request(fix request body)")
	ErrUserRegistrationConflict = errors.New("user already exist")
	ErrInternalServerError      = errors.New("internal server error")
	ErrUncastable               = errors.New("can't cast")
)

var (
	ErrUserConflict     = errors.New("user already exist")
	ErrUserDataConflict = errors.New("invalid login or password")
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
)

// ClientHttpInterface для возможности подмены на тестовый клиент
type ClientHttpInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Validatable для унификации работы со структурами запросов
type Validatable interface {
	Valid() bool
}

type Storable4Server interface {
	Storable4Users
	Storable4Data
}

type Storable4Users interface {
	CreateUser(login, pass string) (string, error)
	CheckUser(login string) (string, string, error)
}

type Storable4Data interface {
	SetData(req UserData, user string) error
	GetData(user string) ([]UserData, error)
	Delete(req DeleteRequest, user string) error
	Update(req UserData, user string) error
}

type Config struct {
	HostAddr   string `yaml:"host"`
	DbLocation string `yaml:"db_location"`
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	Token string `json:"token"`
}

type UserData struct {
	ID      string `json:"id"`
	Data    string `json:"data"`
	Comment string `json:"metadata,omitempty"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type UserDataResponse struct {
	Data []UserData `json:"data"`
}

type ClientStorable interface {
	Set(r UserData) error
	GetAll() ([]UserData, error)
	Update(r UserData) error
	Delete(r DeleteRequest) error
}
