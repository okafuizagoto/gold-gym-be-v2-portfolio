package goldgym

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldItemsEntity "gold-gym-be/internal/entity/items"
	goldQuotaEntity "gold-gym-be/internal/entity/quota"

	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	// "go.opentelemetry.io/otel/trace"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	GetItems(ctx context.Context, goldid int, name string, outcode string, limit, offset int) ([]goldItemsEntity.Item, error)
	GetTotalItems(ctx context.Context, goldid int, name string, outcode string) (int64, error)
	InsertItems(ctx context.Context, tx *gorm.DB, id int, user []goldItemsEntity.InsertItem) error
	GetLastItemCode(ctx context.Context, goldid int, code string) (*string, error)
	UpdateItems(ctx context.Context, items goldItemsEntity.UpdateItems) error
	GetItemGoldID(ctx context.Context, itemID int, outcode string) (int, error)
	DeleteItems(ctx context.Context, goldid, golditemid int, outcode string) error
	WithTransactionItems(ctx context.Context, fn func(tx *gorm.DB) error) error
	EnsureTherapyStock(ctx context.Context, goldid int, outcode string) error
	GetOutletCodesByGoldID(ctx context.Context, goldid int) ([]string, error)
	GetItemByID(ctx context.Context, itemID int) (goldItemsEntity.Item, error)
	UpdateItemPhoto(ctx context.Context, itemID int, filename string, bytes int) error
	ClearItemPhoto(ctx context.Context, itemID int) error
	GetUserStorageUsedKB(ctx context.Context, goldID int) (int, error)
	AddUserStorageUsedKB(ctx context.Context, goldID int, deltaKB int) error
	InsertStorageDeleteHistory(ctx context.Context, h goldQuotaEntity.StorageDeleteHistory) error
}

type UserData interface {
	GetGoldUserByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
}

type RedisData interface {
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	goldgymitems Data
	goldgymuser  UserData
	goldgymredis RedisData
	tracer       opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(goldgymItemsData Data, goldGymUserData UserData, goldGymRedisData RedisData, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		goldgymitems: goldgymItemsData,
		goldgymuser:  goldGymUserData,
		goldgymredis: goldGymRedisData,
		tracer:       tracer,
		logger:       logger,
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
