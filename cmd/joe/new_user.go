package main

import (
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/joe/models"
	"golang.org/x/crypto/ssh/terminal"
	"path"
	"fmt"
	"syscall"
)

func addNewUser() {
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
}
