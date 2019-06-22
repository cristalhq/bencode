package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
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
	switch v := v.(type) {
	case Marshaler:
		raw, err := v.MarshalBencode()
		if err != nil {
			return err
		}
		_, _ = e.Write(raw)
		return nil

	case []byte:
		return e.marshalBytes(v)

	case string:
		return e.marshalString(v)

	case int:
		return e.marshalInt(int64(v))

	case int8:
		return e.marshalInt(int64(v))

	case int16:
		return e.marshalInt(int64(v))

	case int32:
		return e.marshalInt(int64(v))

	case int64:
		return e.marshalInt(int64(v))

	case uint:
		return e.marshalUInt(uint64(v))

	case uint8:
		return e.marshalUInt(uint64(v))

	case uint16:
		return e.marshalUInt(uint64(v))

	case uint32:
		return e.marshalUInt(uint64(v))

	case uint64:
		return e.marshalUInt(uint64(v))

	case bool:
		return e.marshalBool(v)

	case map[string]interface{}:
		return e.marshalDictionary(v)

	case []string:
		return e.marshalStringsSlice(v)

	case []interface{}:
		return e.marshalSlice2(v)

	default:
		val := reflect.ValueOf(v)
		return e.marshal(val)
	}
}

func (e *encoder) marshal(val reflect.Value) error {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.marshalIntRefl(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.marshalUIntRefl(val)

	case reflect.String:
		return e.marshalStringRefl(val)

	case reflect.Slice:
		return e.marshalSlice(val)

	case reflect.Array:
		return e.marshalArray(val)

	case reflect.Map:
		return e.marshalMap(val)

	case reflect.Struct:
		return e.marshalStruct(val)

	case reflect.Ptr:
		return e.Marshal(val.Elem().Interface())

	case reflect.Interface:
		return e.Marshal(val.Elem().Interface())

	default:
		return fmt.Errorf("Unknown kind: %v", val.Kind())
	}
}

func (e *encoder) marshalIntRefl(val reflect.Value) error {
	_ = e.WriteByte('i')
	buf := strconv.AppendInt([]byte{}, val.Int(), 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalUIntRefl(val reflect.Value) error {
	_ = e.WriteByte('i')
	buf := strconv.AppendUint([]byte{}, val.Uint(), 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalInt(val int64) error {
	_ = e.WriteByte('i')
	buf := strconv.AppendInt([]byte{}, val, 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalUInt(val uint64) error {
	_ = e.WriteByte('i')
	buf := strconv.AppendUint([]byte{}, val, 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalBool(val bool) error {
	_ = e.WriteByte('i')
	if val {
		_ = e.WriteByte('1')
	} else {
		_ = e.WriteByte('0')
	}
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalStringRefl(val reflect.Value) error {
	buf := strconv.AppendInt([]byte{}, int64(len(val.String())), 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte(':')
	_, _ = e.Write([]byte(val.String()))
	return nil
}

func (e *encoder) marshalBytes(val []byte) error {
	buf := strconv.AppendInt([]byte{}, int64(len(val)), 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte(':')
	_, _ = e.Write(val)
	return nil
}

func (e *encoder) marshalString(val string) error {
	buf := strconv.AppendInt([]byte{}, int64(len(val)), 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte(':')
	_, _ = e.WriteString(val)
	return nil
}

func (e *encoder) marshalSlice(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		return e.marshalList(val)
	}

	// treat slice like string
	valBytes := val.Bytes()

	_, _ = e.Write(strconv.AppendInt([]byte{}, int64(len(valBytes)), 10))
	_ = e.WriteByte(':')
	_, _ = e.Write(valBytes)
	return nil
}

func (e *encoder) marshalArray(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		return e.marshalList(val)
	}

	buf := strconv.AppendInt([]byte{}, int64(val.Len()), 10)
	_, _ = e.Write(buf)
	_ = e.WriteByte(':')

	for i := 0; i < val.Len(); i++ {
		v := byte(val.Index(i).Uint())
		_ = e.WriteByte(v)
	}
	return nil
}

func (e *encoder) marshalList(val reflect.Value) error {
	_ = e.WriteByte('l')

	for i := 0; i < val.Len(); i++ {
		// array of interface{} values, need to extract unterling type of element
		element := reflect.ValueOf(val.Index(i).Interface())
		if err := e.marshal(element); err != nil {
			return err
		}
	}
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalMap(val reflect.Value) error {
	_ = e.WriteByte('d')

	keys := val.MapKeys()
	rawKeys := make(bytesSlice, len(keys))

	for i, key := range keys {
		if key.Kind() != reflect.String {
			return errors.New("Map can be marshaled only if keys are of type 'string'")
		}
		rawKeys[i] = []byte(key.String())
	}
	sort.Sort(rawKeys)

	for _, rawKey := range rawKeys {
		key := string(rawKey)
		vKey := reflect.ValueOf(key)
		if err := e.marshal(vKey); err != nil {
			return err
		}
		value := val.MapIndex(vKey)
		if err := e.marshal(value); err != nil {
			return err
		}
	}
	return e.WriteByte('e')
}

func (e *encoder) marshalStruct(val reflect.Value) error {
	_ = e.WriteByte('d')

	valType := val.Type()

	fields := positionedFieldsByName{}
	for i := 0; i < val.NumField(); i++ {
		fieldOpt := extractFieldOptions(val, valType.Field(i).Name)
		if len(fieldOpt) == 0 {
			continue
		}
		fields = append(fields, positionedField{[]byte(fieldOpt), i})
	}

	sort.Sort(fields)

	for _, f := range fields {
		if err := e.marshal(reflect.ValueOf(f.name)); err != nil {
			return err
		}
		if err := e.marshal(val.Field(f.pos)); err != nil {
			return err
		}
	}
	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalDictionary(d map[string]interface{}) error {
	_ = e.WriteByte('d')

	for key, data := range d {
		_ = e.marshalString(key)

		if err := e.Marshal(data); err != nil {
			return err
		}
	}

	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalSlice2(v []interface{}) error {
	_ = e.WriteByte('l')

	for _, data := range v {
		if err := e.Marshal(data); err != nil {
			return err
		}
	}

	_ = e.WriteByte('e')
	return nil
}

func (e *encoder) marshalStringsSlice(v []string) error {
	_ = e.WriteByte('l')

	for _, data := range v {
		_ = e.marshalString(data)
	}

	_ = e.WriteByte('e')
	return nil
}
