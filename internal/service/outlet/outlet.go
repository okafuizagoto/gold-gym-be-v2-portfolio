package goldgym

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	jaegerLog "gold-gym-be/pkg/log"

	goldOutletEntity "gold-gym-be/internal/entity/outlet"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	// "go.opentelemetry.io/otel/trace"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	InsertOutlet(ctx context.Context, tx *gorm.DB, outlet []goldOutletEntity.Outlet) error
	WithTransactionOutlet(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetNextSequence(ctx context.Context, tx *gorm.DB, reqs []goldOutletEntity.InsertOutletCounter, outletDataJson []goldOutletEntity.InsertOutletCounterJson) ([]goldOutletEntity.Outlet, error)
	UpdateOutlet(ctx context.Context, updateoutlet goldOutletEntity.UpdateOutlet) error
	DeleteOutlet(ctx context.Context, goldid int, code string) error
	GetOutlet(ctx context.Context, goldid int, name, status string, page, length int) ([]goldOutletEntity.Outlet, error)
	GetTotalOutlet(ctx context.Context, goldid int, name string) (int64, error)
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	goldgymoutlet Data
	tracer        opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(goldgymOutletData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		goldgymoutlet: goldgymOutletData,
		tracer:        tracer,
		logger:        logger,
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
