package storage

import (
	"database/sql"
	"fmt"
	"github.com/eltaljohn/go-db/pkg/invoice"
	"github.com/eltaljohn/go-db/pkg/invoiceheader"
	"github.com/eltaljohn/go-db/pkg/invoiceitem"
)

// MySQLInvoice is used to work with MySQL - invoice
type MySQLInvoice struct {
	db            *sql.DB
	storageHeader invoiceheader.Storage
	storageItems  invoiceitem.Storage
}

// MySQLInvoice returns a new pointer of MySQLInvoice
func NewMySQLInvoice(db *sql.DB, h invoiceheader.Storage, i invoiceitem.Storage) *MySQLInvoice {
	return &MySQLInvoice{db, h, i}
}

// Create implements interface invoice.Storage
func (p *MySQLInvoice) Create(m *invoice.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	err = p.storageHeader.CreateTx(tx, m.Header)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Printf("Factura creada con id: %d \n", m.Header.ID)

	if err := p.storageItems.CreateTx(tx, m.Header.ID, m.Items); err != nil {
		tx.Rollback()
		return err
	}
	fmt.Printf("Items creados: %d \n", len(m.Items))

	return tx.Commit()
}
