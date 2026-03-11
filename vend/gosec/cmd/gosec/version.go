package main

// Version is the build version
var Version string

// GitTag is the git tag of the build
var GitTag string

// BuildDate is the date when the build was created
var BuildDate string

// prepareVersionInfo sets some runtime version when the version value
// was not injected by the build into the binary (e.g. go get).
// This returns currently "(devel)" but not an effective version until
// https://github.com/golang/go/issues/29814 gets resolved.
func prepareVersionInfo() {
	if Version == "" {
		// bi, _ := debug.ReadBuildInfo()
		// Version = bi.Main.Version
		// TODO use the debug information when it will provide more details
		// It seems to panic with Go 1.13.
		Version = "dev"
	}
}
