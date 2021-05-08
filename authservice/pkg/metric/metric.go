package metric

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

const ctxKeyStartReqTime uint8 = 5

var (
	hits    *prometheus.CounterVec
	errors  *prometheus.CounterVec
	timings *prometheus.CounterVec
)

func Destroy() {
	prometheus.Unregister(hits)
	prometheus.Unregister(errors)
	prometheus.Unregister(timings)
}

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
	route := mux.CurrentRoute(r)
	path, _ := route.GetPathTemplate()
	timings.WithLabelValues(r.Method, path, time.Since(timeStart).String()).Inc()
}

func CrateRequestHits(status int, r *http.Request) {
	route := mux.CurrentRoute(r)
	path, _ := route.GetPathTemplate()
	hits.WithLabelValues(strconv.Itoa(status), path).Inc()
}

func CrateRequestError(err error) {
	if err != nil {
		errors.WithLabelValues(err.Error()).Inc()
	}
}
