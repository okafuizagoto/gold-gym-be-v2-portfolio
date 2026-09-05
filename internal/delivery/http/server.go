package http

import (
	"context"
	"net/http"

	"gold-gym-be/internal/config"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/gin-gonic/gin"
)

// GoldGymHandler ...
type GoldGymHandler interface {
	// GetGoldGym(w http.ResponseWriter, r *http.Request)
	// InsertGoldGym(w http.ResponseWriter, r *http.Request)
	// DeleteGoldGym(w http.ResponseWriter, r *http.Request)
	// UpdateGoldGym(w http.ResponseWriter, r *http.Request)

	GetGoldGymGin(c *gin.Context)
	InsertGoldGymGin(c *gin.Context)
	InsertGoldGymKafkaGin(c *gin.Context)
	DeleteGoldGymGin(c *gin.Context)
	UpdateGoldGymGin(c *gin.Context)
	UpdateGoldGymTokoGin(c *gin.Context)
	UpdateGoldGymBuyerGin(c *gin.Context)
	UpdateRegistrationModeGin(c *gin.Context)

	// PrintSelisih(w http.ResponseWriter, r *http.Request)
	// PrintExpiredTerpajang(w http.ResponseWriter, r *http.Request)
	// PrintExpiredTerkumpul(w http.ResponseWriter, r *http.Request)

	// PrintBatch(w http.ResponseWriter, r *http.Request)
	// PrintBatchFull(w http.ResponseWriter, r *http.Request)

	// //Trans Out
	// InsertTransOut(w http.ResponseWriter, r *http.Request)
	// InsertSales(w http.ResponseWriter, r *http.Request)
	// DeleteSalesByPeriod(w http.ResponseWriter, r *http.Request)
	// RemoveSalesByOutcode(w http.ResponseWriter, r *http.Request)
	// InsertBatchData(w http.ResponseWriter, r *http.Request)
}

type GoldGymStockHandler interface {
	GetGoldGymStockGin(c *gin.Context)
	InsertGoldGymStockGin(c *gin.Context)
	UpdateGoldGymStockGin(c *gin.Context)
	DeleteGoldGymStockGin(c *gin.Context)
}

type GoldGymItemsHandler interface {
	GetGoldGymItemsGin(c *gin.Context)
	InsertGoldGymItemsGin(c *gin.Context)
	UpdateGoldGymItemsGin(c *gin.Context)
	DeleteGoldGymItemsGin(c *gin.Context)
}

type GoldGymDiscountHandler interface {
	GetGoldGymDiscountGin(c *gin.Context)
	InsertGoldGymDiscountGin(c *gin.Context)
	UpdateGoldGymDiscountGin(c *gin.Context)
	DeleteGoldGymDiscountGin(c *gin.Context)
}

type GoldGymStorageHandler interface {
	GetGoldGymStorageGin(c *gin.Context)
	DeleteGoldGymStorageGin(c *gin.Context)
}

type GoldGymOutletHandler interface {
	GetGoldGymOutletGin(c *gin.Context)
	InsertGoldGymOutletGin(c *gin.Context)
	UpdateGoldGymOutletGin(c *gin.Context)
	DeleteGoldGymOutletGin(c *gin.Context)
}

type GoldGymAreaHandler interface {
	GetGoldGymAreaGin(c *gin.Context)
	InsertGoldGymAreaGin(c *gin.Context)
}

type GoldGymMejaHandler interface {
	GetGoldGymMejaGin(c *gin.Context)
	InsertGoldGymMejaGin(c *gin.Context)
	UpdateGoldGymMejaGin(c *gin.Context)
}

type GoldGymCustomerHandler interface {
	GetGoldGymCustomerGin(c *gin.Context)
	InsertGoldGymCustomerGin(c *gin.Context)
	UpdateGoldGymCustomerGin(c *gin.Context)
	DeleteGoldGymCustomerGin(c *gin.Context)
}

type GoldGymCustomerTypeHandler interface {
	GetGoldGymCustomerTypeGin(c *gin.Context)
	InsertGoldGymCustomerTypeGin(c *gin.Context)
	UpdateGoldGymCustomerTypeGin(c *gin.Context)
	DeleteGoldGymCustomerTypeGin(c *gin.Context)
}

