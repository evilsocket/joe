package main

import (
	"flag"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/api"
	"github.com/evilsocket/joe/models"
)

func main() {
	flag.Parse()

	setup()
	defer cleanup()

	if ver {
		fmt.Println(api.Version)
		return
	}

	if err := models.Setup(dataPath, confFile); err != nil {
		log.Fatal("%v", err)
	}

	if err, server := api.Setup(); err != nil {
		log.Fatal("%v", err)
	} else {
		server.Run(address)
	}
}
