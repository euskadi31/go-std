// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package std

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var errNilPtr = errors.New("destination pointer is nil") // embedded in descriptive error

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
// nolint: gocyclo
func convertAssign(dest, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		if d, ok := dest.(*string); ok {
			*d = s

			return nil
		}
	case []byte:
		if d, ok := dest.(*string); ok {
			*d = string(s)

			return nil
		}
	}

	var sv reflect.Value

	if d, ok := dest.(*bool); ok {
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool) // nolint: forcetypeassert
		}

		return err // nolint: wrapcheck
	}

	if scanner, ok := dest.(sql.Scanner); ok {
		return scanner.Scan(src) // nolint: wrapcheck
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}

	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}

		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))

		return nil
	}

	// The following conversions use a string value as an intermediate representation
	// to convert between various numeric types.
	//
	// This also allows scanning into user defined types such as "type Int int64".
	// For symmetry, also check for string destination types.
	// nolint: exhaustive
	switch dv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := asString(src)

		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)

			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %w", src, s, dv.Kind(), err)
		}

		dv.SetInt(i64)

		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := asString(src)

		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)

			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %w", src, s, dv.Kind(), err)
		}

		dv.SetUint(u64)

		return nil
	case reflect.Float32, reflect.Float64:
		s := asString(src)

		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)

			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %w", src, s, dv.Kind(), err)
		}

		dv.SetFloat(f64)

		return nil
	case reflect.String:
		switch v := src.(type) {
		case string:
			dv.SetString(v)

			return nil
		case []byte:
			dv.SetString(string(v))

			return nil
		}
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

func strconvErr(err error) error {
	var ne *strconv.NumError
	if errors.As(err, &ne) {
		return ne.Err
	}

	return err
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	c := make([]byte, len(b))
	copy(c, b)

	return c
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}

	rv := reflect.ValueOf(src)

	// nolint: exhaustive
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}

	return fmt.Sprintf("%v", src)
}
