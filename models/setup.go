package models

import (
	"database/sql"
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/log"
	"github.com/joho/godotenv"
	"os"
)

var DB *sql.DB

func Cleanup() {
	if DB != nil {
		DB.Close()
		DB = nil
	}
}

func Setup(dataPath string, confFile string) (err error) {
	defer func() {
		log.Info("loaded %d total queries", NumQueries)
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
