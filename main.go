package main

import (
	"circa/config"
	"circa/handler"
	"circa/server"
	"flag"
	"github.com/common-nighthawk/go-figure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)


func main()  {
	debug := flag.Bool("debug", false, "dev mode")
	jsonLogs := flag.Bool("json-out", false, "json logging")
	configFilePath := flag.String("config", "./config.json", "Config path")
	port := flag.String("port", "8000", "Listen port")
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

	log.Info().Str("config", *configFilePath).Msg("Loading... ")
	err := config.AdjustJsonConfig(r, *configFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't load config file")
		return
	}

	server.Run(r, *port)
}



