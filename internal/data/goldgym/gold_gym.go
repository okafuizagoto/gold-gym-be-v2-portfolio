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
	getSubsWithUser  = "GetSubsWithUser"
	qGetSubsWithUser = `SELECT a.gold_id, c.gold_menuid, a.gold_email, a.gold_nama, a.gold_nomorhp, a.gold_expireddate,
	c.gold_namapaket, c.gold_namalayanan, c.gold_harga, c.gold_listlatihan, c.gold_jumlahpertemuan, c.gold_durasi, c.gold_statuslangganan
	FROM data_peserta a
	LEFT JOIN subscription b
	ON a.gold_id = b.gold_id
	LEFT JOIN subscription_detail c
	ON b.gold_id = c.gold_id
	ORDER BY gold_id`

	insertSubscriptionDetail  = "InsertSubscriptionDetail"
	qInsertSubscriptionDetail = `INSERT INTO subscription_detail (gold_id, gold_menuid, gold_namapaket, gold_namalayanan, gold_harga, gold_jadwal, gold_listlatihan, gold_jumlahpertemuan, gold_durasi, gold_statuslangganan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

var (
	readStmt = []statement{
		{getSubsWithUser, qGetSubsWithUser},
	}
	insertStmt = []statement{
		{insertSubscriptionDetail, qInsertSubscriptionDetail},
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
