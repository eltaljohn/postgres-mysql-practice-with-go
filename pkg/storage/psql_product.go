package storage

import (
	"database/sql"
	"fmt"
	"github.com/eltaljohn/go-db/pkg/product"
)

// psqlMigrateProduct cons to create products table
const (
	psqlMigrateProduct = `
	CREATE TABLE IF NOT EXISTS products(
	id SERIAL NOT NULL,
	name VARCHAR(25) NOT NULL,
	observation VARCHAR(100),
	price INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP
	CONSTRAINT products_id_pk PRIMARY KEY (id)
	)`
	psqlCreateProduct = `INSERT INTO products(name, observation, price, created_at) 
	VALUES ($1, $2, $3, $4) RETURNING id`
	psqlGetAllProduct  = `SELECT id, name, observation, price, created_at, updated_at from products`
	psqlGetProductByID = psqlGetAllProduct + " WHERE id = $1"
	psqlUpdateProduct  = `UPDATE products SET name = $1, observation = $2, price = $3, updated_at = $4 WHERE id = $5`
	psqlDeleteProduct  = "DELETE FROM products WHERE id = $1"
)

// psqlProduct used to work with postgres - product
type psqlProduct struct {
	db *sql.DB
}

// newPsqlProduct returns a new pointer of psqlProduct
func newPsqlProduct(db *sql.DB) *psqlProduct {
	return &psqlProduct{db}
}

// Migrate implements interface product.storage
func (p *psqlProduct) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil
	}

	fmt.Println("Migración de producto ejecutada correctamente")
	return nil
}

// Create implements interface product.storage
func (p *psqlProduct) Create(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlCreateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.Name,
		stringToNull(m.Observations),
		m.Price,
		m.CreatedAt,
	).Scan(&m.ID)
	if err != nil {
		return err
	}

	fmt.Println("Se creó producto correctamente")
	return nil
}

// GetAll implements interface product.storage
func (p *psqlProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(psqlGetAllProduct)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(product.Models, 0)
	for rows.Next() {
		m, err := scanRowProduct(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

// GetByID implements interface product.storage
func (p *psqlProduct) GetByID(id uint) (*product.Model, error) {
	stmt, err := p.db.Prepare(psqlGetProductByID)
	if err != nil {
		return &product.Model{}, err
	}
	defer stmt.Close()

	return scanRowProduct(stmt.QueryRow(id))
}

// Update implements interface product.storage
func (p *psqlProduct) Update(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlUpdateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		m.Name,
		stringToNull(m.Observations),
		m.Price,
		timeToNull(m.UpdatedAt),
		m.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no existe el producto con id: %d", m.ID)
	}

	fmt.Println("Se actualizó el producto correctamente")
	return nil
}

// Delete implements interface product.storage
func (p *psqlProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(psqlDeleteProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no existe el producto con id: %d", id)
	}

	fmt.Println("Se eliminó el producto correctamente")
	return nil
}
