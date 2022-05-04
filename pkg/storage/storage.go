package storage

import (
	"database/sql"
	"fmt"
	"github.com/eltaljohn/go-db/pkg/product"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"time"
)

var (
	db   *sql.DB
	once sync.Once
)

// Driver of storage
type Driver string

// Drivers
const (
	MySQL    Driver = "MYSQL"
	Postgres Driver = "POSTGRES"
)

// New creates the connection with DB
func New(d Driver) {
	switch d {
	case MySQL:
		newMySQLDB()
	case Postgres:
		newPostgresDB()
	}
}

func newPostgresDB() {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			"postgres",
			"password",
			"127.0.0.1",
			"5432",
			"postgres",
		)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Can't open db: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("Can't do ping db: %v", err)
		}

		fmt.Println("Connected to PostgreSQL")
	})
}

func newMySQLDB() {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?tls=false&autocommit=true&allowNativePasswords=true&parseTime=true",
			"root",
			"password",
			"127.0.0.1",
			"3306",
			"mysql",
		)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Can't open db: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("Can't do ping db: %v", err)
		}

		fmt.Println("Connected to MySQL")
	})
}

// Pool return a unique intance of db
func Pool() *sql.DB {
	return db
}

func stringToNull(s string) sql.NullString {
	null := sql.NullString{String: s}
	if null.String != "" {
		null.Valid = true
	}
	return null
}

func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}
	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanRowProduct(s scanner) (*product.Model, error) {
	m := &product.Model{}
	observationNull := sql.NullString{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&m.ID,
		&m.Name,
		&observationNull,
		&m.Price,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return &product.Model{}, err
	}

	m.Observations = observationNull.String
	m.UpdatedAt = updatedAtNull.Time

	return m, nil
}

// DAOProduct factory of product.storage
func DAOProduct(driver Driver) (product.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlProduct(db), nil
	case MySQL:
		return newMySQLProduct(db), nil

	default:
		return nil, fmt.Errorf("driver not implemented")
	}
}
