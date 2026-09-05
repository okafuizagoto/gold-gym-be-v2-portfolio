package selleraccess

import (
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type Data struct {
	db     *gorm.DB
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

func New(db *gorm.DB, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
	return &Data{
		db:     db,
		tracer: tracer,
		logger: logger,
	}
}
