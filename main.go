package main

import (
	// general
	"errors"
	"fmt"
	"net/http"
	"os"

	// Edgecast Client
	"github.com/mre/edgecast"

	// Prometheus for logging/metrics
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// go-kit
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

var (
	// user-defined environment-variables that handle access to the API
	accountID = os.Getenv("EDGECAST_ACCOUNT_ID")
	token     = os.Getenv("EDGECAST_TOKEN")
)

func main() {

	// check if account ID and token were properly specified using the environment variables
	if len(accountID) == 0 || len(token) == 0 {
		fmt.Println(errors.New("error: empty Account-ID or Token!\n-> Please specify using environment variables EDGECAST_ACCOUNT_ID and EDGECAST_TOKEN"))
		os.Exit(1)
	}

	// create new logger on Stderr
	logger := log.NewLogfmtLogger(os.Stderr)

	// Prometheus metrics settings for this service
	fieldKeys := []string{"method", "error"} // label names
	requestCount := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "Edgecast",
		Subsystem: "service_metrics",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "Edgecast",
		Subsystem: "service_metrics",
		Name:      "request_latency_distribution_seconds",
		Help:      "Total duration of requests in seconds.",
	}, fieldKeys)
	requestGauge := kitprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Namespace: "Edgecast",
		Subsystem: "service_metrics",
		Name:      "request_latency_seconds",
		Help:      "Duration of request in seconds.",
	}, fieldKeys)

	// create EdgecastClient that communicates with the Edgecast API
	var svc EdgecastInterface = edgecast.NewEdgecastClient(accountID, token)
	// attach logger to service
	svc = loggingMiddleware{logger, svc}
	// attach instrumenting middleware
	svc = instrumentingMiddleware{requestCount, requestLatency, requestGauge, svc}

	// create the prometheus collector that uses the EdgecastClient and register it to prometheus
	collector := NewEdgecastCollector(&svc)
	prometheus.MustRegister(collector)

	// connect handlers
	http.Handle("/metrics", promhttp.Handler())

	// set up logger and start service on port 80
	_ = logger.Log("msg", "HTTP", "addr", ":80")
	_ = logger.Log("err", http.ListenAndServe(":80", nil))
}
