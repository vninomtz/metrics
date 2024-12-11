package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	path string
	db   *sql.DB
}

func NewRepository(path string) (*Repository, error) {
	var db *sql.DB
	db, err := sql.Open("sqlite3", path)
	if err != nil {
	}

	return &Repository{
		path: path,
		db:   db,
	}, nil
}
func (repo Repository) Close() {
	repo.db.Close()
}
func (repo Repository) SaveView(m map[string]string) error {
	query := `
	INSERT INTO views(page, user_agent, visited) VALUES (?,?,?)
	`
	_, err := repo.db.Exec(query, m["page"], m["agent"], m["visited"])
	if err != nil {
		return err
	}
	return nil
}
