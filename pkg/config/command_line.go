package config

import "flag"

func ReadFromCommandLine(cfg *Config) {
	flag.Parse()
	paths := flag.Args()
	if len(paths) != 0 {
		cfg.Paths = paths
	}
}
