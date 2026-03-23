package lib

import (
	"database/sql"
	"encoding/json"
	"strconv"
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

type JSONFloat struct {
	sql.NullFloat64
}

func (jf *JSONFloat) MarshalJSON() ([]byte, error) {
	if !jf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(jf.Float64)
}

func (jf *JSONFloat) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		jf.Float64 = 0
		jf.Valid = false
		return nil
	}

	f, err := strconv.ParseFloat(string(data), 64)
	if err == nil {
		jf.Float64 = f
		jf.Valid = true
		return nil
	} else {
		return err
	}
}

func JSONFloatNil() JSONFloat {
	v := JSONFloat{}
	v.Valid = false
	return v
}

func JSONFloatGet(value float64) JSONFloat {
	v := JSONFloat{}
	v.Valid = true
	v.Float64 = value
	return v
}

type JSONInt struct {
	sql.NullInt64
}

func (ji *JSONInt) MarshalJSON() ([]byte, error) {
	if !ji.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ji.Int64)
}

func (ji *JSONInt) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		ji.Int64 = 0
		ji.Valid = false
		return nil
	}

	i, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		ji.Int64 = i
		ji.Valid = true
		return nil
	} else {
		return err
	}
}

func JSONIntNil() JSONInt {
	v := JSONInt{}
	v.Valid = false
	v.Int64 = 0
	return v
}

func JSONIntGet(value int64) JSONInt {
	v := JSONInt{}
	v.Valid = true
	v.Int64 = value
	return v
}
