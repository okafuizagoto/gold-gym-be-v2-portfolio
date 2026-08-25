package goldgym

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/storage"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	jaegerLog "gold-gym-be/pkg/log"
)

type (
	// Data ...
	Data struct {
		db   *gorm.DB
		dbr  *sqlx.DB
		fsdb *firestore.Client
		s    *storage.Client
		rdb  *redis.Client
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

// Tambahkan query di dalam const
// getAllUser = "GetAllUser"
// qGetAllUser = "SELECT * FROM users"
const (
	insertCustomer  = "InsertCustomer"
	qInsertCustomer = `INSERT INTO customer
(cust_id, cust_gold_id, cust_outcode, cust_code, cust_type_id, cust_name, cust_phone, cust_address, cust_email, cust_status, cust_created_at, cust_updated_at, cust_created_by)
VALUES %s`
)

var (
	readStmt   = []statement{}
	insertStmt = []statement{}
	updateStmt = []statement{}
	deleteStmt = []statement{}
)

// New ...
// func New(db *gorm.DB, fsdb *db.Client, fs *storage.Client, rdb *redis.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
func New(db *gorm.DB, dbr *sqlx.DB, fsdb *firestore.Client, fs *storage.Client, rdb *redis.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
	var (
		stmts = make(map[string]*sqlx.Stmt)
	)
	d := &Data{
		db:     db,
		dbr:    dbr,
		fsdb:   fsdb,
		s:      fs,
		rdb:    rdb,
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
