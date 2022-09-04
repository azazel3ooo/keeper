package server_repo

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	"github.com/azazel3ooo/keeper/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type ServerStorage struct {
	db *sql.DB
	mu *sync.RWMutex
}

func (s *ServerStorage) Init(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), 0774)
		if err != nil {
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	s.mu = new(sync.RWMutex)
	s.db = db
	err = s.CreateTables()
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerStorage) CreateTables() error {
	stmt := `CREATE TABLE if not exists users (
		"id" TEXT primary key,
		"login" TEXT,
		"pass" TEXT
	);`

	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `CREATE TABLE if not exists storage (
    	"id" TEXT PRIMARY key,
    	"user" TEXT,
    	"data" TEXT,
    	"comment" TEXT
	);`

	_, err = s.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerStorage) CreateUser(login, pass string) (string, error) {
	stmt := `select COUNT(*) from users where login=$1`
	r, err := s.db.Query(stmt, login)
	if err != nil {
		return "", err
	}

	if r.Err() != nil {
		return "", r.Err()
	}
	if r.Next() {
		var c int
		err = r.Scan(&c)
		r.Close()
		if err != nil {
			return "", err
		}
		if c > 0 {
			return "", models.ErrUserConflict
		}
	}

	stmt = `insert into users (id, login, pass) values ($1,$2,$3);`
	id := models.GenerateUserID()

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err = s.db.Exec(stmt, id, login, pass)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *ServerStorage) CheckUser(login string) (id string, pass string, err error) {
	stmt := `select id, pass from users where login=$1`
	r, err := s.db.Query(stmt, login)
	if err != nil {
		return "", "", err
	}
	defer r.Close()

	if r.Err() != nil {
		return "", "", r.Err()
	}
	if r.Next() {
		err = r.Scan(&id, &pass)
	}

	return id, pass, err
}

func (s *ServerStorage) SetData(req models.UserData, user string) error {
	stmt := `insert into storage (id, user, data, comment) values ($1,$2,$3,$4);`

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(stmt, req.ID, user, req.Data, req.Comment)
	return err
}

func (s *ServerStorage) GetData(user string) ([]models.UserData, error) {
	stmt := `select id,data,comment from storage where user=$1;`
	r, err := s.db.Query(stmt, user)
	if err != nil {
		return nil, err
	}
	if r.Err() != nil {
		return nil, err
	}
	defer r.Close()

	var (
		data models.UserData
		res  []models.UserData
	)

	for r.Next() {
		err = r.Scan(&data.ID, &data.Data, &data.Comment)
		if err != nil {
			return nil, err
		}
		res = append(res, data)
	}

	return res, nil
}

func (s *ServerStorage) Delete(req models.DeleteRequest, user string) error {
	stmt := `delete from storage where id=$1 AND user=$2`

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(stmt, req.ID, user)
	return err
}

func (s *ServerStorage) Update(req models.UserData, user string) error {
	stmt := `replace into storage(id, user, data, comment) VALUES ($1,$2,$3,$4);`

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(stmt, req.ID, user, req.Data, req.Comment)
	return err
}
