package metric

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)
const ctxKeyStartReqTime uint8 = 5

type PrometheusMetric struct {
	hits *prometheus.CounterVec
	errors *prometheus.CounterVec
	timings *prometheus.CounterVec
}

func New() *PrometheusMetric{
	metrics := &PrometheusMetric{
		hits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hits",
		}, []string{"status", "path"}),
		errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "errors",
		}, []string{"error"}),
		timings: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "timings",
		}, []string{"method", "URL", "time"}),
	}
	prometheus.MustRegister(metrics.hits, metrics.errors, metrics.timings)
	return metrics
}


func (m *PrometheusMetric)CrateRequestTiming(ctx context.Context, r* http.Request){
	timeStart := ctx.Value(ctxKeyStartReqTime).(time.Time)
	m.timings.WithLabelValues(r.Method, r.URL.String(), time.Since(timeStart).String()).Inc()
}


func (m *PrometheusMetric)CrateRequestHits(status int, r* http.Request){
	m.timings.WithLabelValues(strconv.Itoa(status), r.URL.Path).Inc()
}


func (m *PrometheusMetric)CrateRequestError(err error){
	oldErr := errors.Unwrap(err)
	m.errors.WithLabelValues(oldErr.Error()).Inc()
}

func (m *PrometheusMetric)Test(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		print(5)
		next.ServeHTTP(w, r)
	})
}