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
	} else if newUser != "" {
		addNewUser()
		return
	}

	compileViews := docOutput == ""

	if err := models.Setup(confFile, dataPath, usersPath, compileViews); err != nil {
		log.Fatal("%v", err)
	}

	if docOutput != "" {
		makeDoc()
		return
	}

	if err, server := api.Setup(); err != nil {
		log.Fatal("%v", err)
	} else {
		server.Run(address)
	}
}
