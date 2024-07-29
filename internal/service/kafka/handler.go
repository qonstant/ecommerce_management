package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaService interface {
	Producer(ctx context.Context, topic string, message interface{}) error
}

type service struct {
	kafkaURL      string
	kafkaUsername string
	kafkaPassword string
}

func NewKafkaService(kafkaURL, kafkaUsername, kafkaPassword string) KafkaService {
	return &service{
		kafkaURL:      kafkaURL,
		kafkaUsername: kafkaUsername,
		kafkaPassword: kafkaPassword,
	}
}

func (s *service) Producer(ctx context.Context, topic string, message interface{}) error {
	mechanism, err := scram.Mechanism(scram.SHA256, s.kafkaUsername, s.kafkaPassword)
	if err != nil {
		return err
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP(s.kafkaURL),
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
