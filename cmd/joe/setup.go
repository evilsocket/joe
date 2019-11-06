package main

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/models"
)

func setup() {
	if debug {
		log.Level = log.DEBUG
	} else {
		log.Level = log.INFO
	}
	log.OnFatal = log.ExitOnFatal
}

func cleanup() {
	models.Cleanup()
}
