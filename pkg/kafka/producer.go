package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Message adalah envelope standar untuk semua pesan Kafka
type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// Producer wraps kafka.Writer untuk publish message ke satu topic.
type Producer struct {
	// writer *kafka.Writer
	brokers []string
	writers map[string]*kafka.Writer
}

// NewProducer membuat Producer baru yang terhubung ke brokers dan topic tertentu.
//
//	func NewProducer(brokers []string, topic string) *Producer {
//		w := &kafka.Writer{
//			Addr:                   kafka.TCP(brokers...),
//			Topic:                  topic,
//			Balancer:               &kafka.LeastBytes{},
//			AllowAutoTopicCreation: true,
//		}
//		return &Producer{writer: w}
//	}

// NewProducer membuat Producer baru dengan daftar brokers.
// Topic dibuat on-demand saat Publish dipanggil.
func NewProducer(brokers []string) *Producer {
	return &Producer{
		brokers: brokers,
		writers: make(map[string]*kafka.Writer),
	}
}

// getWriter mengembalikan writer untuk topic tertentu,
// membuat baru jika belum ada.
func (p *Producer) getWriter(topic string) *kafka.Writer {
	if w, ok := p.writers[topic]; ok {
		return w
	}
	w := &kafka.Writer{
		Addr:                   kafka.TCP(p.brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		MaxAttempts:            5,
		WriteTimeout:           10 * time.Second,
		ReadTimeout:            10 * time.Second,
		RequiredAcks:           kafka.RequireOne,
	}
	p.writers[topic] = w
	return w
}

// Publish melakukan JSON marshal pada body lalu mengirimnya ke topic tertentu.
func (p *Producer) Publish(ctx context.Context, topic string, action string, body interface{}) error {
	payloadBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal payload error: %w", err)
	}

	msg := Message{
		Action:  action,
		Payload: json.RawMessage(payloadBytes),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	// return p.getWriter(topic).WriteMessages(ctx, kafka.Message{
	// 	Value: data,
	// })
	// Retry with backoff
	maxRetry := 5
	for i := 0; i < maxRetry; i++ {
		err = p.getWriter(topic).WriteMessages(ctx, kafka.Message{
			Value: data,
		})
		if err == nil {
			return nil
		}
		log.Printf("[KAFKA] [topic=%s] [action=%s] attempt %d failed: %v", topic, action, i+1, err)
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}

	return fmt.Errorf("[KAFKA] publish failed after %d retries: %w", maxRetry, err)
}

// // Publish melakukan JSON marshal pada body lalu mengirimnya ke Kafka.
// func (p *Producer) Publish(ctx context.Context, body interface{}) error {
// 	data, err := json.Marshal(body)
// 	if err != nil {
// 		return err
// 	}
// 	return p.writer.WriteMessages(ctx, kafka.Message{
// 		Value: data,
// 	})
// }

// // Close menutup koneksi writer.
// func (p *Producer) Close() {
// 	_ = p.writer.Close()
// }

// Close menutup semua writer.
func (p *Producer) Close() {
	for _, w := range p.writers {
		_ = w.Close()
	}
}
