package group

import (
	"testing"

	asserterror "github.com/ymz-ncnk/assert/error"

	cln "github.com/cmd-stream/cmd-stream-go/client"
)

type mockStrategyFactory[T any] struct {
	DispatchStrategyFactory[T]
}

func TestWithFactory(t *testing.T) {
	var (
		opts    = Options[any]{}
		factory = &mockStrategyFactory[any]{}
	)
	WithFactory[any](factory)(&opts)
	asserterror.Equal(t, opts.Factory, DispatchStrategyFactory[any](factory))
}

func TestWithReconnect(t *testing.T) {
	opts := Options[any]{}
	WithReconnect[any]()(&opts)
	asserterror.Equal(t, opts.Reconnect, true)
}

func TestWithClient(t *testing.T) {
	var (
		opts      = Options[any]{}
		clientOpt = cln.WithTransport()
	)
	WithClient[any](clientOpt)(&opts)
	asserterror.Equal(t, len(opts.ClientOpts), 1)
}

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options Options[any]
		wantErr bool
	}{
		{
			name:    "valid options",
			options: Options[any]{Factory: &mockStrategyFactory[any]{}},
			wantErr: false,
		},
		{
			name:    "nil factory",
			options: Options[any]{Factory: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asserterror.Equal(t, tt.options.Validate() != nil, tt.wantErr)
		})
	}
}

func TestApply(t *testing.T) {
	var (
		opts    = DefaultOptions[any]()
		factory = &mockStrategyFactory[any]{}
	)
	err := Apply(&opts, WithFactory[any](factory), WithReconnect[any](), nil)
	asserterror.Equal(t, err, nil)
	asserterror.Equal(t, opts.Factory, DispatchStrategyFactory[any](factory))
	asserterror.Equal(t, opts.Reconnect, true)
}

func TestApplyError(t *testing.T) {
	opts := Options[any]{}
	err := Apply(&opts, WithFactory[any](nil))
	asserterror.Equal(t, err.Error(), "factory is nil")
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions[any]()
	if opts.Factory == nil {
		t.Error("DefaultOptions should have a default factory")
	}
}
