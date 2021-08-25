package main

import (
	"circa/manage"
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
	"circa/handler"
	"circa/server"
)

func main() {
	debug := flag.Bool("debug", false, "dev mode")
	jsonLogs := flag.Bool("json-out", false, "json logging")
	configFilePath := flag.String("config", "./config.json", "Config path")
	port := flag.String("port", "8000", "Listen port")
	managePort := flag.String("manage-port", "9991", "Listen port")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if !*jsonLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	handler.RegisterMetrics()
	server.RegisterMetrics()
	storages.RegisterMetrics()

	figure.NewColorFigure("| CIRCA |", "cyberlarge", "yellow", true).Print()

	runner := handler.NewRunner(server.MakeRequest)
	log.Info().Str("config", *configFilePath).Msg("Loading... ")
	err := config.AdjustJsonConfig(runner, *configFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't load config file")
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	circa := server.Run(cancel, runner, *port)
	manageSrv := manage.Run(runner, *managePort)
	<-done
	manageSrv.Shutdown(ctx)
	circa.Shutdown()
}
