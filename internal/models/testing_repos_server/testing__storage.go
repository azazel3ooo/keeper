package testing_repos_server

import (
	"errors"
	"github.com/azazel3ooo/keeper/internal/models"
)

type TestUser struct {
	Log string
	Pas string
}

type TestExample struct {
	User    string
	Data    string
	Comment string
}

type TestUsers map[string]TestUser
type TestData map[string]TestExample

type TestingServerStorage struct {
	users TestUsers
	data  TestData
}

func (t *TestingServerStorage) Init() {
	t.users = make(TestUsers)
	t.data = make(TestData)
}

func (t TestingServerStorage) CreateUser(log, pas string) (string, error) {
	for k, _ := range t.users {
		if log == t.users[k].Log {
			return "", models.ErrUserConflict
		}
	}

	id := models.GenerateUserID()
	t.users[id] = TestUser{Log: log, Pas: pas}

	return id, nil
}

func (t TestingServerStorage) CheckUser(login string) (id string, pass string, err error) {
	for k, _ := range t.users {
		if login == t.users[k].Log {
			return k, t.users[k].Pas, nil
		}
	}

	return "", "", nil
}

func (t TestingServerStorage) SetData(req models.UserData, user string) error {
	t.data[req.ID] = TestExample{
		User:    user,
		Data:    req.Data,
		Comment: req.Comment,
	}

	return nil
}

func (t TestingServerStorage) GetData(user string) ([]models.UserData, error) {
	var res []models.UserData

	for k, v := range t.data {
		if v.User == user {
			tmp := models.UserData{
				ID:      k,
				Data:    v.Data,
				Comment: v.Comment,
			}
			res = append(res, tmp)
		}
	}

	return res, nil
}

func (t TestingServerStorage) Delete(req models.DeleteRequest, user string) error {
	v, ok := t.data[req.ID]
	if !ok {
		return errors.New("unknown ID")
	}

	if v.User == user {
		delete(t.data, req.ID)
	}
	return nil
}

func (t TestingServerStorage) Update(req models.UserData, user string) error {
	t.data[req.ID] = TestExample{
		User:    user,
		Data:    req.Data,
		Comment: req.Comment,
	}
	return nil
}
