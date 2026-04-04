package transport

import (
	"github.com/ymz-ncnk/mok"
)

type (
	ReaderReadFn     func(p []byte) (n int, err error)
	ReaderReadByteFn func() (byte, error)
)

type Reader struct {
	*mok.Mock
}

func NewReader() Reader {
	return Reader{Mock: mok.New("Reader")}
}

func (m Reader) RegisterRead(fn ReaderReadFn) Reader {
	m.Register("Read", fn)
	return m
}

func (m Reader) RegisterReadByte(fn ReaderReadByteFn) Reader {
	m.Register("ReadByte", fn)
	return m
}

func (m Reader) Read(p []byte) (n int, err error) {
	vals, err := m.Call("Read", p)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m Reader) ReadByte() (byte, error) {
	vals, err := m.Call("ReadByte")
	if err != nil {
		panic(err)
	}
	b := vals[0].(byte)
	err, _ = vals[1].(error)
	return b, err
}
