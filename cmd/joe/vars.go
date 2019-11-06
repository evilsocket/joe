package main

import (
	"flag"
	"github.com/evilsocket/islazy/log"
)

var (
	debug    = false
	ver      = false
	address  = "0.0.0.0:8080"
	confFile = "/etc/joe/joe.conf"
	dataPath = "/etc/joe/queries"
)

func init() {
	flag.BoolVar(&ver, "version", ver, "Print version and exit.")
	flag.BoolVar(&debug, "debug", debug, "Enable debug logs.")
	flag.StringVar(&log.Output, "log", log.Output, "Log file path or empty for standard output.")
	flag.StringVar(&address, "address", address, "API address.")
	flag.StringVar(&confFile, "conf", confFile, "Configuration file.")
	flag.StringVar(&dataPath, "data", dataPath, "Data path.")
}
