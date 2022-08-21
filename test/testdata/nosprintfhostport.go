//golangcitest:args -Enosprintfhostport
package testdata

import (
	"fmt"
	"net"
)

func _() {

	_ = fmt.Sprintf("postgres://%s:%s@127.0.0.1/%s", "foo", "bar", "baz")

	_ = fmt.Sprintf("http://api.%s/foo", "example.com")

	_ = fmt.Sprintf("http://api.%s:6443/foo", "example.com")

	_ = fmt.Sprintf("http://%s/foo", net.JoinHostPort("foo", "80"))

	_ = fmt.Sprintf("9invalidscheme://%s:%d", "myHost", 70)

	_ = fmt.Sprintf("gopher://%s/foo", net.JoinHostPort("foo", "80"))

	_ = fmt.Sprintf("telnet+ssl://%s/foo", net.JoinHostPort("foo", "80"))

	_ = fmt.Sprintf("http://%s/foo:bar", net.JoinHostPort("foo", "80"))

	_ = fmt.Sprintf("http://user:password@%s/foo:bar", net.JoinHostPort("foo", "80"))

	_ = fmt.Sprintf("http://example.com:9211")

	_ = fmt.Sprintf("gopher://%s:%d", "myHost", 70) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	_ = fmt.Sprintf("telnet+ssl://%s:%d", "myHost", 23) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	_ = fmt.Sprintf("weird3.6://%s:%d", "myHost", 23) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	_ = fmt.Sprintf("https://user@%s:%d", "myHost", 8443) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	_ = fmt.Sprintf("postgres://%s:%s@%s:5050/%s", "foo", "bar", "baz", "qux") // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	_ = fmt.Sprintf("https://%s:%d", "myHost", 8443) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	_ = fmt.Sprintf("https://%s:9211", "myHost") // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"

	ip := "fd00::1"
	_ = fmt.Sprintf("http://%s:1936/healthz", ip) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"
}
