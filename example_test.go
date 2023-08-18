package bencode_test

import (
	"bytes"
	"fmt"

	"github.com/cristalhq/bencode"
)

func ExampleMarshal() {
	// data to process, most of the types are supported
	var data interface{} = map[string]interface{}{
		"1":     42,
		"hello": "world",
		"foo":   []string{"bar", "baz"},
	}

	buf, err := bencode.Marshal(data)
	checkErr(err)
	fmt.Printf("marshaled: %s\n", string(buf))

	// or via Encoder:
	w := &bytes.Buffer{} // or any other io.Writer
	err = bencode.NewEncoder(w).Encode(data)
	checkErr(err)

	// Output:
	// marshaled: d1:1i42e3:fool3:bar3:baze5:hello5:worlde
}

func ExampleMarshalTo() {
	var data interface{} = map[string]interface{}{
		"1":     42,
		"hello": "world",
		"foo":   []string{"bar", "baz"},
	}

	buf := make([]byte, 0, 128)

	buf, err := bencode.MarshalTo(buf, data)
	checkErr(err)
	fmt.Printf("marshaled: %s\n", string(buf))

	// Output:
	// marshaled: d1:1i42e3:fool3:bar3:baze5:hello5:worlde
}

func ExampleUnmarshal() {
	var data interface{}

	buf := []byte("li1ei42ee")

	err := bencode.Unmarshal(buf, &data)
	checkErr(err)

	// or via Decoder:
	r := bytes.NewBufferString("li1ei42ee") // or any other io.Reader
	err = bencode.NewDecoder(r).Decode(&data)
	checkErr(err)

	fmt.Printf("unmarshaled: %v\n", data)

	// Output:
	// unmarshaled: [1 42]
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
