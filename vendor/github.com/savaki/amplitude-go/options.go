package amplitude

import "time"

type Option func(*Client)

func Interval(v time.Duration) Option {
	return func(c *Client) {
		c.interval = v
	}
}

func OnPublishFunc(fn func(status int, err error)) Option {
	return func(c *Client) {
		c.onPublishFunc = fn
	}
}
