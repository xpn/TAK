module xpnsec.com/reverse-tunnel/v2

go 1.25.6

require (
	github.com/fatih/color v1.18.0
	github.com/gravitational/teleport/api v0.0.0-20260116142645-4b284116367d
	github.com/spf13/cobra v1.10.2
	golang.org/x/crypto v0.47.0
	xpnsec.com/shared/v2 v2.0.0
)

replace xpnsec.com/shared/v2 => ../shared

require (
	github.com/beevik/etree v1.5.0 // indirect
	github.com/charlievieth/strcase v0.0.5 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/go-piv/piv-go/v2 v2.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gravitational/trace v1.5.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/russellhaering/gosaml2 v0.10.0 // indirect
	github.com/russellhaering/goxmldsig v1.5.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
