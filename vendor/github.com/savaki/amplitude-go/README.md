amplitude-go
============

Amplitude client for Go. For additional documentation visit https://amplitude.com/docs or view the godocs.

## Installation

	$ go get github.com/savaki/amplitude-go

## Examples

### Basic Client

Full example of a simple event tracker.

```go
	apiKey := os.Getenv("AMPLITUDE_API_KEY")
	client := amplitude.New(apiKey)
	client.Publish(amplitude.Event{
		UserId:    "123",
		EventType: "sample",
	})
```

