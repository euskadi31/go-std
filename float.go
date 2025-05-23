package std

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Float is a nullable float64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float struct {
	Data  float64
	Valid bool // Valid is true if Float64 is not NULL
}

// NewFloat creates a new Float.
func NewFloat(f float64, valid bool) Float {
	return Float{
		Data:  f,
		Valid: valid,
	}
}

// FloatFrom creates a new Float that will always be valid.
func FloatFrom(f float64) Float {
	return NewFloat(f, true)
}

// FloatFromPtr creates a new Float that be null if f is nil.
func FloatFromPtr(f *float64) Float {
	if f == nil {
		return NewFloat(0, false)
	}

	return NewFloat(*f, true)
}

// Scan implements the Scanner interface.
func (f *Float) Scan(value interface{}) error {
	if value == nil {
		f.Data, f.Valid = 0, false

		return nil
	}

	f.Valid = true

	return convertAssign(&f.Data, value)
}

// Value implements the driver Valuer interface.
func (f Float) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}

	return f.Data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float.
// It also supports unmarshalling a sql.NullFloat64.
func (f *Float) UnmarshalJSON(data []byte) error {
	var (
		err error
		v   interface{}
	)

	if err = json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("json: cannot unmarshal %s into Go value of type null.Float: %w", string(data), err)
	}

	switch x := v.(type) {
	case float64:
		f.Data = float64(x)
	case nil:
		f.Valid = false

		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float", reflect.TypeOf(v).Name())
	}

	f.Valid = err == nil

	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float) UnmarshalText(text []byte) error {
	str := string(text)

	if str == "" || str == "null" {
		f.Valid = false

		return nil
	}

	var err error

	f.Data, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil

	return err // nolint: wrapcheck
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float is null.
func (f Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}

	return []byte(strconv.FormatFloat(f.Data, 'f', -1, 64)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float is null.
func (f Float) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}

	return []byte(strconv.FormatFloat(f.Data, 'f', -1, 64)), nil
}

// SetValid changes this Float's value and also sets it to be non-null.
func (f *Float) SetValid(n float64) {
	f.Data = n
	f.Valid = true
}

// Ptr returns a pointer to this Float's value, or a nil pointer if this Float is null.
func (f Float) Ptr() *float64 {
	if !f.Valid {
		return nil
	}

	return &f.Data
}

// IsZero returns true for invalid Floats, for future omitempty support (Go 1.4?)
// A non-null Float with a 0 value will not be considered zero.
func (f Float) IsZero() bool {
	return !f.Valid
}

// String implements fmt.Stringer interface.
func (f Float) String() string {
	if !f.Valid {
		return ""
	}

	return strconv.FormatFloat(f.Data, 'f', -1, 64)
}
