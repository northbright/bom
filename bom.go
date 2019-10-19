package bom

import (
	"bufio"
	"bytes"
	"io"
)

// Encoding is interface type represents BOM encoding.
// String() implements the fmt.Stringer interface.
// Bytes() will output the bytes of BOM.
type Encoding interface {
	String() string
	Bytes() []byte
}

// encoding is a trivial implmentation of Encoding.
type encoding struct {
	name string
	hex  string
}

// String is the implementation for fmt.Stringer interface.
// It will return the name of BOM encoding.
func (enc *encoding) String() string {
	return enc.name
}

// Bytes returns the bytes of BOM.
func (enc *encoding) Bytes() []byte {
	return []byte(enc.hex)
}

var (
	// UTF32BE is the BOM of UTF-32, big-endian.
	UTF32BE = &encoding{"UTF-32, big-endian", "\x00\x00\xFE\xFF"}
	// UTF32LE is the BOM of UTF-32,little-endian.
	UTF32LE = &encoding{"UTF-32, little-endian", "\xFF\xFE\x00\x00"}
	// UTF8 is the BOM of UTF-8.
	UTF8 = &encoding{"UTF-8", "\xEF\xBB\xBF"}
	// UTF16BE is the BOM of UTF-16,big-endian.
	UTF16BE = &encoding{"UTF-16, big-endian", "\xFE\xFF"}
	// UTF16LE is the BOM of UTF-16,little-endian.
	UTF16LE = &encoding{"UTF-16, little-endian", "\xFF\xFE"}
	// NoBOM is the default BOM if no BOM is detected.
	NoBOM = &encoding{"No BOM", ""}

	encodings = []*encoding{
		UTF32BE,
		UTF32LE,
		UTF8,
		UTF16BE,
		UTF16LE,
		NoBOM,
	}
)

// DetectEncoding detects the BOM encoding by given buffer.
func DetectEncoding(buf []byte) Encoding {
	for _, enc := range encodings {
		if bytes.HasPrefix(buf, enc.Bytes()) {
			return enc
		}
	}
	return NoBOM
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
				return NoBOM, reader, err
			}
			// err == EOF, try to Peek(i - 1) again.
			continue
		}

		enc := DetectEncoding(buf)
		// No BOM detected.
		if enc == NoBOM {
			return NoBOM, reader, nil
		}

		// Make reader advance n bytes which n = len(BOM) to skip BOM.
		buf = make([]byte, len(enc.Bytes()))
		_, err = reader.Read(buf)

		return enc, reader, err
	}

	// No BOM detected.
	return NoBOM, reader, nil
}
