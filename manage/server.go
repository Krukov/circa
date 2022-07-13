package manage

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"circa/config"
)

func Run(c *config.Config, port string) *http.Server {
	cm := newConfigManage(c)
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/api/storage/", cm.Storages)
	http.HandleFunc("/api/rules/", cm.Rules)
	// http.HandleFunc("/api/target/", cm.Target)

	manageSrv := http.Server{Addr: fmt.Sprintf(":%s", port)}
	go func() {
		log.Info().Str("port", port).Msg("Start manage server")
		if err := manageSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn().Err(err).Msg("Can't start manage server")
		}
	}()
	return &manageSrv
}
