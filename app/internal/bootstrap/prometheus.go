package bootstrap

import "github.com/prometheus/client_golang/prometheus"

type Prometheus struct {
	HttpRequestTotal                 *prometheus.CounterVec
	HttpRequestDuration              *prometheus.HistogramVec
	HttpRequestSize                  *prometheus.HistogramVec
	HttpResponseSize                 *prometheus.HistogramVec
	WebsocketMessageHandlingDuration *prometheus.HistogramVec
	WebsocketConnectionsCount        prometheus.Gauge
}

func NewPrometheus() *Prometheus {
	return &Prometheus{
		HttpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "The total number of HTTP requests.",
			},
			[]string{"method", "endpoint", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "The HTTP request latencies in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HttpRequestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "The HTTP request sizes in bytes.",
				Buckets: prometheus.ExponentialBuckets(100, 10, 6),
			},
			[]string{"method", "endpoint"},
		),
		HttpResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "The HTTP response sizes in bytes.",
				Buckets: prometheus.ExponentialBuckets(100, 10, 6),
			},
			[]string{"method", "endpoint"},
		),
		WebsocketMessageHandlingDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "websocket_message_handling_duration_seconds",
				Help:    "The websocket message handling latencies in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"endpoint", "message_type", "user_id"},
		),
		WebsocketConnectionsCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "websocket_connections_count",
				Help: "The websocket connections count.",
			},
		),
	}

}

func (p *Prometheus) Init() {
	prometheus.MustRegister(p.HttpRequestTotal, p.HttpRequestDuration, p.HttpRequestSize, p.HttpResponseSize, p.WebsocketMessageHandlingDuration, p.WebsocketConnectionsCount)
}
