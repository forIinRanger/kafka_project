package adapters

import (
	"context"
	"github.com/IBM/sarama"
	"log"

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
	go func() {
		for {
			select {
			case _, ok := <-prod.Successes():
				if !ok {
					return
				}
				log.Printf("Message sent successfully")
			case _, ok := <-prod.Errors():
				if !ok {
					return
				}
				log.Printf("Error sending message")
			}
		}
	}()
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
	log.Println("Kafka Producer: message sent")
	return nil
}
