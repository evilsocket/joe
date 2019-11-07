package models

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/evilsocket/islazy/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	TokenTTL int    `yaml:"token_ttl"`
}

func SaveUser(user User, fileName string) error {
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum(nil))

	if raw, err := yaml.Marshal(&user); err != nil {
		return err
	} else if err := ioutil.WriteFile(fileName, raw, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func LoadUser(fileName string) (*User, error) {
	log.Debug("loading %s ...", fileName)

	user := &User{}

	if raw, err := ioutil.ReadFile(fileName); err != nil {
		return nil, err
	} else if err = yaml.Unmarshal(raw, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) ValidPassword(password string) bool {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)) == u.Password
}
