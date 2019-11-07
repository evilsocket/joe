package main

import (
	"flag"
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/api"
	"github.com/evilsocket/joe/models"
	"golang.org/x/crypto/ssh/terminal"
	"path"
	"syscall"
)

func main() {
	flag.Parse()

	setup()
	defer cleanup()

	if ver {
		fmt.Println(api.Version)
		return
	} else if newUser != "" {
		userFile := path.Join(usersPath, fmt.Sprintf("%s.yml", newUser))
		if fs.Exists(userFile) {
			fmt.Printf("%s already exists.\n", userFile)
			return
		}

		fmt.Print("Enter Password: ")
		raw, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}

		password := string(raw)
		user := models.User {
			Username: newUser,
			Password: password,
			TokenTTL: tokenTTL,
		}

		if err := models.SaveUser(user, userFile); err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Printf("Saved to %s\n", userFile)
		return
	}

	if err := models.Setup(confFile, dataPath, usersPath); err != nil {
		log.Fatal("%v", err)
	}

	if err, server := api.Setup(); err != nil {
		log.Fatal("%v", err)
	} else {
		server.Run(address)
	}
}
