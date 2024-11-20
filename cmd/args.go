package cmd

import (
	"github.com/alexflint/go-arg"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

var cliOptions Args
var dbOptions *redis.Options

type Args struct {
	Environment   string     `arg:"--environment,env:ENVIRONMENT" help:"Staging environment" default:"debug" placeholder:"debug | release | test |"`
	LogLevel      slog.Level `arg:"--log-level,env:LOG_LEVEL" help:"Log level" default:"INFO" placeholder:"INFO | DEBUG | WARN | ERROR"`
	AllowInsecure bool       `arg:"--allow-insecure,env:ALLOW_INSECURE" help:"Allow insecure minecraft accounts to authenticate" default:"false"`

	// Database options from redis.Options
	Addr     string `arg:"--db-addr,env:DB_ADDR" help:"Database address host:port" default:":6379"`
	Username string `arg:"--db-username,env:DB_USERNAME" help:"Database username" default:""`
	Password string `arg:"--db-password,env:DB_PASSWORD" help:"Database password" default:""`
	DB       int    `arg:"--db-index,env:DB_INDEX" help:"Database number" default:"0"`
}

func init() {
	arg.MustParse(&cliOptions)

	// Copy the CLI options to the database options
	dbOptions.Addr = cliOptions.Addr
	dbOptions.Username = cliOptions.Username
	dbOptions.Password = cliOptions.Password
	dbOptions.DB = cliOptions.DB
}

func Arguments() Args {
	return cliOptions
}

func RedisOptions() *redis.Options {
	return &*dbOptions // Prevent modification of the original struct
}
