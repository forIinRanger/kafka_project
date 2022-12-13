package adapters

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"project/internal/app/sender"
)

type NaiveSenderImpl struct{}

func NewSenderNaive() NaiveSenderImpl {
	return NaiveSenderImpl{}
}
func (s NaiveSenderImpl) Send(ctx context.Context, m sender.MappedMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	fmt.Println(m.Value)
	return nil
}

type KafkaProducer struct {
	sender *kafka.Writer
}

func NewKafkaSender() KafkaProducer {
	kp := kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "my-topic",
	}
	return KafkaProducer{sender: &kp}
}

func (s *KafkaProducer) Send(ctx context.Context, m sender.MappedMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return s.sender.WriteMessages(ctx, kafka.Message{Value: []byte(m.Value)})
}
