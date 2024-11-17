package cmd

import (
	"log/slog"
)

type Args struct {
	Environment   string     `arg:"--environment,env:ENVIRONMENT" help:"Staging environment" default:"debug" placeholder:"debug | release | test |"`
	LogLevel      slog.Level `arg:"--log-level,env:LOG_LEVEL" help:"Log level" default:"INFO" placeholder:"INFO | DEBUG | WARN | ERROR"`
	AllowInsecure bool       `arg:"--allow-insecure,env:ALLOW_INSECURE" help:"Allow insecure minecraft accounts to authenticate" default:"false"`
}
