package http

import (
	"context"
	"gold-gym-be/internal/delivery/http/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// // Handler will initialize mux router and register handler
// func (s *Server) Handler() *mux.Router {
// 	r := mux.NewRouter()
// 	// Jika tidak ditemukan, jangan diubah.
// 	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
// 	// Health Check
// 	r.HandleFunc("", defaultHandler).Methods("GET")
// 	r.HandleFunc("/", defaultHandler).Methods("GET")

// 	// Tambahan Prefix di depan API endpoint
// 	router := r.PathPrefix("/gold-gym").Subrouter()

// 	router.HandleFunc("", defaultHandler).Methods("GET")
// 	router.HandleFunc("/", defaultHandler).Methods("GET")

// 	sub := router.PathPrefix("/v2").Subrouter()

// 	// Routes
// 	goldgym := sub.PathPrefix("/userdata").Subrouter()

// 	goldgym.HandleFunc("", s.Goldgym.GetGoldGym).Methods("GET")
// 	goldgym.HandleFunc("", s.Goldgym.InsertGoldGym).Methods("POST")
// 	goldgym.HandleFunc("", s.Goldgym.UpdateGoldGym).Methods("PUT")
// 	goldgym.HandleFunc("", s.Goldgym.DeleteGoldGym).Methods("DELETE")

// 	goldgym.HandleFunc("/login", s.Auth.LoginUser).Methods("POST")

// 	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
// 	return r
// }

// func defaultHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Example Service API"))
// }

// Handler will initialize Gin router and register handler
func (s *Server) Handler() *gin.Engine {
	r := gin.New()

	// recovery
	if s.Config.Server.Env == "local" {
		r.Use(gin.Recovery())
	} else {
		// r.Use(gin.CustomRecovery(func(c gin.Context, recovered interface{}) {
		// 	s.Logger.For(context.Background()).Error(
		// 		"panic",
		// 		zap.Any("err", recovered),
		// 	)
		// 	c.JSON(500, gin.H{"error": "internal server error"})
		// }))
		// ----------------------------------------------------------
		// r.Use(func(c *gin.Context) {
		// 	defer func() {
		// 		if err := recover(); err != nil {
		// 			log.Printf("panic: %v", err)
		// 			c.JSON(500, gin.H{"error": "internal error"})
		// 		}
		// 	}()
		// 	c.Next()
		// })
		// Untuk semua environment termasuk production
		r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
			// Log dengan stack trace
			log.Printf("[GIN] panic recovered: %v", recovered)
			middleware.AlertPanic(c, recovered)

			// Jangan expose detail error ke client
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			c.Abort()
		}))
	}

	// metrics + access log
	r.Use(middleware.PrometheusMetrics())
	r.Use(middleware.AccessLogger())
	r.Use(middleware.ErrorAlert())

	// timeout
	// 5s terlalu ketat untuk registrasi (4x hashing argon2 + kirim email SMTP
	// ke Gmail bisa >5s), sehingga context deadline exceeded sebelum insert
	// user sempat commit ke DB. Dinaikkan ke 20s.
	r.Use(middleware.Timeout(20 * time.Second))

	// r.Use(middleware.Timeout(5 * time.Second))
	// r.Use(func(c *gin.Context) {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			log.Printf("panic: %v", err)
	// 			c.JSON(500, gin.H{"error": "internal error"})
	// 		}
	// 	}()
	// 	c.Next()
	// })
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		s.Logger.For(context.Background()).Error(
			"panic recovered",
			zap.Any("err", recovered),
		)
		middleware.AlertPanic(c, recovered)

		c.JSON(500, gin.H{
			"error": "internal server error",
		})
	}))
	// Health Check
	// r.GET("/", defaultHandler)
	r.GET("", defaultHandler)
	r.GET("/healthz", s.Health.Check)

	// Tambahan Prefix di depan API endpoint
	router := r.Group("/gold-gym")

	goldgymAuth := router.Group("/v2/auth")
	{
		goldgymAuth.POST("/refreshtoken", s.Auth.RefreshToken)
		goldgymAuth.POST("/logout", s.Auth.Logout)
	}

	// Routes
	goldgym := router.Group("/v2/userdata")
	{
		// Define the routes for GoldGym
		goldgym.GET("", s.Goldgym.GetGoldGymGin)                                                                                                                  // GET
		goldgym.POST("", s.Middleware.RegistrationRateLimit, s.Middleware.LoginRateLimit, s.Middleware.CheckUniqueRequest, s.Goldgym.InsertGoldGymGin)            // POST sync
		goldgym.POST("/kafka", s.Middleware.RegistrationRateLimit, s.Middleware.LoginRateLimit, s.Middleware.CheckUniqueRequest, s.Goldgym.InsertGoldGymKafkaGin) // POST async via Kafka
		goldgym.PUT("", s.Middleware.SensitiveUserActionRateLimit, s.Goldgym.UpdateGoldGymGin)                                                                    // PUT
		goldgym.PUT("/toko", s.Middleware.ValidateToken, s.Goldgym.UpdateGoldGymTokoGin)                                                                          // PUT nama toko sendiri
		goldgym.PUT("/buyer", s.Middleware.ValidateToken, s.Goldgym.UpdateGoldGymBuyerGin)                                                                        // PUT daftar sebagai pembeli (flag)
		goldgym.PUT("/registrationmode", s.Middleware.ValidateToken, s.Goldgym.UpdateRegistrationModeGin)                                                         // PUT mode pendaftaran mandiri (ADMIN)
		goldgym.DELETE("", s.Middleware.SensitiveUserActionRateLimit, s.Goldgym.DeleteGoldGymGin)                                                                 // DELETE

		// Auth routes
		goldgym.POST("/login", s.Middleware.LoginRateLimit, s.Auth.LoginUser) // POST
	}
	goldgymStock := router.Group("/v2/stock")
	{
		// Define the routes for GoldGym
		goldgymStock.GET("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymStock.GetGoldGymStockGin)     // GET
		goldgymStock.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymStock.InsertGoldGymStockGin) // POST
		goldgymStock.PUT("", s.GoldgymStock.UpdateGoldGymStockGin)                                                               // PUT
		goldgymStock.DELETE("", s.GoldgymStock.DeleteGoldGymStockGin)                                                            // DELETE

		// Elastic routes
		elastic := router.Group("/v2/elastic")
		{
			elastic.GET("", s.Elastic.GetElasticGin)       // GET: search or getbyid
			elastic.POST("", s.Elastic.PostElasticGin)     // POST: index document
			elastic.PUT("", s.Elastic.PutElasticGin)       // PUT: update document
			elastic.DELETE("", s.Elastic.DeleteElasticGin) // DELETE: delete document
		}
		// // Auth routes
		// goldgym.POST("/login", s.Auth.LoginUser) // POST
	}
	goldgymItems := router.Group("/v2/items")
	{
		// Define the routes for GoldGym
		goldgymItems.GET("", s.Middleware.ValidateToken, s.GoldgymItems.GetGoldGymItemsGin)                                      // GET
		goldgymItems.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymItems.InsertGoldGymItemsGin) // POST
		goldgymItems.PUT("", s.Middleware.ValidateToken, s.GoldgymItems.UpdateGoldGymItemsGin)                                   // PUT
		goldgymItems.DELETE("", s.Middleware.ValidateToken, s.GoldgymItems.DeleteGoldGymItemsGin)                                // DELETE
	}

	goldgymDiscount := router.Group("/v2/discount")
	{
		goldgymDiscount.GET("", s.Middleware.ValidateToken, s.GoldgymDiscount.GetGoldGymDiscountGin)                                      // GET
		goldgymDiscount.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymDiscount.InsertGoldGymDiscountGin) // POST
		goldgymDiscount.PUT("", s.Middleware.ValidateToken, s.GoldgymDiscount.UpdateGoldGymDiscountGin)                                   // PUT
		goldgymDiscount.DELETE("", s.Middleware.ValidateToken, s.GoldgymDiscount.DeleteGoldGymDiscountGin)                                // DELETE
	}

	// Menu Storage -- daftar & hapus foto (item + bukti pembayaran) milik
	// user yang login. Semua role KECUALI ADMIN (ditolak di handler).
	goldgymStorage := router.Group("/v2/storage")
	{
		goldgymStorage.GET("", s.Middleware.ValidateToken, s.GoldgymStorage.GetGoldGymStorageGin)       // GET
		goldgymStorage.DELETE("", s.Middleware.ValidateToken, s.GoldgymStorage.DeleteGoldGymStorageGin) // DELETE
	}
	goldgymOutlet := router.Group("/v2/outlet")
	{
		// Define the routes for GoldGym
		goldgymOutlet.GET("", s.Middleware.ValidateToken, s.GoldgymOutlet.GetGoldGymOutletGin)                                      // GET
		goldgymOutlet.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymOutlet.InsertGoldGymOutletGin) // POST
		goldgymOutlet.PUT("", s.Middleware.ValidateToken, s.GoldgymOutlet.UpdateGoldGymOutletGin)                                   // PUT
		goldgymOutlet.DELETE("", s.Middleware.ValidateToken, s.GoldgymOutlet.DeleteGoldGymOutletGin)                                // DELETE

		// // Auth routes
		// goldgym.POST("/login", s.Auth.LoginUser) // POST
	}

	goldgymArea := router.Group("/v2/area")
	{
		goldgymArea.GET("", s.Middleware.ValidateToken, s.GoldgymArea.GetGoldGymAreaGin)                                      // GET
		goldgymArea.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymArea.InsertGoldGymAreaGin) // POST
	}

	goldgymMeja := router.Group("/v2/meja")
	{
		goldgymMeja.GET("", s.Middleware.ValidateToken, s.GoldgymMeja.GetGoldGymMejaGin)                                      // GET
		goldgymMeja.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymMeja.InsertGoldGymMejaGin) // POST
		goldgymMeja.PUT("", s.Middleware.ValidateToken, s.GoldgymMeja.UpdateGoldGymMejaGin)                                   // PUT: reservemeja / releasemeja
	}

	goldgymCustomer := router.Group("/v2/cust")
	{
		// Define the routes for GoldGym
		goldgymCustomer.GET("", s.Middleware.ValidateToken, s.GoldgymCustomer.GetGoldGymCustomerGin)                                      // GET
		goldgymCustomer.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymCustomer.InsertGoldGymCustomerGin) // POST
		goldgymCustomer.PUT("", s.Middleware.ValidateToken, s.GoldgymCustomer.UpdateGoldGymCustomerGin)                                   // PUT
		goldgymCustomer.DELETE("", s.Middleware.ValidateToken, s.GoldgymCustomer.DeleteGoldGymCustomerGin)                                // DELETE

		// // Auth routes
		// goldgym.POST("/login", s.Auth.LoginUser) // POST
	}
	goldgymCustomerType := router.Group("/v2/typecust")
	{
		// Define the routes for GoldGym
		goldgymCustomerType.GET("", s.Middleware.ValidateToken, s.GoldgymCustomerType.GetGoldGymCustomerTypeGin)                                      // GET
		goldgymCustomerType.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymCustomerType.InsertGoldGymCustomerTypeGin) // POST
		goldgymCustomerType.PUT("", s.Middleware.ValidateToken, s.GoldgymCustomerType.UpdateGoldGymCustomerTypeGin)                                   // PUT
		goldgymCustomerType.DELETE("", s.Middleware.ValidateToken, s.GoldgymCustomerType.DeleteGoldGymCustomerTypeGin)                                // DELETE

		// // Auth routes
		// goldgym.POST("/login", s.Auth.LoginUser) // POST
	}
	goldgymBooking := router.Group("/v2/booking")
	{
		// Booking slot terapi (outlet type THERAPY)
		goldgymBooking.GET("", s.Middleware.ValidateToken, s.GoldgymBooking.GetGoldGymBookingGin)                                      // GET: slots / buyers
		goldgymBooking.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymBooking.InsertGoldGymBookingGin) // POST: insertbooking
		goldgymBooking.PUT("", s.Middleware.ValidateToken, s.GoldgymBooking.UpdateGoldGymBookingGin)                                   // PUT: paybooking
		goldgymBooking.DELETE("", s.Middleware.ValidateToken, s.GoldgymBooking.DeleteGoldGymBookingGin)                                // DELETE: removebooking (SELLER/ADMIN, UNPAID saja)
	}
	goldgymOrder := router.Group("/v2/order")
	{
		// Pesanan pembeli (BUYER order flow) — outlet non-THERAPY.
		goldgymOrder.GET("", s.Middleware.ValidateToken, s.GoldgymOrder.GetGoldGymOrderGin)                                      // GET: outlets/catalog/buyerorders/sellerorders/orderdetail
		goldgymOrder.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymOrder.InsertGoldGymOrderGin) // POST: insertorder
		goldgymOrder.PUT("", s.Middleware.ValidateToken, s.GoldgymOrder.UpdateGoldGymOrderGin)                                   // PUT: confirm/reject/finish
		goldgymOrder.DELETE("", s.Middleware.ValidateToken, s.GoldgymOrder.DeleteGoldGymOrderGin)                                // DELETE: removevisible (ADMIN)
	}
	goldgymSellerAccess := router.Group("/v2/selleraccess")
	{
		// ADMIN: aktif/nonaktifkan menu Daftar Pembeli & Mode Pembeli milik penjual.
		goldgymSellerAccess.GET("", s.Middleware.ValidateToken, s.GoldgymSellerAccess.GetGoldGymSellerAccessGin)    // GET: list
		goldgymSellerAccess.PUT("", s.Middleware.ValidateToken, s.GoldgymSellerAccess.UpdateGoldGymSellerAccessGin) // PUT: daftarpembeli/modepembeli
	}
	goldgymSale := router.Group("/v2/sales")
	{
		// Define the routes for GoldGym
		goldgymSale.GET("", s.Middleware.ValidateToken, s.GoldgymSale.GetGoldGymSaleGin)                                      // GET
		goldgymSale.POST("", s.Middleware.ValidateToken, s.Middleware.CheckUniqueRequest, s.GoldgymSale.InsertGoldGymSaleGin) // POST
		goldgymSale.PUT("", s.Middleware.ValidateToken, s.GoldgymSale.UpdateGoldGymSaleGin)                                   // PUT
		goldgymSale.DELETE("", s.Middleware.ValidateToken, s.GoldgymSale.DeleteGoldGymSaleGin)                                // DELETE

		// // Auth routes
		// goldgym.POST("/login", s.Auth.LoginUser) // POST
	}
	// // Elastic routes
	// elastic := router.Group("/v2/elastic")
	// {
	// 	elastic.GET("", s.Elastic.GetElasticGin)   // GET: search or getbyid
	// 	elastic.POST("", s.Elastic.PostElasticGin) // POST: index document
	// }

	// GymActivity routes (MongoDB)
	// mongodb
	// activity := router.Group("/v2/activity")
	// {
	// 	activity.GET("", s.GymActivity.GetGymActivityGin)          // GET: list activities
	// 	activity.POST("", s.GymActivity.InsertGymActivityGin)      // POST: insert activity
	// 	activity.PUT("", s.GymActivity.UpdateGymActivityGin)       // PUT: update activity by ?id=
	// 	activity.DELETE("", s.GymActivity.DeleteGymActivityGin)    // DELETE: delete activity by ?id=
	// }

	// GraphQL endpoint — POST for queries/mutations, GET for GraphiQL Playground
	r.POST("/graphql", s.GraphQL.GraphQLGin)
	r.GET("/graphql", s.GraphQL.GraphQLGin)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus metrics endpoint — scraped by Prometheus server
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func defaultHandler(c *gin.Context) {
	c.String(200, "Example Service API")
}

// func notFoundHandler(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		resp   *response.Response
// 		err    error
// 		errRes response.Error
// 	)
// 	resp = &response.Response{}
// 	defer resp.RenderJSON(w, r)

// 	err = errors.New("404 Not Found")

// 	if err != nil {
// 		// Error response handling
// 		errRes = response.Error{
// 			Code:   404,
// 			Msg:    "404 Not Found",
// 			Status: true,
// 		}

// 		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
// 		resp.StatusCode = 404
// 		resp.Error = errRes
// 		return
// 	}
// }
