package client_repo

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/azazel3ooo/keeper/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type ClientStorage struct {
	d *sql.DB
}

func (c *ClientStorage) Init(path string) error {
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

	c.d = db
	err = c.CreateTables()
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientStorage) CreateTables() error {
	stmt := `CREATE TABLE if not exists storage (
    	"id" TEXT PRIMARY key,
    	"data" TEXT,
    	"comment" TEXT
	);`

	_, err := c.d.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientStorage) Set(r models.UserData) error {
	stmt := `insert or replace into storage (id,"data",comment) values($1,$2,$3);`

	_, err := c.d.Exec(stmt, r.ID, r.Data, r.Comment)
	return err
}

func (c *ClientStorage) GetAll() ([]models.UserData, error) {
	stmt := `select id, "data", comment from storage`

	rows, err := c.d.Query(stmt)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var (
		tmp models.UserData
		res []models.UserData
	)
	for rows.Next() {
		err = rows.Scan(&tmp.ID, &tmp.Data, &tmp.Comment)
		if err != nil {
			log.Println(err)
			continue
		}

		res = append(res, tmp)
	}

	return res, nil
}

func (c *ClientStorage) Update(r models.UserData) error {
	stmt := `insert or replace into storage (id,"data",comment) values($1,$2,$3);`

	_, err := c.d.Exec(stmt, r.ID, r.Data, r.Comment)
	return err
}

func (c *ClientStorage) Delete(r models.DeleteRequest) error {
	stmt := `delete from storage where id=$1`

	_, err := c.d.Exec(stmt, r.ID)
	return err
}
