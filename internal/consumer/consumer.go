package consumer

import (
	"context"
	"encoding/json"
	"gold-gym-be/internal/registry"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type CDCEvent struct {
	Payload struct {
		Op     string                 `json:"op"`
		After  map[string]interface{} `json:"after"`
		Before map[string]interface{} `json:"before"`
		Source map[string]interface{} `json:"source"`
	} `json:"payload"`
}

func ConsumeLoop(brokers []string, topic, groupID string, reg *registry.Registry) {
	// topic = "mysql_server.u868654674_gold_gym_bez.data_peserta"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("kafka read error:", err)
			time.Sleep(time.Second)
			continue
		}
		var ev CDCEvent
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			log.Println("unmarshal CDC error:", err)
			continue
		}
		table, _ := ev.Payload.Source["table"].(string)
		if handler, ok := reg.GetHandler(table); ok {
			if err := handler(context.Background(), ev.Payload.Op, ev.Payload.After, ev.Payload.Before); err != nil {
				log.Printf("handler error table=%s op=%s: %v", table, ev.Payload.Op, err)
			}
		} else {
			log.Printf("no handler for table=%s op=%s", table, ev.Payload.Op)
		}
	}
}
