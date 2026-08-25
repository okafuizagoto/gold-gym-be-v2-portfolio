package kafkagoldgym

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	goldCustomerEntity "gold-gym-be/internal/entity/customer"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldSaleEntity "gold-gym-be/internal/entity/sales"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/segmentio/kafka-go"
)

// RepoService adalah interface ke data layer yang digunakan oleh worker.
// Worker memanggil data layer langsung (bypass service) untuk menghindari
// pengiriman email OTP yang ada di service.InsertGoldUser.
type RepoService interface {
	// InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (string, error)
	InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error)
}

type RepoSales interface {
	InsertSales(ctx context.Context, goldid int, sales goldSaleEntity.InsertSales) (string, error)
}

type RepoCustomer interface {
	InsertCustomer(ctx context.Context, id int, items []goldCustomerEntity.Customer) (string, error)
}

type KafkaMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// Worker adalah Kafka consumer untuk memproses pesan yang berasal dari HTTP request.
//
//	type Worker struct {
//		repo   RepoService
//		sales  RepoSales
//		reader *kafka.Reader
//		logger jaegerLog.Factory
//	}
type Worker struct {
	repo     RepoService
	sales    RepoSales
	customer RepoCustomer
	brokers  []string
	groupID  string
	topics   []string
	logger   jaegerLog.Factory
}

// New membuat Worker baru.
// brokers: list Kafka broker address
// topic: nama topic yang di-consume
// groupID: consumer group ID
func New(repo RepoService, sales RepoSales, customer RepoCustomer, brokers []string, groupID string, topics []string, logger jaegerLog.Factory) *Worker {
	// r := kafka.NewReader(kafka.ReaderConfig{
	// 	Brokers:     brokers,
	// 	Topic:       topic,
	// 	GroupID:     groupID,
	// 	StartOffset: kafka.FirstOffset,
	// 	MinBytes:    1,
	// 	MaxBytes:    10e6, // 10MB
	// })
	// return &Worker{
	// 	repo:   repo,
	// 	sales:  sales,
	// 	reader: r,
	// 	logger: logger,
	// }
	return &Worker{
		repo:     repo,
		sales:    sales,
		customer: customer,
		brokers:  brokers,
		groupID:  groupID,
		topics:  topics,
		logger:  logger,
	}
}

// Start menjalankan loop consumer di goroutine. Berhenti saat ctx.Done().
func (w *Worker) Start(ctx context.Context) {
	// log.Println("[KafkaWorker] started, waiting for messages...")
	// defer func() {
	// 	if err := w.reader.Close(); err != nil {
	// 		log.Printf("[KafkaWorker] reader close error: %v", err)
	// 	}
	// 	log.Println("[KafkaWorker] stopped")
	// }()

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 		msg, err := w.reader.FetchMessage(ctx)
	// 		if err != nil {
	// 			if ctx.Err() != nil {
	// 				return
	// 			}
	// 			log.Printf("[KafkaWorker] fetch error: %v", err)
	// 			time.Sleep(time.Second)
	// 			continue
	// 		}

	// 		if err := w.processMessage(ctx, msg); err != nil {
	// 			log.Printf("[KafkaWorker] process error (nack): %v", err)
	// 			// Tidak requeue — commit tetap agar tidak loop terus untuk pesan rusak
	// 		}

	// 		if err := w.reader.CommitMessages(ctx, msg); err != nil {
	// 			log.Printf("[KafkaWorker] commit error: %v", err)
	// 		}
	// 	}
	// }
	log.Printf("[KafkaWorker] starting %d topic consumer(s): %v", len(w.topics), w.topics)

	var wg sync.WaitGroup
	for _, topic := range w.topics {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			w.consumeTopic(ctx, t)
		}(topic)
	}

	wg.Wait()
	log.Println("[KafkaWorker] all consumers stopped")
}

func (w *Worker) consumeTopic(ctx context.Context, topic string) {
	// Recovery di level goroutine worker
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[KafkaWorker] [topic=%s] PANIC recovered: %v", topic, r)
			// restart consumer topic ini
			go w.consumeTopic(ctx, topic)
		}
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     w.brokers,
		Topic:       topic,
		GroupID:     w.groupID + "-" + topic,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("[KafkaWorker] [topic=%s] reader close error: %v", topic, err)
		}
		log.Printf("[KafkaWorker] [topic=%s] stopped", topic)
	}()

	log.Printf("[KafkaWorker] [topic=%s] listening...", topic)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("[KafkaWorker] [topic=%s] fetch error: %v", topic, err)
				time.Sleep(time.Second)
				continue
			}

			if err := w.routeMessage(ctx, topic, msg); err != nil {
				log.Printf("[KafkaWorker] [topic=%s] process error: %v", topic, err)
				// tetap commit agar tidak loop untuk pesan rusak
			}

			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("[KafkaWorker] [topic=%s] commit error: %v", topic, err)
			}
		}
	}
}

// ----------------------------------------
// routeMessage — routing by topic lalu action
// ----------------------------------------

