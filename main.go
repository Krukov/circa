package main

import (
	"circa/resolver"
	"circa/storages"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/common-nighthawk/go-figure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"circa/config"
	"circa/runner"
	"circa/server"
)

func main() {
	debug := flag.Bool("debug", false, "dev mode")
	jsonLogs := flag.Bool("json-out", false, "json logging")
	configPath := flag.String("config", "./circa.json", "Config")
	port := flag.String("port", "8000", "Listen port")
	// managePort := flag.String("manage-port", "", "Listen port")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if !*jsonLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	runner.RegisterMetrics()
	server.RegisterMetrics()
	storages.RegisterMetrics()

	figure.NewColorFigure("| CIRCA |", "cyberlarge", "yellow", true).Print()
	Resolver := resolver.NewResolver()

	log.Info().Str("config", *configPath).Msg("Loading... ")
	conf, err := config.NewConfigFromDSN(*configPath, Resolver)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't load config")
		return
	}
	Runner := runner.NewRunner(conf)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	circa := server.Run(cancel, Runner, *port)
	// var manageSrv *http.Server
	// if *managePort != "" {
	// 	manageSrv = manage.Run(runner, *managePort)
	// }
	<-done
	// if manageSrv != nil {
	// 	manageSrv.Shutdown(ctx)
	// }
	circa.Shutdown()
}
