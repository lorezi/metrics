package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpMetrics struct {
	RequestTotal             *prometheus.CounterVec
	RequestDurationHistogram *prometheus.HistogramVec
}

func NewHttpMetrics() HttpMetrics {
	return HttpMetrics{
		RequestTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "total HTTP requests processed",
		}, []string{"code", "method"}),

		RequestDurationHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Seconds spent serving HTTP requests.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"code", "method"}),
	}
}

// InstrumentHandler instruments any HTTP handler for the request total and request duration metric

func InstrumentHandler(next http.HandlerFunc, _http HttpMetrics) http.HandlerFunc {
	return promhttp.InstrumentHandlerCounter(_http.RequestTotal, promhttp.InstrumentHandlerDuration(_http.RequestDurationHistogram, next))
}
