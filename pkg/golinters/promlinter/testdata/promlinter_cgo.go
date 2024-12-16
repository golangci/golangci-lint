//golangcitest:args -Epromlinter
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var (
	_ = promauto.NewCounterVec(
		prometheus.CounterOpts{ // want `Metric: test_metric_name Error: counter metrics should have "_total" suffix`
			Name: "test_metric_name",
			Help: "test help text",
		}, []string{},
	)

	_ = promauto.NewCounterVec(
		prometheus.CounterOpts{ // want "Metric: test_metric_total Error: no help text"
			Name: "test_metric_total",
		}, []string{},
	)

	_ = promauto.NewCounterVec(
		prometheus.CounterOpts{ // want `Metric: metric_type_in_name_counter_total Error: metric name should not include type 'counter'`
			Name: "metric_type_in_name_counter_total",
			Help: "foo",
		}, []string{},
	)

	_ = prometheus.NewHistogram(prometheus.HistogramOpts{ // want `Metric: test_duration_milliseconds Error: use base unit "seconds" instead of "milliseconds"`
		Name: "test_duration_milliseconds",
		Help: "",
	})
)
