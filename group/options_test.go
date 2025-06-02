package group

import (
	"testing"

	cln "github.com/cmd-stream/cmd-stream-go/client"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options[any]{}
		wantFactory   = RoundRobinStrategyFactory[any]{}
		wantClientOps = []cln.SetOption{}
	)
	ApplyGroup([]SetOption[any]{
		WithFactory[any](wantFactory),
		WithClientOps[any](wantClientOps...),
	}, &o)

	asserterror.EqualDeep(o.Factory, wantFactory, t)
	asserterror.EqualDeep(o.ClientOps, wantClientOps, t)

}
