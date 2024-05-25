package cmd

import (
	"log/slog"
)

type Args struct {
	Port        int    `arg:"-p,--port" help:"Port number" default:"8080"`
	Environment string `arg:"-e,--environment" help:"Environment" default:"debug" placeholder:"debug | release | test |"`

	Verbose    slog.Level `arg:"-v,--verbose" help:"Log level" default:"INFO" placeholder:"INFO | DEBUG | WARN | ERROR"`
	PrettyJson bool       `arg:"--pretty-json" help:"Return pretty JSON responses (CPU intensive)" default:"false"`

	AllowInsecure bool `arg:"--allow-insecure" help:"Allow insecure minecraft accounts to connect" default:"false"`

	RateLimit  int   `arg:"--rate-limit" help:"Rate limit per second" default:"5"`
	RateBurst  int   `arg:"--rate-burst" help:"Rate burst" default:"10"`
	RatePunish int64 `arg:"--rate-punish" help:"Rate punish duration" default:"10"`
}
