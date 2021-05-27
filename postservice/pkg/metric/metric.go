package metric

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

var (
	hits    *prometheus.CounterVec
	errors  *prometheus.CounterVec
	Timings *prometheus.HistogramVec
)

func Destroy() {
	prometheus.Unregister(hits)
	prometheus.Unregister(errors)
	prometheus.Unregister(Timings)
}

func New() {
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})
	errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "errors",
	}, []string{"error"})

	Timings = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "timings",
		Buckets: []float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.1,
			0.15, 0.2, 0.4, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 15, 20},
	}, []string{"method", "URL"})
	prometheus.MustRegister(hits, errors, Timings)
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
