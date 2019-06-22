package bencode

import (
	"bytes"
	"testing"
)

type encodeTestCase struct {
	in  interface{}
	out []byte
	err error
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

func TestEncodeMarshaler(t *testing.T) {
	tcs := []encodeTestCase{
		{
			myBoolType(true),
			[]byte(`1:y`),
			nil,
		},
		{
			myBoolType(true),
			[]byte(`1:y`),
			nil,
		},
		{
			myBoolType(false),
			[]byte(`1:n`),
			nil,
		},
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

func TestEncodeInt(t *testing.T) {
	tcs := []encodeTestCase{
		{
			1,
			[]byte(`i1e`),
			nil,
		},
		{
			-1,
			[]byte(`i-1e`),
			nil,
		},
		{
			0,
			[]byte(`i0e`),
			nil,
		},
		{
			43210,
			[]byte(`i43210e`),
			nil,
		},
		{
			int(1),
			[]byte(`i1e`),
			nil,
		},
		{
			int8(2),
			[]byte(`i2e`),
			nil,
		},
		{
			int16(3),
			[]byte(`i3e`),
			nil,
		},
		{
			int32(4),
			[]byte(`i4e`),
			nil,
		},
		{
			int64(5),
			[]byte(`i5e`),
			nil,
		},
		{
			uint(6),
			[]byte(`i6e`),
			nil,
		},
		{
			uint8(7),
			[]byte(`i7e`),
			nil,
		},
		{
			uint16(8),
			[]byte(`i8e`),
			nil,
		},
		{
			uint32(9),
			[]byte(`i9e`),
			nil,
		},
		{
			uint64(10),
			[]byte(`i10e`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
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
			// TODO
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
	b := true
	s := "well"
	i := 42
	tcs := []encodeTestCase{
		{
			&map[string]string{},
			[]byte{'d', 'e'},
			nil,
		},
		{
			&[]string{},
			[]byte(`le`),
			nil,
		},
		{
			&b,
			[]byte(`i1e`),
			nil,
		},
		{
			&s,
			[]byte(`4:well`),
			nil,
		},
		{
			&i,
			[]byte(`i42e`),
			nil,
		},
		{
			&[]*[]string{
				&[]string{},
				&[]string{},
			},
			[]byte(`llelee`),
			nil,
		},
	}
	testEncodeLoop(t, tcs)
}
