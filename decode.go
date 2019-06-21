package bencode

import (
	"bytes"
	"io"
)

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

// A Decoder reads and decodes Bencode values from an input stream.
type Decoder struct {
	r io.Reader
	d decoder
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the Bencode encoding value from stream
// and stores it in the value pointed to by v.
func (dec *Decoder) Decode(v interface{}) error {
	dec.d.Reset()
	if err := dec.d.Unmarshal(v); err != nil {
		return err
	}
	return nil
}

// we can ignore every error result from bytes.Buffer 'cause it's nil
type decoder struct {
	bytes.Buffer
}

func (d *decoder) Unmarshal(v interface{}) error {
	return nil
}
