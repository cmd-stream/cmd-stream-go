package mock

import (
	"github.com/ymz-ncnk/mok"
)

type (
	WriterWriteByteFn   func(c byte) error
	WriterWriteFn       func(p []byte) (n int, err error)
	WriterWriteStringFn func(s string) (n int, err error)
	WriterFlushFn       func() error
)

type Writer struct {
	*mok.Mock
}

func NewWriter() Writer {
	return Writer{Mock: mok.New("Writer")}
}

func (m Writer) RegisterWriteByte(fn WriterWriteByteFn) Writer {
	m.Register("WriteByte", fn)
	return m
}

func (m Writer) RegisterWrite(fn WriterWriteFn) Writer {
	m.Register("Write", fn)
	return m
}

func (m Writer) RegisterWriteString(fn WriterWriteStringFn) Writer {
	m.Register("WriteString", fn)
	return m
}

func (m Writer) RegisterFlush(fn WriterFlushFn) Writer {
	m.Register("Flush", fn)
	return m
}

func (m Writer) WriteByte(c byte) error {
	vals, err := m.Call("WriteByte", c)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m Writer) Write(p []byte) (n int, err error) {
	vals, err := m.Call("Write", p)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m Writer) WriteString(s string) (n int, err error) {
	vals, err := m.Call("WriteString", s)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m Writer) Flush() error {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}
