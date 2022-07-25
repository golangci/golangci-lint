//golangcitest:args -Ereusestdlibvars
package testdata

import (
	"net/http"
	_ "time"
)

func _200() {
	_ = 200 // ERROR `can use http.StatusOK instead "200"`
}

func _200_1() {
	var w http.ResponseWriter
	w.WriteHeader(200) // ERROR `can use http.StatusOK instead "200"`
}

func sunday() {
	_ = "Sunday" // ERROR `can use time.Sunday.String\(\) instead "Sunday"`
}

func monday() {
	_ = "Monday" // ERROR `can use time.Monday.String\(\) instead "Monday"`
}

func tuesday() {
	_ = "Tuesday" // ERROR `can use time.Tuesday.String\(\) instead "Tuesday"`
}

func wednesday() {
	_ = "Wednesday" // ERROR `can use time.Wednesday.String\(\) instead "Wednesday"`
}

func thursday() {
	_ = "Thursday" // ERROR `can use time.Thursday.String\(\) instead "Thursday"`
}

func friday() {
	_ = "Friday" // ERROR `can use time.Friday.String\(\) instead "Friday"`
}

func saturday() {
	_ = "Saturday" // ERROR `can use time.Saturday.String\(\) instead "Saturday"`
}

func january() {
	_ = "January" // ERROR `can use time.January.String\(\) instead "January"`
}

func february() {
	_ = "February" // ERROR `can use time.February.String\(\) instead "February"`
}

func march() {
	_ = "March" // ERROR `can use time.March.String\(\) instead "March"`
}

func april() {
	_ = "April" // ERROR `can use time.April.String\(\) instead "April"`
}

func may() {
	_ = "May" // ERROR `can use time.May.String\(\) instead "May"`
}

func june() {
	_ = "June" // ERROR `can use time.June.String\(\) instead "June"`
}

func july() {
	_ = "July" // ERROR `can use time.July.String\(\) instead "July"`
}

func august() {
	_ = "August" // ERROR `can use time.August.String\(\) instead "August"`
}

func september() {
	_ = "September" // ERROR `can use time.September.String\(\) instead "September"`
}

func october() {
	_ = "October" // ERROR `can use time.October.String\(\) instead "October"`
}

func november() {
	_ = "November" // ERROR `can use time.November.String\(\) instead "November"`
}

func december() {
	_ = "December" // ERROR `can use time.December.String\(\) instead "December"`
}
