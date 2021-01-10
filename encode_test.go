package bencode

import (
	"bytes"
	"errors"
	"testing"
)

type encodeTestCase struct {
	in  interface{}
	out []byte
	err error
}

func TestEncodeMarshaler(t *testing.T) {
	tcs := []encodeTestCase{
		// {
		// 	myBoolType(true),
		// 	[]byte(`1:y`),
		// 	nil,
		// },
		// {
		// 	myBoolType(true),
		// 	[]byte(`1:y`),
		// 	nil,
		// },
		// {
		// 	myBoolType(false),
		// 	[]byte(`1:n`),
		// 	nil,
		// },
		// {
		// 	myTimeType{now},
		// 	fmt.Sprintf("i%de", now.Unix()),
		// 	nil,
		// },
		{
			errorMarshalType{},
			[]byte(``),
			ErrJustAnError,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestMarshalInt(t *testing.T) {
	tcs := []struct {
		val interface{}
		raw string
	}{
		{1, `i1e`},
		{-1, `i-1e`},
		{0, `i0e`},
		{43210, `i43210e`},
		{int(1), `i1e`},
		{int8(2), `i2e`},
		{int16(3), `i3e`},
		{int32(4), `i4e`},
		{int64(5), `i5e`},
		{uint(6), `i6e`},
		{uint8(7), `i7e`},
		{uint16(8), `i8e`},
		{uint32(9), `i9e`},
		{uint64(10), `i10e`},
	}
	for i, test := range tcs {
		buf, err := Marshal(test.val)
		if err != nil {
			t.Fatalf("[test %d] unexpected err %v", i, err)
		}

		got := string(buf)
		if want := string(test.raw); got != want {
			t.Fatalf("[test %d] got %v want: %v", i, got, want)
		}
	}
}

func TestEncodeBool(t *testing.T) {
	tcs := []encodeTestCase{
		{
			true,
			[]byte(`i1e`),
			nil,
		},
		{
			false,
			[]byte(`i0e`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodeString(t *testing.T) {
	tcs := []encodeTestCase{
		// {
		// 	(*string)(nil),
		// 	[]byte(``),
		// 	nil,
		// },
		{
			[]byte(``),
			[]byte(`0:`),
			nil,
		},
		{
			[]byte(`x`),
			[]byte(`1:x`),
			nil,
		},
		{
			[]byte(`foobarbaz`),
			[]byte(`9:foobarbaz`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodeSliceAsString(t *testing.T) {
	tcs := []encodeTestCase{
		{
			[]byte(nil),
			[]byte{'0', ':'},
			nil,
		},
		{
			[]byte{},
			[]byte{'0', ':'},
			nil,
		},
		{
			[]byte(`test`),
			[]byte(`4:test`),
			nil,
		},
		{
			[]byte{0, 1, 2, 3},
			[]byte{'4', ':', 0, 1, 2, 3},
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodeArrayAsString(t *testing.T) {
	tcs := []encodeTestCase{
		{
			[...]byte{},
			[]byte(`0:`),
			nil,
		},
		{
			[...]byte{'x', 'y'},
			[]byte(`2:xy`),
			nil,
		},
		{
			[...]byte{0, 1, 2},
			[]byte{'3', ':', 0, 1, 2},
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodeSliceAsList(t *testing.T) {
	tcs := []encodeTestCase{
		{
			[]interface{}{},
			[]byte(`le`),
			nil,
		},
		{
			[]interface{}{1, 2, 3},
			[]byte(`li1ei2ei3ee`),
			nil,
		},
		{
			[]interface{}{"foo", "bar", "baz"},
			[]byte(`l3:foo3:bar3:baze`),
			nil,
		},
		{
			[]interface{}{1, "foo", 2, "bar"},
			[]byte(`li1e3:fooi2e3:bare`),
			nil,
		},
		{
			[]interface{}{1, []interface{}{"bar"}},
			[]byte(`li1el3:baree`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodeArrayAsList(t *testing.T) {
	tcs := []encodeTestCase{
		{
			[...]interface{}{},
			[]byte(`le`),
			nil,
		},
		{
			[...]interface{}{1, 2, 3},
			[]byte(`li1ei2ei3ee`),
			nil,
		},
		{
			[...]interface{}{"foo", "bar", "baz"},
			[]byte(`l3:foo3:bar3:baze`),
			nil,
		},
		{
			[...]interface{}{1, "foo", 2, "bar"},
			[]byte(`li1e3:fooi2e3:bare`),
			nil,
		},
		{
			[...]interface{}{1, [...]interface{}{"bar"}},
			[]byte(`li1el3:baree`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodeStringsSlice(t *testing.T) {
	tcs := []encodeTestCase{
		{
			[]string(nil),
			[]byte(`le`),
			nil,
		},
		{
			[]string{},
			[]byte(`le`),
			nil,
		},
		{
			[]string{"foo"},
			[]byte(`l3:fooe`),
			nil,
		},
		{
			[]string{"foo", "barbaz"},
			[]byte(`l3:foo6:barbaze`),
			nil,
		},
		{
			[]string{"foo", "barbaz", "go"},
			[]byte(`l3:foo6:barbaz2:goe`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}

func TestEncodePointer(t *testing.T) {
	// b := true
	// s := "well"
	// i := 42
	// tcs := []encodeTestCase{
	// 	{
	// 		&map[string]string{},
	// 		[]byte{'d', 'e'},
	// 		nil,
	// 	},
	// 	{
	// 		&[]string{},
	// 		[]byte(`le`),
	// 		nil,
	// 	},
	// 	{
	// 		&b,
	// 		[]byte(`i1e`),
	// 		nil,
	// 	},
	// 	{
	// 		&s,
	// 		[]byte(`4:well`),
	// 		nil,
	// 	},
	// 	{
	// 		&i,
	// 		[]byte(`i42e`),
	// 		nil,
	// 	},
	// 	{
	// 		&[]*[]string{
	// 			&[]string{},
	// 			&[]string{},
	// 		},
	// 		[]byte(`llelee`),
	// 		nil,
	// 	},
	// }
	// testEncodeLoop(t, tcs)
}

func testEncodeLoop(t *testing.T, tcs []encodeTestCase) {
	t.Helper()

	for i, test := range tcs {
		var b bytes.Buffer
		err := NewEncoder(&b).Encode(test.in)

		if test.err == nil && err != nil {
			t.Fatalf("[test %d] unexpected err %v", i, err)
		}

		if test.err != nil && err == nil {
			t.Fatalf("[test %d] expected err", i)
		}

		output := b.Bytes()
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("[test %d] got %v expected: %v", i, test.out, string(output))
		}
	}
}

type myBoolType bool

var ErrJustAnError = errors.New("oops")

type errorMarshalType struct{}

// MarshalBencode implements Marshaler.MarshalBencode
func (emt errorMarshalType) MarshalBencode() ([]byte, error) {
	return nil, ErrJustAnError
}

// UnmarshalBencode implements Unmarshaler.UnmarshalBencode
func (emt errorMarshalType) UnmarshalBencode([]byte) error {
	return ErrJustAnError
}
