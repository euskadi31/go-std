Go Standard Library ![Last release](https://img.shields.io/github/release/euskadi31/go-std.svg)
===================

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-std)](https://goreportcard.com/report/github.com/euskadi31/go-std)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/euskadi31/go-std/master.svg)](https://travis-ci.org/euskadi31/go-std) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-std/master.svg)](https://coveralls.io/github/euskadi31/go-std?branch=master) |


go-std is a library with reasonable options for dealing with nullable SQL and JSON values.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this library in place of `sql.NullXXX`.
All types also implement: `encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`, `json.Unmarshaler` and `fmt.Stringer`.

Types
-----

- `std.Bool`: Nullable bool
- `std.Float`: Nullable float64
- `std.String`: Nullable string
- `std.Int`: Nullable int64
- `std.Uint`: Nullable uint64
- `std.Time`: Nullable Time
- `std.DateTime`: Nullable Time with ISO8601 format
- `std.Date`: Nullable Time with ISO8601 (yyyy-mm-dd) format


License
-------

go-std is licensed under [the MIT license](LICENSE.md).
