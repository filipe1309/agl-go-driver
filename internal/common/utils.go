package common

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 sql.NullInt64

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

// Scan implements the Scanner interface for NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}

	return nil
}

// ToInt64 returns the int64 value
func (ni *NullInt64) ToInt64() int64 {
	return ni.Int64
}

// ValidNullInt64 is a one-liner that returns a valid NullInt64 type
func ValidNullInt64(i int64) NullInt64 {
	return NullInt64{
		Int64: i,
		Valid: true,
	}
}

// InvalidNullInt64 is a one-liner that returns an invalid NullInt64 type
func InvalidNullInt64(i int64) NullInt64 {
	return NullInt64{
		Int64: i,
		Valid: false,
	}
}
