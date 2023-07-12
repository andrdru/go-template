package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/andrdru/go-template/cmd/app"
)

type (
	flags struct {
		script *string
		isHelp *bool
	}
)

const (
	// serviceName name of service
	// redefine here or with ldflags
	// go build -ldflags="-X 'main.serviceName=my_service'"
	serviceName = "service"

	scriptTest = "script-test"
)

func main() {
	f := initFlags()
	if *f.isHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logger := initLogger()

	var (
		code int
	)

	switch *f.script {
	case "":
		code = app.Run(logger)
	default:
		logger.Error().Msgf("unknown script: %s", *f.script)
		os.Exit(1)

	case scriptTest:
		//code = somescript.Run(logger)
	}

	os.Exit(code)
}

func initFlags() (fv flags) {
	fv.isHelp = flag.Bool("help", false, "Print help and exit")
	fv.script = flag.String("script", "", "Run in script mode. One of: test-script")

	flag.Parse()
	return fv
}

func initLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	return zerolog.New(os.Stdout).
		With().Timestamp().
		Str("service", serviceName).
		Logger()
}
