package bencode

import (
	"bytes"
	"reflect"
	"testing"
)

type unmarshalTestCase struct {
	input   string
	want    any
	wantErr bool
}

func TestUnmarshalInt(t *testing.T) {
	tcs := []unmarshalTestCase{
		{`i1e`, int64(1), false},
		{`i-1e`, -int64(1), false},
		{`i0e`, int64(0), false},
		{`i43210e`, int64(43210), false},
		{`i1e`, int64(1), false},
		{`i2e`, int64(2), false},
		{`i3e`, int64(3), false},
		{`i4e`, int64(4), false},
		{`i5e`, int64(5), false},
		{`i6e`, int64(6), false},
		{`i7e`, int64(7), false},
		{`i8e`, int64(8), false},
		{`i9e`, int64(9), false},
		{`i10e`, int64(10), false},
	}
	testLoopUnmarshal(t, tcs)
}

func TestUnmarshalFloat(t *testing.T) {
	tcs := []unmarshalTestCase{
		// we don't know is it a float, so decoding as int64
		{`i4638387860618067575e`, int64(4638387860618067575), false},
		{`i-4576640212951153351e`, int64(-4576640212951153351), false},
	}
	testLoopUnmarshal(t, tcs)
}

func TestUnmarshalBool(t *testing.T) {
	tcs := []unmarshalTestCase{
		// we don't know is it a bool, so decoding as int64
		{`i1e`, int64(1), false},
		{`i0e`, int64(0), false},
	}
	testLoopUnmarshal(t, tcs)
}

func TestUnmarshalString(t *testing.T) {
	tcs := []unmarshalTestCase{
		{"0:", []byte(""), false},
		{"1:x", []byte("x"), false},
		{"9:foobarbaz", []byte("foobarbaz"), false},
	}
	testLoopUnmarshal(t, tcs)
}

func TestUnmarshalByteSliceAsString(t *testing.T) {
	tcs := []unmarshalTestCase{
		{
			string([]byte{'0', ':'}),
			[]byte(""),
			false,
		},
		{
			string([]byte{'0', ':'}),
			[]byte{},
			false,
		},
		{
			string([]byte(`4:test`)),
			[]byte(`test`),
			false,
		},
		{
			string([]byte{'4', ':', 0, 1, 2, 3}),
			[]byte{0, 1, 2, 3},
			false,
		},
	}
	testLoopUnmarshal(t, tcs)
}

func TestUnmarshalSlice(t *testing.T) {
	tcs := []unmarshalTestCase{
		{
			`le`,
			[]any{},
			false,
		},
		{
			`li1ei2ei3ee`,
			[]any{int64(1), int64(2), int64(3)},
			false,
		},
		{
			`l3:foo3:bar3:baze`,
			[]any{[]byte("foo"), []byte("bar"), []byte("baz")},
			false,
		},
		{
			`li1e3:fooi2e3:bare`,
			[]any{int64(1), []byte("foo"), int64(2), []byte("bar")},
			false,
		},
		{
			`li1el3:baree`,
			[]any{int64(1), []any{[]byte("bar")}},
			false,
		},
		{
			`l3:fooe`,
			[]any{[]byte("foo")},
			false,
		},
		{
			`l3:foo6:barbaze`,
			[]any{[]byte("foo"), []byte("barbaz")},
			false,
		},
		{
			`l3:foo6:barbaz2:goe`,
			[]any{[]byte("foo"), []byte("barbaz"), []byte("go")},
			false,
		},
		{
			`l3:fooi20ee`,
			[]any{[]byte("foo"), int64(20)},
			false,
		},
		{
			`li90ei20ee`,
			[]any{int64(90), int64(20)},
			false,
		},
		{
			`ll3:foo3:barei20ee`,
			[]any{
				[]any{
					[]byte("foo"),
					[]byte("bar"),
				},
				int64(20),
			},
			false,
		},
		{
			`l5:024683:acee`,
			[]any{
				[]byte{'0', '2', '4', '6', '8'},
				[]byte{'a', 'c', 'e'},
			},
			false,
		},
		{
			`ld1:ai0e1:bi1eed1:ci2e1:di3eee`,
			[]any{
				map[string]any{
					"a": int64(0),
					"b": int64(1),
				},
				map[string]any{
					"c": int64(2),
					"d": int64(3),
				},
			},
			false,
		},
	}
	testLoopUnmarshal(t, tcs)
}

