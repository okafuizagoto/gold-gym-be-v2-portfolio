package goldgym

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	jaegerLog "gold-gym-be/pkg/log"
)

type (
	// Data ...
	Data struct {
		// db   *sqlx.DB
		db   *gorm.DB
		dbr  *sqlx.DB
		stmt *map[string]*sqlx.Stmt

		tracer opentracing.Tracer
		logger jaegerLog.Factory
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

const (
	insertItems  = "InsertItems"
	qInsertItems = `INSERT INTO items
(item_gold_id, item_outcode, item_code, item_name, item_type, item_pack, item_brand, item_description, item_status, item_created_at, item_updated_at)
VALUES %s`
)

var (
	readStmt   = []statement{}
	insertStmt = []statement{
		// {insertItems, qInsertItems},
	}
	updateStmt = []statement{}
	deleteStmt = []statement{}
)

// New ...
func New(db *gorm.DB, dbr *sqlx.DB, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
	var (
		stmts = make(map[string]*sqlx.Stmt)
	)
	d := &Data{
		db:     db,
		dbr:    dbr,
		tracer: tracer,
		logger: logger,
		stmt:   &stmts,
	}
	d.InitStmt()
	return d
}

func (d *Data) InitStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize select statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range insertStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize insert statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range updateStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize update statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range deleteStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize delete statement key %v, err : %v", v.key, err)
		}
	}

	*d.stmt = stmts
}

// Close will cleanup prepared statements
func (d *Data) Close() {
	if d.stmt == nil {
		return
	}

	for k, stmt := range *d.stmt {
		if stmt != nil {
			if err := stmt.Close(); err != nil {
				log.Printf("[DB] failed to close stmt %s: %v", k, err)
			}
		}
	}
}
