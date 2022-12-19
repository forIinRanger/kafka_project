package adapters

import (
	"context"
	"github.com/IBM/sarama"
	"grpcservice/internal/app"
	"log"
	"runtime/debug"
)

const PartitionName = "my_topic"

type KafkaConsumer struct {
	ctx      context.Context
	Consumer sarama.PartitionConsumer
	service  *app.GrpcService
}

//nolint:all
func NewKafkaConsumer(ctx context.Context, addrs []string, serv *app.GrpcService) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	cons, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}
	partCons, err := cons.ConsumePartition(PartitionName, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	ctx1, _ := context.WithCancel(ctx)
	return &KafkaConsumer{
		ctx:      ctx1,
		Consumer: partCons,
		service:  serv,
	}, nil
}

func (kc *KafkaConsumer) StartCatching() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Kafka Consumer panic: %v", r)
			log.Printf("Stack:\n%s", debug.Stack())
		}
	}()
	for {
		select {
		case <-kc.ctx.Done():
			return kc.ctx.Err()
		case msg, ok := <-kc.Consumer.Messages():
			if !ok {
				log.Println("Kafka Consumer: messages channel closed")
				return nil
			}
			log.Println("Kafka Consumer: message received")
			log.Printf("Msg.Value: '%v'", string(msg.Value))
			err := kc.service.ProcessMessage(kc.ctx, string(msg.Value))
			if err != nil {
				log.Printf("Processing error: %v", err)
			}
		case err, ok := <-kc.Consumer.Errors():
			if !ok {
				log.Println("Kafka Consumer: errors channel closed")
				return nil
			}
			if err != nil {
				log.Printf("Error got from consumer error's channel: %v", err)
			}
		}

	}
}

func (kc *KafkaConsumer) StopCatching() error {
	return kc.Consumer.Close()
}
