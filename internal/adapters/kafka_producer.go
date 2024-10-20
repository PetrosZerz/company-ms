package adapters

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducerImpl struct {
	producer *kafka.Producer
}

func NewKafkaProducer(broker string) (*KafkaProducerImpl, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaProducerImpl{producer: p}, nil
}

func (kp *KafkaProducerImpl) Produce(topic string, message []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		return err
	}

	// Wait for the delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}

func (kp *KafkaProducerImpl) Close() {
	kp.producer.Close()
}
