package auth

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	// "go.opentelemetry.io/otel/trace"
)

type RedisData interface {
	GetFromRedis(ctx context.Context, key string, dest interface{}) (err error)
	DeleteFromRedis(ctx context.Context, key string) error
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	redis  RedisData
	tracer opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(redisData RedisData, tracer opentracing.Tracer, logger jaegerLog.Factory) *Service {
	// Assign variable dari parameter ke object
	return &Service{
		redis:  redisData,
		tracer: tracer,
		logger: logger,
	}
}

func (s Service) checkPermission(ctx context.Context, _permissions ...string) error {
	claims := ctx.Value(entity.ContextKey("claims"))
	if claims != nil {
		actions := claims.(entity.ContextValue).Get("permissions").(map[string]interface{})
		for _, action := range actions {
			permissions := action.([]interface{})
			for _, permission := range permissions {
				for _, _permission := range _permissions {
					if permission.(string) == _permission {
						return nil
					}
				}
			}
		}
	}
	return errors.New("401 unauthorized")
}
