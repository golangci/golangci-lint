package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	// NAME is the name of this notifier library as reported to the Rollbar API.
	NAME = "go-rollbar"

	// VERSION is the version number of this notifier library as reported to the
	// Rollbar API.
	VERSION = "0.4.0"

	// CRIT is the critical Rollbar severity level as reported to the Rollbar
	// API.
	CRIT = "critical"

	// ERR is the error Rollbar severity level as reported to the Rollbar API.
	ERR = "error"

	// WARN is the warning Rollbar severity level as reported to the Rollbar API.
	WARN = "warning"

	// INFO is the info Rollbar severity level as reported to the Rollbar API.
	INFO = "info"

	// DEBUG is the debug Rollbar severity level as reported to the Rollbar API.
	DEBUG = "debug"

	// FILTERED is the text that replaces all sensitive values in items sent to
	// the Rollbar API.
	FILTERED = "[FILTERED]"
)

var (
	// Token is the Rollbar access token under which all items will be reported.
	// If Token is blank, no errors will be reported to Rollbar.
	Token = ""

	// Environment is the environment under which all items will be reported.
	Environment = "development"

	// Platform is the platform reported for all Rollbar items. The default is
	// the running operating system (darwin, freebsd, linux, etc.) but it can
	// also be application specific (client, heroku, etc.).
	Platform = runtime.GOOS

	// Endpoint is the URL destination for all Rollbar item POST requests.
	Endpoint = "https://api.rollbar.com/api/1/item/"

	// Buffer is the maximum number of errors that will be queued for sending.
	// When the buffer is full, new errors are dropped on the floor until the API
	// can catch up.
	Buffer = 1000

	// FilterFields is a regular expression that matches field names that should
	// not be sent to Rollbar. Values for these fields are replaced with
	// "[FILTERED]".
	FilterFields = regexp.MustCompile("password|secret|token")

	// ErrorWriter is the destination for errors encountered while POSTing items
	// to Rollbar. By default, this is stderr. This can be nil.
	ErrorWriter io.Writer = os.Stderr

	// CodeVersion is the optional code version reported to the Rollbar API for
	// all items.
	CodeVersion = ""

	bodyChannel chan map[string]interface{}
	waitGroup   sync.WaitGroup
	postErrors  chan error
	nilErrTitle = "<nil>"
)

// Field is a custom data field used to report arbitrary data to the Rollbar
// API.
type Field struct {
	Name string
	Data interface{}
}

// -- Setup

func init() {
	bodyChannel = make(chan map[string]interface{}, Buffer)
	postErrors = make(chan error, Buffer)

	go func() {
		var err error
		for body := range bodyChannel {
			err = post(body)
			if err != nil {
				if len(postErrors) == cap(postErrors) {
					<-postErrors
				}
				postErrors <- err
			}
			waitGroup.Done()
		}
		close(postErrors)
	}()
}

// -- Error reporting

func Errorf(level string, format string, args ...interface{}) {
	ErrorWithStackSkip(level, fmt.Errorf(format, args...), 1)
}

// Error asynchronously sends an error to Rollbar with the given severity
// level. You can pass, optionally, custom Fields to be passed on to Rollbar.
func Error(level string, err error, fields ...*Field) {
	ErrorWithStackSkip(level, err, 1, fields...)
}

// ErrorWithStackSkip asynchronously sends an error to Rollbar with the given
// severity level and a given number of stack trace frames skipped. You can
// pass, optionally, custom Fields to be passed on to Rollbar.
func ErrorWithStackSkip(level string, err error, skip int, fields ...*Field) {
	stack := BuildStack(2 + skip)
	ErrorWithStack(level, err, stack, fields...)
}

// ErrorWithStack asynchronously sends and error to Rollbar with the given
// stacktrace and (optionally) custom Fields to be passed on to Rollbar.
func ErrorWithStack(level string, err error, stack Stack, fields ...*Field) {
	buildAndPushError(level, err, stack, fields...)
}

// RequestError asynchronously sends an error to Rollbar with the given
// severity level and request-specific information. You can pass, optionally,
// custom Fields to be passed on to Rollbar.
func RequestError(level string, r *http.Request, err error, fields ...*Field) {
	RequestErrorWithStackSkip(level, r, err, 1, fields...)
}

// RequestErrorWithStackSkip asynchronously sends an error to Rollbar with the
// given severity level and a given number of stack trace frames skipped, in
// addition to extra request-specific information. You can pass, optionally,
// custom Fields to be passed on to Rollbar.
func RequestErrorWithStackSkip(level string, r *http.Request, err error, skip int, fields ...*Field) {
	stack := BuildStack(2 + skip)
	RequestErrorWithStack(level, r, err, stack, fields...)
}

// RequestErrorWithStack asynchronously sends an error to Rollbar with the
// given severity level, request-specific information provided by the given
// http.Request, and a custom Stack. You You can pass, optionally, custom
// Fields to be passed on to Rollbar.
func RequestErrorWithStack(level string, r *http.Request, err error, stack Stack, fields ...*Field) {
	buildAndPushError(level, err, stack, append(fields, &Field{Name: "request", Data: errorRequest(r)})...)
}

