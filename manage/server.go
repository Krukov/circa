package manage

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"circa/handler"
)

func Run(h *handler.Runner, port string) *http.Server {
	runnerHandlers := newRunnerHandler(h)
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/api/handlers", runnerHandlers.GetAllHandlers)
	http.HandleFunc("/api/route", runnerHandlers.GetHandlers)
	manageSrv := http.Server{Addr: fmt.Sprintf(":%s", port)}
	go func() {
		log.Info().Str("port", port).Msg("Start manage server")
		if err := manageSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn().Err(err).Msg("Can't start manage server")
		}
	}()
	return &manageSrv
}
