package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type service struct{}

type KafkaService interface {
	Producer(ctx context.Context, topic string, message interface{}) error
}

func NewKafkaService() KafkaService {
	return &service{}
}

func (s *service) Producer(ctx context.Context, topic string, message interface{}) error {
	mechanism, err := scram.Mechanism(scram.SHA256, os.Getenv("UPSTASH_KAFKA_REST_USERNAME"), os.Getenv("UPSTASH_KAFKA_REST_PASSWORD"))
	if err != nil {
		return err
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP(os.Getenv("UPSTASH_KAFKA_REST_URL")),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}
	defer writer.Close()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Value: messageBytes,
	})

	if err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}

	log.Println("Message successfully written to Kafka topic", topic)
	return nil
}
