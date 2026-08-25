package goldgym

import (
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"

	jaegerLog "gold-gym-be/pkg/log"
)

type (
	// Data ...
	Data struct {
		rdb *redis.Client

		tracer opentracing.Tracer
		logger jaegerLog.Factory
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

// New ...
// func New(db *gorm.DB, fsdb *db.Client, fs *storage.Client, rdb *redis.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
func New(rdb *redis.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
	d := &Data{
		rdb: rdb,
	}

	return d
}
