package cmd

import (
	"log/slog"
)

type Args struct {
	Port        int    `env:"PORT" arg:"-p,--port" help:"Port number" default:"80"`
	Environment string `env:"ENVIRONMENT" arg:"-e,--environment" help:"Environment" default:"debug" placeholder:"debug | release | test |"`

	Verbose    slog.Level `env:"VERBOSE" arg:"-v,--verbose" help:"Log level" default:"INFO" placeholder:"INFO | DEBUG | WARN | ERROR"`
	PrettyJson bool       `env:"PRETTY_JSON" arg:"--pretty-json" help:"Return pretty JSON responses (CPU intensive)" default:"false"`

	AllowInsecure bool `env:"ALLOW_INSECURE" arg:"--allow-insecure" help:"Allow insecure minecraft accounts to connect" default:"false"`
}
