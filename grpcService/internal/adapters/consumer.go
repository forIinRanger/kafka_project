package adapters

import (
	"context"
	"github.com/IBM/sarama"
	"grpcservice/internal/app"
)

const PartitionName = "my_topic"

type KafkaConsumer struct {
	ctx      context.Context
	cancel   context.CancelFunc
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
	ctx1, cancel := context.WithCancel(ctx)
	return &KafkaConsumer{
		ctx:      ctx1,
		cancel:   cancel,
		Consumer: partCons,
		service:  serv,
	}, nil
}

func (kc *KafkaConsumer) StartCatching() error {
	for {
		select {
		case <-kc.ctx.Done():
			return kc.ctx.Err()
		case msg := <-kc.Consumer.Messages():
			err := kc.service.ProcessMessage(kc.ctx, string(msg.Value))
			if err != nil {
				return err
			}
		case err := <-kc.Consumer.Errors():
			if err != nil {
				return err
			}
		}

	}
}

func (kc *KafkaConsumer) StopCatching() error {
	kc.cancel()
	return kc.Consumer.Close()
}
