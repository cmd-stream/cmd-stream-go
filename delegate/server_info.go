package delegate

import (
	mus "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
	"github.com/mus-format/mus-stream-go/raw"
)

var byteSliceMUS = ord.NewSliceSer(raw.Byte)

// ServerInfo allows the client to identify a compatible server.
type ServerInfo []byte

// ServerInfoMUS is a ServerInfo MUS serializer.
var ServerInfoMUS = serverInfoMUS{}

type serverInfoMUS struct{}

func (s serverInfoMUS) Marshal(info ServerInfo, w mus.Writer) (n int,
	err error,
) {
	return byteSliceMUS.Marshal(info, w)
}

func (s serverInfoMUS) Unmarshal(r mus.Reader) (info ServerInfo, n int,
	err error,
) {
	return byteSliceMUS.Unmarshal(r)
}

func (s serverInfoMUS) Size(info ServerInfo) (size int) {
	return byteSliceMUS.Size(info)
}

func (s serverInfoMUS) Skip(r mus.Reader) (n int, err error) {
	return byteSliceMUS.Skip(r)
}
