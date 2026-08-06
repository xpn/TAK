package connection

type FlagSet interface {
	StringVarP(*string, string, string, string, string)
	StringVar(*string, string, string, string)
}

// Options contains the connection settings shared by the TAK command-line tools.
type Options struct {
	Proxy       string
	ClientCert  string
	ClientKey   string
	ClusterName string
}

func (o *Options) AddProxyFlag(flags FlagSet) {
	flags.StringVarP(&o.Proxy, "proxy", "x", "", "Proxy address (host:port)")
}

func (o *Options) AddClientCredentialFlags(flags FlagSet) {
	flags.StringVarP(&o.ClientCert, "client-cert", "c", "", "Path to client certificate")
	flags.StringVarP(&o.ClientKey, "client-key", "k", "", "Path to client key")
}

func (o *Options) AddClusterNameFlag(flags FlagSet) {
	flags.StringVar(&o.ClusterName, "cluster-name", "", "Teleport cluster name")
}
