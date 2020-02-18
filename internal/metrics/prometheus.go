package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MonService is a struct to care out about metrics server
type MonService struct {
	Port string
}

// ServeMetrics run the client for Prometheus
func (m *MonService) ServeMetrics() {
	go func() {
		err := http.ListenAndServe(":"+m.Port, promhttp.Handler())
		if err != nil {
			log.Fatal("Unable to run the metrics server", err)
		}
	}()
}
