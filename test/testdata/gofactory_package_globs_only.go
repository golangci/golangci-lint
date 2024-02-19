//golangcitest:args -Egofactory
//golangcitest:config_path configs/go_factory_package_globs_only.yml
package testdata

import (
	"net/http"
	neturl "net/url"
)

var (
	nestedGlobalRequest    = http.Request{}
	nestedGlobalRequestPtr = &http.Request{}

	blockedGlobalURL    = neturl.URL{}  // want `Use factory for url.URL`
	blockedGlobalURLPtr = &neturl.URL{} // want `Use factory for url.URL`
)

func Blocked() {
	_ = http.Request{}
	_ = &http.Request{}

	_ = neturl.URL{}  // want `Use factory for url.URL`
	_ = &neturl.URL{} // want `Use factory for url.URL`
}

type URL struct {
	Scheme      string
	Opaque      string
	User        *neturl.Userinfo
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
	_ = neturl.URL(URL{}) // want `Use factory for url.URL`

	uPtr, _ := neturl.Parse("")
	u := *uPtr
	_ = URL(u)
}
