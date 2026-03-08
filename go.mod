module github.com/cmd-stream/cmd-stream-go

go 1.23.0

replace github.com/cmd-stream/testkit-go => ../testkit-go

require (
	github.com/cmd-stream/core-go v0.0.0-20260308140214-8371a5716599
	github.com/cmd-stream/delegate-go v0.0.0-20260308140853-81d810b4c662
	github.com/cmd-stream/handler-go v0.0.0-20260308162944-3e190d687853
	github.com/cmd-stream/testkit-go v0.0.0-20251122130859-27b372d2b32f
	github.com/cmd-stream/transport-go v0.0.0-20260308162028-cda69b948c47
	github.com/mus-format/mus-stream-go v0.8.0
	github.com/ymz-ncnk/assert v0.0.0-20260108210721-155bc9aa4282
	github.com/ymz-ncnk/mok v0.2.2
)

require (
	github.com/mus-format/common-go v0.0.0-20260225152706-590b1bf7cb37 // indirect
	github.com/mus-format/dts-stream-go v0.10.0 // indirect
	github.com/ymz-ncnk/jointwork-go v0.0.0-20240428103805-1ee224bde88a // indirect
	github.com/ymz-ncnk/multierr-go v0.0.0-20230813140901-5e9302c2e02a // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
)
