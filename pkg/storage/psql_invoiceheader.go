package storage

import (
	"database/sql"
	"fmt"
	"github.com/eltaljohn/go-db/pkg/invoiceheader"
)

// psqlMigrateInvoiceHeader cons to create invoice_headers table
const (
	psqlMigrateInvoiceHeader = `CREATE TABLE IF NOT EXISTS invoice_headers(
	id SERIAL NOT NULL,
	client VARCHAR(25) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP,
	CONSTRAINT invoice_headers_id_pk PRIMARY KEY (id)
)`
	psqlCreateInvoiceHeader = `INSERT INTO invoice_headers(client) VALUES ($1) RETURNING id, created_at`
)

// PsqlInvoiceHeader used to work with postgres - invoice_headers
type PsqlInvoiceHeader struct {
	db *sql.DB
}

// NewPsqlInvoiceHeader returns a new pointer of PsqlInvoiceHeader
func NewPsqlInvoiceHeader(db *sql.DB) *PsqlInvoiceHeader {
	return &PsqlInvoiceHeader{db}
}

// Migrate implements interface invoiceHeader.storage
func (p *PsqlInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateInvoiceHeader)
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

func (p *PsqlInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(psqlCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(m.Client).Scan(&m.ID, &m.CreateAt)
}
