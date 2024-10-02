package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
)

type MetricsMiddleware struct {
	Prom *bootstrap.Prometheus
}

func NewMetricsMiddleware(prom *bootstrap.Prometheus) *MetricsMiddleware {
	return &MetricsMiddleware{
		Prom: prom,
	}
}

func (m *MetricsMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		route := r.URL.Path
		method := r.Method

		size := r.ContentLength
		m.Prom.HttpRequestSize.WithLabelValues(method, route).Observe(float64(size))

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(recorder, r)

		duration := time.Since(startTime).Seconds()
		m.Prom.HttpRequestDuration.WithLabelValues(method, route).Observe(duration)
		m.Prom.HttpRequestTotal.WithLabelValues(method, route, strconv.Itoa(recorder.statusCode)).Inc()

		m.Prom.HttpResponseSize.WithLabelValues(method, route).Observe(float64(recorder.size))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}
