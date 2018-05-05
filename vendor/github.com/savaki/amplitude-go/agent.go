package amplitude

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	DefaultQueueSize = 250
	ApiEndpoint      = "https://api.amplitude.com/httpapi"
)

type Event struct {
	UserId             string                 `json:"user_id,omitempty"`
	DeviceId           string                 `json:"device_id,omitempty"`
	EventType          string                 `json:"event_type,omitempty"`
	Time               time.Time              `json:"-"`
	TimeInMillis       int64                  `json:"timestamp,omitempty"`
	EventProperties    map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
	AppVersion         string                 `json:"app_version,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	OSName             string                 `json:"os_name,omitempty"`
	OSVersion          string                 `json:"os_version,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	DeviceType         string                 `json:"device_type,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	Country            string                 `json:"country,omitempty"`
	Region             string                 `json:"region,omitempty"`
	City               string                 `json:"city,omitempty"`
	DMA                string                 `json:"dma,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Revenue            float64                `json:"revenu,omitempty"`
	Lat                float64                `json:"lat,omitempty"`
	Lon                float64                `json:"lon,omitempty"`
	Ip                 string                 `json:"ip,omitempty"`
	IDFA               string                 `json:"idfa,omitempty"`
	ADID               string                 `json:"adid,omitempty"`
}

type Client struct {
	cancel        func()
	ctx           context.Context
	apiKey        string
	ch            chan Event
	flush         chan chan struct{}
	queueSize     int
	interval      time.Duration
	onPublishFunc func(status int, err error)
}

func New(apiKey string, options ...Option) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		cancel:        cancel,
		ctx:           ctx,
		apiKey:        apiKey,
		ch:            make(chan Event, DefaultQueueSize),
		flush:         make(chan chan struct{}),
		queueSize:     DefaultQueueSize,
		interval:      time.Second * 15,
		onPublishFunc: func(status int, err error) {},
	}

	for _, opt := range options {
		opt(client)
	}

	go client.start()

	return client
}

func (c *Client) Publish(e Event) error {
	if !e.Time.IsZero() {
		e.TimeInMillis = e.Time.UnixNano() / int64(time.Millisecond)
	}

	select {
	case c.ch <- e:
		return nil
	default:
		return fmt.Errorf("Unable to send event, queue is full.  Use a larger queue size or create more workers.")
	}
}

func (c *Client) Event(e map[string]interface{}) error {
	if _, ok := e["user_id"]; ok {
		return fmt.Errorf("missing required parameter: user_id")
	}
	if _, ok := e["event_type"]; ok {
		return fmt.Errorf("missing required parameter: event_type")
	}

	return c.Publish(Event{
		UserId:    fmt.Sprintf("%v", e["user_id"]),
		EventType: fmt.Sprintf("%v", e["event_type"]),
	})
}

func (c *Client) start() {
	timer := time.NewTimer(c.interval)

	bufferSize := 256
	buffer := make([]Event, bufferSize)
	index := 0

	for {
		timer.Reset(c.interval)

		select {
		case <-c.ctx.Done():
			return

		case <-timer.C:
			if index > 0 {
				c.publish(buffer[0:index])
				index = 0
			}

		case v := <-c.ch:
			buffer[index] = v
			index++
			if index == bufferSize {
				c.publish(buffer[0:index])
				index = 0
			}

		case v := <-c.flush:
			if index > 0 {
				c.publish(buffer[0:index])
				index = 0
			}
			v <- struct{}{}
		}
	}
}

func (c *Client) publish(events []Event) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("api_key", c.apiKey)
	params.Set("event", string(data))

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := ctxhttp.PostForm(ctx, http.DefaultClient, ApiEndpoint, params)
	if resp != nil {
		defer resp.Body.Close()
	}
	c.onPublishFunc(resp.StatusCode, err)

	return err
}

func (c *Client) Flush() {
	ch := make(chan struct{})
	defer close(ch)

	c.flush <- ch
	<-ch
}

func (c *Client) Close() {
	c.Flush()
	c.cancel()
}
