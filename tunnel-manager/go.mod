module xpnsec.com/teleport-tunnel-manager/v2

go 1.25.5

require (
	github.com/cosmos/gogoproto v1.7.2
	github.com/gorilla/websocket v1.5.3
	github.com/gravitational/teleport/api v0.0.0-20260114180212-534806aeae0d
	github.com/spf13/cobra v1.10.2
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	xpnsec.com/shared/v2 v2.0.0
)

replace xpnsec.com/shared/v2 => ../shared

require (
	github.com/beevik/etree v1.5.0 // indirect
	github.com/charlievieth/strcase v0.0.5 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/go-piv/piv-go/v2 v2.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gravitational/trace v1.5.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/russellhaering/gosaml2 v0.10.0 // indirect
	github.com/russellhaering/goxmldsig v1.5.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/term v0.37.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
)
