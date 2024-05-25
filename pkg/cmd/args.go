package cmd

import "net/netip"

type Args struct {
	Host netip.Addr `arg:"-h,--host" help:"Host address, supports v4 and v6" default:"127.0.0.1"`
	Port int        `arg:"-p,--port" help:"Port number" default:"8080"`

	Verbose string `arg:"-v,--verbose" help:"Log level" default:"info" placeholder:"INFO | DEBUG | WARN | ERROR"`

	AllowInsecure bool `arg:"--allow-insecure" help:"Allow insecure minecraft accounts to connect" default:"false"`
}

var DefaultArgs = &Args{
	Host:          netip.MustParseAddr("127.0.0.1"),
	Port:          8080,
	Verbose:       "info",
	AllowInsecure: false,
}
