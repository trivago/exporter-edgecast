package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mre/edgecast"
)

/*
 * instrumentingMiddleware wraps a given EdgecastInterface an creates metrics for its invoked functions
 * The following metrics are created per function:
 * - requestCount:					incremented on every invocation of that function
 * - requestLatency:				time in seconds that function took from invocation to return
 * - requestLatencyDistribution:	histogram distribution of all invocations so far including phi-quantiles, total, sum
 */
type instrumentingMiddleware struct {
	requestCount               metrics.Counter   // positive/incrementing only value
	requestLatencyDistribution metrics.Histogram // bucket sampling
	requestLatency             metrics.Gauge     // positive and negative counting value
	next                       EdgecastInterface
}

func (mw instrumentingMiddleware) Bandwidth(platform int) (bandwidthData *edgecast.BandwidthData, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Bandwidth", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatencyDistribution.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.requestLatency.With(lvs...).Set(time.Since(begin).Seconds())
	}(time.Now())

	bandwidthData, err = mw.next.Bandwidth(platform) // hand request to logged service
	return
}

func (mw instrumentingMiddleware) Connections(platform int) (connectionData *edgecast.ConnectionData, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Connections", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatencyDistribution.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.requestLatency.With(lvs...).Set(time.Since(begin).Seconds())
	}(time.Now())

	connectionData, err = mw.next.Connections(platform) // hand request to logged service
	return
}

func (mw instrumentingMiddleware) CacheStatus(platform int) (cacheStatusData *edgecast.CacheStatusData, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CacheStatus", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatencyDistribution.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.requestLatency.With(lvs...).Set(time.Since(begin).Seconds())
	}(time.Now())

	cacheStatusData, err = mw.next.CacheStatus(platform) // hand request to logged service
	return
}

func (mw instrumentingMiddleware) StatusCodes(platform int) (statusCodeData *edgecast.StatusCodeData, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "StatusCodes", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatencyDistribution.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.requestLatency.With(lvs...).Set(time.Since(begin).Seconds())
	}(time.Now())

	statusCodeData, err = mw.next.StatusCodes(platform) // hand request to logged service
	return
}
