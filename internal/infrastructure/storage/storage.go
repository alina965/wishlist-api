package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbURL string) (*Storage, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("ping error: %w", err)
	}

	return &Storage{db: db}, nil
}

func (storage *Storage) Close() error {
	return storage.db.Close()
}

func (storage *Storage) GetDB() *sql.DB {
	return storage.db
}
