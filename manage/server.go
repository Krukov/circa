package manage

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func Run(port string) *http.Server {
	http.Handle("/metrics", promhttp.Handler())
	manageSrv := http.Server{Addr: fmt.Sprintf(":%s", port)}
	go func() {
		log.Info().Str("port", port).Msg("Start manage server")
		if err := manageSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn().Err(err).Msg("Can't start manage server")
		}
	}()
	return &manageSrv
}
