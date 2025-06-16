//golangcitest:args -Eworkflowcheck
package testdata

import (
	"crypto/rand"
	"fmt"
	"log"
	mathrand "math/rand"
	"net/http"
	"os"
	"time"
)

var BadVar time.Time // want BadVar:"declared non-deterministic"

func AccessesStdout() { // want AccessesStdout:"accesses non-deterministic var os.Stdout"
	os.Stdout.Write([]byte("Hello"))
}

func AccessesStdoutTransitively() { // want AccessesStdoutTransitively:"calls non-deterministic function testdata.AccessesStdout"
	AccessesStdout()
}

func CallsOtherStdoutCall() { // want CallsOtherStdoutCall:"calls non-deterministic function fmt.Println"
	fmt.Println()
}

func AccessesBadVar() { // want AccessesBadVar:"accesses non-deterministic var testdata.BadVar"
	BadVar.Day()
}

func AccessesCryptoRandom() { // want AccessesCryptoRandom:"accesses non-deterministic var crypto/rand.Reader"
	rand.Reader.Read(nil)
}

func AccessesCryptoRandomTransitively() { // want AccessesCryptoRandomTransitively:"calls non-deterministic function crypto/rand.Read"
	rand.Read(nil)
}

func CallsTime() { // want CallsTime:"calls non-deterministic function time.Now"
	time.Now()
}

func CallsTimeTransitively() { // want CallsTimeTransitively:"calls non-deterministic function testdata.CallsTime"
	CallsTime()
}

func CallsOtherTimeCall() { // want CallsOtherTimeCall:"calls non-deterministic function time.Until"
	time.Until(time.Time{})
}

func MultipleCalls() { // want MultipleCalls:"calls non-deterministic function time.Now, calls non-deterministic function testdata.CallsTime"
	time.Now()
	CallsTime()
}

func CallsLog() { // want CallsLog:"calls non-deterministic function log.Println"
	log.Println()
}

func CallsMathRandom() { // want CallsMathRandom:"calls non-deterministic function math/rand.Int"
	mathrand.Int()
}

func CallsHTTP() { // want CallsHTTP:"calls non-deterministic function net/http.Get"
	http.Get("http://example.com")
}

func NotSafeFmtCall() { // want NotSafeFmtCall:"calls non-deterministic function time.Now"
	fmt.Sprintf("%d", time.Now())
}

func SafeOnlyDoesAddition() {
	return 2 + 2
}