func (w *Worker) routeMessage(ctx context.Context, topic string, msg kafka.Message) (err error) {
	// Recovery di level per message
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[KafkaWorker] [topic=%s] PANIC in processMessage: %v", topic, r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	var km KafkaMessage
	if err := json.Unmarshal(msg.Value, &km); err != nil {
		log.Printf("[KafkaWorker] [topic=%s] [action=unknown] unmarshal error: %v, raw: %s", topic, err, string(msg.Value))
		return nil // skip pesan rusak
	}

	log.Printf("[KafkaWorker] [topic=%s] [action=%s] received", topic, km.Action)

	// Routing berdasarkan action di dalam pesan, bukan nama topic. Nama topic
	// berbeda per environment (staging memakai prefix "staging." dan suffix
	// "_insert"), sehingga match nama topic persis membuat pesan di-skip
	// diam-diam — nota penjualan tidak pernah tersimpan ke database.
	switch km.Action {
	case "insert_user":
		return w.processInsertUser(ctx, km.Payload)
	case "insert_sales":
		return w.processInsertSales(ctx, km.Payload)
	case "insert_customers_bulk":
		return w.processInsertCustomersBulk(ctx, km.Payload)
	default:
		log.Printf("[KafkaWorker] [topic=%s] [action=%s] unknown action, skipping", topic, km.Action)
		return nil
	}
}

// func (w *Worker) processMessage(ctx context.Context, msg kafka.Message) error {
// 	// var user goldEntity.GetGoldUsers
// 	// if err := json.Unmarshal(msg.Value, &user); err != nil {
// 	// 	log.Printf("[KafkaWorker] unmarshal error: %v, raw: %s", err, string(msg.Value))
// 	// 	return err
// 	// }

// 	// result, err := w.repo.InsertGoldUser(ctx, user)
// 	// if err != nil {
// 	// 	log.Printf("[KafkaWorker] InsertGoldUser error: %v", err)
// 	// 	return err
// 	// }

// 	// log.Printf("[KafkaWorker] InsertGoldUser success, result: %s", result)
// 	// return nil
// 	// Validasi topic → panggil function yang sesuai
// 	switch msg.Topic {
// 	case "goldgym.http.insert_user":
// 		return w.processInsertUser(ctx, msg)
// 	case "goldgym.http.sales_insert":
// 		return w.processInsertSales(ctx, msg)
// 	// case "gold-gym.stock.insert":
// 	//     return w.processInsertStock(ctx, msg)
// 	// case "gold-gym.customer.insert":
// 	//     return w.processInsertCustomer(ctx, msg)
// 	default:
// 		log.Printf("[KafkaWorker] unknown topic: %s", msg.Topic)
// 		return nil // skip, jangan return error agar tidak loop
// 	}
// }

// ----------------------------------------
// Processor — user
// ----------------------------------------

func (w *Worker) processInsertUser(ctx context.Context, payload json.RawMessage) error {
	log.Printf("[KafkaWorker] [topic=goldgym.user] [action=insert_user] processing...")

	var user goldEntity.GetGoldUsers
	if err := json.Unmarshal(payload, &user); err != nil {
		log.Printf("[KafkaWorker] [topic=goldgym.user] [action=insert_user] unmarshal error: %v", err)
		return err
	}

	result, err := w.repo.InsertGoldUser(ctx, user)
	if err != nil {
		log.Printf("[KafkaWorker] [topic=goldgym.user] [action=insert_user] error: %v", err)
		return err
	}

	log.Printf("[KafkaWorker] [topic=goldgym.user] [action=insert_user] success: %v", result)
	return nil
}

// ----------------------------------------
// Processor — sales
// ----------------------------------------

func (w *Worker) processInsertSales(ctx context.Context, payload json.RawMessage) error {
	log.Printf("[KafkaWorker] [topic=goldgym.sales] [action=insert_sales] processing...")

	var sale goldSaleEntity.InsertSales
	if err := json.Unmarshal(payload, &sale); err != nil {
		log.Printf("[KafkaWorker] [topic=goldgym.sales] [action=insert_sales] unmarshal error: %v", err)
		return err
	}

	result, err := w.sales.InsertSales(ctx, sale.Header.SaleGoldID, sale)
	if err != nil {
		log.Printf("[KafkaWorker] [topic=goldgym.sales] [action=insert_sales] error: %v", err)
		return err
	}

	log.Printf("[KafkaWorker] [topic=goldgym.sales] [action=insert_sales] success: %v", result)
	return nil
}

// ----------------------------------------
// Processor — customer (bulk insert)
// ----------------------------------------

func (w *Worker) processInsertCustomersBulk(ctx context.Context, payload json.RawMessage) error {
	log.Printf("[KafkaWorker] [action=insert_customers_bulk] processing...")

	var bulk goldCustomerEntity.BulkInsertCustomer
	if err := json.Unmarshal(payload, &bulk); err != nil {
		log.Printf("[KafkaWorker] [action=insert_customers_bulk] unmarshal error: %v", err)
		return err
	}
	if len(bulk.Items) == 0 {
		return nil
	}

	result, err := w.customer.InsertCustomer(ctx, bulk.GoldID, bulk.Items)
	if err != nil {
		log.Printf("[KafkaWorker] [action=insert_customers_bulk] error: %v", err)
		return err
	}

	log.Printf("[KafkaWorker] [action=insert_customers_bulk] success (%d items): %v", len(bulk.Items), result)
	return nil
}
