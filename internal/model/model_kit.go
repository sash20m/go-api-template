package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString is a helper wrapper around sql.NullString that handles the null case of the value when marshalling to JSON.
type NullString sql.NullString

// Scan implements the sql.Scanner interface so database/sql (and sqlx) can scan NULLs and strings into this type.
func (ns *NullString) Scan(value any) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*ns = NullString(s)
	return nil
}

// Value implements the driver.Valuer interface so this type can be used in query parameters.
func (ns NullString) Value() (driver.Value, error) {
	s := sql.NullString(ns)
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// NullInt64 is a helper wrapper around sql.NullInt64 that handles the null case of the value when marshalling to JSON.
type NullInt64 sql.NullInt64

// Scan implements the sql.Scanner interface so database/sql (and sqlx) can scan NULLs and integers into this type.
func (ni *NullInt64) Scan(value any) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	*ni = NullInt64(i)
	return nil
}

// Value implements the driver.Valuer interface so this type can be used in query parameters.
func (ni NullInt64) Value() (driver.Value, error) {
	i := sql.NullInt64(ni)
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	var s *int64
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s != nil {
		ni.Valid = true
		ni.Int64 = *s
	} else {
		ni.Valid = false
	}
	return nil
}
