//args: -Epromlinter
package testdata

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	_ = promauto.NewCounterVec(
		prometheus.CounterOpts{ // ERROR `Metric: test_metric_name Error: counter metrics should have "_total" suffix`
			Name: "test_metric_name",
			Help: "test help text",
		}, []string{},
	)

	_ = promauto.NewCounterVec(
		prometheus.CounterOpts{ // ERROR "Metric: test_metric_total Error: no help text"
			Name: "test_metric_total",
		}, []string{},
	)

	_ = promauto.NewCounterVec(
		prometheus.CounterOpts{ // ERROR `Metric: metric_type_in_name_counter_total Error: metric name should not include type 'counter'`
			Name: "metric_type_in_name_counter_total",
			Help: "foo",
		}, []string{},
	)

	_ = prometheus.NewHistogram(prometheus.HistogramOpts{ // ERROR `Metric: test_duration_milliseconds Error: use base unit "seconds" instead of "milliseconds"`
		Name: "test_duration_milliseconds",
		Help: "",
	})
)
