package bom_test

import (
	"bytes"
	"io"
	"log"

	"github.com/northbright/bom"
)

func ExampleSkip() {
	type TestData struct {
		str         string
		expectedEnc bom.Encoding
	}

	testDatas := []TestData{
		TestData{"\x00\x00\xFE\xFF", bom.UTF32BE},         // UTF-32, big-endian
		TestData{"\xFF\xFE\x00\x00", bom.UTF32LE},         // UTF-32, little-endian
		TestData{"\xFE\xFF", bom.UTF16BE},                 // UTF-16, big-endian
		TestData{"\xFF\xFE", bom.UTF16LE},                 // UTF-16, little-endian
		TestData{"\xEF\xBB\xBF", bom.UTF8},                // UTF-8
		TestData{"No BOM", bom.NoBOM},                     // UTF-8 string without BOM
		TestData{"\x00", bom.NoBOM},                       // only 1 byte
		TestData{"\xFE\xFF\x67\x0D\x52\xA1", bom.UTF16BE}, // UTF-16 big-endian string: "服务"
		TestData{"\xFF\xFE\x0D\x67\xA1\x52", bom.UTF16LE}, // UTF-16 little-endian string: "服务"
		TestData{"\xEF\xBB\xBFHello World!", bom.UTF8},    // UTF-8 string with BOM
	}

	for _, data := range testDatas {
		// Create a bytes.Buffer(io.Reader).
		b := bytes.NewBufferString(data.str)

		// Create a bufio.Reader and try to skip BOM.
		enc, r, err := bom.Skip(b)
		if err != nil {
			log.Printf("bom.Skip() error: %v\n", err)
			return
		}

		if enc != data.expectedEnc {
			log.Printf("encoding: %s is not expected: %s\n", enc, data.expectedEnc)
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
			log.Printf("encoding: %s, read %v bytes\n", enc, n)
		} else {
			log.Printf("encoding: %s, read %v bytes, buf: %X\n", enc, n, buf[:n])
		}
	}

	// Output:
}
