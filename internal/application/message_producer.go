package application

type MessageProducer interface {
	Produce(topic string, message []byte) error
}
