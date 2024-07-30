package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type KafkaService interface {
	Producer(ctx context.Context, topic string, message interface{}) error
}

type service struct {
	kafkaURL      string
	kafkaUsername string
	kafkaPassword string
	schemaURL     string
}

func NewKafkaService(kafkaURL, kafkaUsername, kafkaPassword, schemaURL string) KafkaService {
	return &service{
		kafkaURL:      kafkaURL,
		kafkaUsername: kafkaUsername,
		kafkaPassword: kafkaPassword,
		schemaURL:     schemaURL,
	}
}

func (s *service) Producer(ctx context.Context, topic string, message interface{}) error {
	// Create Kafka Producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": s.kafkaURL,
		"sasl.mechanism":    "SCRAM-SHA-256",
		"security.protocol": "SASL_SSL",
		"sasl.username":     s.kafkaUsername,
		"sasl.password":     s.kafkaPassword,
	})
	if err != nil {
		return err
	}
	defer p.Close()

	// Create Schema Registry Client
	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		s.schemaURL,
		s.kafkaUsername,
		s.kafkaPassword))
	if err != nil {
		return err
	}

	// Create Avro Serializer
	ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		return err
	}

	// Serialize the payload
	payload, err := ser.Serialize(topic, message)
	if err != nil {
		return err
	}

	// Produce the message
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("test"), // You might want to modify this key as needed
		Value:          payload,
	}, deliveryChan)

	if err != nil {
		return err
	}

	// Wait for delivery confirmation
	ev := <-deliveryChan
	m := ev.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	log.Printf("Message produced to topic %s [%d] at offset %d\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	return nil
}
