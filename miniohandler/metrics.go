package miniohandler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var httpMetrics = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: "golang",
		Name:      "my_histogram",
		Help:      "This is my histogram",
	})

// PrometheusHandler
func PrometheusHandler() {
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "5000"
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", metricsPort), mux)

	log.Info().Msgf("metrics server is listening on port: %s", metricsPort)
}
