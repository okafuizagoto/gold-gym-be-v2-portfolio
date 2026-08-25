package goldgym

import (
	"context"
	jaegerLog "gold-gym-be/pkg/log"

	firebaseEntity "gold-gym-be/internal/entity/firebase"
	goldStockEntity "gold-gym-be/internal/entity/stock"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcStock interface {
	// LoginUser(ctx context.Context, _user, _password string, _host string) (auth.Token, map[string]interface{}, error)
	GetOneStockProduct(ctx context.Context, stockcode string, stocknmame string, stockid string) (goldStockEntity.GetOneStock, error)
	// InsertStockSales(ctx context.Context, stock goldStockEntity.InsertStockData) (string, error)
	GetAllStockHeader(ctx context.Context) ([]goldStockEntity.GetOneStock, error)
	GetAllStockHeaderToRedis(ctx context.Context) ([]goldStockEntity.GetOneStock, error)
	GetFromFirebase(ctx context.Context, userID string) (*firebaseEntity.User, error)
	CreateUser(ctx context.Context, user firebaseEntity.User) (string, error)
	GetStockByName(ctx context.Context, goldid int, stockname string) ([]goldStockEntity.GetOneStock, error)
	InsertStockSalesNew(ctx context.Context, goldid int, stock []goldStockEntity.InsertStock) (string, error)
	GetStock(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldStockEntity.GetOneStock, goldStockEntity.MetadataPaginationDetail, error)
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcStock IgoldgymSvcStock
		tracer          opentracing.Tracer
		logger          jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(isst IgoldgymSvcStock, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcStock: isst,
		tracer:          tracer,
		logger:          logger,
	}
}
