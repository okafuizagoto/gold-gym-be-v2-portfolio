package kafka

import (
	"fmt"
	"log"
	"net"

	"github.com/segmentio/kafka-go"
)

func EnsureTopics(brokers []string, topics []string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("controller error: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(
		controller.Host,
		fmt.Sprintf("%d", controller.Port),
	))
	if err != nil {
		return fmt.Errorf("controller dial error: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, t := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             t,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
		// ignore jika topic sudah ada
		log.Printf("[KAFKA] EnsureTopics (mungkin sudah ada): %v", err)
	}

	log.Printf("[KAFKA] EnsureTopics OK: %v", topics)
	return nil
}
