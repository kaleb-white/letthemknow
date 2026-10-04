package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "modernc.org/sqlite"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

const SOURCE string = "db"

var dataSourceName string = os.Getenv("DATA_SOURCE_NAME")
var db *sql.DB
var once sync.Once

func GetConnection() (*sql.Conn, error) {
	once.Do(func() {
		log.Log(SOURCE, fmt.Sprintf("Opening db with DSN %s...", dataSourceName), log.DEBUG)
		var err error
		db, err = sql.Open("sqlite", dataSourceName)
		if err != nil {
			log.Log(SOURCE, fmt.Sprintf("Error while opening db: %s", err.Error()), log.ERROR)
			panic(fmt.Sprintf("Error while opening db: %s", err.Error()))
		}
	})

	err := db.Ping()
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while opening db: %s", err.Error()), log.ERROR)
		return nil, err
	}

	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while fetching connection: %s", err.Error()), log.ERROR)
		return nil, err
	}
	
	log.Log(SOURCE, "Db conn fetch successful", log.DEBUG)
	return conn, nil
}

func CloseConnection(db *sql.DB) error {
	// Handle connection failures gracefully
	if db == nil {
		return nil
	}

	err := db.Close()
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while closing db: %s", log.ERROR))
		return err
	}
	return nil
}
