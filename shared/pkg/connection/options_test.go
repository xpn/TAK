package connection

import "testing"

type testFlagSet struct {
	values map[string]string
}

func (f testFlagSet) StringVarP(value *string, name, _ string, defaultValue, _ string) {
	*value = f.value(name, defaultValue)
}

func (f testFlagSet) StringVar(value *string, name, defaultValue, _ string) {
	*value = f.value(name, defaultValue)
}

func (f testFlagSet) value(name, defaultValue string) string {
	if value, ok := f.values[name]; ok {
		return value
	}
	return defaultValue
}

func TestSharedFlags(t *testing.T) {
	var options Options
	flags := testFlagSet{values: map[string]string{
		"proxy":        "proxy.example.com:443",
		"client-cert":  "client.crt",
		"client-key":   "client.key",
		"cluster-name": "example.com",
	}}
	options.AddProxyFlag(flags)
	options.AddClientCredentialFlags(flags)
	options.AddClusterNameFlag(flags)

	if options.Proxy != "proxy.example.com:443" {
		t.Fatalf("unexpected proxy: %q", options.Proxy)
	}
	if options.ClientCert != "client.crt" {
		t.Fatalf("unexpected client certificate: %q", options.ClientCert)
	}
	if options.ClientKey != "client.key" {
		t.Fatalf("unexpected client key: %q", options.ClientKey)
	}
	if options.ClusterName != "example.com" {
		t.Fatalf("unexpected cluster name: %q", options.ClusterName)
	}
}
