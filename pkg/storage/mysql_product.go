package storage

import (
	"database/sql"
	"fmt"
	"github.com/eltaljohn/go-db/pkg/product"
)

const (
	mySQLMigrateProduct = `
	CREATE TABLE IF NOT EXISTS products(
	id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	name VARCHAR(25) NOT NULL,
	observation VARCHAR(100),
	price INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP
	)`
	mySQLCreateProduct  = `INSERT INTO products(name, observation, price, created_at) VALUES (?, ?, ?, ?)`
	mySQLGetAllProduct  = `SELECT id, name, observation, price, created_at, updated_at from products`
	mySQLGetProductByID = mySQLGetAllProduct + " WHERE id = ?"
	mySQLUpdateProduct  = `UPDATE products SET name = ?, observation = ?, price = ?, updated_at = ? WHERE id = ?`
	mySQLDeleteProduct  = "DELETE FROM products WHERE id = ?"
)

// mySQLProduct used to work with mySQL - product
type mySQLProduct struct {
	db *sql.DB
}

// NewMySQLProduct returns a new pointer of mySQLProduct
func newMySQLProduct(db *sql.DB) *mySQLProduct {
	return &mySQLProduct{db}
}

// Migrate implements interface product.storage
func (p *mySQLProduct) Migrate() error {
	stmt, err := p.db.Prepare(mySQLMigrateProduct)
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
func (p *mySQLProduct) Create(m *product.Model) error {
	stmt, err := p.db.Prepare(mySQLCreateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		m.Name,
		stringToNull(m.Observations),
		m.Price,
		m.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = uint(id)

	fmt.Printf("Se creó producto correctamente con ID: %d\n", m.ID)
	return nil
}

// GetAll implements interface product.storage
func (p *mySQLProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(mySQLGetAllProduct)
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
func (p *mySQLProduct) GetByID(id uint) (*product.Model, error) {
	stmt, err := p.db.Prepare(mySQLGetProductByID)
	if err != nil {
		return &product.Model{}, err
	}
	defer stmt.Close()

	return scanRowProduct(stmt.QueryRow(id))
}

// Update implements interface product.storage
func (p *mySQLProduct) Update(m *product.Model) error {
	stmt, err := p.db.Prepare(mySQLUpdateProduct)
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
func (p *mySQLProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(mySQLDeleteProduct)
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
