package lib

import (
	globals "GoDojoGo/deff"
	"context"
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func DbConnect() (*sqlx.DB, error) {
	// Load .env (optional in production)
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	dsn := globals.Settings.MYSQL_DSN
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

func DbConnectX(ctx context.Context) (*sqlx.DB, error) {
	// Load .env (optional in production)
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	dsn := globals.Settings.MYSQL_DSN
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

func EmptyIsNull(str string) sql.NullString {
	var value sql.NullString

	if strings.Trim(str, " ") == "" {
		value.Scan(nil)
	} else {
		value.Scan(str)
	}
	return value
}