func buildError(level string, err error, stack Stack, fields ...*Field) map[string]interface{} {
	title := nilErrTitle
	if err != nil {
		title = err.Error()
	}

	body := buildBody(level, title)
	data := body["data"].(map[string]interface{})
	errBody := errorBody(err, stack)
	data["body"] = errBody

	for _, field := range fields {
		data[field.Name] = field.Data
	}

	return body
}

func buildAndPushError(level string, err error, stack Stack, fields ...*Field) {
	push(buildError(level, err, stack, fields...))
}

// -- Message reporting

// Message asynchronously sends a message to Rollbar with the given severity
// level.
func Message(level string, msg string) {
	body := buildBody(level, msg)
	data := body["data"].(map[string]interface{})
	data["body"] = messageBody(msg)

	push(body)
}

// -- Misc.

// PostErrors returns a channel that receives all errors encountered while
// POSTing items to the Rollbar API.
func PostErrors() <-chan error {
	return postErrors
}

// Wait will block until the queue of errors / messages is empty. This allows
// you to ensure that errors / messages are sent to Rollbar before exiting an
// application.
func Wait() {
	waitGroup.Wait()
}

// Build the main JSON structure that will be sent to Rollbar with the
// appropriate metadata.
func buildBody(level, title string) map[string]interface{} {
	timestamp := time.Now().Unix()
	hostname, _ := os.Hostname()

	data := map[string]interface{}{
		"environment": Environment,
		"title":       title,
		"level":       level,
		"timestamp":   timestamp,
		"platform":    Platform,
		"language":    "go",
		"server": map[string]interface{}{
			"host": hostname,
		},
		"notifier": map[string]interface{}{
			"name":    NAME,
			"version": VERSION,
		},
	}
	if CodeVersion != "" {
		data["code_version"] = CodeVersion
	}

	return map[string]interface{}{
		"access_token": Token,
		"data":         data,
	}
}

// errorBody generates a Rollbar error body with a given stack trace.
func errorBody(err error, stack Stack) map[string]interface{} {
	message := nilErrTitle
	if err != nil {
		message = err.Error()
	}

	errBody := map[string]interface{}{
		"trace": map[string]interface{}{
			"frames": stack,
			"exception": map[string]interface{}{
				"class":   errorClass(err),
				"message": message,
			},
		},
	}
	return errBody
}

// errorRequest extracts details from a Request in a format that Rollbar
// accepts.
func errorRequest(r *http.Request) map[string]interface{} {
	cleanQuery := filterParams(r.URL.Query())

	return map[string]interface{}{
		"url":     r.URL.String(),
		"method":  r.Method,
		"headers": flattenValues(r.Header),

		// GET params
		"query_string": url.Values(cleanQuery).Encode(),
		"GET":          flattenValues(cleanQuery),

		// POST / PUT params
		"POST":    flattenValues(filterParams(r.Form)),
		"user_ip": r.RemoteAddr,
	}
}

// filterParams filters sensitive information like passwords from being sent to
// Rollbar.
func filterParams(values map[string][]string) map[string][]string {
	for key := range values {
		if FilterFields.Match([]byte(key)) {
			values[key] = []string{FILTERED}
		}
	}

	return values
}

func flattenValues(values map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range values {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
		}
	}

	return result
}

// Build a message inner-body for the given message string.
func messageBody(s string) map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]interface{}{
			"body": s,
		},
	}
}

func errorClass(err error) string {
	if err == nil {
		return nilErrTitle
	}

	class := reflect.TypeOf(err).String()
	if class == "" {
		return "panic"
	} else if class == "*errors.errorString" {
		checksum := adler32.Checksum([]byte(err.Error()))
		return fmt.Sprintf("{%x}", checksum)
	} else {
		return strings.TrimPrefix(class, "*")
	}
}

// -- POST handling

// Queue the given JSON body to be POSTed to Rollbar.
func push(body map[string]interface{}) {
	if len(bodyChannel) < Buffer {
		waitGroup.Add(1)
		bodyChannel <- body
	} else {
		stderr("buffer full, dropping error on the floor")
	}
}

// POST the given JSON body to Rollbar synchronously.
func post(body map[string]interface{}) error {
	if len(Token) == 0 {
		stderr("empty token")
		return nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		stderr("failed to encode payload: %s", err.Error())
		return err
	}

	resp, err := http.Post(Endpoint, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		stderr("POST failed: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		stderr("received response: %s", resp.Status)
		return ErrHTTPError(resp.StatusCode)
	}

	return nil
}

// -- stderr
func stderr(format string, args ...interface{}) {
	if ErrorWriter != nil {
		format = "Rollbar error: " + format + "\n"
		fmt.Fprintf(ErrorWriter, format, args...)
	}
}
