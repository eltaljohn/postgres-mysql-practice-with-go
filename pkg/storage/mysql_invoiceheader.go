package storage

import (
	"database/sql"
	"fmt"
	"github.com/eltaljohn/go-db/pkg/invoiceheader"
)

// psqlMigrateInvoiceHeader cons to create invoice_headers table
const (
	mySQLMigrateInvoiceHeader = `CREATE TABLE IF NOT EXISTS invoice_headers(
	id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	client VARCHAR(25) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP
)`
	mySQLCreateInvoiceHeader = `INSERT INTO invoice_headers(client) VALUES (?)`
)

// MYSQLInvoiceHeader used to work with MySQL - invoice_headers
type MYSQLInvoiceHeader struct {
	db *sql.DB
}

// NewMYSQLInvoiceHeader returns a new pointer of MYSQLInvoiceHeader
func NewMYSQLInvoiceHeader(db *sql.DB) *MYSQLInvoiceHeader {
	return &MYSQLInvoiceHeader{db}
}

// Migrate implements interface invoiceHeader.storage
func (p *MYSQLInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(mySQLMigrateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil
	}

	fmt.Println("Migración de InvoiceHeader ejecutada correctamente")
	return nil
}

func (p *MYSQLInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(mySQLCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(m.Client)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	m.ID = uint(id)
	return nil
}
