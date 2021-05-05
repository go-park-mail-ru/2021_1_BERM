package metric

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

const ctxKeyStartReqTime uint8 = 5

var (
	hits    *prometheus.CounterVec
	errors   *prometheus.CounterVec
	timings *prometheus.CounterVec
)

func New() {
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})
	errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "errors",
	}, []string{"error"})
	timings = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "timings",
	}, []string{"method", "URL", "time"})
	prometheus.MustRegister(hits, errors, timings)
}

func CrateRequestTiming(ctx context.Context, r *http.Request) {
	timeStart := ctx.Value(ctxKeyStartReqTime).(time.Time)
	timings.WithLabelValues(r.Method, r.URL.String(), time.Since(timeStart).String()).Inc()
}

func CrateRequestHits(status int, r *http.Request) {
	hits.WithLabelValues(strconv.Itoa(status), r.URL.Path).Inc()
}

func  CrateRequestError(err error) {
	if err != nil {
		errors.WithLabelValues(err.Error()).Inc()
	}
}
