package herbbucketsgo

import (
	"bytes"
	"io"
)

type LenReader interface {
	io.Reader
	Len() int
}

func ByteReader(data []byte) LenReader {
	return bytes.NewBuffer(data)
}

type PlainLenReader struct {
	io.Reader
	Length int
}

func NewReader(r io.Reader, length int) LenReader {
	return &PlainLenReader{
		Reader: r,
		Length: length,
	}
}
func (r *PlainLenReader) Len() int {
	return r.Length
}
