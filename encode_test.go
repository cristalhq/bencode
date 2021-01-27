package bencode

import (
	"testing"
)

type marshalTestCase struct {
	val     interface{}
	want    string
	wantErr bool
}

func TestMarshalInt(t *testing.T) {
	tcs := []marshalTestCase{
		{42, `i42e`, false},
		{-42, `i-42e`, false},
		{0, `i0e`, false},
		{43210, `i43210e`, false},
		{int(1), `i1e`, false},
		{int8(2), `i2e`, false},
		{int16(3), `i3e`, false},
		{int32(4), `i4e`, false},
		{int64(5), `i5e`, false},
		{uint(6), `i6e`, false},
		{uint8(7), `i7e`, false},
		{uint16(8), `i8e`, false},
		{uint32(9), `i9e`, false},
		{uint64(10), `i10e`, false},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalFloat(t *testing.T) {
	tcs := []marshalTestCase{
		{123.456, `i4638387860618067575e`, false},
		{-456.1234, `i-4576640212951153351e`, false},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalBool(t *testing.T) {
	tcs := []marshalTestCase{
		{true, `i1e`, false},
		{false, `i0e`, false},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalString(t *testing.T) {
	tcs := []marshalTestCase{
		{"", "0:", false},
		{[]byte(""), "0:", false},
		{[]byte("x"), "1:x", false},
		{[]byte("foobarbaz"), "9:foobarbaz", false},
		{"foobarbaz", "9:foobarbaz", false},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalByteSliceAsString(t *testing.T) {
	tcs := []marshalTestCase{
		{
			[]byte(nil),
			string([]byte{'0', ':'}),
			false,
		},
		{
			[]byte{},
			string([]byte{'0', ':'}),
			false,
		},
		{
			[]byte(`test`),
			string([]byte(`4:test`)),
			false,
		},
		{
			[]byte{0, 1, 2, 3},
			string([]byte{'4', ':', 0, 1, 2, 3}),
			false,
		},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalByteArrayAsString(t *testing.T) {
	tcs := []marshalTestCase{
		{[...]byte{}, string([]byte(`0:`)), false},
		{[...]byte{'x', 'y'}, string([]byte{'2', ':', 'x', 'y'}), false},
		{[...]byte{0, 1, 2}, string([]byte{'3', ':', 0, 1, 2}), false},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalSlice(t *testing.T) {
	tcs := []marshalTestCase{
		{[]interface{}{}, `le`, false},
		{[]interface{}{1, 2, 3}, `li1ei2ei3ee`, false},
		{[]interface{}{"foo", "bar", "baz"}, `l3:foo3:bar3:baze`, false},
		{[]interface{}{1, "foo", 2, "bar"}, `li1e3:fooi2e3:bare`, false},
		{[]interface{}{1, []interface{}{"bar"}}, `li1el3:baree`, false},
		{[]string(nil), `le`, false},
		{[]string{}, `le`, false},
		{[]string{"foo"}, `l3:fooe`, false},
		{[]string{"foo", "barbaz"}, `l3:foo6:barbaze`, false},
		{[]string{"foo", "barbaz", "go"}, `l3:foo6:barbaz2:goe`, false},
		{(*[]interface{})(nil), ``, false},
		{[]interface{}{"foo", 20}, `l3:fooi20ee`, false},
		{[]interface{}{90, 20}, `li90ei20ee`, false},
		{[]interface{}{[]interface{}{"foo", "bar"}, 20}, `ll3:foo3:barei20ee`, false},
		{
			[]map[string]int{
				{"a": 0, "b": 1},
				{"c": 2, "d": 3},
			},
			`ld1:ai0e1:bi1eed1:ci2e1:di3eee`, false,
		},
		{
			[][]byte{
				[]byte{'0', '2', '4', '6', '8'},
				[]byte{'a', 'c', 'e'},
			},
			`l5:024683:acee`, false,
		},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalArray(t *testing.T) {
	tcs := []marshalTestCase{
		{[...]interface{}{}, `le`, false},
		{[...]interface{}{1, 2, 3}, `li1ei2ei3ee`, false},
		{[...]interface{}{"foo", "bar", "baz"}, `l3:foo3:bar3:baze`, false},
		{[...]interface{}{1, "foo", 2, "bar"}, `li1e3:fooi2e3:bare`, false},
		{[...]interface{}{1, [...]interface{}{"bar"}}, `li1el3:baree`, false},
		{[...]int{0, 1, 2}, "li0ei1ei2ee", false},
		{[...]float32{10, 20, 30}, "li1092616192ei1101004800ei1106247680ee", false},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalMap(t *testing.T) {
	tcs := []marshalTestCase{
		{map[string]string{}, `de`, false},
		{map[string]int{"1": 2, "4": 5}, `d1:1i2e1:4i5ee`, false},
		{
			map[string]int{"1": 1, "3": 3, "123": 123},
			"d1:1i1e3:123i123e1:3i3ee", false,
		},
		{
			map[string]string{
				"publisher":          "bob",
				"publisher-webpage":  "www.example.com",
				"publisher.location": "home",
			},
			"d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee",
			false,
		},
		{
			map[string]interface{}{"1": "one"},
			"d1:13:onee", false,
		},
		{
			map[string]interface{}{"1": "one", "two": "2"},
			"d1:13:one3:two1:2e", false,
		},
		{
			map[string]interface{}{"1": func() {}},
			"d1:13:one3:two1:2e", true,
		},
		{
			map[string][]int{
				"a": {0, 1},
				"b": {2, 3},
			},
			`d1:ali0ei1ee1:bli2ei3eee`, false,
		},
	}
	testLoopMarshal(t, tcs)
}

func TestMarshalPointer(t *testing.T) {
	b := true
	s := "well"
	i := 42

	tcs := []marshalTestCase{
		{&map[string]string{}, "de", false},
		{&[]string{}, "le", false},
		{&b, "i1e", false},
		{&s, "4:well", false},
		{&i, "i42e", false},
		{
			&[]*[]string{
				&[]string{},
				&[]string{},
			},
			"llelee", false,
		},
	}
	testLoopMarshal(t, tcs)
}

func testLoopMarshal(t *testing.T, tcs []marshalTestCase) {
	t.Helper()

	for i, test := range tcs {
		buf, err := Marshal(test.val)
		if err != nil {
			if test.wantErr {
				continue
			}
			t.Fatalf("[test %d] unexpected err %v", i+1, err)
		}

		got := string(buf)
		if want := string(test.want); got != want {
			t.Fatalf("[test %d] got %v want: %v", i+1, got, want)
		}
	}
}

var marshalBenchData = map[string]interface{}{
	"announce": ("udp://tracker.publicbt.com:80/announce"),
	"announce-list": []interface{}{
		[]interface{}{("udp://tracker.publicbt.com:80/announce")},
		[]interface{}{[]byte("udp://tracker.openbittorrent.com:80/announce")},
		[]interface{}{
			"udp://tracker.openbittorrent.com:80/announce",
			"udp://tracker.openbittorrent.com:80/announce",
		},
	},
	"comment": []byte("Debian CD from cdimage.debian.org"),
	"info": map[string]interface{}{
		"name":         []byte("debian-8.8.0-arm64-netinst.iso"),
		"length":       170917888,
		"piece length": 262144,
	},
}

func Benchmark_Marshal(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, err := Marshal(marshalBenchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_MarshalTo(b *testing.B) {
	dst := make([]byte, 0, 1<<12)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, err := MarshalTo(dst, marshalBenchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
