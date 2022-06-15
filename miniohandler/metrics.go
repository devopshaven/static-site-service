package miniohandler

import (
	"fmt"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
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

// NewMetricsHandler
func NewMetricsHandler() {
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "5000"
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/ready", ReadyHandler)
	mux.HandleFunc("/health", ReadyHandler)

	go func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", metricsPort), mux)
	}()

	log.Info().Msgf("metrics server is listening on port: %s", metricsPort)
}

func ReadyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Server", "DOH-Static-Site-Service")

	enc := jsoniter.NewEncoder(w)
	enc.Encode(map[string]interface{}{
		"success": true,
	})
}
