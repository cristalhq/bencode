package bencode

// Marshaler is the interface implemented by types that
// can marshal themselves into valid Bencode.
type Marshaler interface {
	MarshalBencode() ([]byte, error)
}

// Marshal returns bencode encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	var e encoder
	if err := e.marshal(v); err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

// Unmarshaler is the interface implemented by types
// that can unmarshal a Bencode description of themselves.
type Unmarshaler interface {
	UnmarshalBencode([]byte) error
}

// Unmarshal parses the bencoded data and stores the result
// in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	var d decoder
	if err := d.Unmarshal(v); err != nil {
		return err
	}
	return nil
}
