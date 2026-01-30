package lib

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func DbConnect(dsn string) (*sqlx.DB, error) {

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

func GetString(db *sqlx.DB, query string, args ...any) (string, error) {
	v, err := GetValue(db, query, args...)
	if err == nil {
		if vs, ok := v.(sql.NullString); ok {
			if vs.Valid {
				return vs.String, nil
			} else {
				return "", errors.New("null value")
			}
		} else {
			return "", errors.New("type is not string")
		}
	} else {
		return "", err
	}
}

func GetFloat(db *sqlx.DB, query string, args ...any) (float64, error) {
	v, err := GetValue(db, query, args...)
	if err == nil {
		if vs, ok := v.(sql.NullFloat64); ok {
			if vs.Valid {
				return vs.Float64, nil
			} else {
				return 0, errors.New("null value")
			}
		} else {
			return 0, errors.New("type is not float")
		}
	} else {
		return 0, err
	}
}

func GetInt(db *sqlx.DB, query string, args ...any) (int64, error) {
	v, err := GetValue(db, query, args...)
	if err == nil {
		if vs, ok := v.(sql.NullInt64); ok {
			if vs.Valid {
				return vs.Int64, nil
			} else {
				return 0, errors.New("null value")
			}
		} else {
			return 0, errors.New("type is not float")
		}
	} else {
		return 0, err
	}
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
