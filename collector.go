package main

import (
	"sync"

	"github.com/mre/edgecast"
	"github.com/prometheus/client_golang/prometheus"
)

// EdgecastInterface to be used for logging and instrumenting middleware
type EdgecastInterface interface {
	Bandwidth(int) (*edgecast.BandwidthData, error)
	Connections(int) (*edgecast.ConnectionData, error)
	CacheStatus(int) (*edgecast.CacheStatusData, error)
	StatusCodes(int) (*edgecast.StatusCodeData, error)
}

// EdgecastCollector needs an edgecast client that implements the given interface to fetch metrics from edgecast API
type EdgecastCollector struct {
	ec EdgecastInterface
}

const (
	// NAMESPACE declaration for all exposed metrics in Prometheus
	NAMESPACE = "Edgecast"
)

var (
	// Platforms maps all possible media-types/platforms to it's IDs used in a request
	Platforms = map[int]string{
		2:  "flash",
		3:  "http_large",
		8:  "http_small",
		14: "adn",
	}

	// Prepared Description of all fetchable metrics
	bandwidth = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "bandwidth_bps"), "Current amount of bandwidth usage per platform (bits per second).", []string{"platform"}, nil,
	)
	cachestatus = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "cachestatus"), "Breakdown of the cache statuses currently being returned for requests to CDN account.", []string{"platform", "CacheStatus"}, nil,
	)
	connections = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "connections"), "Total active connections per second per platform.", []string{"platform"}, nil,
	)
	statuscodes = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "statuscodes"), "Breakdown of the HTTP status codes currently being returned for requests to CDN account.", []string{"platform", "StatusCode"}, nil,
	)
)

// NewEdgecastCollector constructs a new EdgecastCollector using a given edgecast-client that implements the EdgecastInterface
func NewEdgecastCollector(client *EdgecastInterface) *EdgecastCollector {
	return &EdgecastCollector{ec: *client}
}

// Describe describes all exported metrics
//- implements function of interface prometheus.Collector
func (col EdgecastCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- bandwidth
	ch <- cachestatus
	ch <- connections
	ch <- statuscodes
}

// Collect is called by Prometheus Server
// - concurrently fetches metrics for all possible platforms and exposes them in Prometheus format
//- implements function of interface prometheus.Collector
func (col EdgecastCollector) Collect(ch chan<- prometheus.Metric) {
	var collectWaitGroup sync.WaitGroup
	for p := range Platforms { // for each possible platform concurrently
		collectWaitGroup.Add(1)
		go col.metrics(ch, &collectWaitGroup, p) // fetch all possible metrics concurrently
	}
	collectWaitGroup.Wait()
}

// metrics() concurrently fetches all possible metric types for a given platform
func (col EdgecastCollector) metrics(ch chan<- prometheus.Metric, collectWaitgroup *sync.WaitGroup, platform int) {
	var metricsWaitGroup sync.WaitGroup
	metricsWaitGroup.Add(4) // 4 goroutines per platform for the 4 possible metric types
	go col.bandwidth(ch, &metricsWaitGroup, platform)
	go col.connections(ch, &metricsWaitGroup, platform)
	go col.cachestatus(ch, &metricsWaitGroup, platform)
	go col.statuscodes(ch, &metricsWaitGroup, platform)
	metricsWaitGroup.Wait() // wait for metric-fetching to finish
	collectWaitgroup.Done() // DONE fetching and exposing metrics for this platform
}

// bandwidth() fetches bandwidth metrics from API and pushes them to the channel as a new prometheus const metric
func (col EdgecastCollector) bandwidth(ch chan<- prometheus.Metric, metricsWaitGroup *sync.WaitGroup, platform int) {
	defer metricsWaitGroup.Done()

	bw, err := col.ec.Bandwidth(platform)
	if err == nil {
		bwBps := bw.Bps
		bwPlatform := Platforms[bw.Platform]
		ch <- prometheus.MustNewConstMetric(bandwidth, prometheus.GaugeValue, bwBps, []string{bwPlatform}...)
	}
}

// connections() fetches connection metrics from API and pushes them to the channel as a new prometheus const metric
func (col EdgecastCollector) connections(ch chan<- prometheus.Metric, metricsWaitGroup *sync.WaitGroup, platform int) {
	defer metricsWaitGroup.Done()

	con, err := col.ec.Connections(platform)
	if err == nil {
		conCon := con.Connections
		conPlatform := Platforms[con.Platform]
		ch <- prometheus.MustNewConstMetric(connections, prometheus.GaugeValue, conCon, []string{conPlatform}...)
	}
}

// cachestatus() fetches cachestatus metrics from API and pushes them to the channel as a new prometheus const metric
func (col EdgecastCollector) cachestatus(ch chan<- prometheus.Metric, metricsWaitGroup *sync.WaitGroup, platform int) {
	defer metricsWaitGroup.Done()

	cs, err := col.ec.CacheStatus(platform)
	if err == nil {
		csList := *cs
		var val float64
		var labelVals []string
		for c := range csList {
			val = float64(csList[c].Connections)
			labelVals = []string{Platforms[platform], csList[c].CacheStatus}
			ch <- prometheus.MustNewConstMetric(cachestatus, prometheus.GaugeValue, val, labelVals...)
		}

	}

}

// statuscodes() fetches statuscodes metrics from API and pushes them to the channel as a new prometheus const metric
func (col EdgecastCollector) statuscodes(ch chan<- prometheus.Metric, metricsWaitGroup *sync.WaitGroup, platform int) {
	defer metricsWaitGroup.Done()

	sc, err := col.ec.StatusCodes(platform)
	if err == nil {
		scList := *sc
		var val float64
		var labelVals []string
		for s := range scList {
			val = float64(scList[s].Connections)
			labelVals = []string{Platforms[platform], scList[s].StatusCode}
			ch <- prometheus.MustNewConstMetric(statuscodes, prometheus.GaugeValue, val, labelVals...)
		}
	}
}
