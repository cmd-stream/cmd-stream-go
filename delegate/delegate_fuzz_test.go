package delegate

import (
	"bufio"
	"bytes"
	"testing"
)

func FuzzServerInfoUnmarshal(f *testing.F) {
	f.Add([]byte("test server info"))
	f.Add(make([]byte, DefaultServerInfoMaxLength+1))

	f.Fuzz(func(t *testing.T, data []byte) {
		reader := bufio.NewReader(bytes.NewReader(data))
		_, _, err := ServerInfoValidMUS.Unmarshal(reader)
		if err != nil {
			// Expected errors are fine, we just want to ensure no panics.
			return
		}
	})
}
