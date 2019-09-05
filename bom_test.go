package bom_test

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/northbright/bom"
)

func ExampleSkip() {
	strs := []string{
		"\x00\x00\xFE\xFF", // UTF-32, big-endian
		"\xFF\xFE\x00\x00", // UTF-32, little-endian
		"\xFE\xFF",         // UTF-16, big-endian
		"\xFF\xFE",         // UTF-16, little-endian
		"\xEF\xBB\xBF",     // UTF-8
		"No BOM",           // UTF-8 string without BOM
		"\x00",             // only 1 byte
		"\xFE\xFF\x67\x0D\x52\xA1", // UTF-16 big-endian string: "服务"
		"\xFF\xFE\x0D\x67\xA1\x52", // UTF-16 little-endian string: "服务"
		"\xEF\xBB\xBFHello World!", // UTF-8 string with BOM
	}

	for _, str := range strs {
		// Create a bytes.Buffer(io.Reader).
		b := bytes.NewBufferString(str)

		// Create a bufio.Reader and try to skip BOM.
		encoding, r, err := bom.Skip(b)
		if err != nil {
			log.Printf("bom.Skip() error: %v\n", err)
			return
		}

		// Read data to buffer after skip BOM.
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Read() error: %v\n", err)
			return
		}

		if n == 0 {
			fmt.Printf("encoding:%v, read %v bytes\n", encoding.Name, n)
		} else {
			fmt.Printf("encoding:%v, read %v bytes, buf: %X\n", encoding.Name, n, buf[:n])
		}
	}

	// Output:
	//encoding:UTF-32,big-endian, read 0 bytes
	//encoding:UTF-32,little-endian, read 0 bytes
	//encoding:UTF-16,big-endian, read 0 bytes
	//encoding:UTF-16,little-endian, read 0 bytes
	//encoding:UTF-8, read 0 bytes
	//encoding:Unkown, read 6 bytes, buf: 4E6F20424F4D
	//encoding:Unkown, read 1 bytes, buf: 00
	//encoding:UTF-16,big-endian, read 4 bytes, buf: 670D52A1
	//encoding:UTF-16,little-endian, read 4 bytes, buf: 0D67A152
	//encoding:UTF-8, read 12 bytes, buf: 48656C6C6F20576F726C6421
}
