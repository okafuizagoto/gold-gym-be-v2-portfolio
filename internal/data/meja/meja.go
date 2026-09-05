package goldgym

import (
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	jaegerLog "gold-gym-be/pkg/log"
)

type Data struct {
	db  *gorm.DB
	dbr *sqlx.DB

	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

func New(db *gorm.DB, dbr *sqlx.DB, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
	return &Data{
		db:     db,
		dbr:    dbr,
		tracer: tracer,
		logger: logger,
	}
}
