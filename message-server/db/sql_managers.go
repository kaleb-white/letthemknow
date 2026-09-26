package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

const SOURCE_SQL_HELPERS string = "sql_managers"

type PreparedSql struct {
	OperationName				 	 string
	Err 									 error
	ErrorDuringPreparation bool
	Prepared							 bool
	RawSql                 string
	PreparedSql            *sql.Stmt
}

func (p *PreparedSql) tryAndPrepare(ctx context.Context, conn *sql.Conn) (*sql.Stmt, error) {
	log.Log(SOURCE_SQL_HELPERS, fmt.Sprintf("Preparing operation %s", p.OperationName), log.INFO)
	log.Log(SOURCE_SQL_HELPERS, fmt.Sprintf("Sql: %s", p.RawSql), log.DEBUG)
	preparedSql, err := conn.PrepareContext(ctx, p.RawSql)

	if err != nil {
		log.Log(SOURCE_SQL_HELPERS, fmt.Sprintf("Error while preparing sql for operation %s: %s", p.OperationName, err.Error()))
		return nil, err
	}
	
	return preparedSql, nil
}

func (p *PreparedSql) GetPreparedStatement(ctx context.Context, conn *sql.Conn) (*sql.Stmt, error) {
	switch {
	case p.ErrorDuringPreparation: 
		preparedSql, err := p.tryAndPrepare(ctx, conn)
		if err != nil {
			p.Err = err
			return nil, err
		}
		return preparedSql, nil 
	case p.Prepared:
		return p.PreparedSql, nil
	default:
		preparedSql, err := p.tryAndPrepare(ctx, conn)
		if err != nil {
			p.Err = err
			p.ErrorDuringPreparation = true
			p.Prepared = false
			return nil, err
		}
		p.ErrorDuringPreparation = false
		p.Prepared = true
		p.PreparedSql = preparedSql

		return preparedSql, nil
	}
}

// Implements io.Closer
func (p *PreparedSql) Close() error {
	if !p.Prepared {
		return nil
	} 
	err := p.PreparedSql.Close()
	return err
}

