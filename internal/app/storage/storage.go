package storage

import (
	"database/sql"
	_ "github.com/lib/pq"

	//"gorm.io/driver/postgres"
	//"gorm.io/gorm"
)

type Storage struct {
	config *Config
	db *sql.DB
	EventRepository *EventRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

// Open connection and ping db
func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// Close connection
func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Event() *EventRepository {
	if s.EventRepository != nil {
		return s.EventRepository
	}

	s.EventRepository = &EventRepository{
		store: s,
	}

	return s.EventRepository
}