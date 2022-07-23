package manage

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"circa/config"
)

func Run(c *config.Config, port string) *http.Server {
	cm := newConfigManage(c)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/api/storage/", cm.Storages)
	router.HandleFunc("/api/rules/", cm.Rules)
	router.HandleFunc("/api/sync/", cm.Sync)
	
	logger := log.With().
		Str("port", port).
		Str("managment_api", "true").
		Logger()

	manageSrv := http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      logging(logger)(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Info().Str("port", port).Msg("Start manage server")
		if err := manageSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn().Err(err).Msg("Can't start manage server")
		}
	}()
	return &manageSrv
}

func logging(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Info().Str("method", r.Method).Str("path", r.URL.Path).Msg("#> Request")
			}()
			next.ServeHTTP(w, r)
		})
	}
}
