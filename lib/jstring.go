package lib

import (
	"database/sql"
	"encoding/json"
)

type JSONNullString struct {
	sql.NullString
}

func (ns *JSONNullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *JSONNullString) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		ns.String = ""
		ns.Valid = false
		return nil
	}

	// Otherwise expect a string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ns.String = s
	ns.Valid = true
	return nil
}
