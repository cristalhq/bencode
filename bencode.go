package bencode

import (
	"bytes"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into valid Bencode.
type Marshaler interface {
	MarshalBencode() ([]byte, error)
}

// Marshal returns bencode encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalTo returns bencode encoding of v written to dst.
func MarshalTo(dst []byte, v interface{}) ([]byte, error) {
	enc := &Encoder{buf: dst}
	if err := enc.marshal(v); err != nil {
		return nil, err
	}
	return enc.buf, nil
}

// Unmarshaler is the interface implemented by types
// that can unmarshal a Bencode description of themselves.
type Unmarshaler interface {
	UnmarshalBencode([]byte) error
}

// Unmarshal parses the bencoded data and stores the result
// in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	d := NewDecodeBytes(data)
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}

// A is a Bencode array.
//
// Example:
//
//	bencode.A{"hello", "world", 3.14159, bencode.D{{"foo", 12345}}}
type A []interface{}

// D is an ordered representation of a Bencode document.
//
// Example usage:
//
//	bencode.D{{"hello", "world"}, {"foo", "bar"}, {"pi", 3.14159}}
type D []e

// e represents a Bencode element for a D. It is usually used inside a D.
type e struct {
	K string
	V interface{}
}

// M is an unordered representation of a Bencode document.
//
// Example usage:
//
//	bencode.M{"hello": "world", "foo": "bar", "pi": 3.14159}
type M map[string]interface{}
