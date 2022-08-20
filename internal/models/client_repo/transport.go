package client_repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/azazel3ooo/keeper/internal/models"
	"io"
	"log"
	"net/http"
)

// GetToken производит авторизацию\регистрацию пользователя и обновляет токен клиента
func (c *Client) GetToken(r models.UserRequest, addr string) error {
	s, err := json.Marshal(r)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(s)

	req, err := http.NewRequest(http.MethodPost, addr, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cl.Do(req)
	if err != nil {
		log.Println(err)
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var res models.UserResponse

		err = json.Unmarshal(b, &res)
		if err != nil {
			return err
		}
		c.UpdateToken(res.Token)

	case http.StatusForbidden:
		log.Println(string(b))
		return models.ErrForbidden

	case http.StatusBadRequest:
		log.Println(string(b))
		return models.ErrBadRequest

	case http.StatusConflict:
		log.Println(string(b))
		return models.ErrUserRegistrationConflict

	case http.StatusInternalServerError:
		log.Println(string(b))
		return models.ErrInternalServerError

	default:
		return errors.New("unexpected status code " + resp.Status)
	}

	return nil
}

// ActionToServer отправляет запрос с необходимым действием на сервер
func (c Client) ActionToServer(r models.Validatable, addr, method string) error {
	s, err := json.Marshal(r)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(s)

	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	resp, err := c.cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusForbidden:
		return models.ErrForbidden

	case http.StatusUnauthorized:
		return models.ErrExpiredToken

	case http.StatusBadRequest:
		return models.ErrBadRequest

	case http.StatusOK:
		return nil
	}

	return nil
}

// GetActualData получает все записи клиента с сервера
func (c Client) GetActualData() ([]models.UserData, error) {
	req, err := http.NewRequest(http.MethodGet, c.ActionAddr(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)

	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var res models.UserDataResponse
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &res)
		if err != nil {
			return nil, err
		}
		return res.Data, nil

	case http.StatusForbidden:
		return nil, models.ErrForbidden

	case http.StatusUnauthorized:
		return nil, models.ErrExpiredToken

	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError

	default:
		return nil, errors.New("unknown status " + resp.Status)
	}
}
