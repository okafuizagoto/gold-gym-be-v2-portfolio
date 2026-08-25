package graphql

import (
	"context"
	"net/http"

	"gold-gym-be/internal/entity/auth/v2"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	jaegerLog "gold-gym-be/pkg/log"
	"gold-gym-be/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/handler"
	"github.com/opentracing/opentracing-go"
)

// IgoldgymSvc defines all goldgym service methods exposed via GraphQL
type IgoldgymSvc interface {
	// Query
	GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	GetGoldUserByEmail(ctx context.Context, email string) (string, error)
	GetGoldUserDataByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
	GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error)
	GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error)
	GetSubscriptionHeaderTotalHarga(ctx context.Context, email string) (goldEntity.SubscriptionHeaderPayment, error)

	// Insert
	InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error)
	InsertSubscriptionUser(ctx context.Context, subs goldEntity.InsertSubsAll) (string, error)
	InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) (string, error, response.Response)

	// Update
	UpdateDataPeserta(ctx context.Context, subs goldEntity.UpdatePassword) (string, error)
	UpdateNama(ctx context.Context, subs goldEntity.UpdateNama) (string, error)
	UpdateKartu(ctx context.Context, subs goldEntity.UpdateKartu) (string, error)
	UpdateSubscriptionDetail(ctx context.Context, subs goldEntity.UpdateSubs) (string, error)
	UpdateValidationOTP(ctx context.Context, otp string, email string) (string, error)
	UpdateOTP(ctx context.Context, email string) (string, error)
	UpdateOTPSubscription(ctx context.Context, id string) (string, error)
	UpdatePayment(ctx context.Context, otp string, email string) (string, error, response.Response)

	// Auth
	LoginUser(ctx context.Context, email string, password string, host, device string) (auth.Token, map[string]interface{}, string, error)
	Logout(ctx context.Context, subs goldEntity.Logout) (string, error)
}

// Handler wraps the graphql-go HTTP handler with Gin
type Handler struct {
	goldgymSvc  IgoldgymSvc
	tracer      opentracing.Tracer
	logger      jaegerLog.Factory
	httpHandler http.Handler
}

// New creates a new GraphQL Handler with schema built from svc
func New(svc IgoldgymSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	schema := buildSchema(svc)

	h := gqlhandler.New(&gqlhandler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, // GraphiQL Playground aktif di GET /graphql
	})

	return &Handler{
		goldgymSvc:  svc,
		tracer:      tracer,
		logger:      logger,
		httpHandler: h,
	}
}

// GraphQLGin is the Gin endpoint that wraps the graphql-go HTTP handler
func (h *Handler) GraphQLGin(c *gin.Context) {
	h.httpHandler.ServeHTTP(c.Writer, c.Request)
}

// buildSchema assembles the complete GraphQL schema
func buildSchema(svc IgoldgymSvc) graphql.Schema {
	queryType := buildQueryType(svc)
	mutationType := buildMutationType(svc)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic("failed to build GraphQL schema: " + err.Error())
	}
	return schema
}
