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
var (
	db *sql.DB
	once sync.Once
	err error
)

func OpenDB() {
	once.Do(func() {
		log.Log(SOURCE, fmt.Sprintf("Opening db with DSN %s...", dataSourceName), log.DEBUG)
		db, err = sql.Open("sqlite", dataSourceName)
		if err != nil {
			log.Log(SOURCE, fmt.Sprintf("Error while opening db: %s", err.Error()), log.ERROR)
			panic(fmt.Sprintf("Error while opening db: %s", err.Error()))
		}

		log.Log(SOURCE, "Pinging db...", log.DEBUG)
		err = db.Ping()
	})
}

func GetConnection() (*sql.Conn, error) {
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while opening db: %s", err.Error()), log.ERROR)
		return nil, err
	}
	
	log.Log(SOURCE, "Fetching conn...", log.DEBUG)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Log(SOURCE, fmt.Sprintf("Error while fetching connection: %s", err.Error()), log.ERROR)
		return nil, err
	}
	return conn, nil
}

func CloseDB() error {
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
