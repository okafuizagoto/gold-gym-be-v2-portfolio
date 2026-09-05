package goldgym

import (
	"context"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"

	goldItemsEntity "gold-gym-be/internal/entity/items"
)

type IgoldgymSvcItems interface {
	GetItems(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldItemsEntity.Item, goldItemsEntity.MetadataPaginationDetail, error)
	InsertItems(ctx context.Context, id int, items []goldItemsEntity.InsertItem, applyAllOutlets bool) (goldItemsEntity.InsertItemsResult, error)
	UpdateItems(ctx context.Context, items goldItemsEntity.UpdateItems) (string, error)
	DeleteItems(ctx context.Context, goldid, golditemid int, outcode string) (string, error)
	SaveItemPhoto(ctx context.Context, itemID int, originalName string, mimeType string, content []byte, goldID int, isAdmin bool) (goldItemsEntity.Item, error)
	GetItemPhotoFile(ctx context.Context, itemID int) (string, []byte, error)
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcItems IgoldgymSvcItems
		tracer          opentracing.Tracer
		logger          jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(isim IgoldgymSvcItems, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcItems: isim,
		tracer:          tracer,
		logger:          logger,
	}
}
