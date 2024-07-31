package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

// Credentials for Kafka service
type Credentials struct {
	KafkaURL      string
	KafkaUsername string
	KafkaPassword string
	SchemaURL     string
}

type KafkaService interface {
	Producer(ctx context.Context, topic string, message interface{}) error
}

type service struct {
	credentials Credentials
}

func NewKafkaService(credentials Credentials) KafkaService {
	return &service{
		credentials: credentials,
	}
}

func (s *service) Producer(ctx context.Context, topic string, message interface{}) error {
	// Create Kafka Producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": s.credentials.KafkaURL,
		"sasl.mechanism":    "SCRAM-SHA-256",
		"security.protocol": "SASL_SSL",
		"sasl.username":     s.credentials.KafkaUsername,
		"sasl.password":     s.credentials.KafkaPassword,
	})
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		return err
	}
	defer p.Close()

	// Create Schema Registry Client
	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		s.credentials.SchemaURL,
		s.credentials.KafkaUsername,
		s.credentials.KafkaPassword))
	if err != nil {
		log.Printf("Failed to create Schema Registry client: %v", err)
		return err
	}

	// Create Avro Serializer
	ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		log.Printf("Failed to create Avro serializer: %v", err)
		return err
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	// Serialize the payload
	payload, err := ser.Serialize(topic, message)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return err
	}

	// Produce the message
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("test"),
		Value:          payload,
	}, deliveryChan)

	if err != nil {
		log.Printf("Failed to produce message: %v", err)
		return err
	}

	// Wait for delivery confirmation
	ev := <-deliveryChan
	m := ev.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.Printf("Failed to deliver message: %v", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}

	log.Printf("Message produced to topic %s [%d] at offset %d\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	return nil
}
