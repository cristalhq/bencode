package bencode

import (
	"bytes"
	"io"
)

// An Encoder writes Bencode values to an output stream.
type Encoder struct {
	w io.Writer
	e encoder
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the Bencode encoding of v to the stream.
func (enc *Encoder) Encode(v interface{}) error {
	enc.e.Reset()
	if err := enc.e.Marshal(v); err != nil {
		return err
	}
	_, err := enc.w.Write(enc.e.Bytes())
	return err
}

// we can ignore every error result from bytes.Buffer 'cause it's nil
type encoder struct {
	bytes.Buffer
}

func (e *encoder) Marshal(v interface{}) error {
	return nil
}
