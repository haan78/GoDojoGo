package lib

import (
	"errors"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func DbConnect() (*sqlx.DB, error) {
	// Load .env (optional in production)
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return nil, errors.New("no connection string in .env file")
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetValue(db *sqlx.DB, query string, args ...any) (any, error) {
	var value any

	err := db.Get(&value, query, args...)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func GenericQuery[T any](db *sqlx.DB, query string, args ...any) ([]T, error) {
	var result []T
	if err := db.Select(&result, query, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func GenericRow[T any](db *sqlx.DB, query string, args ...any) (T, error) {
	var result []T
	var empty T
	if err := db.Select(&result, query, args...); err != nil {
		return empty, err
	}
	if len(result) > 0 {
		return result[0], nil
	} else {
		return empty, errors.New("no record found")
	}
}
