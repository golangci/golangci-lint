//golangcitest:args -Eresponsewriterlint
package testdata

import (
	"errors"
	"fmt"
	"net/http"
)

type notAResponseWriter struct{}

func (narw notAResponseWriter) Write(in []byte) (int, error) {
	// actually do nothing
	return 42, errors.New("no")
}

func (narw notAResponseWriter) WriteHeader(code int) {
	// also do nothing
}

func (narw notAResponseWriter) Header() http.Header {
	return http.Header{}
}

type rwlRandom struct{}

func rwlExampleOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("some header", "value")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`boys in the yard`))
}

func rwlExampleTwo(s string, r *http.Request) error {
	fmt.Printf("this is a thing")

	return nil
}

func rwlFakeWriter(w notAResponseWriter) {
	w.Header().Add("something", "other")
	w.WriteHeader(420)
	_, _ = w.Write([]byte(`fooled ya`))
}

func (b rwlRandom) method(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("some header", "value")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`boys in the yard`))
}

func (b *rwlRandom) methodPointer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("some header", "value")
	_, _ = w.Write([]byte(`boys in the yard`))
	w.WriteHeader(http.StatusOK) // want "function methodPointer: http.ResponseWriter.Write is called before http.ResponseWriter.WriteHeader. Headers are already sent, this has no effect."
}

func rwlExampleThree(bloe http.ResponseWriter, r *http.Request) {
	_, _ = bloe.Write([]byte(`hellyea`)) // want "function rwlExampleThree: Multiple calls to http.ResponseWriter.Write in the same function body. This is most probably a bug."

	bloe.WriteHeader(http.StatusBadRequest)          // want "function rwlExampleThree: Multiple calls to http.ResponseWriter.WriteHeader in the same function body. This is most probably a bug."
	_, _ = bloe.Write([]byte(`hellyelamdflmda`))     // want "function rwlExampleThree: Multiple calls to http.ResponseWriter.Write in the same function body. This is most probably a bug."
	bloe.WriteHeader(http.StatusInternalServerError) // want "function rwlExampleThree: Multiple calls to http.ResponseWriter.WriteHeader in the same function body. This is most probably a bug."

	bloe.Header().Set("help", "somebody")     // want "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.Write. This has no effect." "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.Write. This has no effect." "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.WriteHeader. This has no effect." "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.WriteHeader. This has no effect."
	bloe.Header().Set("dddd", "someboddaady") // want "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.Write. This has no effect." "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.Write. This has no effect." "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.WriteHeader. This has no effect." "function rwlExampleThree: http.ResponseWriter.Header called after calling http.ResponseWriter.WriteHeader. This has no effect."
}
