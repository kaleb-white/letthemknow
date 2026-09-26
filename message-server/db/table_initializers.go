package db

import (
	"context"
	"fmt"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

const SOURCE_INITIALIZER = "db_table_initializer"

func InitializeTable(ctx context.Context, sql string, tableName string) error {
	log.Log(SOURCE_INITIALIZER, fmt.Sprintf("Initializing table: %s", tableName), log.INFO)

	// Get connection
	conn, err := GetConnection()
	if err != nil {
		log.Log(SOURCE_INITIALIZER, fmt.Sprintf("Error while getting conn %s: %s", tableName, err.Error()), log.ERROR)
	}
	defer conn.Close()

	// Run initialization
	_, err = conn.ExecContext(ctx, sql)
	if err != nil {
		log.Log(SOURCE_INITIALIZER, fmt.Sprintf("Error during initialization of table %s: %s", tableName, err.Error()), log.ERROR)
	}
	return nil
}
