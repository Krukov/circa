package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/common-nighthawk/go-figure"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"circa/config"
	"circa/handler"
	"circa/server"
)


func main()  {
	debug := flag.Bool("debug", false, "dev mode")
	jsonLogs := flag.Bool("json-out", false, "json logging")
	configFilePath := flag.String("config", "./config.json", "Config path")
	port := flag.String("port", "8000", "Listen port")
	managePort := flag.String("manage-port", "9000", "Listen port")
	flag.Parse()


	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if !*jsonLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	figure.NewColorFigure("| CIRCA |", "cyberlarge", "yellow", true).Print()

	r := handler.NewRunner(server.MakeRequest)
	handler.RegisterMetrics()
	log.Info().Str("config", *configFilePath).Msg("Loading... ")
	err := config.AdjustJsonConfig(r, *configFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't load config file")
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	circa := server.Run(cancel, r, *port)
	http.Handle("/metrics", promhttp.Handler())
	manageSrv := http.Server{Addr: fmt.Sprintf(":%s", *managePort)}
	go func () {
		if err := manageSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn().Err(err).Msg("Can't start manage server")
		}
	}()
	<- done
	manageSrv.Shutdown(ctx)
	circa.Shutdown()


}



