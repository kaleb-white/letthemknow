package utils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/kaleb-white/letthemknow/log"
)

const SOURCE string = "db"

var dataSourceName string = os.Getenv("DATASOURCENAME")

func GetConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while opening db: %s", err.Error()), log.ERROR)
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while opening db: %s", err.Error()), log.ERROR)
		return db, err
	}

	return db, nil
}
