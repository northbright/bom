package bom

import (
	"bufio"
	"bytes"
	"io"
)

// Encoding represents the BOM encoding.
type Encoding struct {
	// Name is the name of BOM encoding.
	Name string
	// Data contains BOM bytes.
	Data []byte
}

var (
	// Encodings represent the unicode BOM Encodings.
	Encodings = []Encoding{
		Encoding{"UTF-32,big-endian", []byte{0x00, 0x00, 0xFE, 0xFF}},
		Encoding{"UTF-32,little-endian", []byte{0xFF, 0xFE, 0x00, 0x00}},
		Encoding{"UTF-8", []byte{0xEF, 0xBB, 0xBF}},
		Encoding{"UTF-16,big-endian", []byte{0xFE, 0xFF}},
		Encoding{"UTF-16,little-endian", []byte{0xFF, 0xFE}},
	}
)

// DetectEncoding detects the BOM encoding by given buffer.
func DetectEncoding(buf []byte) Encoding {
	for _, enc := range Encodings {
		if bytes.HasPrefix(buf, enc.Data) {
			return enc
		}
	}
	return Encoding{"Unkown", []byte{}}
}

// Skip tries to detect BOM in original reader and return a new bufio.Reader with BOM skipped.
// After Skip is called, the original reader was advanced N bytes because of using bufio.Reader.
// N is the buffer size of bufio.Reader.
func Skip(r io.Reader) (Encoding, *bufio.Reader, error) {
	reader := bufio.NewReader(r)

	for i := 4; i >= 2; i-- {
		buf, err := reader.Peek(i)

		if err != nil {
			if err != io.EOF {
				return Encoding{"Unkown", []byte{}}, reader, err
			}
			// err == EOF, try to Peek(i - 1) again.
			continue
		}

		encoding := DetectEncoding(buf)
		// No BOM detected.
		l := len(encoding.Data)
		if l <= 0 {
			return encoding, reader, nil
		}

		// Make reader advance n bytes which n = len(BOM) to skip BOM.
		buf = make([]byte, len(encoding.Data))
		_, err = reader.Read(buf)

		return encoding, reader, err
	}

	// No BOM detected.
	return Encoding{"Unkown", []byte{}}, reader, nil
}
