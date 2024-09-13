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

	RateLimit    int `env:"RATE_LIMIT" arg:"--rate-limit" help:"Maximum number of requests allowed in the duration" default:"5"`
	RateDuration int `env:"RATE_DURATION" arg:"--rate-duration" help:"Time frame in milliseconds in which requests can be made until n is reached" default:"10000"`
	RateBurst    int `env:"RATE_BURST" arg:"--rate-burst" help:"How many requests can be handled simultaneously from one IP" default:"2"`
}
