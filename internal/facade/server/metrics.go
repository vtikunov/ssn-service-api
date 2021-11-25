package server

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ozonmp/ssn-service-api/internal/facade/config"
)

func createMetricsServer(cfg *config.Config) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return metricsServer
}
