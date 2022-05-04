package invoice

import (
	"github.com/eltaljohn/go-db/pkg/invoiceheader"
	"github.com/eltaljohn/go-db/pkg/invoiceitem"
)

// Model of invoice
type Model struct {
	Header *invoiceheader.Model
	Items  invoiceitem.Models
}

// Storage interface that must implement a db storage
type Storage interface {
	Create(*Model) error
}

// Service of invoice
type Service struct {
	storage Storage
}

// NewService returns a service pointer
func NewService(s Storage) *Service {
	return &Service{s}
}

// Create creates a new invoice
func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}
