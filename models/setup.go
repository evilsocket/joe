package models

import (
	"database/sql"
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/log"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var (
	DB         = (*sql.DB)(nil)
	Queries    = sync.Map{}
	NumQueries = 0
	Users      = sync.Map{}
	NumUsers   = 0
)

func FindQuery(name string) *Query {
	if q, found := Queries.Load(name); found {
		return q.(*Query)
	}
	return nil
}

func Cleanup() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			fmt.Println(err)
		}
		DB = nil
	}
}

func Setup(confFile, dataPath, usersPath string) (err error) {
	defer func() {
		log.Info("users:%d queries:%d", NumUsers, NumQueries)
	}()

	if err := godotenv.Load(confFile); err != nil {
		return fmt.Errorf("error while loading %s: %v", confFile, err)
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)

	log.Debug("connecting to %s database at %s ...", dbDriver, dbURL)

	if DB, err = sql.Open(dbDriver, dbURL); err != nil {
		return
	} else if err = DB.Ping(); err != nil {
		return
	}

	log.Info("loading users from %s ...", usersPath)
	err = fs.Glob(usersPath, "*.yml", func(fileName string) error {
		if user, err := LoadUser(fileName); err != nil {
			return fmt.Errorf("error while loading %s: %v", fileName, err)
		} else {
			Users.Store(user.Username, user)
			NumUsers++
		}
		return nil
	})
	if err != nil {
		return err
	}
	Users.Store("anonymous", &User{})

	log.Info("loading data from %s ...", dataPath)
	return fs.Glob(dataPath, "*.yml", func(fileName string) error {
		if query, err := LoadQuery(fileName); err != nil {
			return fmt.Errorf("error while loading %s: %v", fileName, err)
		} else {
			Queries.Store(query.Name, query)
			NumQueries++
		}
		return nil
	})
}