// func TestUnmarshalArray(t *testing.T) {
// 	tcs := []unmarshalTestCase{
// 		{
// 			`le`,
// 			[...]any{},
// 			false,
// 		},
// 		{
// 			`li1ei2ei3ee`,
// 			[...]any{1, 2, 3},
// 			false,
// 		},
// 		{
// 			[...]any{"foo", "bar", "baz"},
// 			`l3:foo3:bar3:baze`,
// 			false,
// 		},
// 		{
// 			[...]any{1, "foo", 2, "bar"},
// 			`li1e3:fooi2e3:bare`,
// 			false,
// 		},
// 		{
// 			[...]any{1, [...]any{"bar"}},
// 			`li1el3:baree`,
// 			false,
// 		},
// 		{
// 			[...]int{0, 1, 2},
// 			"li0ei1ei2ee",
// 			false,
// 		},
// 		{
// 			[...]float32{10, 20, 30},
// 			"li1092616192ei1101004800ei1106247680ee",
// 			false,
// 		},
// 	}
// 	testLoopUnmarshal(t, tcs)
// }

func TestUnmarshalMap(t *testing.T) {
	tcs := []unmarshalTestCase{
		{
			`de`,
			map[string]any{},
			false,
		},
		{
			`d1:1i2e1:4i5ee`,
			map[string]any{"1": int64(2), "4": int64(5)},
			false,
		},
		{
			"d1:1i1e3:123i123e1:3i3ee",
			map[string]any{
				"1":   int64(1),
				"3":   int64(3),
				"123": int64(123),
			},
			false,
		},
		{
			"d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee",
			map[string]any{
				"publisher":          []byte("bob"),
				"publisher-webpage":  []byte("www.example.com"),
				"publisher.location": []byte("home"),
			},
			false,
		},
		{
			"d1:13:onee",
			map[string]any{"1": []byte("one")},
			false,
		},
		{
			"d1:13:one3:two1:2e",
			map[string]any{
				"1":   []byte("one"),
				"two": []byte("2"),
			},
			false,
		},
		{
			`d1:ali0ei1ee1:bli2ei3eee`,
			map[string]any{
				"a": []any{int64(0), int64(1)},
				"b": []any{int64(2), int64(3)},
			},
			false,
		},
	}
	testLoopUnmarshal(t, tcs)
}

func testLoopUnmarshal(t *testing.T, tcs []unmarshalTestCase) {
	t.Helper()

	for i, test := range tcs {
		var got any
		err := Unmarshal([]byte(test.input), &got)
		if err != nil {
			if test.wantErr {
				continue
			}
			t.Fatalf("[test %d] unexpected err %v", i+1, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("[test %d] got %v want: %v", i+1, got, test.want)
		}
	}
}

var unmarshalBenchData = []byte("d4:infod6:lengthi170917888e12:piece lengthi262144e4:name30:debian-8.8.0-arm64-netinst.isoe8:announce38:udp://tracker.publicbt.com:80/announce13:announce-listll38:udp://tracker.publicbt.com:80/announceel44:udp://tracker.openbittorrent.com:80/announceee7:comment33:Debian CD from cdimage.debian.orge")

func Benchmark_Unmarshal(b *testing.B) {
	var res any
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		err := NewDecodeBytes(unmarshalBenchData).Decode(&res)
		if err != nil {
			b.Fatal(err)
		}
		if res == nil {
			b.Fatal("is nil")
		}
	}
}

func Benchmark_UnmarshalReader(b *testing.B) {
	var res any
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r := bytes.NewReader(unmarshalBenchData)
		err := NewDecoder(r).Decode(&res)
		if err != nil {
			b.Fatal(err)
		}
		if res == nil {
			b.Fatal("is nil")
		}
	}
}
