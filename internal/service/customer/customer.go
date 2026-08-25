package goldgym

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	goldCustomerEntity "gold-gym-be/internal/entity/customer"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	// "go.opentelemetry.io/otel/trace"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	WithTransactionCustomer(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetCustomer(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerEntity.Customer, error)
	GetTotalCustomer(ctx context.Context, goldid int, name string, outcode string) (int64, error)
	GetLastCustomerCode(ctx context.Context, goldid int, code string) (*string, error)
	InsertCustomer(ctx context.Context, tx *gorm.DB, id int, items []goldCustomerEntity.Customer) error
	UpdateCustomer(ctx context.Context, cust goldCustomerEntity.Customer) error
	DeleteCustomer(ctx context.Context, goldid, goldcustomerid int, outcode string) error
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	goldgymcustomer Data
	tracer          opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(goldgymCustomerData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		goldgymcustomer: goldgymCustomerData,
		tracer:          tracer,
		logger:          logger,
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
