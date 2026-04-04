package testkit

import (
	"encoding/json"

	"github.com/cmd-stream/cmd-stream-go/core"
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-stream-go/ord"
	"github.com/mus-format/mus-stream-go/typed"
)

func AsyncResult(seq core.Seq, result core.Result) core.AsyncResult {
	return core.AsyncResult{
		Seq:       seq,
		BytesRead: CalcResultSize(seq, result),
		Result:    result,
	}
}

func CalcCmdSize(seq core.Seq, cmd Cmd) (size int) {
	return cmdSize(seq, CmdDTM, cmd)
}

func CalcMultiCmdSize(seq core.Seq, cmd MultiCmd) (size int) {
	return cmdSize(seq, MultiCmdDTM, cmd)
}

func CalcResultSize(seq core.Seq, result core.Result) (size int) {
	size = core.SeqMUS.Size(seq)
	bs, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	return size + ord.ByteSlice.Size(bs)
}

func cmdSize[T any](seq core.Seq, dtm com.DTM, t T) (size int) {
	size = core.SeqMUS.Size(seq)
	size += typed.DTMSer.Size(dtm)
	bs, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return size + ord.ByteSlice.Size(bs)
}
