//golangcitest:args -Egofactory
//golangcitest:config_path configs/go_factory_package_globs_only.yml
package gofactory

import (
	"net/http"
	"net/url"
)

var (
	nestedGlobalRequest    = http.Request{}
	nestedGlobalRequestPtr = &http.Request{}

	blockedGlobalURL    = url.URL{}  // want `Use factory for url.URL`
	blockedGlobalURLPtr = &url.URL{} // want `Use factory for url.URL`
)

func Blocked() {
	_ = http.Request{}
	_ = &http.Request{}

	_ = url.URL{}  // want `Use factory for url.URL`
	_ = &url.URL{} // want `Use factory for url.URL`
}

type URL struct {
	Scheme      string
	Opaque      string
	User        *url.Userinfo
	Host        string
	Path        string
	RawPath     string
	OmitHost    bool
	ForceQuery  bool
	RawQuery    string
	Fragment    string
	RawFragment string
}

func Casting() {
	_ = url.URL(URL{}) // want `Use factory for url.URL`

	uPtr, _ := url.Parse("")
	u := *uPtr
	_ = URL(u)
}
