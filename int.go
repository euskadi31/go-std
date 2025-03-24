package std

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Int is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	Data  int64
	Valid bool // Valid is true if Int64 is not NULL
}

// NewInt creates a new Int.
func NewInt(i int64, valid bool) Int {
	return Int{
		Data:  i,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int64) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}

	return NewInt(*i, true)
}

// Scan implements the Scanner interface.
func (i *Int) Scan(value interface{}) error {
	if value == nil {
		i.Data, i.Valid = 0, false

		return nil
	}

	i.Valid = true

	return convertAssign(&i.Data, value)
}

// Value implements the driver Valuer interface.
func (i Int) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	return i.Data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error {
	var (
		err error
		v   interface{}
	)

	if err = json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("json: cannot unmarshal %s into Go value of type null.Int: %w", string(data), err)
	}

	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Data)
	case nil:
		i.Valid = false

		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}

	i.Valid = err == nil

	return err // nolint: wrapcheck
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false

		return nil
	}

	var err error

	i.Data, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil

	return err // nolint: wrapcheck
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}

	return []byte(strconv.FormatInt(i.Data, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}

	return []byte(strconv.FormatInt(i.Data, 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Data = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}

	return &i.Data
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}

// String implements fmt.Stringer interface.
func (i Int) String() string {
	if !i.Valid {
		return ""
	}

	return strconv.Itoa(int(i.Data))
}
