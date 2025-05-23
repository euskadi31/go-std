package std

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	floatJSON     = []byte(`1.2345`)
	nullFloatJSON = []byte(`{"Float64":1.2345,"Valid":true}`)
)

func TestFloatFrom(t *testing.T) {
	f := FloatFrom(1.2345)
	assertFloat(t, f, "FloatFrom()")

	zero := FloatFrom(0)
	if !zero.Valid {
		t.Error("FloatFrom(0)", "is invalid, but should be valid")
	}
}

func TestFloatFromPtr(t *testing.T) {
	n := float64(1.2345)
	iptr := &n
	f := FloatFromPtr(iptr)
	assertFloat(t, f, "FloatFromPtr()")

	null := FloatFromPtr(nil)
	assertNullFloat(t, null, "FloatFromPtr(nil)")
}

func TestUnmarshalFloat(t *testing.T) {
	var f Float
	err := json.Unmarshal(floatJSON, &f)
	assert.NoError(t, err)
	assertFloat(t, f, "float json")

	var null Float
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assertNullFloat(t, null, "null json")

	var badType Float
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat(t, badType, "wrong type json")

	var invalid Float
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalFloat(t *testing.T) {
	var f Float
	err := f.UnmarshalText([]byte("1.2345"))
	assert.NoError(t, err)
	assertFloat(t, f, "UnmarshalText() float")

	var blank Float
	err = blank.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assertNullFloat(t, blank, "UnmarshalText() empty float")

	var null Float
	err = null.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assertNullFloat(t, null, `UnmarshalText() "null"`)
}

func TestMarshalFloat(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := json.Marshal(f)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewFloat(0, false)
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalFloatText(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := f.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewFloat(0, false)
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloatPointer(t *testing.T) {
	f := FloatFrom(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloatIsZero(t *testing.T) {
	f := FloatFrom(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloatSetValid(t *testing.T) {
	change := NewFloat(0, false)
	assertNullFloat(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat(t, change, "SetValid()")
}

func TestFloatScan(t *testing.T) {
	var f Float
	err := f.Scan(1.2345)
	assert.NoError(t, err)
	assertFloat(t, f, "scanned float")

	var null Float
	err = null.Scan(nil)
	assert.NoError(t, err)
	assertNullFloat(t, null, "scanned null")
}

func TestFloatString(t *testing.T) {
	f := FloatFrom(1.2345)
	assert.Equal(t, "1.2345", f.String())

	null := Float{}
	assert.Equal(t, "", null.String())
}

func TestFloatValue(t *testing.T) {
	i := FloatFrom(1.2345)

	v, err := i.Value()
	assert.NoError(t, err)
	assert.Equal(t, float64(1.2345), v.(float64))

	null := Float{}

	v, err = null.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func assertFloat(t *testing.T, f Float, from string) {
	if f.Data != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Data, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat(t *testing.T, f Float, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func BenchmarkFloatString(b *testing.B) {
	f := FloatFrom(1567.5346456)

	for n := 0; n < b.N; n++ {
		_ = f.String()
	}
}
