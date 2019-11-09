package main

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/doc"
)

func makeDoc() {
	if docFormat == "markdown" {
		if err := doc.ToMarkdown(address, docOutput); err != nil {
			log.Fatal("%v", err)
		}
	} else {
		log.Fatal("documentation format '%s' not supported", docFormat)
	}
}
