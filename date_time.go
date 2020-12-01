package std

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// ISO8601 format
const dateTimeFormat = "2006-01-02T15:04:05-0700"

var nullType = []byte("null")

// DateTime is a nullable time.Time with ISO8601 format. It supports SQL and JSON serialization.
// It will marshal to null if null.
// swagger:strfmt date-time
type DateTime struct {
	Data  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (t *DateTime) Scan(value interface{}) error {
	var err error

	switch x := value.(type) {
	case time.Time:
		t.Data = x
	case nil:
		t.Data = time.Time{}
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("std: cannot scan type %T into std.DateTime: %v", value, value)
	}

	t.Valid = err == nil

	return err
}

// Value implements the driver Valuer interface.
func (t DateTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	return t.Data, nil
}

// NewDateTime creates a new DateTime.
func NewDateTime(t time.Time, valid bool) DateTime {
	return DateTime{
		Data:  t,
		Valid: valid,
	}
}

// DateTimeFrom creates a new Time that will always be valid.
func DateTimeFrom(t time.Time) DateTime {
	return NewDateTime(t, !t.IsZero())
}

// DateTimeFromPtr creates a new DateTime that will be null if t is nil.
func DateTimeFromPtr(t *time.Time) DateTime {
	if t == nil {
		return NewDateTime(time.Time{}, false)
	}

	return NewDateTime(*t, true)
}

// MarshalText implement the json.Marshaler interface
func (t DateTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}

	return []byte(t.Data.Format(dateTimeFormat)), nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t DateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullType, nil
	}

	b, _ := t.MarshalText()

	if len(b) == 0 {
		return nullType, nil
	}

	dt := []byte{}
	dt = append(dt, 0x22) // 0x22 => "
	dt = append(dt, b...)
	dt = append(dt, 0x22) // 0x22 => "

	return dt, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *DateTime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	return t.UnmarshalText(b)
}

// UnmarshalText allows ISO8601Time to implement the TextUnmarshaler interface
func (t *DateTime) UnmarshalText(b []byte) error {
	str := string(b)
	var err error

	if str == "" || str == "null" {
		t.Data = time.Time{}
		t.Valid = false

		return nil
	}

	t.Data, err = time.Parse(dateTimeFormat, str)

	if err != nil {
		t.Valid = false
	} else {
		t.Valid = true
	}

	return err
}

// SetValid changes this Time's value and sets it to be non-null.
func (t *DateTime) SetValid(v time.Time) {
	t.Data = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t DateTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}

	return &t.Data
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t DateTime) IsZero() bool {
	return !t.Valid
}

// String implements fmt.Stringer interface
func (t DateTime) String() string {
	if !t.Valid {
		return ""
	}

	return t.Data.Format(dateTimeFormat)
}
