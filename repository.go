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

	rep := &Repository{
		path: path,
		db:   db,
	}

	err = rep.Setup()
	return rep, err

}

func (repo *Repository) Setup() error {
	query := `
	PRAGMA journal_mode = WAL;
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_size_limit = 67108864; -- 64 megabytes
	PRAGMA mmap_size = 134217728; -- 128 megabytes
	PRAGMA cache_size = 2000;
	PRAGMA busy_timeout = 5000;
	`
	_, err := repo.db.Exec(query)
	if err != nil {
		return err
	}
	schema := `CREATE TABLE IF NOT EXISTS views(path text, ip text, agent text, created text)`
	_, err = repo.db.Exec(schema)
	return err
}
func (repo Repository) Close() {
	repo.db.Close()
}
func (repo *Repository) SaveView(v View) error {
	query := `
	INSERT INTO views(path, ip, agent, created) VALUES (?,?,?,?)
	`
	_, err := repo.db.Exec(query, v.Path, v.IP, v.Agent, v.Created.UTC().String())
	if err != nil {
		return err
	}
	return nil
}
