package bencode

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
)

// An Encoder writes Bencode values to an output stream.
type Encoder struct {
	w   io.Writer
	buf []byte
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return NewEncoderWithBuffer(w, make([]byte, 0, 512))
}

// NewEncoderWithBuffer returns a new encoder that writes to w.
func NewEncoderWithBuffer(w io.Writer, buf []byte) *Encoder {
	return &Encoder{
		w:   w,
		buf: buf,
	}
}

// Encode writes the Bencode encoding of v to the stream.
func (e *Encoder) Encode(v any) error {
	e.buf = e.buf[:0]
	if err := e.marshal(v); err != nil {
		return fmt.Errorf("bencode: encode failed: %w", err)
	}
	_, err := e.w.Write(e.buf)
	return err
}

func (e *Encoder) marshal(v any) error {
	switch v := v.(type) {
	case []byte:
		e.marshalBytes(v)
	case string:
		e.marshalString(v)

	case M:
		return e.marshalDictionary(v)
	case D:
		return e.marshalDictionaryNew(v)
	case A:
		return e.marshalSlice(v)

	case map[string]any:
		return e.marshalDictionary(v)

	case []any:
		return e.marshalSlice(v)

	case int, int8, int16, int32, int64:
		e.marshalIntGen(v)
	case uint, uint8, uint16, uint32, uint64:
		e.marshalIntGen(v)

	case float32:
		e.marshalInt(int64(math.Float32bits(v)))
	case float64:
		e.marshalInt(int64(math.Float64bits(v)))

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
		e.buf = append(e.buf, raw...)

	default:
		return e.marshalReflect(reflect.ValueOf(v))
	}
	return nil
}

func (e *Encoder) writeInt(n int64) {
	var bs [20]byte // max_str_len( math.MaxInt64, math.MinInt64 ) base 10
	buf := strconv.AppendInt(bs[0:0], n, 10)
	e.buf = append(e.buf, buf...)
}

func (e *Encoder) marshalBytes(b []byte) {
	// manual inline of writeInt
	var bs [20]byte // max_str_len( math.MaxInt64, math.MinInt64 ) base 10
	buf := strconv.AppendInt(bs[0:0], int64(len(b)), 10)
	buf = append(buf, ':')
	e.buf = append(e.buf, buf...)
	e.buf = append(e.buf, b...)
}

func (e *Encoder) marshalString(s string) {
	// manual inline of writeInt
	var bs [20]byte // max_str_len( math.MaxInt64, math.MinInt64 ) base 10
	buf := strconv.AppendInt(bs[0:0], int64(len(s)), 10)
	buf = append(buf, ':')
	e.buf = append(e.buf, buf...)
	e.buf = append(e.buf, s...)
}

func (e *Encoder) marshalIntGen(val any) {
	var num int64
	switch val := val.(type) {
	case int64:
		num = val
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
	}
	e.marshalInt(num)
}

func (e *Encoder) marshalInt(num int64) {
	e.buf = append(e.buf, 'i')
	e.writeInt(num)
	e.buf = append(e.buf, 'e')
}

func (e *Encoder) marshalReflect(val reflect.Value) error {
	switch val.Kind() {
	case reflect.Slice:
		return e.marshalSliceReflect(val)
	case reflect.Array:
		return e.marshalArrayReflect(val)

	case reflect.Map:
		return e.marshalMap(val)

	case reflect.Struct:
		return e.marshalStruct(val)

	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		return e.marshal(val.Elem().Interface())

	default:
		return fmt.Errorf("unknown kind: %q", val)
	}
}

func (e *Encoder) marshalSliceReflect(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind == reflect.Uint8 {
		e.marshalBytes(val.Bytes())
		return nil
	}
	return e.marshalList(val)
}

func (e *Encoder) marshalArrayReflect(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		return e.marshalList(val)
	}

	e.writeInt(int64(val.Len()))
	buf := make([]byte, 1+val.Len())
	buf[0] = ':'

	for i := 1; i <= val.Len(); i++ {
		buf[i] = byte(val.Index(i - 1).Uint())
	}
	e.buf = append(e.buf, buf...)
	return nil
}

