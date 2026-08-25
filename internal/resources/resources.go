package resources

import (
	"gold-gym-be/internal/service/goldgym" // sesuaikan path

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

type BootResources struct {
	DBLocal      *gorm.DB
	DBProd       *gorm.DB
	Redis        *redis.Client
	GoldSvcLocal goldgym.Service
	GoldSvcProd  goldgym.Service
	Tracer       opentracing.Tracer
	Logger       *zap.Logger
}
