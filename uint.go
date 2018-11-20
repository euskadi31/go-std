package std

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Uint is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint struct {
	Data  uint64
	Valid bool // Valid is true if Uint64 is not NULL
}

// NewUint creates a new Uint
func NewUint(i uint64, valid bool) Uint {
	return Uint{
		Data:  i,
		Valid: valid,
	}
}

// UintFrom creates a new Uint that will always be valid.
func UintFrom(i uint64) Uint {
	return NewUint(i, true)
}

// UintFromPtr creates a new Uint that be null if i is nil.
func UintFromPtr(i *uint64) Uint {
	if i == nil {
		return NewUint(0, false)
	}

	return NewUint(*i, true)
}

// Scan implements the Scanner interface.
func (i *Uint) Scan(value interface{}) error {
	if value == nil {
		i.Data, i.Valid = 0, false
		return nil
	}

	i.Valid = true

	return convertAssign(&i.Data, value)
}

// Value implements the driver Valuer interface.
func (i Uint) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	return i.Data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Uint.
// It also supports unmarshalling a sql.NullInt64.
func (i *Uint) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Data)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint", reflect.TypeOf(v).Name())
	}

	i.Valid = err == nil

	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false

		return nil
	}

	var err error

	i.Data, err = strconv.ParseUint(string(text), 10, 64)
	i.Valid = err == nil

	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint is null.
func (i Uint) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}

	return []byte(strconv.FormatUint(i.Data, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint is null.
func (i Uint) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}

	return []byte(strconv.FormatUint(i.Data, 10)), nil
}

// SetValid changes this Uint's value and also sets it to be non-null.
func (i *Uint) SetValid(n uint64) {
	i.Data = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint's value, or a nil pointer if this Uint is null.
func (i Uint) Ptr() *uint64 {
	if !i.Valid {
		return nil
	}

	return &i.Data
}

// IsZero returns true for invalid Uints, for future omitempty support (Go 1.4?)
// A non-null Uint with a 0 value will not be considered zero.
func (i Uint) IsZero() bool {
	return !i.Valid
}

// String implements fmt.Stringer interface
func (i Uint) String() string {
	if !i.Valid {
		return ""
	}

	return strconv.FormatUint(i.Data, 10)
}
