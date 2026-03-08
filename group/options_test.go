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
		WithFactory(wantFactory),
		WithClient[any](wantClientOps...),
	}, &o)

	asserterror.EqualDeep(t, o.Factory, wantFactory)
	asserterror.EqualDeep(t, o.ClientOps, wantClientOps)
}
