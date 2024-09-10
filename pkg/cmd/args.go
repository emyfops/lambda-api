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

	RateLimit    int `arg:"--rate-limit" help:"Maximum number of requests allowed in the duration" default:"5"`
	RateDuration int `arg:"--rate-duration" help:"Time frame in milliseconds in which requests can be made until n is reached" default:"10000"`
	RateBurst    int `arg:"--rate-burst" help:"How many requests can be handled simultaneously from one IP" default:"2"`
}