type GoldGymSaleHandler interface {
	GetGoldGymSaleGin(c *gin.Context)
	InsertGoldGymSaleGin(c *gin.Context)
	UpdateGoldGymSaleGin(c *gin.Context)
	DeleteGoldGymSaleGin(c *gin.Context)
}

// GoldGymBookingHandler ...
type GoldGymBookingHandler interface {
	GetGoldGymBookingGin(c *gin.Context)
	InsertGoldGymBookingGin(c *gin.Context)
	UpdateGoldGymBookingGin(c *gin.Context)
	DeleteGoldGymBookingGin(c *gin.Context)
}

// GoldGymOrderHandler ...
type GoldGymOrderHandler interface {
	GetGoldGymOrderGin(c *gin.Context)
	InsertGoldGymOrderGin(c *gin.Context)
	UpdateGoldGymOrderGin(c *gin.Context)
	DeleteGoldGymOrderGin(c *gin.Context)
}

// GoldGymSellerAccessHandler ...
type GoldGymSellerAccessHandler interface {
	GetGoldGymSellerAccessGin(c *gin.Context)
	UpdateGoldGymSellerAccessGin(c *gin.Context)
}

// AuthHandler ...
type AuthHandler interface {
	// LoginUser(w http.ResponseWriter, r *http.Request)
	LoginUser(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type MiddlewareHandler interface {
	CheckUniqueRequest(c *gin.Context)
	ValidateToken(c *gin.Context)
	RegistrationRateLimit(c *gin.Context)
	LoginRateLimit(c *gin.Context)
	SensitiveUserActionRateLimit(c *gin.Context)
}

type HealthHandler interface {
	Check(c *gin.Context)
}

type ElasticHandler interface {
	GetElasticGin(c *gin.Context)
	PostElasticGin(c *gin.Context)
	PutElasticGin(c *gin.Context)
	DeleteElasticGin(c *gin.Context)
}

// GymActivityHandler defines the HTTP interface for gymactivity (MongoDB) endpoints
type GymActivityHandler interface {
	GetGymActivityGin(c *gin.Context)
	InsertGymActivityGin(c *gin.Context)
	UpdateGymActivityGin(c *gin.Context)
	DeleteGymActivityGin(c *gin.Context)
}

// GraphQLHandler defines the interface for GraphQL endpoint
type GraphQLHandler interface {
	GraphQLGin(c *gin.Context)
}

// Server ...
type Server struct {
	Goldgym             GoldGymHandler
	GoldgymStock        GoldGymStockHandler
	GoldgymItems        GoldGymItemsHandler
	GoldgymDiscount     GoldGymDiscountHandler
	GoldgymStorage      GoldGymStorageHandler
	GoldgymOutlet       GoldGymOutletHandler
	GoldgymArea         GoldGymAreaHandler
	GoldgymMeja         GoldGymMejaHandler
	GoldgymCustomer     GoldGymCustomerHandler
	GoldgymCustomerType GoldGymCustomerTypeHandler
	GoldgymSale         GoldGymSaleHandler
	GoldgymBooking      GoldGymBookingHandler
	GoldgymOrder        GoldGymOrderHandler
	GoldgymSellerAccess GoldGymSellerAccessHandler
	Auth                AuthHandler
	Middleware          MiddlewareHandler
	// mongodb
	// GymActivity GymActivityHandler
	Elastic ElasticHandler
	GraphQL GraphQLHandler

	engine *gin.Engine
	server *http.Server

	Health HealthHandler

	Logger jaegerLog.Factory
	Config *config.Config
}

// Serve is serving HTTP gracefully on port x ...
// func (s *Server) Serve(port string) error {
// 	handler := cors.AllowAll().Handler(s.Handler())
// 	return grace.Serve(port, handler)
// }

// func (s *Server) Serve(port string) error {
// 	handler := s.Handler()         // Change this to use Gin handler instead of mux
// 	return handler.Run(":" + port) // Gin's way to serve the app
// }

func (s *Server) Serve(port string) error {
	s.engine = s.Handler()

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: s.engine,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	return s.server.Shutdown(ctx)
}
