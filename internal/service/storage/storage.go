package goldgym

import (
	"context"
	goldStorageEntity "gold-gym-be/internal/entity/storage"
	jaegerLog "gold-gym-be/pkg/log"

	itemsService "gold-gym-be/internal/service/items"
	salesService "gold-gym-be/internal/service/sales"

	"github.com/opentracing/opentracing-go"
)

// Data ...
type Data interface {
	ListItemPhotos(ctx context.Context, goldID int) ([]goldStorageEntity.StorageEntry, error)
	ListPaymentProofs(ctx context.Context, goldID int) ([]goldStorageEntity.StorageEntry, error)
	GetUserStorageUsedKB(ctx context.Context, goldID int) (int, error)
}

// Service ...
// Delete-nya sengaja TIDAK mengulang logika file/DB items & payment_proof --
// dispatch ke Service milik modul items/sales yang sudah punya
// DeleteItemPhoto/DeletePaymentProofPhoto (menyimpan tempat file & cleanup
// kuota di satu tempat, bukan diduplikasi di sini).
type Service struct {
	storage  Data
	itemsSvc itemsService.Service
	salesSvc salesService.Service
	tracer   opentracing.Tracer
	logger   jaegerLog.Factory
}

// New ...
func New(storageData Data, itemsSvc itemsService.Service, salesSvc salesService.Service, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		storage:  storageData,
		itemsSvc: itemsSvc,
		salesSvc: salesSvc,
		tracer:   tracer,
		logger:   logger,
	}
}
