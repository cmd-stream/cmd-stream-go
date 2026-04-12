package delegate

import (
	com "github.com/mus-format/common-go"
	mus "github.com/mus-format/mus-stream-go"
	bslopts "github.com/mus-format/mus-stream-go/options/byte_slice"
	"github.com/mus-format/mus-stream-go/ord"
)

// DefaultServerInfoMaxLength is the default maximum length of the server info.
const DefaultServerInfoMaxLength = 1024

// var byteSliceMUS = ord.NewSliceSer(raw.Byte)

// ServerInfo allows the client to identify a compatible server.
type ServerInfo []byte

// ServerInfoValidMUS is a ServerInfo MUS serializer with length validation.
var ServerInfoValidMUS = newServerInfoValidMUS()

func newServerInfoValidMUS() serverInfoValidMUS {
	return serverInfoValidMUS{
		byteSliceMUS: ord.NewValidByteSliceSer(
			bslopts.WithLenValidator(com.ValidatorFn[int](
				func(i int) error {
					if i > DefaultServerInfoMaxLength {
						return ErrTooLargeServerInfo
					}
					return nil
				},
			)),
		),
	}
}

type serverInfoValidMUS struct {
	byteSliceMUS mus.Serializer[[]byte]
}

// Marshal encodes ServerInfo to the writer.
func (s serverInfoValidMUS) Marshal(info ServerInfo, w mus.Writer) (n int,
	err error,
) {
	return s.byteSliceMUS.Marshal(info, w)
}

// Unmarshal decodes ServerInfo from the reader.
func (s serverInfoValidMUS) Unmarshal(r mus.Reader) (info ServerInfo, n int,
	err error,
) {
	return s.byteSliceMUS.Unmarshal(r)
}

// Size returns the encoded size of the ServerInfo.
func (s serverInfoValidMUS) Size(info ServerInfo) (size int) {
	return s.byteSliceMUS.Size(info)
}

// Skip skips the ServerInfo in the reader.
func (s serverInfoValidMUS) Skip(r mus.Reader) (n int, err error) {
	return s.byteSliceMUS.Skip(r)
}
