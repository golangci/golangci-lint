//golangcitest:args -Eiface
//golangcitest:config_path testdata/duplicate.yml
package testdata

type Pinger interface { // want "interface Pinger contains duplicate methods or type constraints from another interface, causing redundancy"
	Ping() error
}

type Healthcheck interface { // want "interface Healthcheck contains duplicate methods or type constraints from another interface, causing redundancy"
	Ping() error
}
