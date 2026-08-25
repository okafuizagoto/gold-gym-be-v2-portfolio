package boot

import (
	// "context"

	"context"
	"encoding/json"
	"gold-gym-be/docs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "gold-gym-be/internal/data/auth"

	// "gold-gym-be/pkg/firebaseclient"

	"gold-gym-be/pkg/tracing"
	"log"

	"gold-gym-be/internal/config"
	jaegerLog "gold-gym-be/pkg/log"

	// Log "gold-gym-be/pkg/logs"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// "golang.org/x/net/trace"
	// "go.opentelemetry.io/otel/trace"
	// "gold-gym-be/pkg/trace"

	es "github.com/elastic/go-elasticsearch/v8"

	elasticData "gold-gym-be/internal/data/elastic"
	goldgymData "gold-gym-be/internal/data/goldgym"

	// gymactivityData "gold-gym-be/internal/data/gymactivity"
	goldgymGrpcHandler "gold-gym-be/internal/delivery/grpc/goldgym"
	goldgymServer "gold-gym-be/internal/delivery/http"
	elasticHandler "gold-gym-be/internal/delivery/http/elastic"
	goldgymHandler "gold-gym-be/internal/delivery/http/goldgym"
	graphqlHandler "gold-gym-be/internal/delivery/http/graphql"

	authHandler "gold-gym-be/internal/delivery/http/auth"
	authService "gold-gym-be/internal/service/auth"

	// gymactivityHandler "gold-gym-be/internal/delivery/http/gymactivity"
	elasticService "gold-gym-be/internal/service/elastic"
	goldgymService "gold-gym-be/internal/service/goldgym"

	// gymactivityService "gold-gym-be/internal/service/gymactivity"

	middlewareHandler "gold-gym-be/internal/delivery/http/middleware"
	middlewareService "gold-gym-be/internal/service/middleware"

	healthHandler "gold-gym-be/internal/delivery/http/health"

	goldgymStockData "gold-gym-be/internal/data/stock"
	goldgymStockHandler "gold-gym-be/internal/delivery/http/stock"
	goldgymStockService "gold-gym-be/internal/service/stock"

	// kafka new
	kafkagoldgymWorker "gold-gym-be/internal/worker/kafkagoldgym"
	"gold-gym-be/pkg/kafka"

	goldgymDiscountData "gold-gym-be/internal/data/discount"
	goldgymItemsData "gold-gym-be/internal/data/items"
	goldgymDiscountHandler "gold-gym-be/internal/delivery/http/discount"
	goldgymItemsHandler "gold-gym-be/internal/delivery/http/items"
	goldgymDiscountService "gold-gym-be/internal/service/discount"
	goldgymItemsService "gold-gym-be/internal/service/items"

	goldgymOutletData "gold-gym-be/internal/data/outlet"
	goldgymOutletHandler "gold-gym-be/internal/delivery/http/outlet"
	goldgymOutletService "gold-gym-be/internal/service/outlet"

	goldgymCustomerData "gold-gym-be/internal/data/customer"
	goldgymCustomerHandler "gold-gym-be/internal/delivery/http/customer"
	goldgymCustomerService "gold-gym-be/internal/service/customer"

	goldgymCustomerTypeData "gold-gym-be/internal/data/customertype"
	goldgymCustomerTypeHandler "gold-gym-be/internal/delivery/http/customertype"
	goldgymCustomerTypeService "gold-gym-be/internal/service/customertype"

	goldgymSaleData "gold-gym-be/internal/data/sales"
	goldgymSaleHandler "gold-gym-be/internal/delivery/http/sales"
	goldgymSaleService "gold-gym-be/internal/service/sales"

	goldgymBookingData "gold-gym-be/internal/data/booking"
	goldgymBookingHandler "gold-gym-be/internal/delivery/http/booking"
	goldgymBookingService "gold-gym-be/internal/service/booking"

	goldgymOrderData "gold-gym-be/internal/data/order"
	goldgymOrderHandler "gold-gym-be/internal/delivery/http/order"
	goldgymOrderService "gold-gym-be/internal/service/order"

	goldgymRedisData "gold-gym-be/internal/data/redis"

	pb "gold-gym-be/proto"
	"net"

	// mongodb ---------------------------------------------------
	// "go.mongodb.org/mongo-driver/v2/mongo"
	// mongoOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
	// mongodb ---------------------------------------------------

	"google.golang.org/grpc"
	// goldgymStockData "gold-gym-be/internal/data/stock"
	// pushNotifData "gold-gym-be/internal/data/pushnotif"
	// pushNotifHandler "gold-gym-be/internal/delivery/http/pushnotif"
	// pushNotifService "gold-gym-be/internal/service/pushnotif"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
func HTTP() error {
	var (
		// 	ctx = context.Background()
		// firebase
		// cred map[string]string
		// firebase
		cfg *config.Config // Configuration object
	)
	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	// firebase
	cfg, _ = config.Get()

	rdb := newRedisClient(cfg.Redis)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("[REDIS] Failed to connect: %v", err)
	}
	defer rdb.Close()

	// t, err := trace.New(ctx, cfg.Trace.Exporter)
	// if err != nil {
	// 	log.Fatalf("[CONFIG] Failed to initialize tracer: %v", err)
	// }
	// defer t.Shutdown(ctx)

	// Open MySQL DB Connection
	db, dbr, err := openDatabases(cfg)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	// // Open MySQL DB Connection
	// dbprod, dbrprod, err := openDatabasesProd(cfg)
	// if err != nil {
	// 	log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	// }

	// firebase
	// // Open MySQL DB Connection
	// f, err := firebaseclient.NewClient(cfg, cred)
	// if err != nil {
	// 	log.Fatalf("[FIREBASE] Failed to initialize firebase client: %v", err)
	// }
	// fs := f.StorageClient

	// ctx := context.Background()

	// firebaseApp, err := openFirebaseClient(ctx, cfg.Firebase, cred)
	// if err != nil {
	// 	log.Fatalf("[FIREBASE] Failed to initialize firebase client: %v", err)
	// }

	// fsdb, err := openFirestoreClient(ctx, firebaseApp)
	// if err != nil {
	// 	log.Fatalf("[FIRESTORE] Failed to initialize Firestore client: %v", err)
	// }
	// defer fsdb.Close()
	// firebase

	// fsdb, err := openFirebaseDatabaseClient(ctx, firebaseApp)
	// if err != nil {
	// 	log.Fatalf("[FIREBASE] Failed to initialize Realtime Database client: %v", err)
	// }

	// Firebase Client Init
	// fcmCredB2BPelapak, err := firebaseclient.NewClient(cfg.Firebase.FcmProjectIDB2BPelapak, cred)
	// if err != nil {
	// 	log.Fatalf("[FIREBASE] Failed to initialize firebase client: %v", err)
	// }
	// fcmB2BPelapak := fcmCredB2BPelapak.MessagingClient

	//
	docs.SwaggerInfo.Host = cfg.Swagger.Host
	docs.SwaggerInfo.Schemes = cfg.Swagger.Schemes

	// Set logger used for jaeger
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	zapLogger := logger.With(zap.String("service", "goldgym"))
	zlogger := jaegerLog.NewFactory(zapLogger)
	// loggers := Log.NewLogrusLogger()
	// Set tracer for service
	tracer, closer := tracing.Init("goldgym", zlogger)
	defer closer.Close()

	// httpc := httpclient.NewClient(tracer)
	// ad := auth.New(httpc, cfg.API.Auth)

	sdrd := goldgymRedisData.New(rdb, tracer, zlogger)

	sdst := goldgymStockData.New(db, dbr, nil, nil, nil, tracer, zlogger)
	ssst := goldgymStockService.New(sdst, tracer, zlogger)
	shst := goldgymStockHandler.New(ssst, tracer, zlogger)

	// Diganti dengan domain yang anda buat
	sd := goldgymData.New(db, dbr, tracer, zlogger)
	// ss := goldgymService.New(sd, ad, tracer, zlogger)
	ss := goldgymService.New(sd, sdrd, tracer, zlogger)

	sdsls := goldgymSaleData.New(db, dbr, nil, nil, nil, tracer, zlogger)
	sssls := goldgymSaleService.New(sdsls, tracer, zlogger)

	// // kafka new
	// ----- Kafka HTTP Producer + Consumer Worker -----
	// kafkaProd := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topics.HTTPInsert)
	// Sesudah (baru) — multi topic
	// kafkaProd := kafka.NewProducer(cfg.Kafka.Brokers)
	// defer kafkaProd.Close()

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	// kafkaWorker := kafkagoldgymWorker.New(ss, sssls, cfg.Kafka.Brokers, cfg.Kafka.Topics.HTTPInsert, cfg.Kafka.GroupID+"-http", zlogger)
	// go kafkaWorker.Start(ctx)
	// -------------------------------------------------------
	// Kafka
	// -------------------------------------------------------

	// 1. Pastikan semua topic sudah ada sebelum producer/consumer start
	topics := []string{
		cfg.Kafka.Topics.HTTPInsert,
		cfg.Kafka.Topics.Sales,
	}
	if err := kafka.EnsureTopics(cfg.Kafka.Brokers, topics); err != nil {
		log.Printf("[KAFKA] warning EnsureTopics: %v", err)
	}

	// 2. Init producer (multi topic)
	kafkaProd := kafka.NewProducer(cfg.Kafka.Brokers)
	defer kafkaProd.Close()

	// 3. Context untuk graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// customer service dibuat di sini juga karena worker Kafka membutuhkannya
	// untuk memproses bulk insert customer (action insert_customers_bulk)
	sdct := goldgymCustomerData.New(db, dbr, nil, nil, nil, tracer, zlogger)
	ssct := goldgymCustomerService.New(sdct, tracer, zlogger)

	// 4. Init worker — 1 worker, spawn 1 goroutine per topic otomatis
	kafkaWorker := kafkagoldgymWorker.New(
		ss,
		sssls,
		ssct,
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		topics,
		zlogger,
	)
	go kafkaWorker.Start(ctx)

	// -------------------------------------------------------
	// -------------------------------------------------

	// sh := goldgymHandler.New(ss, ssst, kafkaProd, tracer, zlogger)
	sh := goldgymHandler.New(ss, kafkaProd, tracer, zlogger)

	// discount service dikonstruksi di sini (sebelum sales handler) karena
	// sales handler butuh ssdc.RedeemVoucher untuk validasi voucher SEBELUM
	// insert sales diantre ke Kafka (lihat insert_gold_gym_sales_gin.go).
	sddc := goldgymDiscountData.New(db, tracer, zlogger)
	ssdc := goldgymDiscountService.New(sddc, tracer, zlogger)
	shdc := goldgymDiscountHandler.New(ssdc, tracer, zlogger)

	shsls := goldgymSaleHandler.New(sssls, ssdc, kafkaProd, tracer, zlogger)

	sdbk := goldgymBookingData.New(db, tracer, zlogger)
	ssbk := goldgymBookingService.New(sdbk, tracer, zlogger)
	shbk := goldgymBookingHandler.New(ssbk, kafkaProd, tracer, zlogger)

	sdor := goldgymOrderData.New(db, tracer, zlogger)
	ssor := goldgymOrderService.New(sdor, tracer, zlogger)
	shor := goldgymOrderHandler.New(ssor, kafkaProd, tracer, zlogger)

	sdim := goldgymItemsData.New(db, dbr, tracer, zlogger)
	ssim := goldgymItemsService.New(sdim, sd, sdrd, tracer, zlogger)
	shim := goldgymItemsHandler.New(ssim, tracer, zlogger)

	sdot := goldgymOutletData.New(db, dbr, tracer, zlogger)
	ssot := goldgymOutletService.New(sdot, tracer, zlogger)
	shot := goldgymOutletHandler.New(ssot, tracer, zlogger)

	shct := goldgymCustomerHandler.New(ssct, kafkaProd, tracer, zlogger)

	sdctp := goldgymCustomerTypeData.New(db, dbr, nil, nil, nil, tracer, zlogger)
	ssctp := goldgymCustomerTypeService.New(sdctp, tracer, zlogger)
	shctp := goldgymCustomerTypeHandler.New(ssctp, tracer, zlogger)

	// Elasticsearch
	esClient, err := es.NewClient(es.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	})
	if err != nil {
		log.Fatalf("[ES] Failed to create Elasticsearch client: %v", err)
	}
	sed := elasticData.New(esClient, tracer, zlogger)
	ses := elasticService.New(sed, tracer, zlogger)
	seh := elasticHandler.New(ses, tracer, zlogger)

	// // MongoDB
	// // mongodb
	// // --------------------------------------------------------------------------------
	// mongoClient, err := openMongoClient(cfg.MongoDB.URI)
	// if err != nil {
	// 	log.Fatalf("[MONGO] Failed to connect to MongoDB: %v", err)
	// }
	// defer func() {
	// 	if err := mongoClient.Disconnect(context.Background()); err != nil {
	// 		log.Printf("[MONGO] Disconnect error: %v", err)
	// 	}
	// }()
	// mongoDB := mongoClient.Database(cfg.MongoDB.Database)
	// gymactivityRepo := gymactivityData.New(mongoDB, tracer, zlogger)
	// gymactivitySvc := gymactivityService.New(gymactivityRepo, tracer, zlogger)
	// gymactivityH := gymactivityHandler.New(gymactivitySvc, tracer, zlogger)
	// // --------------------------------------------------------------------------------

	//middleware
	ms := middlewareService.New(sd, tracer, zlogger)
	mh := middlewareHandler.New(ms, ss, ssst, tracer, zlogger)
	middlewareHandler.InitDiscordAlert(db, cfg.Discord.Webhook)

	hh := healthHandler.New(db)

	// sdprod := goldgymData.New(dbprod, tracer, zlogger)
	// ssprod := goldgymService.New(sdprod, tracer, zlogger)

	sa := authService.New(sdrd, tracer, zlogger)
	sha := authHandler.New(ss, sa, tracer, zlogger)
	gqlH := graphqlHandler.New(ss, tracer, zlogger)
	// sh := goldgymHandler.New(ss, tracer, zlogger)

	// gRPC handler
	grpcHandler := goldgymGrpcHandler.NewHandler(ss, tracer, zlogger)

	// sdpn := pushNotifData.New(fcmB2BPelapak, loggers)
	// sspn := pushNotifService.New(sdpn, t.Tracer, loggers)
	// spnh := pushNotifHandler.New(sspn, loggers)

	// // // // ----- kafka -----
	// res := &resources.BootResources{
	// 	DBLocal:      db,
	// 	DBProd:       dbprod,
	// 	Redis:        rdb,
	// 	GoldSvcLocal: ss,
	// 	GoldSvcProd:  ssprod,
	// 	Tracer:       tracer,
	// 	Logger:       logger,
	// }
	// reg := registry.New(res)
	// go StartKafkaConsumers(cfg.Kafka, reg)
	// // // ----------------------------------------------------------------------------------------------------

	// // kafka new ( di comment jika kafka new jalan )
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()
	// --------------------------------------------------------------------------------------

	// config.PrepareWatchPath()
	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	err := config.Init()
	// 	if err != nil {
	// 		log.Printf("[VIPER] Error get config file, %v", err)
	// 	}
	// 	// firebase
	// 	cfg, _ = config.Get()

	// 	// reload local db
	// 	newDB, newDBR, err := openDatabases(cfg)
	// 	if err != nil {
	// 		log.Fatalf("[DB] Failed to initialize local database: %v", err)
	// 	} else {
	// 		db = newDB
	// 		dbr = newDBR
	// 		log.Println("[DB] local db reloaded")
	// 	}

	// 	// // reload prod db
	// 	// masterNewProduction, masterNewProductionR, err := openDatabasesProd(cfg)
	// 	// if err != nil {
	// 	// 	log.Fatalf("[DB] Failed to initialize production database: %v", err)
	// 	// } else {
	// 	// 	dbprod = masterNewProduction
	// 	// 	dbrprod = masterNewProductionR
	// 	// 	log.Println("[DB] prod db reloaded")
	// 	// }
	// })

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("[VIPER] config changed")
		log.Println("[VIPER] please restart service to apply DB changes")
	})

	defer func() {
		sd.Close()

		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// // // ----- kafka -----
	// // // prepare BootResources
	// // res := &BootResources{
	// // 	DBLocal:      db,
	// // 	DBProd:       dbprod,
	// // 	Redis:        rdb,
	// // 	GoldSvcLocal: ss,
	// // 	GoldSvcProd:  ssprod,
	// // }

	// // // // START Kafka Consumers (background)
	// // // StartKafkaConsumers(res)
	// // // // ----- kafka -----

	// // // daftarkan ke registry
	// // registry := consumer.NewRegistry(goldConsumer, stockConsumer, authConsumer)

	// // // jalankan consumer loop (background worker)
	// // go consumer.StartKafkaConsumers(cfg.Kafka, registry)

	// // potongan terakhir di func HTTP()
	// res := &resources.BootResources{
	// 	DBLocal:      db,
	// 	DBProd:       dbprod,
	// 	Redis:        rdb,
	// 	GoldSvcLocal: ss,
	// 	GoldSvcProd:  ssprod,
	// 	Tracer:       tracer,
	// 	Logger:       logger,
	// }
	// reg := registry.New(res)
	// fmt.Printf("test res %v", res)
	// go StartKafkaConsumers(cfg.Kafka, reg)
	// // // ----------------------------------------------------------------------------------------------------

	s := goldgymServer.Server{
		Goldgym:             sh,
		GoldgymStock:        shst,
		GoldgymItems:        shim,
		GoldgymDiscount:     shdc,
		GoldgymOutlet:       shot,
		GoldgymCustomer:     shct,
		GoldgymCustomerType: shctp,
		GoldgymSale:         shsls,
		GoldgymBooking:      shbk,
		GoldgymOrder:        shor,
		Auth:                sha,
		Middleware:          mh,
		Health:              hh,
		Logger:              zlogger,
		// mongodb
		// GymActivity: gymactivityH,
		Elastic: seh,
		GraphQL: gqlH,
		Config:  cfg,
		// PushNotification: spnh,
	}

	//force shutdown db connection when server stopped
	go func() {
		log.Printf("[HTTP] Starting HTTP server on port %s", cfg.Server.Port)
		if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
			// return err
			log.Fatalf("[HTTP] serve error: %v", err)
		}
	}()

	// Start gRPC server
	grpcLis, err := net.Listen("tcp", ":"+cfg.Server.GrpcPort)
	if err != nil {
		log.Fatalf("[GRPC] Failed to listen on port %s: %v", cfg.Server.GrpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoldGymServiceServer(grpcServer, grpcHandler)

	go func() {
		log.Printf("[GRPC] Starting gRPC server on port %s", cfg.Server.GrpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("[GRPC] serve error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down servers...")

	// Graceful shutdown for both HTTP and gRPC
	_ = s.Shutdown(context.Background())
	grpcServer.GracefulStop()
	log.Println("servers shutdown complete")

	return nil
}

func openDatabases(cfg *config.Config) (master *gorm.DB, masterDB *sqlx.DB, err error) {
	master, masterDB, err = openConnectionPool(cfg.Database.Master)
	if err != nil {
		return nil, nil, err
	}

	return master, masterDB, err
}

func openDatabasesProd(cfg *config.Config) (master *gorm.DB, masterDB *sqlx.DB, err error) {
	master, masterDB, err = openConnectionPool(cfg.Database.Master)
	if err != nil {
		return nil, nil, err
	}

	return master, masterDB, err
}

func openConnectionPool(connString string) (*gorm.DB, *sqlx.DB, error) {
	// Add MySQL connection parameters to improve connection handling
	// parseTime=true: Parse TIME/DATE/DATETIME to time.Time
	// loc=Local: Use local timezone
	// charset=utf8mb4: Use UTF8MB4 charset
	// timeout=10s: Connection timeout
	// readTimeout=30s: Read timeout
	// writeTimeout=30s: Write timeout
	// connString += "?parseTime=true&loc=Asia%2FJakarta&charset=utf8mb4&timeout=10s&readTimeout=30s&writeTimeout=30s"

	// connString += "?loc=Asia%2FJakarta&parseTime=true"
	connString += "?parseTime=true&loc=Asia%2FJakarta&charset=utf8mb4&timeout=10s&readTimeout=30s&writeTimeout=30s&time_zone=%27%2B07:00%27"

	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	// Connection pool settings optimized to prevent "connection reset by peer"
	sqlDB.SetMaxOpenConns(25)                 // Max connections in pool
	sqlDB.SetMaxIdleConns(10)                 // Max idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Max lifetime (reduced from 30min to 5min)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Max idle time (reduced from 10min to 5min)

	// Ping to validate initial connection
	if err = sqlDB.Ping(); err != nil {
		return nil, nil, err
	}

	sqlxDB := sqlx.NewDb(sqlDB, "mysql")

	return db, sqlxDB, err
}

// func openConnectionPool(driver string, connString string) (db *sqlx.DB, err error) {
// 	db, err = sqlx.Open(driver, connString)
// 	if err != nil {
// 		return db, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		return db, err
// 	}

// 	return db, err
// }

// openMongoClient creates a MongoDB client and verifies the connection with a ping
// mongodb
// -------------------------------------------------------------------------------------
// func openMongoClient(uri string) (*mongo.Client, error) {
// 	opts := mongoOptions.Client().ApplyURI(uri)
// 	client, err := mongo.Connect(opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := client.Ping(ctx, nil); err != nil {
// 		return nil, err
// 	}

// 	log.Println("[MONGO] Connected to MongoDB")
// 	return client, nil
// }
// -------------------------------------------------------------------------------------

func newRedisClient(cred config.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cred.Host,
		Password: cred.Password,
		DB:       0,
	})
	return client
}

func openFirebaseClient(ctx context.Context, cfg config.FirebaseConfig, cred map[string]string) (*firebase.App, error) {
	credBytes, err := json.Marshal(cred)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(credBytes)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID:     cfg.ProjectID,
		DatabaseURL:   cfg.DatabaseURL,
		StorageBucket: cfg.StorageBucket,
	}, opt)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func openFirestoreClient(ctx context.Context, app *firebase.App) (*firestore.Client, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func openFirebaseDatabaseClient(ctx context.Context, app *firebase.App) (*db.Client, error) {
	client, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
