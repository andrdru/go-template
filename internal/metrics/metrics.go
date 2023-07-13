package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "template"
	subsystem = "example"
)

var (
	databases = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "databases",
			Help:      "databases query metrics",
			Buckets:   []float64{.001, .005, .01, .025, .05, .075, .1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"database", "name", "error"})
)

// HistogramObserverDB .
func HistogramObserverDB(database string, name string, errFunc func() string) prometheus.Observer {
	return databases.With(map[string]string{
		"database": database,
		"name":     name,
		"error":    errFunc(),
	})
}