func (e *Encoder) marshalList(val reflect.Value) error {
	if val.Len() == 0 {
		e.buf = append(e.buf, "le"...)
		return nil
	}

	e.buf = append(e.buf, 'l')
	for i := 0; i < val.Len(); i++ {
		if err := e.marshal(val.Index(i).Interface()); err != nil {
			return err
		}
	}
	e.buf = append(e.buf, 'e')
	return nil
}

func (e *Encoder) marshalMap(val reflect.Value) error {
	rawKeys := val.MapKeys()
	if len(rawKeys) == 0 {
		e.buf = append(e.buf, "de"...)
		return nil
	}

	keys := make([]string, len(rawKeys))

	for i, key := range rawKeys {
		if key.Kind() != reflect.String {
			return errors.New("map can be marshaled only if keys are of type 'string'")
		}
		keys[i] = key.String()
	}

	sortStrings(keys)

	e.buf = append(e.buf, 'd')
	for _, key := range keys {
		e.marshalString(key)

		value := val.MapIndex(reflect.ValueOf(key))
		if err := e.marshal(value.Interface()); err != nil {
			return err
		}
	}
	e.buf = append(e.buf, 'e')
	return nil
}

func (e *Encoder) marshalStruct(x reflect.Value) error {
	dict := make(dictStruct, 0, x.Type().NumField())

	dict, err := walkStruct(dict, x)
	if err != nil {
		return err
	}

	sort.Sort(dict)

	e.buf = append(e.buf, 'd')
	for _, def := range dict {
		e.marshalString(def.Key)
		if err := e.marshal(def.Value.Interface()); err != nil {
			return err
		}
	}
	e.buf = append(e.buf, 'e')
	return nil
}

type dictStruct []dictPair

type dictPair struct {
	Key   string
	Value reflect.Value
}

func (d dictStruct) Len() int           { return len(d) }
func (d dictStruct) Less(i, j int) bool { return d[i].Key < d[j].Key }
func (d dictStruct) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func walkStruct(dict dictStruct, v reflect.Value) (dictStruct, error) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		strField := t.Field(i)
		field := v.FieldByIndex(strField.Index)

		if !field.CanInterface() || isNil(field) {
			continue
		}

		tag, ok := fieldTag(strField, field)
		if !ok {
			continue
		}

		if tag == "" && strField.Anonymous &&
			strField.Type.Kind() == reflect.Struct {

			var err error
			dict, err = walkStruct(dict, field)
			if err != nil {
				return nil, err
			}
		} else {
			dict = append(dict, dictPair{Key: tag, Value: field})
		}
	}
	return dict, nil
}

func (e *Encoder) marshalDictionaryNew(dict D) error {
	if len(dict) == 0 {
		e.buf = append(e.buf, "de"...)
		return nil
	}

	// TODO(cristaloleg): maybe reuse sort from util.go ?
	sort.Slice(dict, func(i, j int) bool {
		return dict[i].K < dict[j].K
	})

	e.buf = append(e.buf, 'd')
	for _, pair := range dict {
		e.marshalString(pair.K)
		if err := e.marshal(pair.V); err != nil {
			return err
		}
	}
	e.buf = append(e.buf, 'e')
	return nil
}

func (e *Encoder) marshalDictionary(dict map[string]any) error {
	if len(dict) == 0 {
		e.buf = append(e.buf, "de"...)
		return nil
	}

	// less than `strSliceLen` keys in dict? - take from pool
	var keys []string
	if len(dict) <= strSliceLen {
		strArr := getStrArray()
		defer putStrArray(strArr)
		keys = strArr[:0:len(dict)]
	} else {
		keys = make([]string, 0, len(dict))
	}

	for key := range dict {
		keys = append(keys, key)
	}

	sortStrings(keys)

	e.buf = append(e.buf, 'd')
	for _, key := range keys {
		e.marshalString(key)
		if err := e.marshal(dict[key]); err != nil {
			return err
		}
	}
	e.buf = append(e.buf, 'e')
	return nil
}

func (e *Encoder) marshalSlice(v []any) error {
	if len(v) == 0 {
		e.buf = append(e.buf, "le"...)
		return nil
	}

	e.buf = append(e.buf, 'l')
	for _, data := range v {
		if err := e.marshal(data); err != nil {
			return err
		}
	}
	e.buf = append(e.buf, 'e')
	return nil
}
