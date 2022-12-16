package adapters

import (
	"context"
	"github.com/IBM/sarama"

	"httpservice/internal/app/sender"
)

type KafkaProducer struct {
	Producer sarama.AsyncProducer
}

func NewKafkaSender(addrs []string, config *sarama.Config) (KafkaProducer, error) {
	prod, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		return KafkaProducer{}, err
	}
	return KafkaProducer{Producer: prod}, nil
}

func (s *KafkaProducer) Send(ctx context.Context, m sender.MappedMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	msg := &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.StringEncoder(m.Value),
	}
	s.Producer.Input() <- msg
	return nil
}
