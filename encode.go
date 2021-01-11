package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"sync"
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

// NewEncoderWithBuffer returns a new encoder that writes to w.
func NewEncoderWithBuffer(w io.Writer, buf *bytes.Buffer) *Encoder {
	return &Encoder{
		w: w,
		e: encoder{*buf},
	}
}

// Encode writes the Bencode encoding of v to the stream.
func (enc *Encoder) Encode(v interface{}) error {
	enc.e.Reset()
	if err := enc.e.marshal(v); err != nil {
		return err
	}
	_, err := enc.w.Write(enc.e.Bytes())
	return err
}

// we can ignore every error result from bytes.Buffer 'cause it's nil
type encoder struct {
	bytes.Buffer
}

func (e *encoder) marshal(v interface{}) error {
	switch v := v.(type) {
	case []byte:
		e.marshalBytes(v)
	case string:
		e.marshalString(v)

	case map[string]interface{}:
		return e.marshalDictionary(v)

	case []interface{}:
		return e.marshalSlice(v)

	case int, int8, int16, int32, int64:
		e.marshalIntGen(v)
	case uint, uint8, uint16, uint32, uint64:
		e.marshalIntGen(v)

	case bool:
		var n int64
		if v {
			n = 1
		}
		e.marshalInt(n)

	case Marshaler:
		raw, err := v.MarshalBencode()
		if err != nil {
			return err
		}
		e.Write(raw)

	default:
		return e.marshalReflect(reflect.ValueOf(v))
	}
	return nil
}

func (e *encoder) writeInt(n int64) {
	var bs [20]byte // max_str_len( math.MaxInt64, math.MinInt64 ) base 10
	buf := strconv.AppendInt(bs[0:0], n, 10)
	e.Write(buf)
}

func (e *encoder) marshalBytes(b []byte) error {
	e.writeInt(int64(len(b)))
	e.WriteByte(':')
	e.Write(b)
	return nil
}

func (e *encoder) marshalString(s string) error {
	e.writeInt(int64(len(s)))
	e.WriteByte(':')
	e.WriteString(s)
	return nil
}

func (e *encoder) marshalIntGen(val interface{}) error {
	var num int64
	switch val := val.(type) {
	case int64:
		num = int64(val)
	case int32:
		num = int64(val)
	case int16:
		num = int64(val)
	case int8:
		num = int64(val)
	case int:
		num = int64(val)
	case uint64:
		num = int64(val)
	case uint32:
		num = int64(val)
	case uint16:
		num = int64(val)
	case uint8:
		num = int64(val)
	case uint:
		num = int64(val)
	default:
		return fmt.Errorf("unknown int type %T", val)
	}
	e.marshalInt(num)
	return nil
}

func (e *encoder) marshalInt(num int64) error {
	e.WriteByte('i')
	e.writeInt(num)
	e.WriteByte('e')
	return nil
}

func (e *encoder) marshalReflect(val reflect.Value) error {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.marshalIntGen(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.marshalIntGen(val.Uint())

	case reflect.String:
		e.marshalString(val.String())

	case reflect.Slice:
		return e.marshalSliceReflect(val)
	case reflect.Array:
		return e.marshalArrayReflect(val)

	case reflect.Map:
		return e.marshalMap(val)

	case reflect.Struct:
		return e.marshalStruct(val)

	case reflect.Ptr:
		return e.marshal(val.Elem().Interface())

	case reflect.Interface:
		return e.marshal(val.Elem().Interface())

	case reflect.Bool:
		return errors.New("reflect.Bool")
	case reflect.Uintptr:
		return errors.New("reflect.Uintptr")
	case reflect.Float32:
		return errors.New("reflect.Float32")
	case reflect.Float64:
		return errors.New("reflect.Float64")
	case reflect.Complex64:
		return errors.New("reflect.Complex64")
	case reflect.Complex128:
		return errors.New("reflect.Complex128")
	case reflect.Chan:
		return errors.New("reflect.Chan")
	case reflect.Func:
		return errors.New("reflect.Func")
	case reflect.UnsafePointer:
		return errors.New("reflect.UnsafePointer")

	default:
		return fmt.Errorf("Unknown kind: %q", val)
	}
	return nil
}

func (e *encoder) marshalSliceReflect(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind == reflect.Uint8 {
		return e.marshalBytes(val.Bytes())
	}
	return e.marshalList(val)
}

func (e *encoder) marshalArrayReflect(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		return e.marshalList(val)
	}

	e.writeInt(int64(val.Len()))
	e.WriteByte(':')

	for i := 0; i < val.Len(); i++ {
		v := byte(val.Index(i).Uint())
		e.WriteByte(v)
	}
	return nil
}

func (e *encoder) marshalList(val reflect.Value) error {
	e.WriteByte('l')

	for i := 0; i < val.Len(); i++ {
		if err := e.marshal(val.Index(i).Interface()); err != nil {
			return err
		}
	}

	e.WriteByte('e')
	return nil
}

func (e *encoder) marshalMap(val reflect.Value) error {
	e.WriteByte('d')

	rawKeys := val.MapKeys()
	if len(rawKeys) == 0 {
		e.WriteByte('e')
		return nil
	}
	keys := make([]string, len(rawKeys))

	for i, key := range rawKeys {
		if key.Kind() != reflect.String {
			return errors.New("Map can be marshaled only if keys are of type 'string'")
		}
		keys[i] = key.String()
	}

	sortStrings(keys)

	for _, key := range rawKeys {
		key := string(key.String())
		e.marshalString(key)

		vKey := reflect.ValueOf(key)
		value := val.MapIndex(vKey)
		if err := e.marshalReflect(value); err != nil {
			return err
		}
	}

	e.WriteByte('e')
	return nil
}

func (e *encoder) marshalStruct(val reflect.Value) error {
	// TODO
	return nil
}

func (e *encoder) marshalDictionary(dict map[string]interface{}) error {
	e.WriteByte('d')

	if len(dict) == 0 {
		e.WriteByte('e')
		return nil
	}

	keys := strslicePool.Get().([]string)
	defer func() {
		keys = keys[:0]
		strslicePool.Put(keys)
	}()

	for key := range dict {
		keys = append(keys, key)
	}

	sortStrings(keys)

	for _, key := range keys {
		e.marshalString(key)
		if err := e.marshal(dict[key]); err != nil {
			return err
		}
	}

	e.WriteByte('e')
	return nil
}

func (e *encoder) marshalSlice(v []interface{}) error {
	e.WriteByte('l')

	for _, data := range v {
		if err := e.marshal(data); err != nil {
			return err
		}
	}

	e.WriteByte('e')
	return nil
}

var strslicePool = sync.Pool{
	New: func() interface{} {
		return make([]string, 0, 20)
	},
}

func sortStrings(ss []string) {
	if len(ss) <= 20 {
		// for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
		// 	ss[j], ss[j-1] = ss[j-1], ss[j]
		// }
		// below is the code above, but (almost) without bound checks

		for i := 1; i < len(ss); i++ {
			j := i
			for {
				a, b := ss[j], ss[j-1]
				if j > 0 && a < b && j < len(ss) {
					a, b = b, a
					ss[j] = a
					ss[j-1] = b // one bound check
					j--
				}
				break
			}
		}
	} else {
		sort.Strings(ss)
	}
}
