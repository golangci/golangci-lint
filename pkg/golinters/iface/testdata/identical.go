//golangcitest:args -Eiface
//golangcitest:config_path testdata/identical.yml
package testdata

type Pinger interface { // want "interface Pinger contains identical methods or type constraints from another interface, causing redundancy"
	Ping() error
}

type Healthcheck interface { // want "interface Healthcheck contains identical methods or type constraints from another interface, causing redundancy"
	Ping() error
}
