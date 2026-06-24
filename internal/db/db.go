package db

import (
	"database/sql"
	"log"
	"sync"

	"github.com/go-sql-driver/mysql"
)

type mySQLStore struct {
	db *sql.DB
	mu sync.Mutex
}

func NewMySQLStorage(cfg mysql.Config) *mySQLStore {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MySQL")
	return &mySQLStore{
		db: db,
	}
}

func (s *mySQLStore) GetDB() *sql.DB {
	return s.db
}

func (s *mySQLStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.db != nil {
		log.Println("Closing MySQL connection")
		return s.db.Close()
	}
	return nil
}