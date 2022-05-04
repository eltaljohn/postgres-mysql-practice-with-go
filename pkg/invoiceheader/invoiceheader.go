package invoiceheader

import (
	"database/sql"
	"time"
)

// Model of invoiceheader
type Model struct {
	ID        uint
	Client    string
	CreateAt  time.Time
	UpdatedAt time.Time
}

type Storage interface {
	Migrate() error
	CreateTx(tx *sql.Tx, model *Model) error
}

// Service of invoiceheader
type Service struct {
	storage Storage
}

// NewService returns a pointer of Service
func NewService(s Storage) *Service {
	return &Service{s}
}

// Migrate is used to migrate invoiceheader
func (s Service) Migrate() error {
	return s.storage.Migrate()
}
